package statistics

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/prommetrics"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/statistics/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// imagePullCollector is a collector that watches for image pull events in a Kubernetes pod.
//
// TODO(Izzette): Use a context.Context to handle cancellation instead of the image pull collector.
type imagePullCollector struct {
	// options are the options used to configure the imagePullCollector.
	options *options.Options

	// canceled is an atomic boolean that indicates whether the collector has been canceled.
	canceled *atomic.Bool
	// cancelChan is a channel used to signal cancellation of the collector.
	cancelChan chan string

	// statisticEventLoop is the [github.com/Izzette/go-safeconcurrency/types.EventLoop] used to handle image pull
	// statistic states.
	statisticEventLoop types.ImagePullStatisticEventLoop

	// pod is the Kubernetes pod for which image pull events are being collected.
	pod *corev1.Pod
}

// imagePullCollectorFactory is a function type that creates a new imagePullCollector instance.
// It is used to allow mocking in tests and to provide a clear contract for the collector's behavior.
type imagePullCollectorFactory func(
	*options.Options, types.ImagePullStatisticEventLoop, *corev1.Pod,
) types.ImagePullCollector

// newImagePullCollector creates (but does not start) a new imagePullCollector instance.
func newImagePullCollector(
	options *options.Options,
	statisticEventLoop types.ImagePullStatisticEventLoop,
	pod *corev1.Pod,
) *imagePullCollector {
	return &imagePullCollector{
		options:            options,
		canceled:           &atomic.Bool{},
		cancelChan:         make(chan string),
		statisticEventLoop: statisticEventLoop,
		pod:                pod,
	}
}

// Run starts the imagePullCollector and begins watching for image pull events.
// It does not start a new goroutine and will block until the image pull is complete or the collector is canceled.
//
// Run implements [types.ImagePullCollector.Run].
func (c *imagePullCollector) Run(clientset *kubernetes.Clientset) {
	logger := c.Logger()

	logger.Debug().Msg("Started ImagePullCollector ...")
	prommetrics.ImagePullCollectorRoutines.Inc()

	defer func() {
		_, err := c.statisticEventLoop.ImagePullDelete(context.TODO(), c.pod)
		if err != nil {
			logger.Error().Err(err).Msg("Error cleaning up image pull statistic")
		}

		logger.Debug().Msg("Stopped ImagePullCollector.")
		prommetrics.ImagePullCollectorRoutines.Dec()
	}()

	for {
		if c.Watch(clientset) {
			return
		}

		logger.Warn().Msg("Watch ended, restarting. Some events may be lost.")
		prommetrics.ImagePullCollectorRestarts.Inc()
	}
}

// HandleWatchEvent processes a watch event and returns true if the watch should be stopped.
//
// HandleWatchEvent implements [types.ImagePullCollector.HandleWatchEvent].
func (c *imagePullCollector) HandleWatchEvent(watchEvent watch.Event) bool {
	logger := c.Logger()

	var (
		event   *corev1.Event
		isEvent bool
	)

	if watchEvent.Type == watch.Error {
		logger.Error().Msgf("Watch event error: %+v", watchEvent)
		prommetrics.ImagePullCollectorErrors.Inc()

		return true
	} else if event, isEvent = watchEvent.Object.(*corev1.Event); !isEvent {
		logger.Panic().Msgf("Watch event is not an Event: %+v", watchEvent)
	}

	c.HandleEvent(watchEvent.Type, event)

	return false
}

// HandleEvent handles a Kubernetes event and publishes it to the statistic event loop if it is an image pull event.
//
// HandleEvent implements [types.ImagePullCollector.HandleEvent].
func (c *imagePullCollector) HandleEvent(
	eventType watch.EventType,
	event *corev1.Event,
) {
	logger := c.Logger()

	if eventType != watch.Added {
		logger.Debug().Msgf("Ignoring non-Added event: %+v", eventType)

		return
	}

	switch event.Reason {
	case "Pulling", "Pulled":
		_, err := c.statisticEventLoop.ImagePullUpdate(context.TODO(), c.pod, event)
		if err != nil {
			logger.Error().Err(err).Any("event", event).Msg("Error publishing ImagePull event")
		}
	default:
		logger.Debug().Msgf("Ignoring non-ImagePull event: %+v", event.Reason)
	}
}

// Watch performs a Watch on the Kubernetes API for image pull events related to the pod.
//
// Watch implements [types.ImagePullCollector.Watch].
func (c *imagePullCollector) Watch(clientset *kubernetes.Clientset) bool {
	logger := c.Logger()

	// TODO: use a ("k8s.io/client-go/tools/watch").RetryWatcher to allow fetching
	// existing events.
	watchOpts := c.WatchOptions()

	watcher, err :=
		clientset.CoreV1().Events(c.pod.Namespace).Watch(context.TODO(), watchOpts)
	if err != nil {
		watchOptsJSON, marshalErr := json.Marshal(watchOpts)
		if marshalErr != nil {
			watchOptsJSON = []byte("null")
		}

		logger.Panic().
			Str("watch_namespace", c.pod.Namespace).
			RawJSON("watch_opts", watchOptsJSON).
			AnErr("watch_opts_marshal_err", marshalErr).
			Err(err).Msg("Error starting watcher.")
	}
	defer watcher.Stop()

	for {
		select {
		case reason := <-c.cancelChan:
			logger.Debug().Msgf("Received cancel event: %s", reason)

			return true
		case watchEvent, watcherOpen := <-watcher.ResultChan():
			if !watcherOpen {
				return false
			}

			shouldBreak := c.HandleWatchEvent(watchEvent)
			prommetrics.ImagePullWatchEvents.
				With(prometheus.Labels{"event_type": string(watchEvent.Type)}).
				Inc()

			if shouldBreak {
				return false
			}
		}
	}
}

// WatchOptions builds the list options for the watch request.
//
// WatchOptions implements [types.ImagePullCollector.WatchOptions].
func (c *imagePullCollector) WatchOptions() metav1.ListOptions {
	timeOut := c.options.KubeWatchTimeout
	watchOps := metav1.ListOptions{
		TimeoutSeconds: &timeOut,
		Watch:          true,
		Limit:          c.options.KubeWatchMaxEvents,
		FieldSelector: fields.Set(
			map[string]string{
				"involvedObject.uid": string(c.pod.UID),
			},
		).AsSelector().String(),
	}

	return watchOps
}

// Cancel cancels the image pull collector.
// Should be run in goroutine to avoid blocking.
//
// Cancel implements [types.ImagePullCollector.Cancel].
func (c *imagePullCollector) Cancel(reason string) {
	logger := c.Logger()

	// Sleep for a bit to allow any pending events to flush.
	// This is a workaround for the fact that the Kubernetes Watch API does not guarantee that all events related to a pod
	// will be delivered before the pod is deleted.
	//
	// TODO(Izzette): surely there's a better way to do this?
	time.Sleep(time.Second * time.Duration(c.options.ImagePullCancelDelay))

	if !c.canceled.Swap(true) {
		logger.Debug().Msgf("Canceling collector: %s", reason)

		c.cancelChan <- reason

		close(c.cancelChan)
	} else {
		logger.Debug().Msgf("Duplicate collector cancel received: %s", reason)
	}
}

// Logger returns a Logger scoped to the image pull collector.
//
// TODO(Izzette): Use a context.Context to propagate the Logger fields.
//
// Logger implements [types.ImagePullCollector.Logger].
func (c *imagePullCollector) Logger() *zerolog.Logger {
	logger := log.With().
		Str("subsystem", "image_pull_collector").
		Str("kube_namespace", c.pod.Namespace).
		Str("pod_uid", string(c.pod.UID)).
		Logger()

	return &logger
}
