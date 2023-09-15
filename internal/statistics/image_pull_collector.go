package statistics

import (
	"context"
	"reflect"
	"regexp"
	"sync/atomic"
	"time"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/prommetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type imagePullCollector struct {
	canceled   *atomic.Bool
	cancelChan chan string
	eh         *StatisticEventHandler
	namespace  string
	podUID     types.UID
}

func newImagePullCollector(
	eventHandler *StatisticEventHandler,
	namespace string,
	podUID types.UID,
) imagePullCollector {
	return imagePullCollector{
		canceled:   &atomic.Bool{},
		cancelChan: make(chan string),
		eh:         eventHandler,
		namespace:  namespace,
		podUID:     podUID,
	}
}

// Should be run in goroutine to avoid blocking.
func (c imagePullCollector) cancel(reason string) {
	logger := c.logger()

	// Sleep for a bit to allow any pending events to flush.
	time.Sleep(time.Second * time.Duration(c.eh.options.ImagePullCancelDelay))

	if !c.canceled.Swap(true) {
		logger.Debug().Msgf("Canceling collector: %s", reason)
		c.cancelChan <- reason
		close(c.cancelChan)
	} else {
		logger.Debug().Msgf("Duplicate collector cancel received: %s", reason)
	}
}

func (c imagePullCollector) logger() *zerolog.Logger {
	logger := log.With().
		Str("subsystem", "image_pull_collector").
		Str("kube_namespace", c.namespace).
		Str("pod_uid", string(c.podUID)).
		Logger()

	return &logger
}

func (c imagePullCollector) parseContainerName(
	fieldRef string,
) (bool, string) {
	r := regexp.MustCompile(`^spec\.((?:initC|c)ontainers)\{(.*)\}$`)

	matches := r.FindStringSubmatch(fieldRef)
	logger := c.logger()
	if matches == nil {
		logger.Error().
			Str("field_ref", fieldRef).
			Msg("Failed to find container name")

		return false, ""
	}

	return matches[1] == "initContainers", matches[2]
}

type imagePullingEvent struct {
	containerName string
	initContainer bool
	event         *corev1.Event
	collector     *imagePullCollector
}

func (ev imagePullingEvent) podUID() types.UID {
	return ev.collector.podUID
}

func (ev imagePullingEvent) logger() *zerolog.Logger {
	logger := ev.collector.logger().
		With().
		Str("image_pull_event_type", "image_pulling").
		Str("container_name", ev.containerName).
		Bool("init_container", ev.initContainer).
		Logger()

	return &logger
}

func (ev imagePullingEvent) handle(statistic *podStatistic) bool {
	logger := ev.logger()

	var containerStatistic *containerStatistic
	if ev.initContainer {
		var ok bool
		containerStatistic, ok = statistic.initContainers[ev.containerName]
		if !ok {
			logger.Error().Msgf(
				"Init container statistic does not exist for %s", ev.containerName,
			)

			return false
		}
	} else {
		var ok bool
		containerStatistic, ok = statistic.containers[ev.containerName]
		if !ok {
			logger.Error().Msgf(
				"Container statistic does not exist for %s", ev.containerName,
			)

			return false
		}
	}
	imagePullStatistic := &containerStatistic.imagePull

	if !imagePullStatistic.finishedAt.IsZero() {
		logger.Debug().Str("container_name", ev.containerName).
			Msg("Skipping event for initialized pod")
	} else if imagePullStatistic.startedAt.IsZero() {
		imagePullStatistic.startedAt = ev.event.FirstTimestamp.Time
	}

	return false
}

type imagePulledEvent struct {
	containerName string
	initContainer bool
	event         *corev1.Event
	collector     *imagePullCollector
}

func (ev imagePulledEvent) podUID() types.UID {
	return ev.collector.podUID
}

func (ev imagePulledEvent) logger() *zerolog.Logger {
	logger := ev.collector.logger().
		With().
		Str("image_pull_event_type", "image_pulled").
		Str("container_name", ev.containerName).
		Bool("init_container", ev.initContainer).
		Logger()

	return &logger
}

func (ev imagePulledEvent) handle(statistic *podStatistic) bool {
	logger := ev.logger()

	var containerStatistic *containerStatistic
	if ev.initContainer {
		var ok bool
		containerStatistic, ok = statistic.initContainers[ev.containerName]
		if !ok {
			logger.Error().Msgf(
				"Init container statistic does not exist for %s", ev.containerName,
			)

			return false
		}
	} else {
		var ok bool
		containerStatistic, ok = statistic.containers[ev.containerName]
		if !ok {
			logger.Error().Msgf(
				"Container statistic does not exist for %s", ev.containerName,
			)

			return false
		}
	}
	imagePullStatistic := &containerStatistic.imagePull

	if imagePullStatistic.finishedAt.IsZero() {
		imagePullStatistic.finishedAt = ev.event.LastTimestamp.Time
	}

	imagePullStatistic.log(ev.event.Message)

	return false
}

func (c imagePullCollector) handleEvent(
	eventType watch.EventType,
	event *corev1.Event,
) statisticEvent {
	logger := c.logger()

	if eventType != watch.Added {
		logger.Debug().Msgf("Ignoring non-Added event: %+v", eventType)

		return nil
	}

	switch event.Reason {
	case "Pulling", "Pulled":
		fieldPath := event.InvolvedObject.FieldPath
		initContainer, containerName := c.parseContainerName(fieldPath)

		if containerName == "" {
			return nil
		}
		if event.Reason == "Pulling" {
			return &imagePullingEvent{
				initContainer: initContainer,
				containerName: containerName,
				event:         event,
				collector:     &c,
			}
		} else if event.Reason == "Pulled" {
			return &imagePulledEvent{
				initContainer: initContainer,
				containerName: containerName,
				event:         event,
				collector:     &c,
			}
		}
	default:
		logger.Debug().Msgf("Ignoring non-ImagePull event: %+v", event.Reason)
	}

	return nil
}

func (c imagePullCollector) handleWatchEvent(watchEvent watch.Event) bool {
	logger := c.logger()

	var event *corev1.Event
	var isEvent bool
	if watchEvent.Type == watch.Error {
		logger.Error().Msgf("Watch event error: %+v", watchEvent)
		prommetrics.ImagePullCollectorErrors.Inc()

		return true
	} else if event, isEvent = watchEvent.Object.(*corev1.Event); !isEvent {
		logger.Panic().Msgf("Watch event is not an Event: %+v", watchEvent)
	} else if statisticEvent :=
		c.handleEvent(watchEvent.Type, event); statisticEvent != nil {
		logger.Debug().Msgf(
			"Publish event: %s", reflect.TypeOf(statisticEvent).String(),
		)
		c.eh.Publish(statisticEvent)
	}

	return false
}

func (c imagePullCollector) watch(clientset *kubernetes.Clientset) bool {
	logger := c.logger()

	watchOps := c.watchOptions()
	watcher, err :=
		clientset.CoreV1().Events(c.namespace).Watch(context.Background(), watchOps)
	if err != nil {
		logger.Panic().Err(err).Msg("Error starting watcher.")
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
			prommetrics.EventsProcessed.
				With(prometheus.Labels{"event_type": string(watchEvent.Type)}).
				Inc()
			if shouldBreak {
				return false
			}
		}
	}
}

func (c imagePullCollector) watchOptions() metav1.ListOptions {
	timeOut := c.eh.options.KubeWatchTimeout
	sendInitialEvents := true
	watchOps := metav1.ListOptions{
		TimeoutSeconds:    &timeOut,
		SendInitialEvents: &sendInitialEvents,
		Watch:             true,
		Limit:             c.eh.options.KubeWatchMaxEvents,
		FieldSelector: fields.Set(
			map[string]string{
				"involvedObject.uid": string(c.podUID),
			},
		).AsSelector().String(),
	}

	return watchOps
}

func (c imagePullCollector) Run(clientset *kubernetes.Clientset) {
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
