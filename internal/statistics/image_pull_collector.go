package statistics

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/prommetrics"
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

	// statisticEventLoop is the [github.com/Izzette/go-safeconcurrency/types.EventLoop] used to handle pod and image pull
	statisticEventLoop *StatisticEventLoop

	// pod is the Kubernetes pod for which image pull events are being collected.
	pod *corev1.Pod
}

// newImagePullCollector creates (but does not start) a new imagePullCollector instance.
func newImagePullCollector(
	options *options.Options,
	eventHandler *StatisticEventLoop,
	pod *corev1.Pod,
) *imagePullCollector {
	return &imagePullCollector{
		options:            options,
		canceled:           &atomic.Bool{},
		cancelChan:         make(chan string),
		statisticEventLoop: eventHandler,
		pod:                pod,
	}
}

// Run starts the imagePullCollector and begins watching for image pull events.
// It does not start a new goroutine and will block until the image pull is complete or the collector is canceled.
func (c *imagePullCollector) Run(clientset *kubernetes.Clientset) {
	logger := c.logger()

	logger.Debug().Msg("Started ImagePullCollector ...")
	prommetrics.ImagePullCollectorRoutines.Inc()
	defer func() {
		logger.Debug().Msg("Stopped ImagePullCollector.")
		prommetrics.ImagePullCollectorRoutines.Dec()
	}()

	for {
		if c.watch(clientset) {
			return
		}

		logger.Warn().Msg("Watch ended, restarting. Some events may be lost.")
		prommetrics.ImagePullCollectorRestarts.Inc()
	}
}

// handleWatchEvent processes a watch event and returns true if the watch should be stopped.
func (c *imagePullCollector) handleWatchEvent(watchEvent watch.Event) bool {
	logger := c.logger()

	var event *corev1.Event
	var isEvent bool
	if watchEvent.Type == watch.Error {
		logger.Error().Msgf("Watch event error: %+v", watchEvent)
		prommetrics.ImagePullCollectorErrors.Inc()

		return true
	} else if event, isEvent = watchEvent.Object.(*corev1.Event); !isEvent {
		logger.Panic().Msgf("Watch event is not an Event: %+v", watchEvent)
	}

	c.handleEvent(watchEvent.Type, event)

	return false
}

// handleEvent handles a Kubernetes event and publishes it to the statistic event loop if it is an image pull event.
func (c *imagePullCollector) handleEvent(
	eventType watch.EventType,
	event *corev1.Event,
) {
	logger := c.logger()

	if eventType != watch.Added {
		logger.Debug().Msgf("Ignoring non-Added event: %+v", eventType)

		return
	}

	switch event.Reason {
	case "Pulling", "Pulled":
		if _, err := c.statisticEventLoop.ImagePullUpdate(context.TODO(), c.pod, event); err != nil {
			logger.Error().Err(err).Any("event", event).Msg("Error publishing ImagePull event")
		}
	default:
		logger.Debug().Msgf("Ignoring non-ImagePull event: %+v", event.Reason)
	}
}

// watch performs a watch on the Kubernetes API for image pull events related to the pod.
func (c *imagePullCollector) watch(clientset *kubernetes.Clientset) bool {
	logger := c.logger()

	// TODO: use a ("k8s.io/client-go/tools/watch").RetryWatcher to allow fetching
	// existing events.
	watchOpts := c.watchOptions()
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

			shouldBreak := c.handleWatchEvent(watchEvent)
			prommetrics.ImagePullWatchEvents.
				With(prometheus.Labels{"event_type": string(watchEvent.Type)}).
				Inc()
			if shouldBreak {
				return false
			}
		}
	}
}

// watchOptions builds the list options for the watch request.
func (c *imagePullCollector) watchOptions() metav1.ListOptions {
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

// cancel cancels the image pull collector.
// Should be run in goroutine to avoid blocking.
func (c *imagePullCollector) cancel(reason string) {
	logger := c.logger()

	// Sleep for a bit to allow any pending events to flush.
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

// logger returns a logger scoped to the image pull collector.
//
// TODO(Izzette): Use a context.Context to propagate the logger fields.
func (c *imagePullCollector) logger() *zerolog.Logger {
	logger := log.With().
		Str("subsystem", "image_pull_collector").
		Str("kube_namespace", c.pod.Namespace).
		Str("pod_uid", string(c.pod.UID)).
		Logger()

	return &logger
}
