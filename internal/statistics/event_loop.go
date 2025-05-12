package statistics

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/prommetrics"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/statistics/state"
	safeconcurrencytypes "github.com/Izzette/go-safeconcurrency/api/types"
	"github.com/Izzette/go-safeconcurrency/eventloop"
	"github.com/Izzette/go-safeconcurrency/eventloop/snapshot"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// PodStatisticEventLoop loops over pod statistic events sent by collectors to track and update metrics.
type PodStatisticEventLoop struct {
	safeconcurrencytypes.EventLoop[*state.PodStatistics]
	watcherChan <-chan struct{}
}

// NewStatisticEventLoop creates a new StatisticEventHandler which filters out events for the provided
// initialSyncBlacklist Pod UIDs.
func NewStatisticEventLoop(options *options.Options) *PodStatisticEventLoop {
	s := state.NewPodStatistics([]types.UID{})
	snapshot := snapshot.NewCopyable[*state.PodStatistics](s)
	if options.StatisticEventQueueLength < 0 {
		log.Panic().Msg("StatisticEventQueueLength must be greater than 0")
	}
	buffer := uint(options.StatisticEventQueueLength)

	return &PodStatisticEventLoop{
		EventLoop: eventloop.NewBuffered[*state.PodStatistics](snapshot, buffer),
	}
}

// Start starts the event loop and begins watching the state of the event loop.
// Start implements [eventloop.EventLoop.Start].
func (el *PodStatisticEventLoop) Start() {
	el.EventLoop.Start()
	// EventLoop.Start() will panic if the event loop is already started, so we can be sure to do this assignment only
	// once.
	el.watcherChan = eventloop.WatchState(context.TODO(), el.EventLoop, el.watcher)
}

// Close closes the event loop and waits for the watcher to finish.
// Close implements [eventloop.EventLoop.Close].
func (el *PodStatisticEventLoop) Close() {
	el.EventLoop.Close()
	// Wait for the watcher to finish too.
	<-el.watcherChan
}

// Send sends an event to the event loop and tracks the time it took to process the event.
// Send implements [eventloop.EventLoop.Send].
func (el *PodStatisticEventLoop) Send(
	ctx context.Context,
	event safeconcurrencytypes.Event[*state.PodStatistics],
) (safeconcurrencytypes.GenerationID, error) {
	start := time.Now()

	labels := prometheus.Labels{"event_loop": "pod"}

	// Wrap the event to track the event metrics.
	tevent := &trackEvent[*state.PodStatistics]{event, &sync.Once{}, labels}
	gen, err := el.EventLoop.Send(ctx, tevent)
	if err == nil {
		// Only increment the queue depth if the event was successfully sent to the event loop.
		tevent.incrementQueueDepth()
	}

	dur := time.Since(start)
	prommetrics.StatisticEventPublish.With(labels).Observe(dur.Seconds())

	//nolint:wrapcheck
	return gen, err
}

// PodUpdate sends an event to update the pod statistic for a pod based on the latest Kubernetes Pod.
func (el *PodStatisticEventLoop) PodUpdate(
	ctx context.Context,
	pod *corev1.Pod,
) (safeconcurrencytypes.GenerationID, error) {
	return el.Send(ctx, &podUpdateEvent{
		pod:       pod,
		eventTime: time.Now(),
		output:    metricOutput,
	})
}

// PodDelete sends an event to stop tracking the pod statistic for a pod after it has been deleted from the Kubernetes
// API.
func (el *PodStatisticEventLoop) PodDelete(
	ctx context.Context,
	podUID types.UID,
) (safeconcurrencytypes.GenerationID, error) {
	return el.Send(ctx, &podDeleteEvent{
		podUID: podUID,
	})
}

// PodResync sends an event to resync the event loop if the Kubernetes Watch API times out and events are lost.
func (el *PodStatisticEventLoop) PodResync(
	ctx context.Context,
	blacklistUIDs []types.UID,
) (safeconcurrencytypes.GenerationID, error) {
	return el.Send(ctx, &resyncEvent{
		blacklistUIDs: blacklistUIDs,
	})
}

// watcher watches the state of the event loop and updates the prometheus metrics.
func (el *PodStatisticEventLoop) watcher(
	ctx context.Context,
	s safeconcurrencytypes.StateSnapshot[*state.PodStatistics],
) bool {
	prommetrics.PodsTracked.Set(float64(s.State().Len()))

	return true
}

// ImagePullStatisticEventLoop loops over image pull statistic events sent by collectors to track and update metrics for
// image pulls.
type ImagePullStatisticEventLoop struct {
	safeconcurrencytypes.EventLoop[*state.ImagePullStatistics]
	watcherChan <-chan struct{}
}

// NewImagePullStatisticEventLoop creates a new ImagePullStatisticEventLoop.
func NewImagePullStatisticEventLoop(options *options.Options) *ImagePullStatisticEventLoop {
	s := state.NewImagePullStatistics()
	snapshot := snapshot.NewCopyable[*state.ImagePullStatistics](s)
	if options.StatisticEventQueueLength < 0 {
		log.Panic().Msg("StatisticEventQueueLength must be greater than 0")
	}
	buffer := uint(options.StatisticEventQueueLength)

	return &ImagePullStatisticEventLoop{
		EventLoop: eventloop.NewBuffered[*state.ImagePullStatistics](snapshot, buffer),
	}
}

// Start starts the event loop and begins watching the state of the event loop.
func (el *ImagePullStatisticEventLoop) Start() {
	el.EventLoop.Start()
	// EventLoop.Start() will panic if the event loop is already started, so we can be sure to do this assignment only
	// once.
	el.watcherChan = eventloop.WatchState(context.TODO(), el.EventLoop, el.watcher)
}

// Close closes the event loop and waits for the watcher to finish.
func (el *ImagePullStatisticEventLoop) Close() {
	el.EventLoop.Close()
	// Wait for the watcher to finish too.
	<-el.watcherChan
}

// Send sends an event to the event loop and tracks the time it took to process the event.
func (el *ImagePullStatisticEventLoop) Send(
	ctx context.Context,
	event safeconcurrencytypes.Event[*state.ImagePullStatistics],
) (safeconcurrencytypes.GenerationID, error) {
	start := time.Now()

	labels := prometheus.Labels{"event_loop": "image_pull"}

	// Wrap the event to track the event metrics.
	tevent := &trackEvent[*state.ImagePullStatistics]{event, &sync.Once{}, labels}
	gen, err := el.EventLoop.Send(ctx, tevent)
	if err == nil {
		// Only increment the queue depth if the event was successfully sent to the event loop.
		tevent.incrementQueueDepth()
	}

	dur := time.Since(start)
	prommetrics.StatisticEventPublish.With(labels).Observe(dur.Seconds())

	//nolint:wrapcheck
	return gen, err
}

// ImagePullUpdate sends an event to update the image pull statistic for a pod from the latest Kubernetes Event for
// image pulling related events.
func (el *ImagePullStatisticEventLoop) ImagePullUpdate(
	ctx context.Context,
	pod *corev1.Pod,
	k8sEvent *corev1.Event,
) (safeconcurrencytypes.GenerationID, error) {
	return el.Send(ctx, &imagePullUpdateEvent{
		pod:      pod,
		k8sEvent: k8sEvent,
		output:   metricOutput,
	})
}

// ImagePullDelete sends an event to delete the image pull statistic for a pod after the image pull is no longer being
// tracked.
func (el *ImagePullStatisticEventLoop) ImagePullDelete(
	ctx context.Context,
	podUID types.UID,
) (safeconcurrencytypes.GenerationID, error) {
	return el.Send(ctx, &deleteImagePullEvent{
		podUID: podUID,
	})
}

// watcher watches the state of the event loop and updates the prometheus metrics.
func (el *ImagePullStatisticEventLoop) watcher(
	ctx context.Context,
	s safeconcurrencytypes.StateSnapshot[*state.ImagePullStatistics],
) bool {
	prommetrics.ImagePullTracked.Set(float64(s.State().Len()))

	return true
}

// trackEvent wraps an event to track the event metrics.
type trackEvent[StateT any] struct {
	// event is the event to dispatch.
	event safeconcurrencytypes.Event[StateT]
	// queueDepthIncrementOnce is used to ensure that the queue depth is only incremented once per event.
	// It allows the event to increment the queue depth when starting in the racy-case that the publisher did not yet
	// increment it, allowing the queue depth to never be negative.
	queueDepthIncrementOnce *sync.Once
	// labels are the labels to use for the prometheus metrics.
	labels prometheus.Labels
}

// Dispatch implements [safeconcurrencytypes.Event.Dispatch].
func (e *trackEvent[StateT]) Dispatch(gen safeconcurrencytypes.GenerationID, statisticState StateT) StateT {
	// Increment in-case the publisher did not yet increment the queue depth.
	// This prevents the queue depth from going negative in racey edge cases.
	e.incrementQueueDepth()
	prommetrics.StatisticEventQueueDepth.With(e.labels).Dec()

	start := time.Now()
	statisticState = e.event.Dispatch(gen, statisticState)
	dur := time.Since(start)
	prommetrics.StatisticEventProcessing.With(e.labels).Observe(dur.Seconds())

	return statisticState
}

// incrementQueueDepth increments the queue depth metric.
func (e *trackEvent[StateT]) incrementQueueDepth() {
	e.queueDepthIncrementOnce.Do(prommetrics.StatisticEventQueueDepth.With(e.labels).Inc)
}

// podUpdateEvent is used to update the pod statistic for a pod from the latest Kubernetes Pod.
type podUpdateEvent struct {
	pod       *corev1.Pod
	eventTime time.Time
	output    io.Writer
}

// Dispatch implements [safeconcurrencytypes.Event.Dispatch].
func (e *podUpdateEvent) Dispatch(
	_ safeconcurrencytypes.GenerationID,
	podStatistics *state.PodStatistics,
) *state.PodStatistics {
	if podStatistics.IsBlacklisted(e.pod.UID) {
		return podStatistics
	}

	statistic, ok := podStatistics.Get(e.pod.UID)
	if !ok {
		statistic = state.NewPodStatistic(e.eventTime, e.pod)
	}

	statistic = statistic.Update(e.eventTime, e.pod)
	podStatistics = podStatistics.Set(e.pod.UID, statistic)

	statistic.Report(e.output)

	return podStatistics
}

// podDeleteEvent is used to delete the pod statistic for a pod after it has been deleted from the Kubernetes API.
type podDeleteEvent struct {
	podUID types.UID
}

// Dispatch implements [safeconcurrencytypes.Event.Dispatch].
func (e *podDeleteEvent) Dispatch(
	_ safeconcurrencytypes.GenerationID,
	podStatistics *state.PodStatistics,
) *state.PodStatistics {
	return podStatistics.Delete(e.podUID)
}

// resyncEvent is used to resync the event loop if the Kubernetes Watch API times out, and events are lost.
// resyncEvent implements [safeconcurrencytypes.Event].
type resyncEvent struct {
	blacklistUIDs []types.UID
}

// Dispatch implements [safeconcurrencytypes.Event.Dispatch].
func (e *resyncEvent) Dispatch(
	_ safeconcurrencytypes.GenerationID,
	podStatistics *state.PodStatistics,
) *state.PodStatistics {
	// newBlacklist contains UIDs that weren't previously present in the tracked pods.
	newBlacklist := make([]types.UID, 0, len(e.blacklistUIDs))
	statistics := make(map[types.UID]*state.PodStatistic, len(e.blacklistUIDs))
	for _, uid := range e.blacklistUIDs {
		statistic, ok := podStatistics.Get(uid)
		if !ok {
			newBlacklist = append(newBlacklist, uid)

			continue
		}
		statistics[uid] = statistic
	}

	podStatistics = state.NewPodStatistics(newBlacklist)
	for uid, statistic := range statistics {
		podStatistics = podStatistics.Set(uid, statistic)
	}

	return podStatistics
}

// imagePullUpdateEvent is used to update the image pull statistic for a pod from the latest Kubernetes Event for image
// pulling related events.
type imagePullUpdateEvent struct {
	pod      *corev1.Pod
	k8sEvent *corev1.Event
	output   io.Writer
}

// Dispatch implements [safeconcurrencytypes.Event.Dispatch].
func (e *imagePullUpdateEvent) Dispatch(
	_ safeconcurrencytypes.GenerationID,
	statisticState *state.ImagePullStatistics,
) *state.ImagePullStatistics {
	containerName, err := e.getContainerName()
	if err != nil {
		// e.logWith returns the zerolog.Event, so we can chain the calls.
		//nolint:zerologlint
		e.logWith(log.Err(err)).Send()

		return statisticState
	}

	podImagePullStatistic, podOk := statisticState.Get(e.pod.UID)
	if !podOk {
		// e.logWith returns the zerolog.Event, so we can chain the calls.
		//nolint:zerologlint
		e.logWith(log.Trace()).Msg("pod image pull statistic not found")
		podImagePullStatistic = state.NewPodImagePullStatistic(e.pod)
	}
	// e.logWith returns the zerolog.Event, so we can chain the calls.
	//nolint:zerologlint
	e.logWith(log.Trace()).Msgf("pod image pull statistic: %#+v", podImagePullStatistic)

	containerImagePullStatistic, containerOk := podImagePullStatistic.Get(containerName)
	if !containerOk {
		// e.logWith returns the zerolog.Event, so we can chain the calls.
		//nolint:zerologlint
		e.logWith(log.Panic()).Msgf("container %#v in image pull statistic not found", containerName)
	}
	containerImagePullStatistic = containerImagePullStatistic.Update(e.k8sEvent)
	containerImagePullStatistic.Report(e.output, e.k8sEvent.Message)

	podImagePullStatistic = podImagePullStatistic.Set(containerImagePullStatistic)
	statisticState = statisticState.Set(e.pod.UID, podImagePullStatistic)

	return statisticState
}

// getContainerName parses the container name from the fieldRef of the Kubernetes Event.
func (e *imagePullUpdateEvent) getContainerName() (string, error) {
	fieldRef := e.k8sEvent.InvolvedObject.FieldPath

	matches := fieldPathContainerRegex.FindStringSubmatch(fieldRef)
	if matches == nil {
		return "", newParseContainerNameError(fieldRef)
	} else if len(matches) != 1+fieldPathContainerRegex.NumSubexp() {
		log.Panic().Msgf("regex %#+v does not produce the expected number of groups", fieldPathContainerRegex)
	}

	return matches[1], nil
}

// TODO(Izzette): replace with [zerolog.Context].
func (e *imagePullUpdateEvent) logWith(ev *zerolog.Event) *zerolog.Event {
	return ev.
		Str("kube_namespace", e.pod.Namespace).
		Str("pod_name", e.pod.Name)
}

// deleteImagePullEvent is used to delete the image pull statistic for a pod after it has been deleted from the
// Kubernetes API.
type deleteImagePullEvent struct {
	podUID types.UID
}

// Dispatch implements [safeconcurrencytypes.Event.Dispatch].
func (e *deleteImagePullEvent) Dispatch(
	_ safeconcurrencytypes.GenerationID,
	statisticState *state.ImagePullStatistics,
) *state.ImagePullStatistics {
	return statisticState.Delete(e.podUID)
}

// fieldPathContainerRegex is used to parse the container name from the fieldRef of the Kubernetes Event.
var fieldPathContainerRegex = regexp.MustCompile(`^spec\.(?:initC|c)ontainers\{(.*)\}$`)

// parseContainerNameError is used to indicate that the container name could not be parsed from the involved object's
// field-path of the Kubernetes Event.
type parseContainerNameError struct {
	// fieldRef is the fieldRef of the Kubernetes Event.
	fieldRef string
}

// newParseContainerNameError creates a new parseContainerNameError.
func newParseContainerNameError(fieldRef string) *parseContainerNameError {
	return &parseContainerNameError{fieldRef}
}

// Error implements the error interface for parseContainerNameError.
func (e *parseContainerNameError) Error() string {
	return fmt.Sprintf("failed to parse container name from %#v", e.fieldRef)
}
