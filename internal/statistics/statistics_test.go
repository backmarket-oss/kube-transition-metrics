package statistics

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func setupLoggerToBuffer() *bytes.Buffer {
	buf := &bytes.Buffer{}
	metricOutput = buf

	return buf
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
	podStat := podStatistic{}
	podStat.initialize(pod)

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

type MockTimeSource struct {
	mockedTime time.Time
}

func (mts MockTimeSource) Now() time.Time {
	return mts.mockedTime
}

func TestPodStatisticUpdate(t *testing.T) {
	zerolog.DurationFieldInteger = false
	zerolog.DurationFieldUnit = time.Second

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

	stat := podStatistic{}
	stat.initialize(pod)

	// Stub time source to fix time to constant for testing
	stat.timeSource = MockTimeSource{created.Add(3 * time.Second)}

	// Stub StatisticEventHandler
	eh := &StatisticEventHandler{
		options: &options.Options{},
	}
	stat.imagePullCollector = newImagePullCollector(eh, "", pod.UID)

	// Update the pod statistic for the "new" state
	stat.update(pod)

	assert.NotZero(t, stat.scheduledTimestamp, "scheduledTimestamp was not set")
	assert.NotZero(
		t, stat.initializedTimestamp, "initializedTimestamp was not set")
	assert.NotEmpty(t, stat.containers, "containers map was not populated")

	// Check that the imagePullCollector would have been canceled for the right
	// reasons upon pod initialization.
	select {
	case s := <-stat.imagePullCollector.cancelChan:
		assert.Equal(
			t, "pod_initialized", s,
			"ImagePullCollector cancel channel received erroneous cancel reason")
		assert.True(
			t, stat.imagePullCollector.canceled.Load(),
			"ImagePullCollector cancel chan written to without setting canceled true")
	case <-time.NewTimer(time.Second).C:
		assert.Fail(
			t, "ImagePullCollector cancel channel was not written to within 1 second")
	}

	decoder := json.NewDecoder(buf)
	statisticLogs := make([]map[string]interface{}, 0)
	for {
		var document interface{}
		if err := decoder.Decode(&document); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			t.Errorf("Invalid JSON output")
		}

		if mapDocument, ok := document.(map[string]interface{}); ok {
			if _, ok := mapDocument["kube_transition_metric_type"]; ok {
				statisticLogs = append(statisticLogs, mapDocument)
			}
		} else {
			t.Errorf("Log document is not a map")
		}
	}

	if !assert.Len(
		t, statisticLogs, 2, "Not the correct number of statistic logs") {
		return
	}

	sharedAssertations := func(log map[string]interface{}) {
		assert.Equal(t, "test-namespace", log["kube_namespace"])
		assert.Equal(t, "test-pod", log["pod_name"])
	}

	sharedAssertations(statisticLogs[0])

	assert.Equal(t,
		"pod", statisticLogs[0]["kube_transition_metric_type"],
		"first log metric is not of type pod")

	assert.IsType(t,
		make(map[string]interface{}),
		statisticLogs[0]["kube_transition_metrics"],
		"key kube_transition_metrics is not a JSON object")
	metrics, _ :=
		statisticLogs[0]["kube_transition_metrics"].(map[string]interface{})
	assert.InDelta(t,
		2*time.Second.Seconds(), metrics["creation_to_running_seconds"], 1e-5,
		"Initialized latency is not correct")
	assert.InDelta(t,
		time.Second.Seconds(), metrics["creation_to_initializing_seconds"], 1e-5,
		"Scheduled latency is not correct")

	sharedAssertations(statisticLogs[1])
	assert.Equal(t,
		"container", statisticLogs[1]["kube_transition_metric_type"],
		"second log metric is not of type container")
	assert.IsType(t,
		make(map[string]interface{}), statisticLogs[1]["kube_transition_metrics"],
		"key kube_transition_metrics is not a JSON object")
	metrics, _ =
		statisticLogs[1]["kube_transition_metrics"].(map[string]interface{})
	assert.Equal(t,
		false, metrics["init_container"], "Container should not be an init container")
	assert.InDelta(t,
		2*time.Second.Seconds(), metrics["initialized_to_running_seconds"], 1e-5,
		"Container runnning latency is not correct")
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
	podStat := podStatistic{}
	podStat.initialize(pod)
	containerStat := podStat.containers["test-container"]
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
