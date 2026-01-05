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

	options     *options.Options
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
		options:   options,
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
		options:   el.options,
		output:    metricOutput,
	})
}

// PodDelete sends an event to stop tracking the pod statistic for a pod after it has been deleted from the Kubernetes
// API.
func (el *PodStatisticEventLoop) PodDelete(
	ctx context.Context,
	pod *corev1.Pod,
) (safeconcurrencytypes.GenerationID, error) {
	return el.Send(ctx, &podDeleteEvent{
		options: el.options,
		pod:     pod,
	})
}

// PodResync sends an event to resync the event loop if the Kubernetes Watch API times out and events are lost.
func (el *PodStatisticEventLoop) PodResync(
	ctx context.Context,
	blacklistUIDs []types.UID,
) (safeconcurrencytypes.GenerationID, error) {
	return el.Send(ctx, &resyncEvent{
		blacklistUIDs: blacklistUIDs,
		output:        metricOutput,
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

	options     *options.Options
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
		options:   options,
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
		options:  el.options,
	})
}

// ImagePullDelete sends an event to delete the image pull statistic for a pod after the image pull is no longer being
// tracked.
func (el *ImagePullStatisticEventLoop) ImagePullDelete(
	ctx context.Context,
	pod *corev1.Pod,
) (safeconcurrencytypes.GenerationID, error) {
	return el.Send(ctx, &deleteImagePullEvent{
		options: el.options,
		pod:     pod,
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
	options   *options.Options
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

	if !statistic.Partial() {
		log.Trace().Str("pod_uid", string(e.pod.UID)).Msg("Pod statistic is already complete, skipping update")

		return podStatistics
	}

	statistic = statistic.Update(e.eventTime, e.pod)
	podStatistics = podStatistics.Set(e.pod.UID, statistic)

	// Emit the pod and container statistics for the pod.
	if e.options.EmitPartialStatistics || !statistic.Partial() {
		statistic.Report(e.output, e.pod)
	}

	return podStatistics
}

// podDeleteEvent is used to delete the pod statistic for a pod after it has been deleted from the Kubernetes API.
type podDeleteEvent struct {
	options *options.Options
	pod     *corev1.Pod
	output  io.Writer
}

// Dispatch implements [safeconcurrencytypes.Event.Dispatch].
func (e *podDeleteEvent) Dispatch(
	_ safeconcurrencytypes.GenerationID,
	podStatistics *state.PodStatistics,
) *state.PodStatistics {
	statistic, ok := podStatistics.Get(e.pod.UID)
	if !ok {
		return podStatistics
	}

	// Emit statistics for pods that are deleted before they become Ready.
	if statistic.Partial() {
		statistic.Report(e.output, e.pod)
	}

	return podStatistics.Delete(e.pod.UID)
}

// resyncEvent is used to resync the event loop if the Kubernetes Watch API times out, and events are lost.
// resyncEvent implements [safeconcurrencytypes.Event].
type resyncEvent struct {
	blacklistUIDs []types.UID
	output        io.Writer
}

// Dispatch implements [safeconcurrencytypes.Event.Dispatch].
func (e *resyncEvent) Dispatch(
	_ safeconcurrencytypes.GenerationID,
	podStatistics *state.PodStatistics,
) *state.PodStatistics {
	blacklistSet := make(map[types.UID]struct{}, len(e.blacklistUIDs))
	for _, uid := range e.blacklistUIDs {
		blacklistSet[uid] = struct{}{}
	}

	// newBlacklist contains UIDs that weren't previously present in the tracked pods.
	newBlacklist := make([]types.UID, 0, len(e.blacklistUIDs))
	for _, uid := range e.blacklistUIDs {
		if _, ok := podStatistics.Get(uid); !ok {
			// This pod has appeared in the resync set, but was not previously tracked, we've missed some events for it, and
			// need to blacklist it to prevent emitting inaccurate statistics.
			newBlacklist = append(newBlacklist, uid)
		}
	}

	newPodStatistics := state.NewPodStatistics(newBlacklist)

	for uid, statistic := range podStatistics.All() {
		if _, ok := blacklistSet[uid]; ok {
			// This pod was previously tracked, and is in the resync set (still in cluster), we can keep tracking it.
			newPodStatistics = newPodStatistics.Set(uid, statistic)
		} else if statistic.Partial() {
			// This pod was previously tracked, but is not in the resync set (not in cluster), we can stop tracking it.

			// TODO(Izzette): We no longer have the pod object, so we can't emit the statistics.
			// We should probably be storing the last pod object in the pod statistic, so we can emit the statistics in this
			// case.
			//
			// statistic.Report(e.output, &corev1.Pod{})

			// TODO(Izzette): it would be cool to use the pod statistic logger to at least include the pod name and namespace.
			log.Warn().
				Str("pod_uid", string(uid)).
				Msg("Pod was previously tracked, but is not in the resync set (not in cluster), statistics have been lost")
		}
	}

	podStatistics = newPodStatistics

	return podStatistics
}

// imagePullUpdateEvent is used to update the image pull statistic for a pod from the latest Kubernetes Event for image
// pulling related events.
type imagePullUpdateEvent struct {
	options  *options.Options
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
		e.logWith(log.Trace()).Msg("Pod image pull statistic not found")
		podImagePullStatistic = state.NewPodImagePullStatistic(e.pod)
	}
	// e.logWith returns the zerolog.Event, so we can chain the calls.
	//nolint:zerologlint
	e.logWith(log.Trace()).Any("pod_image_pull_statistic", podImagePullStatistic).Msg("Pod image pull statistic")

	containerImagePullStatistic, containerOk := podImagePullStatistic.Get(containerName)
	if !containerOk {
		// e.logWith returns the zerolog.Event, so we can chain the calls.
		//nolint:zerologlint
		e.logWith(log.Panic()).Msgf("container %#v in image pull statistic not found", containerName)
	}

	if !containerImagePullStatistic.Partial() {
		log.Trace().
			Str("pod_uid", string(e.pod.UID)).
			Str("container_name", containerName).
			Msg("Container image pull statistic is already complete, skipping update")

		return statisticState
	}

	containerImagePullStatistic = containerImagePullStatistic.Update(e.k8sEvent)

	if e.options.EmitPartialStatistics || !containerImagePullStatistic.Partial() {
		containerImagePullStatistic.Report(e.output, e.pod, e.k8sEvent.Message)
	}

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
	options *options.Options
	pod     *corev1.Pod
	output  io.Writer
}

// Dispatch implements [safeconcurrencytypes.Event.Dispatch].
func (e *deleteImagePullEvent) Dispatch(
	_ safeconcurrencytypes.GenerationID,
	statisticState *state.ImagePullStatistics,
) *state.ImagePullStatistics {
	imagePullStatistic, ok := statisticState.Get(e.pod.UID)
	if !ok {
		log.Trace().Str("pod_uid", string(e.pod.UID)).Msg("Pod image pull statistic not found for deletion")

		return statisticState
	}

	for _, container := range imagePullStatistic.Containers() {
		if container.Partial() {
			container.Report(e.output, e.pod, "premature deletion of pod")
		}
	}

	return statisticState.Delete(e.pod.UID)
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
