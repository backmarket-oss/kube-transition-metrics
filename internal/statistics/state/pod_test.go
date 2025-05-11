package state

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/logging"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// configureLogging sets up the global logging configuration for the tests.
// It forbids parallel execution of tests that use this function.
func configureLogging(t *testing.T) {
	t.Helper()

	t.Cleanup(func() {
		logging.Unconfigure()
	})
	logging.Configure()
}

func TestNewPodStatistic(t *testing.T) {
	configureLogging(t)

	// Define a test pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
		},
	}

	// Call the function
	podStat := NewPodStatistic(time.Now(), pod)

	// Assertions
	assert.Equal(t, pod.Name, podStat.name, "Pod name does not match")
	assert.Equal(t,
		pod.Namespace, podStat.namespace, "Pod namespace does not match")
	assert.WithinDuration(t,
		pod.CreationTimestamp.Time, podStat.creationTimestamp, time.Millisecond,
		"Pod creation timestamps do not match")
	assert.Empty(t, podStat.initContainers, "Expected no init containers")
	assert.Empty(t, podStat.containers, "Expected no containers")
}

func newTestingPod(created time.Time) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			CreationTimestamp: metav1.NewTime(created),
			Name:              "test-pod",
			Namespace:         "test-namespace",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "test-container",
					Image: "test-image",
				},
			},
		},
		Status: corev1.PodStatus{
			Conditions: []corev1.PodCondition{
				{
					Type:               corev1.PodScheduled,
					Status:             corev1.ConditionTrue,
					LastTransitionTime: metav1.NewTime(created.Add(time.Second)),
				},
				{
					Type:               corev1.PodInitialized,
					Status:             corev1.ConditionTrue,
					LastTransitionTime: metav1.NewTime(created.Add(2 * time.Second)),
				},
			},
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name: "test-container",
					State: corev1.ContainerState{
						Running: &corev1.ContainerStateRunning{},
					},
					Ready: true,
				},
			},
		},
	}
}

func checkBasicPodStatisticFields(t *testing.T, stat *PodStatistic) {
	t.Helper()
	assert.False(t, stat.scheduledTimestamp.IsZero(), "scheduledTimestamp was not set")
	assert.False(
		t, stat.initializedTimestamp.IsZero(), "initializedTimestamp was not set")
	assert.NotEmpty(t, stat.containers, "containers map was not populated")
}

func decodeMetrics(t *testing.T, buf *bytes.Buffer) []map[string]interface{} {
	t.Helper()
	decoder := json.NewDecoder(buf)
	statisticLogs := make([]map[string]interface{}, 0)
	for {
		var document interface{}
		if err := decoder.Decode(&document); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			t.Errorf("Invalid JSON output")
		}

		if !assert.IsType(t, make(map[string]interface{}), document, "Log document is not an object") {
			continue
		}
		mapDocument, _ := document.(map[string]interface{})
		if !assert.IsType(t, make(map[string]interface{}), mapDocument["kube_transition_metrics"],
			"kube_transition_metric key of log document is not an object") {
			continue
		}
		mapMetrics, _ := mapDocument["kube_transition_metrics"].(map[string]interface{})
		statisticLogs = append(statisticLogs, mapMetrics)
	}

	return statisticLogs
}

func TestPodStatisticUpdate(t *testing.T) {
	configureLogging(t)

	format := "2006-01-02T15:04:05Z07:00"
	created, err := time.Parse(format, "2023-08-28T00:00:00Z")
	if err != nil {
		panic(err)
	}
	pod := newTestingPod(created)

	now := pod.CreationTimestamp.Add(3 * time.Second)
	stat := NewPodStatistic(now, pod)

	buf := &bytes.Buffer{}

	// Update the pod statistic for the "new" state
	stat = stat.Update(now, pod)

	checkBasicPodStatisticFields(t, stat)

	stat.Report(buf, pod)
	statisticLogs := decodeMetrics(t, buf)

	if !assert.Len(
		t, statisticLogs, 2, "Not the correct number of statistic logs") {
		return
	}

	sharedAssertations := func(log map[string]interface{}) {
		assert.Equal(t, "test-namespace", log["kube_namespace"])
		assert.Equal(t, "test-pod", log["pod_name"])
	}

	metrics := statisticLogs[0]
	sharedAssertations(metrics)
	assert.Equal(t,
		"pod", metrics["type"],
		"first log metric is not of type pod")
	if assert.IsType(t, make(map[string]interface{}), metrics["pod"]) {
		podMetrics, _ := metrics["pod"].(map[string]interface{})
		assert.InDelta(t,
			2*time.Second.Seconds(), podMetrics["creation_to_initialized_seconds"], 1e-5,
			"Initialized latency is not correct")
		assert.InDelta(t,
			time.Second.Seconds(), podMetrics["creation_to_scheduled_seconds"], 1e-5,
			"Scheduled latency is not correct")
	}

	metrics = statisticLogs[1]
	sharedAssertations(metrics)
	assert.Equal(t,
		"container", metrics["type"],
		"second log metric is not of type container")
	if assert.IsType(t, make(map[string]interface{}), metrics["container"]) {
		containerMetrics, _ := metrics["container"].(map[string]interface{})
		assert.Equal(t,
			false, containerMetrics["init_container"], "Container should not be an init container")
		assert.InDelta(t,
			2*time.Second.Seconds(), containerMetrics["initialized_to_running_seconds"], 1e-5,
			"Container runnning latency is not correct")
	}
}

func TestContainerStatisticUpdate(t *testing.T) {
	configureLogging(t)

	// 1. Setup a sample container status
	status := corev1.ContainerStatus{
		Name: "test-container",
		State: corev1.ContainerState{
			Running: &corev1.ContainerStateRunning{},
		},
		Ready: true,
		Started: func() *bool {
			b := true

			return &b
		}(),
	}

	// 2. Initialize containerStatistic and update
	pod := &corev1.Pod{Spec: corev1.PodSpec{
		Containers: []corev1.Container{
			{Name: "test-container"},
		},
	}}
	now := time.Now()
	podStat := NewPodStatistic(now, pod)
	containerStat, ok := podStat.containers.Get("test-container")
	assert.True(t, ok)
	assert.True(t, containerStat.runningTimestamp.IsZero())
	assert.True(t, containerStat.startedTimestamp.IsZero())
	assert.True(t, containerStat.readyTimestamp.IsZero())
	assert.NotNil(t, status.State.Running)
	assert.NotNil(t, status.Started)
	assert.True(t, *status.Started)
	assert.True(t, status.Ready)
	assert.NotZero(t, now)
	containerStat = containerStat.Update(now, status, podStat)

	// 3. Validate updated fields
	assert.NotZero(t,
		containerStat.runningTimestamp, "runningTimestamp was not set")
	assert.NotZero(t,
		containerStat.startedTimestamp, "startedTimestamp was not set")
	assert.NotZero(t,
		containerStat.readyTimestamp, "readyTimestamp was not set")
}
