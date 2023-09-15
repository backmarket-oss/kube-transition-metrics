package statistics

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
)

type podStatistic struct {
	Initialized bool

	Name       string
	Namespace  string
	TimeSource timeSource

	ImagePullCollector imagePullCollector

	// The timestamp for when the pod was created, same as timestamp of when pod
	// was in Pending and containers were Waiting.
	CreationTimestamp time.Time

	// The timestamp for when the pod was scheduled.
	ScheduledTimestamp time.Time

	// The timestamp for when the pod was initialized.
	InitializedTimestamp time.Time

	// The timestamp for when the pod first turned Ready.
	ReadyTimestamp time.Time

	InitContainers map[string]*containerStatistic
	Containers     map[string]*containerStatistic
}

func (s *podStatistic) initialize(pod *corev1.Pod) {
	if s.Initialized {
		return
	}
	s.Initialized = true
	s.TimeSource = realTimeSource{}
	s.Name = pod.Name
	s.Namespace = pod.Namespace
	s.CreationTimestamp = pod.CreationTimestamp.Time
	s.InitContainers = make(map[string]*containerStatistic)
	s.Containers = make(map[string]*containerStatistic)

	for _, container := range pod.Spec.InitContainers {
		s.InitContainers[container.Name] = newContainerStatistic(s, true, container)
	}
	for _, container := range pod.Spec.Containers {
		s.Containers[container.Name] = newContainerStatistic(s, false, container)
	}
}

func (s podStatistic) logger() zerolog.Logger {
	return log.With().
		Str("kube_namespace", s.Namespace).
		Str("pod_name", s.Name).
		Logger()
}

func (s podStatistic) event() *zerolog.Event {
	event := zerolog.Dict()

	if !s.ScheduledTimestamp.IsZero() {
		event.Float64(
			"scheduled_latency",
			s.ScheduledTimestamp.Sub(s.CreationTimestamp).Seconds())
	}
	if !s.InitializedTimestamp.IsZero() {
		event.Float64(
			"initialized_latency",
			s.InitializedTimestamp.Sub(s.CreationTimestamp).Seconds())
	}
	if !s.ReadyTimestamp.IsZero() {
		event.Float64(
			"ready_latency",
			s.ReadyTimestamp.Sub(s.CreationTimestamp).Seconds())
	}

	return event
}

func (s podStatistic) report() {
	logger := s.logger()

	eventLogger := logger.With().
		Str("kube_transition_metric_type", "pod").
		Dict("kube_transition_metrics", s.event()).Logger()
	eventLogger.Info().Msg("")

	for _, containerStatistics := range s.InitContainers {
		containerStatistics.report()
	}
	for _, containerStatistics := range s.Containers {
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
			if s.ScheduledTimestamp.IsZero() {
				s.ScheduledTimestamp = condition.LastTransitionTime.Time
			}
		case corev1.PodInitialized:
			// Pod Initialized occursafter all images pulled, no need to continue to
			// track

			if s.InitializedTimestamp.IsZero() {
				go s.ImagePullCollector.cancel("pod_initialized")
				s.InitializedTimestamp = condition.LastTransitionTime.Time
			}
		case corev1.PodReady:
			if s.ReadyTimestamp.IsZero() {
				s.ReadyTimestamp = condition.LastTransitionTime.Time
			}
		}
	}

	s.updateContainers(pod)
}

func (s *podStatistic) updateContainers(pod *corev1.Pod) {
	now := s.TimeSource.Now()

	logger := s.logger()

	for _, containerStatus := range pod.Status.InitContainerStatuses {
		containerStatistic, ok := s.InitContainers[containerStatus.Name]
		if !ok {
			logger.Error().Msgf(
				"Init container statistic does not exist for %s", containerStatus.Name,
			)

			continue
		}
		containerStatistic.update(now, containerStatus)
	}

	for _, containerStatus := range pod.Status.ContainerStatuses {
		containerStatistic, ok := s.Containers[containerStatus.Name]
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
