package main

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"sync/atomic"
	"time"

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

var (
	ErrNilWatchEvent = errors.New("nil watch event")
)

type ImagePullCollector struct {
	canceled   *atomic.Bool
	cancelChan chan string
	eh         *StatisticEventHandler
	namespace  string
	podUID     types.UID
}

func NewImagePullCollector(
	event_handler *StatisticEventHandler,
	namespace string,
	pod_uid types.UID,
) ImagePullCollector {
	return ImagePullCollector{
		canceled:   &atomic.Bool{},
		cancelChan: make(chan string),
		eh:         event_handler,
		namespace:  namespace,
		podUID:     pod_uid,
	}
}

// Should be run in goroutine to avoid blocking.
func (c ImagePullCollector) cancel(reason string) {
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

func (c ImagePullCollector) logger() *zerolog.Logger {
	logger := log.With().
		Str("subsystem", "image_pull_collector").
		Str("kube_namespace", c.namespace).
		Str("pod_uid", string(c.podUID)).
		Logger()

	return &logger
}

func (c ImagePullCollector) parseContainerName(
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

type ImagePullingEvent struct {
	containerName string
	initContainer bool
	event         *corev1.Event
	collector     *ImagePullCollector
}

func (ev ImagePullingEvent) PodUID() types.UID {
	return ev.collector.podUID
}

func (ev ImagePullingEvent) logger() *zerolog.Logger {
	logger := ev.collector.logger().
		With().
		Str("image_pull_event_type", "image_pulling").
		Str("container_name", ev.containerName).
		Bool("init_container", ev.initContainer).
		Logger()

	return &logger
}

func (ev ImagePullingEvent) Handle(statistic *PodStatistic) bool {
	logger := ev.logger()

	var container_statistic *ContainerStatistic
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

type ImagePulledEvent struct {
	containerName string
	initContainer bool
	event         *corev1.Event
	collector     *ImagePullCollector
}

func (ev ImagePulledEvent) PodUID() types.UID {
	return ev.collector.podUID
}

func (ev ImagePulledEvent) logger() *zerolog.Logger {
	logger := ev.collector.logger().
		With().
		Str("image_pull_event_type", "image_pulled").
		Str("container_name", ev.containerName).
		Bool("init_container", ev.initContainer).
		Logger()

	return &logger
}

func (ev ImagePulledEvent) Handle(statistic *PodStatistic) bool {
	logger := ev.logger()

	var container_statistic *ContainerStatistic
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

func (c ImagePullCollector) handleEvent(
	event_type watch.EventType,
	event *corev1.Event,
) StatisticEvent {
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
			return &ImagePullingEvent{
				initContainer: init_container,
				containerName: container_name,
				event:         event,
				collector:     &c,
			}
		} else if event.Reason == "Pulled" {
			return &ImagePulledEvent{
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

// Ignores any cancel events to avoid blocking PodCollector when a watch error
// occurs.
func (c ImagePullCollector) stubCancel() {
	reason := <-c.cancelChan

	logger := c.logger()
	logger.Warn().Msgf(
		"Received cancel event after abnormal shutdown: %+v", reason)
}

func (c ImagePullCollector) handleWatchEvent(watch_event watch.Event) bool {
	logger := c.logger()

	var event *corev1.Event
	var ok bool
	if watch_event.Type == watch.Error {
		logger.Error().Msgf("Watch event error: %+v", watch_event)

		return true
	} else if event, ok = watch_event.Object.(*corev1.Event); !ok {
		logger.Error().Msgf("Watch event is not an Event: %+v", watch_event)

		return true
	} else if statistic_event :=
		c.handleEvent(watch_event.Type, event); statistic_event != nil {
		logger.Debug().Msgf(
			"Publish event: %s", reflect.TypeOf(statistic_event).String(),
		)
		c.eh.EventChan <- statistic_event
	}

	return false
}

func (c ImagePullCollector) watchUntilEnd(watcher watch.Interface) (bool, error) {
	logger := c.logger()

	should_break := false
	for !should_break {
		select {
		case watch_event := <-watcher.ResultChan():
			// TODO: This shouldn't normally happen, but it may be caused by the UID
			// tracked being removed.
			if watch_event.Object == nil {
				return true, ErrNilWatchEvent
			}

			should_break = c.handleWatchEvent(watch_event)
			EVENTS_PROCESSED.
				With(prometheus.Labels{"event_type": string(watch_event.Type)}).
				Inc()
		case reason := <-c.cancelChan:
			logger.Debug().Msgf("Received cancel event: %s", reason)

			return true, nil
		}
	}

	return false, nil
}

func (c ImagePullCollector) watchOptions() metav1.ListOptions {
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

func (c ImagePullCollector) Run(clientset *kubernetes.Clientset) {
	logger := c.logger()

	logger.Debug().Msg("Started ImagePullCollector ...")
	IMAGE_PULL_COLLECTOR_ROUTINES.Inc()
	defer func() {
		logger.Debug().Msg("Stopped ImagePullCollector.")
		IMAGE_PULL_COLLECTOR_ROUTINES.Dec()
	}()

	watch_ops := c.watchOptions()
	for {
		watcher, err :=
			clientset.CoreV1().Events(c.namespace).Watch(context.Background(), watch_ops)
		if err != nil {
			logger.Panic().Err(err).Msg("Error starting watcher.")
			IMAGE_PULL_COLLECTOR_ERRORS.Inc()

			go c.stubCancel()

			return
		}
		cancel, err := c.watchUntilEnd(watcher)
		if err != nil {
			logger.Error().Err(err).Msg(err.Error())
			IMAGE_PULL_COLLECTOR_ERRORS.Inc()
		}
		if cancel {
			// In non-error cancel conditions, the cancel channel has already been read
			// and closed.  Otherwise, we need to stub future writes to the cancel
			// channel to avoid blocking the StatisticEventHandler routine.
			if err != nil {
				go c.stubCancel()
			}

			return
		}

		logger.Warn().Msg("Watch ended, restarting. Some events may be lost.")
		IMAGE_PULL_COLLECTOR_RESTARTS.Inc()
	}
}
