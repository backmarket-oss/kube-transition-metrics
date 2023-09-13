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
	event_handler *StatisticEventHandler,
	namespace string,
	pod_uid types.UID,
) imagePullCollector {
	return imagePullCollector{
		canceled:   &atomic.Bool{},
		cancelChan: make(chan string),
		eh:         event_handler,
		namespace:  namespace,
		podUID:     pod_uid,
	}
}

// Should be run in goroutine to avoid blocking.
func (c imagePullCollector) cancel(reason string) {
	logger := c.logger()

	// Sleep for a bit to allow any pending events to flush.
	time.Sleep(time.Second * 3)

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
	field_ref string,
) (bool, string) {
	r := regexp.MustCompile(`^spec\.((?:initC|c)ontainers)\{(.*)\}$`)

	matches := r.FindStringSubmatch(field_ref)
	logger := c.logger()
	if matches == nil {
		logger.Error().
			Str("field_ref", field_ref).
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

func (ev imagePullingEvent) PodUID() types.UID {
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

func (ev imagePullingEvent) Handle(statistic *podStatistic) bool {
	logger := ev.logger()

	var container_statistic *containerStatistic
	if ev.initContainer {
		var ok bool
		container_statistic, ok = statistic.InitContainers[ev.containerName]
		if !ok {
			logger.Error().Msgf(
				"Init container statistic does not exist for %s", ev.containerName,
			)

			return false
		}
	} else {
		var ok bool
		container_statistic, ok = statistic.Containers[ev.containerName]
		if !ok {
			logger.Error().Msgf(
				"Container statistic does not exist for %s", ev.containerName,
			)

			return false
		}
	}
	image_pull_statistic := &container_statistic.imagePull

	if !image_pull_statistic.finishedAt.IsZero() {
		logger.Debug().Str("container_name", ev.containerName).
			Msg("Skipping event for initialized pod")
	} else if image_pull_statistic.startedAt.IsZero() {
		image_pull_statistic.startedAt = ev.event.FirstTimestamp.Time
	}

	return false
}

type imagePulledEvent struct {
	containerName string
	initContainer bool
	event         *corev1.Event
	collector     *imagePullCollector
}

func (ev imagePulledEvent) PodUID() types.UID {
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

func (ev imagePulledEvent) Handle(statistic *podStatistic) bool {
	logger := ev.logger()

	var container_statistic *containerStatistic
	if ev.initContainer {
		var ok bool
		container_statistic, ok = statistic.InitContainers[ev.containerName]
		if !ok {
			logger.Error().Msgf(
				"Init container statistic does not exist for %s", ev.containerName,
			)

			return false
		}
	} else {
		var ok bool
		container_statistic, ok = statistic.Containers[ev.containerName]
		if !ok {
			logger.Error().Msgf(
				"Container statistic does not exist for %s", ev.containerName,
			)

			return false
		}
	}
	image_pull_statistic := &container_statistic.imagePull

	if image_pull_statistic.finishedAt.IsZero() {
		image_pull_statistic.finishedAt = ev.event.LastTimestamp.Time
	}

	image_pull_statistic.log(ev.event.Message)

	return false
}

func (c imagePullCollector) handleEvent(
	event_type watch.EventType,
	event *corev1.Event,
) statisticEvent {
	logger := c.logger()

	if event_type != watch.Added {
		logger.Debug().Msgf("Ignoring non-Added event: %+v", event_type)

		return nil
	}

	switch event.Reason {
	case "Pulling", "Pulled":
		field_path := event.InvolvedObject.FieldPath
		init_container, container_name := c.parseContainerName(field_path)

		if container_name == "" {
			return nil
		}
		if event.Reason == "Pulling" {
			return &imagePullingEvent{
				initContainer: init_container,
				containerName: container_name,
				event:         event,
				collector:     &c,
			}
		} else if event.Reason == "Pulled" {
			return &imagePulledEvent{
				initContainer: init_container,
				containerName: container_name,
				event:         event,
				collector:     &c,
			}
		}
	default:
		logger.Debug().Msgf("Ignoring non-ImagePull event: %+v", event.Reason)
	}

	return nil
}

func (c imagePullCollector) handleWatchEvent(watch_event watch.Event) bool {
	logger := c.logger()

	var event *corev1.Event
	var is_event bool
	if watch_event.Type == watch.Error {
		logger.Error().Msgf("Watch event error: %+v", watch_event)
		prommetrics.IMAGE_PULL_COLLECTOR_ERRORS.Inc()

		return true
	} else if event, is_event = watch_event.Object.(*corev1.Event); !is_event {
		logger.Panic().Msgf("Watch event is not an Event: %+v", watch_event)
	} else if statistic_event :=
		c.handleEvent(watch_event.Type, event); statistic_event != nil {
		logger.Debug().Msgf(
			"Publish event: %s", reflect.TypeOf(statistic_event).String(),
		)
		c.eh.EventChan <- statistic_event
	}

	return false
}

func (c imagePullCollector) watch(clientset *kubernetes.Clientset) bool {
	logger := c.logger()

	watch_ops := c.watchOptions()
	watcher, err :=
		clientset.CoreV1().Events(c.namespace).Watch(context.Background(), watch_ops)
	if err != nil {
		logger.Panic().Err(err).Msg("Error starting watcher.")
	}
	defer watcher.Stop()

	for {
		select {
		case watch_event, watcher_open := <-watcher.ResultChan():
			if !watcher_open {
				return false
			}

			should_break := c.handleWatchEvent(watch_event)
			prommetrics.EVENTS_PROCESSED.
				With(prometheus.Labels{"event_type": string(watch_event.Type)}).
				Inc()
			if should_break {
				return false
			}
		case reason := <-c.cancelChan:
			logger.Debug().Msgf("Received cancel event: %s", reason)

			return true
		}
	}
}

func (c imagePullCollector) watchOptions() metav1.ListOptions {
	time_out := int64(60)
	send_initial_events := true
	watch_ops := metav1.ListOptions{
		TimeoutSeconds:    &time_out,
		SendInitialEvents: &send_initial_events,
		Watch:             true,
		Limit:             100,
		FieldSelector: fields.Set(
			map[string]string{
				"involvedObject.uid": string(c.podUID),
			},
		).AsSelector().String(),
	}

	return watch_ops
}

func (c imagePullCollector) Run(clientset *kubernetes.Clientset) {
	logger := c.logger()

	logger.Debug().Msg("Started ImagePullCollector ...")
	prommetrics.IMAGE_PULL_COLLECTOR_ROUTINES.Inc()
	defer func() {
		logger.Debug().Msg("Stopped ImagePullCollector.")
		prommetrics.IMAGE_PULL_COLLECTOR_ROUTINES.Dec()
	}()

	for {
		if c.watch(clientset) {
			return
		}

		logger.Warn().Msg("Watch ended, restarting. Some events may be lost.")
		prommetrics.IMAGE_PULL_COLLECTOR_RESTARTS.Inc()
	}
}