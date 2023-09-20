package statistics

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
)

type podStatistic struct {
	initialized bool

	name       string
	namespace  string
	timeSource timeSource

	imagePullCollector imagePullCollector

	// The timestamp for when the pod was created, same as timestamp of when pod
	// was in Pending and containers were Waiting.
	creationTimestamp time.Time

	// The timestamp for when the pod was scheduled.
	scheduledTimestamp time.Time

	// The timestamp for when the pod was initialized.
	initializedTimestamp time.Time

	// The timestamp for when the pod first turned Ready.
	readyTimestamp time.Time

	initContainers map[string]*containerStatistic
	containers     map[string]*containerStatistic
}

func (s *podStatistic) initialize(pod *corev1.Pod) {
	if s.initialized {
		return
	}
	s.initialized = true
	s.timeSource = realTimeSource{}
	s.name = pod.Name
	s.namespace = pod.Namespace
	s.creationTimestamp = pod.CreationTimestamp.Time
	s.initContainers = make(map[string]*containerStatistic)
	s.containers = make(map[string]*containerStatistic)

	for _, container := range pod.Spec.InitContainers {
		s.initContainers[container.Name] = newContainerStatistic(s, true, container)
	}
	for _, container := range pod.Spec.Containers {
		s.containers[container.Name] = newContainerStatistic(s, false, container)
	}
}

func (s podStatistic) logger() zerolog.Logger {
	return log.With().
		Str("kube_namespace", s.namespace).
		Str("pod_name", s.name).
		Logger()
}

func (s podStatistic) event() *zerolog.Event {
	event := zerolog.Dict()

	if !s.scheduledTimestamp.IsZero() {
		event.Float64(
			"scheduled_latency",
			s.scheduledTimestamp.Sub(s.creationTimestamp).Seconds())
	}
	if !s.initializedTimestamp.IsZero() {
		event.Float64(
			"initialized_latency",
			s.initializedTimestamp.Sub(s.creationTimestamp).Seconds())
	}
	if !s.readyTimestamp.IsZero() {
		event.Float64(
			"ready_latency",
			s.readyTimestamp.Sub(s.creationTimestamp).Seconds())
	}

	return event
}

func (s podStatistic) report() {
	logger := s.logger()

	eventLogger := logger.Output(metricOutput).With().
		Str("kube_transition_metric_type", "pod").
		Dict("kube_transition_metrics", s.event()).Logger()
	eventLogger.Log().Msg("")

	for _, containerStatistics := range s.initContainers {
		containerStatistics.report()
	}
	for _, containerStatistics := range s.containers {
		containerStatistics.report()
	}
}

func (s *podStatistic) update(pod *corev1.Pod) {
	logger := s.logger()

	for _, condition := range pod.Status.Conditions {
		logger.Debug().Msgf("Pod condition: %+v", condition)
		if condition.Status != corev1.ConditionTrue {
			continue
		}

		// TODO: include core/v1.ContainersReady and core/v1.DisruptionTarget
		//nolint:exhaustive
		switch condition.Type {
		case corev1.PodScheduled:
			if s.scheduledTimestamp.IsZero() {
				s.scheduledTimestamp = condition.LastTransitionTime.Time
			}
		case corev1.PodInitialized:
			// Pod Initialized occursafter all images pulled, no need to continue to
			// track

			if s.initializedTimestamp.IsZero() {
				go s.imagePullCollector.cancel("pod_initialized")
				s.initializedTimestamp = condition.LastTransitionTime.Time
			}
		case corev1.PodReady:
			if s.readyTimestamp.IsZero() {
				s.readyTimestamp = condition.LastTransitionTime.Time
			}
		}
	}

	s.updateContainers(pod)
}

func (s *podStatistic) updateContainers(pod *corev1.Pod) {
	now := s.timeSource.Now()

	logger := s.logger()

	for _, containerStatus := range pod.Status.InitContainerStatuses {
		containerStatistic, ok := s.initContainers[containerStatus.Name]
		if !ok {
			logger.Error().Msgf(
				"Init container statistic does not exist for %s", containerStatus.Name,
			)

			continue
		}
		containerStatistic.update(now, containerStatus)
	}

	for _, containerStatus := range pod.Status.ContainerStatuses {
		containerStatistic, ok := s.containers[containerStatus.Name]
		if !ok {
			logger.Error().Msgf(
				"Container statistic does not exist for %s", containerStatus.Name,
			)

			continue
		}
		containerStatistic.update(now, containerStatus)
	}

	s.report()
}
