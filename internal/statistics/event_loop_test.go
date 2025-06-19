package statistics

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/statistics/state"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/testhelpers"
	"github.com/Izzette/go-safeconcurrency/eventloop"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
)

// newTestingPod creates a partial pod (missing PodReady condition) for testing.
func newTestingPod(created time.Time) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			UID:               "test-uid",
			Name:              "test-pod",
			Namespace:         "test-namespace",
			CreationTimestamp: metav1.NewTime(created),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "test-container", Image: "test-image"},
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
		},
	}
}

// newTestingCompletePod creates a pod with all conditions and container statuses set, so
// [state.PodStatistic.Partial] returns false.
func newTestingCompletePod(created time.Time) *corev1.Pod {
	started := true

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			UID:               "test-uid",
			Name:              "test-pod",
			Namespace:         "test-namespace",
			CreationTimestamp: metav1.NewTime(created),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "test-container", Image: "test-image"},
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
				{
					Type:               corev1.PodReady,
					Status:             corev1.ConditionTrue,
					LastTransitionTime: metav1.NewTime(created.Add(3 * time.Second)),
				},
			},
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name: "test-container",
					State: corev1.ContainerState{
						Running: &corev1.ContainerStateRunning{},
					},
					Started: &started,
					Ready:   true,
				},
			},
		},
	}
}

func TestNewStatisticEventLoop(t *testing.T) {
	ctx := t.Context()
	opts := &options.Options{
		StatisticEventQueueLength: 1,
		LogLevel:                  zerolog.FatalLevel,
	}
	testhelpers.ConfigureLogging(t, opts)

	statisticEventLoop := NewStatisticEventLoop(opts, io.Discard)
	defer statisticEventLoop.Close()

	statisticEventLoop.Start()

	gen, err := statisticEventLoop.PodResync(ctx, []apimachinerytypes.UID{"test-uid"})
	require.NoError(t, err, "Expected no error when sending resync event")
	state, err := eventloop.WaitForGeneration(ctx, statisticEventLoop, gen)
	require.NoError(t, err, "Expected no error when waiting for resync event to be processed")
	assert.NotNil(t, state, "Expected state to be not nil")

	assert.True(t, state.State().IsBlacklisted("test-uid"), "Expected test-uid to be blacklisted")
	testUIDStatistic, ok := state.State().Get("test-uid")
	assert.Nil(t, testUIDStatistic, "Expected test-uid statistic to be nil")
	assert.False(t, ok, "Expected test-uid statistic to not exist")
	assert.Zero(t, state.State().Len(), "Expected number of tracked statistics to be 0")
}

func TestNewImagePullStatisticEventLoop(t *testing.T) {
	ctx := t.Context()
	opts := &options.Options{
		StatisticEventQueueLength: 1,
		LogLevel:                  zerolog.FatalLevel,
	}
	testhelpers.ConfigureLogging(t, opts)

	imagePullEventLoop := NewImagePullStatisticEventLoop(opts, io.Discard)
	defer imagePullEventLoop.Close()

	imagePullEventLoop.Start()

	pod := newTestingPod(time.Now())

	gen, err := imagePullEventLoop.ImagePullDelete(ctx, pod)
	require.NoError(t, err, "Expected no error when sending delete event for untracked pod")
	s, err := eventloop.WaitForGeneration(ctx, imagePullEventLoop, gen)
	require.NoError(t, err, "Expected no error when waiting for delete event to be processed")
	assert.NotNil(t, s, "Expected state to be not nil")
	assert.Equal(t, 0, s.State().Len(), "Expected no image pull statistics")
}

func TestPodResync(t *testing.T) {
	testhelpers.ConfigureLogging(t, &options.Options{})

	buf := &bytes.Buffer{}

	podStatistics := state.NewPodStatistics([]apimachinerytypes.UID{})
	assert.False(t, podStatistics.IsBlacklisted("test-uid"), "Expected test-uid to not be blacklisted")
	assert.Zero(t, podStatistics.Len(), "Expected number of tracked statistics to be 0")

	ev := &resyncEvent{
		blacklistUIDs: []apimachinerytypes.UID{"test-uid"},
		output:        buf,
	}
	nextPodStatistics := ev.Dispatch(0, podStatistics)
	assert.NotNil(t, nextPodStatistics, "Expected nextPodStatistics to be not nil")
	assert.True(t, nextPodStatistics.IsBlacklisted("test-uid"), "Expected test-uid to be blacklisted")
	assert.Zero(t, nextPodStatistics.Len(), "Expected number of tracked statistics to be 0")
	assert.Empty(t, buf.String(), "Expected no output")
}

func TestPodResyncKeepsTrackedPods(t *testing.T) {
	testhelpers.ConfigureLogging(t, &options.Options{})

	created := time.Now()
	pod := newTestingPod(created)

	podStatistics := state.NewPodStatistics([]apimachinerytypes.UID{})
	statistic := state.NewPodStatistic(created, pod)
	podStatistics = podStatistics.Set("test-uid", statistic)

	ev := &resyncEvent{
		blacklistUIDs: []apimachinerytypes.UID{"test-uid"}, // pod is still in cluster
		output:        io.Discard,
	}
	nextStats := ev.Dispatch(0, podStatistics)

	assert.Equal(t, 1, nextStats.Len(), "Expected tracked pod to be kept")
	nextStatistic, ok := nextStats.Get("test-uid")
	assert.True(t, ok, "Expected pod statistic to be found")
	assert.Same(t, statistic, nextStatistic, "Expected pod statistic pointer to be unchanged")
}

func TestPodResyncRemovesLostPods(t *testing.T) {
	testhelpers.ConfigureLogging(t, &options.Options{})

	created := time.Now()
	pod := newTestingPod(created)

	podStatistics := state.NewPodStatistics([]apimachinerytypes.UID{})
	statistic := state.NewPodStatistic(created, pod)
	podStatistics = podStatistics.Set("test-uid", statistic)

	ev := &resyncEvent{
		blacklistUIDs: []apimachinerytypes.UID{}, // pod is NOT in cluster anymore
		output:        io.Discard,
	}
	nextStats := ev.Dispatch(0, podStatistics)

	assert.Equal(t, 0, nextStats.Len(), "Expected lost pod to be removed")
	_, ok := nextStats.Get("test-uid")
	assert.False(t, ok, "Expected pod statistic to not be found")
}

func TestPodUpdateAddsNewPod(t *testing.T) {
	opts := &options.Options{}
	testhelpers.ConfigureLogging(t, opts)

	created := time.Now()
	pod := newTestingPod(created)

	podStatistics := state.NewPodStatistics([]apimachinerytypes.UID{})

	updateEvent := &podUpdateEvent{
		pod:       pod,
		eventTime: created,
		options:   opts,
		output:    io.Discard,
	}
	nextStats := updateEvent.Dispatch(0, podStatistics)

	assert.Equal(t, 1, nextStats.Len(), "Expected 1 pod statistic")
	statistic, ok := nextStats.Get("test-uid")
	assert.True(t, ok, "Expected pod statistic to be found")
	assert.NotNil(t, statistic, "Expected pod statistic to not be nil")
	assert.True(t, statistic.Partial(), "Expected pod statistic to be partial")
}

func TestPodUpdateBlacklisted(t *testing.T) {
	opts := &options.Options{}
	testhelpers.ConfigureLogging(t, opts)

	created := time.Now()
	pod := newTestingPod(created)

	podStatistics := state.NewPodStatistics([]apimachinerytypes.UID{"test-uid"})

	updateEvent := &podUpdateEvent{
		pod:       pod,
		eventTime: created,
		options:   opts,
		output:    io.Discard,
	}
	nextStats := updateEvent.Dispatch(0, podStatistics)

	assert.Same(t, podStatistics, nextStats, "Expected PodStatistics to be unchanged for blacklisted pod")
	assert.Equal(t, 0, nextStats.Len(), "Expected no pod statistics (blacklisted)")
	assert.True(t, nextStats.IsBlacklisted("test-uid"), "Expected test-uid to still be blacklisted")
}

func TestPodUpdateSkipsCompletePod(t *testing.T) {
	opts := &options.Options{}
	testhelpers.ConfigureLogging(t, opts)

	created := time.Now()
	pod := newTestingCompletePod(created)

	statistic := state.NewPodStatistic(created.Add(3*time.Second), pod)
	require.False(t, statistic.Partial(), "Expected pod statistic to not be partial")

	podStatistics := state.NewPodStatistics([]apimachinerytypes.UID{})
	podStatistics = podStatistics.Set("test-uid", statistic)

	updateEvent := &podUpdateEvent{
		pod:       pod,
		eventTime: created.Add(4 * time.Second),
		options:   opts,
		output:    io.Discard,
	}
	nextStats := updateEvent.Dispatch(0, podStatistics)

	assert.Same(t, podStatistics, nextStats, "Expected PodStatistics to be unchanged for complete pod")
}

func TestPodUpdateEmitsPartialStatistics(t *testing.T) {
	opts := &options.Options{EmitPartialStatistics: true}
	testhelpers.ConfigureLogging(t, opts)

	created := time.Now()
	pod := newTestingPod(created)

	podStatistics := state.NewPodStatistics([]apimachinerytypes.UID{})

	output := testhelpers.NewMetricWriter(t)
	updateEvent := &podUpdateEvent{
		pod:       pod,
		eventTime: created,
		options:   opts,
		output:    output,
	}
	nextStats := updateEvent.Dispatch(0, podStatistics)

	require.Equal(t, 1, nextStats.Len(), "Expected 1 pod statistic")

	metrics := testhelpers.DecodeMetricOutput(t, output)
	require.Len(t, metrics, 2, "Expected 2 metrics (pod + container) for partial pod")

	metricTypes := make([]string, 0, len(metrics))

	for _, metric := range metrics {
		assert.Equal(t, true, metric["partial"], "Expected partial=true")
		metricType, ok := metric["type"].(string)
		require.True(t, ok, "Expected type to be a string")

		metricTypes = append(metricTypes, metricType)
	}

	assert.ElementsMatch(t, []string{"pod", "container"}, metricTypes, "Expected pod and container metric types")
}

func TestPodUpdateSuppressesPartialStatistics(t *testing.T) {
	opts := &options.Options{EmitPartialStatistics: false}
	testhelpers.ConfigureLogging(t, opts)

	created := time.Now()
	pod := newTestingPod(created)

	podStatistics := state.NewPodStatistics([]apimachinerytypes.UID{})

	output := testhelpers.NewMetricWriter(t)
	updateEvent := &podUpdateEvent{
		pod:       pod,
		eventTime: created,
		options:   opts,
		output:    output,
	}
	updateEvent.Dispatch(0, podStatistics)

	metrics := testhelpers.DecodeMetricOutput(t, output)
	assert.Empty(t, metrics, "Expected no metrics for partial pod when EmitPartialStatistics=false")
}

func TestPodDeleteRemovesTrackedPod(t *testing.T) {
	opts := &options.Options{}
	testhelpers.ConfigureLogging(t, opts)

	created := time.Now()
	pod := newTestingPod(created)

	podStatistics := state.NewPodStatistics([]apimachinerytypes.UID{})
	statistic := state.NewPodStatistic(created, pod)
	require.True(t, statistic.Partial(), "Expected pod statistic to be partial")
	podStatistics = podStatistics.Set("test-uid", statistic)

	output := testhelpers.NewMetricWriter(t)
	deleteEvent := &podDeleteEvent{
		options: opts,
		pod:     pod,
		output:  output,
	}
	nextStats := deleteEvent.Dispatch(0, podStatistics)

	assert.Equal(t, 0, nextStats.Len(), "Expected no pod statistics after delete")

	metrics := testhelpers.DecodeMetricOutput(t, output)
	require.Len(t, metrics, 2, "Expected 2 metrics (pod + container) for partial pod on deletion")

	metricTypes := make([]string, 0, len(metrics))

	for _, metric := range metrics {
		assert.Equal(t, true, metric["partial"], "Expected partial=true on deletion")
		metricType, ok := metric["type"].(string)
		require.True(t, ok, "Expected type to be a string")

		metricTypes = append(metricTypes, metricType)
	}

	assert.ElementsMatch(t, []string{"pod", "container"}, metricTypes, "Expected pod and container metric types")
}

func TestPodDeleteSkipsUntrackedPod(t *testing.T) {
	opts := &options.Options{}
	testhelpers.ConfigureLogging(t, opts)

	pod := newTestingPod(time.Now())
	podStatistics := state.NewPodStatistics([]apimachinerytypes.UID{})

	ev := &podDeleteEvent{
		options: opts,
		pod:     pod,
		output:  io.Discard,
	}
	nextStats := ev.Dispatch(0, podStatistics)

	assert.Same(t, podStatistics, nextStats, "Expected PodStatistics to be unchanged for untracked pod")
}

func TestImagePullUpdateAddsStatistic(t *testing.T) {
	opts := &options.Options{}
	testhelpers.ConfigureLogging(t, opts)

	created := time.Now()
	pod := newTestingPod(created)

	k8sEvent := &corev1.Event{
		InvolvedObject: corev1.ObjectReference{
			FieldPath: "spec.containers{test-container}",
		},
		Reason:        "Pulling",
		LastTimestamp: metav1.NewTime(created.Add(time.Second)),
	}

	statisticState := state.NewImagePullStatistics()

	imagePullUpdate := &imagePullUpdateEvent{
		options:  opts,
		pod:      pod,
		k8sEvent: k8sEvent,
		output:   io.Discard,
	}
	nextState := imagePullUpdate.Dispatch(0, statisticState)

	assert.Equal(t, 1, nextState.Len(), "Expected 1 image pull statistic")
	podStat, ok := nextState.Get("test-uid")
	require.True(t, ok, "Expected pod image pull statistic to be found")
	containerStat, ok := podStat.Get("test-container")
	require.True(t, ok, "Expected container image pull statistic to be found")
	assert.True(t, containerStat.Partial(), "Expected container to be partial (only Pulling received)")
}

func TestImagePullUpdateInvalidFieldPath(t *testing.T) {
	opts := &options.Options{}
	testhelpers.ConfigureLogging(t, opts)

	pod := newTestingPod(time.Now())

	k8sEvent := &corev1.Event{
		InvolvedObject: corev1.ObjectReference{
			FieldPath: "invalid-field-path",
		},
		Reason: "Pulling",
	}

	statisticState := state.NewImagePullStatistics()

	imagePullUpdate := &imagePullUpdateEvent{
		options:  opts,
		pod:      pod,
		k8sEvent: k8sEvent,
		output:   io.Discard,
	}
	nextState := imagePullUpdate.Dispatch(0, statisticState)

	assert.Same(t, statisticState, nextState, "Expected state to be unchanged for invalid field path")
}

func TestImagePullUpdateEmitsPartialStatistics(t *testing.T) {
	opts := &options.Options{EmitPartialStatistics: true}
	testhelpers.ConfigureLogging(t, opts)

	created := time.Now()
	pod := newTestingPod(created)

	k8sEvent := &corev1.Event{
		InvolvedObject: corev1.ObjectReference{FieldPath: "spec.containers{test-container}"},
		Reason:         "Pulling",
		LastTimestamp:  metav1.NewTime(created.Add(time.Second)),
	}

	statisticState := state.NewImagePullStatistics()

	output := testhelpers.NewMetricWriter(t)
	imagePullUpdate := &imagePullUpdateEvent{
		options:  opts,
		pod:      pod,
		k8sEvent: k8sEvent,
		output:   output,
	}
	imagePullUpdate.Dispatch(0, statisticState)

	metrics := testhelpers.DecodeMetricOutput(t, output)
	require.Len(t, metrics, 1, "Expected 1 image_pull metric for partial container")
	assert.Equal(t, "image_pull", metrics[0]["type"], "Expected type=image_pull")
	assert.Equal(t, true, metrics[0]["partial"], "Expected partial=true")
}

func TestImagePullUpdateSuppressesPartialStatistics(t *testing.T) {
	opts := &options.Options{EmitPartialStatistics: false}
	testhelpers.ConfigureLogging(t, opts)

	created := time.Now()
	pod := newTestingPod(created)

	k8sEvent := &corev1.Event{
		InvolvedObject: corev1.ObjectReference{FieldPath: "spec.containers{test-container}"},
		Reason:         "Pulling",
		LastTimestamp:  metav1.NewTime(created.Add(time.Second)),
	}

	statisticState := state.NewImagePullStatistics()

	output := testhelpers.NewMetricWriter(t)
	imagePullUpdate := &imagePullUpdateEvent{
		options:  opts,
		pod:      pod,
		k8sEvent: k8sEvent,
		output:   output,
	}
	imagePullUpdate.Dispatch(0, statisticState)

	metrics := testhelpers.DecodeMetricOutput(t, output)
	assert.Empty(t, metrics, "Expected no metrics for partial image pull when EmitPartialStatistics=false")
}

func TestImagePullUpdateSkipsCompleteStatistic(t *testing.T) {
	opts := &options.Options{}
	testhelpers.ConfigureLogging(t, opts)

	created := time.Now()
	pod := newTestingPod(created)

	pullingEvent := &corev1.Event{
		InvolvedObject: corev1.ObjectReference{FieldPath: "spec.containers{test-container}"},
		Reason:         "Pulling",
		LastTimestamp:  metav1.NewTime(created.Add(time.Second)),
	}
	pulledEvent := &corev1.Event{
		InvolvedObject: corev1.ObjectReference{FieldPath: "spec.containers{test-container}"},
		Reason:         "Pulled",
		LastTimestamp:  metav1.NewTime(created.Add(2 * time.Second)),
	}

	statisticState := state.NewImagePullStatistics()
	statisticState = (&imagePullUpdateEvent{
		options: opts, pod: pod, k8sEvent: pullingEvent, output: io.Discard,
	}).Dispatch(0, statisticState)
	statisticState = (&imagePullUpdateEvent{
		options: opts, pod: pod, k8sEvent: pulledEvent, output: io.Discard,
	}).Dispatch(0, statisticState)

	podStat, ok := statisticState.Get("test-uid")
	require.True(t, ok, "Expected pod image pull statistic to be found")
	containerStat, ok := podStat.Get("test-container")
	require.True(t, ok, "Expected container image pull statistic to be found")
	require.False(t, containerStat.Partial(), "Expected container to be complete")

	laterEvent := &corev1.Event{
		InvolvedObject: corev1.ObjectReference{FieldPath: "spec.containers{test-container}"},
		Reason:         "Pulled",
		LastTimestamp:  metav1.NewTime(created.Add(3 * time.Second)),
	}

	output := testhelpers.NewMetricWriter(t)
	nextState := (&imagePullUpdateEvent{
		options: opts, pod: pod, k8sEvent: laterEvent, output: output,
	}).Dispatch(0, statisticState)

	assert.Same(t, statisticState, nextState, "Expected state to be unchanged for complete statistic")
	metrics := testhelpers.DecodeMetricOutput(t, output)
	assert.Empty(t, metrics, "Expected no metrics for already-complete image pull statistic")
}

func TestDeleteImagePullRemovesTrackedPod(t *testing.T) {
	opts := &options.Options{}
	testhelpers.ConfigureLogging(t, opts)

	created := time.Now()
	pod := newTestingPod(created)

	// Seed a partial image pull statistic (only Pulling received).
	pullingEvent := &corev1.Event{
		InvolvedObject: corev1.ObjectReference{FieldPath: "spec.containers{test-container}"},
		Reason:         "Pulling",
		LastTimestamp:  metav1.NewTime(created.Add(time.Second)),
	}
	statisticState := state.NewImagePullStatistics()
	statisticState = (&imagePullUpdateEvent{
		options: opts, pod: pod, k8sEvent: pullingEvent, output: io.Discard,
	}).Dispatch(0, statisticState)
	require.Equal(t, 1, statisticState.Len(), "Expected 1 image pull statistic before delete")

	output := testhelpers.NewMetricWriter(t)
	ev := &deleteImagePullEvent{
		options: opts,
		pod:     pod,
		output:  output,
	}
	nextState := ev.Dispatch(0, statisticState)

	assert.Equal(t, 0, nextState.Len(), "Expected no image pull statistics after deletion")

	metrics := testhelpers.DecodeMetricOutput(t, output)
	require.Len(t, metrics, 1, "Expected 1 image_pull metric for partial container on deletion")
	assert.Equal(t, "image_pull", metrics[0]["type"], "Expected type=image_pull")
	assert.Equal(t, true, metrics[0]["partial"], "Expected partial=true")
}

func TestDeleteImagePullSkipsUntrackedPod(t *testing.T) {
	opts := &options.Options{}
	testhelpers.ConfigureLogging(t, opts)

	pod := newTestingPod(time.Now())
	statisticState := state.NewImagePullStatistics()

	ev := &deleteImagePullEvent{
		options: opts,
		pod:     pod,
		output:  io.Discard,
	}
	nextState := ev.Dispatch(0, statisticState)

	assert.Same(t, statisticState, nextState, "Expected state to be unchanged for untracked pod")
}
