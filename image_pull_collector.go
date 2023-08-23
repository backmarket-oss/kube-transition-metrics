package main

import (
	"context"
	"reflect"
	"regexp"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type ImagePullCollector struct {
	canceled   *atomic.Bool
	cancelChan chan string
	eh         *StatisticEventHandler
	namespace  string
	podUID     types.UID
}

func NewImagePullCollector(eh *StatisticEventHandler, namespace string, pod_uid types.UID) ImagePullCollector {
	return ImagePullCollector{
		canceled:   &atomic.Bool{},
		cancelChan: make(chan string),
		eh:         eh,
		namespace:  namespace,
		podUID:     pod_uid,
	}
}

// Should be run in goroutine to avoid blocking.
func (c ImagePullCollector) cancel(reason string) {
	logger := c.logger()

	// Sleep for a bit to allow any pending events to flush.
	time.Sleep(time.Second * 3)

	if !c.canceled.Load() {
		logger.Debug().Msgf("Canceling collector: %s", reason)
		c.canceled.Store(true)
		c.cancelChan <- reason
		close(c.cancelChan)
	} else {
		logger.Warn().Msgf("Duplicate collector cancel received: %s", reason)
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

func (c ImagePullCollector) parseContainerName(field_ref string) (bool, string) {
	r := regexp.MustCompile(`^spec\.((?:initC|c)ontainers)\{(.*)\}$`)

	matches := r.FindStringSubmatch(field_ref)
	logger := c.logger()
	if matches == nil {
		logger.Error().Str("field_ref", field_ref).Msg("Failed to find container name")

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
		Str("image_pull_event_type", "image_pulling_event").
		Str("container_name", ev.containerName).
		Bool("init_container", ev.initContainer).
		Logger()

	return &logger
}

func (ev ImagePullingEvent) Handle(statistic *PodStatistic) bool {
	logger := ev.logger()

	var image_pull_statistic *ImagePullStatistic
	if ev.initContainer {
		image_pull_statistic = &statistic.InitContainers[ev.containerName].imagePull
	} else {
		image_pull_statistic = &statistic.Containers[ev.containerName].imagePull
	}

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

func (ev ImagePulledEvent) Handle(statistic *PodStatistic) bool {
	var image_pull_statistic *ImagePullStatistic
	if ev.initContainer {
		image_pull_statistic = &statistic.InitContainers[ev.containerName].imagePull
	} else {
		image_pull_statistic = &statistic.Containers[ev.containerName].imagePull
	}

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

// Ignores any cancel events to avoid blocking PodCollector when a watch error occurs.
func (c ImagePullCollector) stubCancel() {
	logger := c.logger()
	reason := <-c.cancelChan

	logger.Warn().Msgf("Received cancel event after shutdown: %+v", reason)
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
	} else if statistic_event := c.handleEvent(watch_event.Type, event); statistic_event != nil {
		logger.Info().Msgf("Publish event: %q", reflect.TypeOf(statistic_event))
		c.eh.EventChan <- statistic_event
	}

	return false
}

func (c ImagePullCollector) Run(clientset *kubernetes.Clientset) {
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

	logger := c.logger()

	logger.Debug().Msg("Started ImagePullCollector ...")
	defer logger.Debug().Msg("Stopped ImagePullCollector.")

	for {
		watcher, err :=
			clientset.CoreV1().Events(c.namespace).Watch(context.Background(), watch_ops)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error starting watcher.")

			go c.stubCancel()

			return
		}
		should_break := false
		for !should_break {
			select {
			case watch_event := <-watcher.ResultChan():
				should_break = c.handleWatchEvent(watch_event)
			case reason := <-c.cancelChan:
				logger.Debug().Msgf("Received cancel event: %s", reason)

				return
			}
		}

		logger.Warn().Msg("Watch ended, restarting. Some events may be lost.")
	}
}
