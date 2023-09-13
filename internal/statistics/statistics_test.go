package statistics

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func setupLoggerToBuffer() *bytes.Buffer {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)

	return &buf
}

func TestNewPodStatistic(t *testing.T) {
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
	podStat := PodStatistic{}
	podStat.Initialize(pod)

	// Assertions
	assert.Equal(t, pod.Name, podStat.Name, "Pod name does not match")
	assert.Equal(t,
		pod.Namespace, podStat.Namespace, "Pod namespace does not match")
	assert.WithinDuration(t,
		pod.CreationTimestamp.Time, podStat.CreationTimestamp, time.Millisecond,
		"Pod creation timestamps do not match")
	assert.Empty(t, podStat.InitContainers, "Expected no init containers")
	assert.Empty(t, podStat.Containers, "Expected no containers")
}

type MockTimeSource struct {
	mockedTime time.Time
}

func (mts MockTimeSource) Now() time.Time {
	return mts.mockedTime
}

//nolint:funlen
func TestPodStatisticUpdate(t *testing.T) {
	// Redirect logger to buffer
	buf := setupLoggerToBuffer()

	// 1. Setup a sample pod
	format := "2006-01-02T15:04:05Z07:00"
	created, err := time.Parse(format, "2023-08-28T00:00:00Z")
	if err != nil {
		panic(err)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			CreationTimestamp: metav1.NewTime(created),
			Name:              "test-pod",
			Namespace:         "test-namespace",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "test-container"},
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

	// 2. Initialize podStatistic and update
	stat := PodStatistic{}
	stat.Initialize(pod)
	stat.TimeSource = MockTimeSource{created.Add(3 * time.Second)}
	stat.update(pod)

	// 3. Validate updated fields
	assert.NotZero(t, stat.ScheduledTimestamp, "scheduledTimestamp was not set")
	assert.NotZero(t, stat.InitializedTimestamp, "initializedTimestamp was not set")
	assert.NotEmpty(t, stat.Containers, "containers map was not populated")

	decoder := json.NewDecoder(buf)
	statistic_logs := make([]map[string]interface{}, 0)
	for {
		var document interface{}
		if err := decoder.Decode(&document); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			t.Errorf("Invalid JSON output")
		}

		if map_document, ok := document.(map[string]interface{}); ok {
			if _, ok := map_document["kube_transition_metric_type"]; ok {
				statistic_logs = append(statistic_logs, map_document)
			}
		} else {
			t.Errorf("Log document is not a map")
		}
	}

	assert.Len(t, statistic_logs, 2, "Not the correct number of statistic logs")

	shared_assertations := func(log map[string]interface{}) {
		assert.Equal(t, "test-namespace", log["kube_namespace"])
		assert.Equal(t, "test-pod", log["pod_name"])
	}

	shared_assertations(statistic_logs[0])
	assert.Equal(t,
		"pod", statistic_logs[0]["kube_transition_metric_type"],
		"first log metric is not of type pod")
	assert.IsType(t,
		make(map[string]interface{}),
		statistic_logs[0]["kube_transition_metrics"],
		"key kube_transition_metrics is not a JSON object")
	metrics, _ :=
		statistic_logs[0]["kube_transition_metrics"].(map[string]interface{})
	assert.InDelta(t,
		2*time.Second.Seconds(), metrics["initialized_latency"], 1e-5,
		"Initialized latency is not correct")
	assert.InDelta(t,
		time.Second.Seconds(), metrics["scheduled_latency"], 1e-5,
		"Scheduled latency is not correct")

	shared_assertations(statistic_logs[1])
	assert.Equal(t,
		"container", statistic_logs[1]["kube_transition_metric_type"],
		"second log metric is not of type container")
	assert.IsType(t,
		make(map[string]interface{}), statistic_logs[1]["kube_transition_metrics"],
		"key kube_transition_metrics is not a JSON object")
	metrics, _ =
		statistic_logs[1]["kube_transition_metrics"].(map[string]interface{})
	assert.Equal(t,
		false, metrics["init_container"], "Container should not be an init container")
	assert.InDelta(t,
		3*time.Second.Seconds(), metrics["ready_latency"], 1e-5,
		"Container ready latency is not correct")
}

func TestContainerStatisticUpdate(t *testing.T) {
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
	podStat := PodStatistic{}
	podStat.Initialize(pod)
	containerStat := podStat.Containers["test-container"]
	assert.True(t, containerStat.runningTimestamp.IsZero())
	assert.True(t, containerStat.startedTimestamp.IsZero())
	assert.True(t, containerStat.readyTimestamp.IsZero())
	assert.NotNil(t, status.State.Running)
	assert.NotNil(t, status.Started)
	assert.True(t, *status.Started)
	assert.True(t, status.Ready)
	now := time.Now()
	assert.NotZero(t, now)
	containerStat.update(now, status)

	// 3. Validate updated fields
	assert.NotZero(t,
		containerStat.runningTimestamp, "runningTimestamp was not set")
	assert.NotZero(t,
		containerStat.startedTimestamp, "startedTimestamp was not set")
	assert.NotZero(t,
		containerStat.readyTimestamp, "readyTimestamp was not set")
}
