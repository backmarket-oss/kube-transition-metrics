package statistics

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
)

type PodStatistic struct {
	Initialized bool

	Name       string
	Namespace  string
	TimeSource TimeSource

	ImagePullCollector ImagePullCollector

	// The timestamp for when the pod was created, same as timestamp of when pod
	// was in Pending and containers were Waiting.
	CreationTimestamp time.Time

	// The timestamp for when the pod was scheduled.
	ScheduledTimestamp time.Time

	// The timestamp for when the pod was initialized.
	InitializedTimestamp time.Time

	// The timestamp for when the pod first turned Ready.
	ReadyTimestamp time.Time

	InitContainers map[string]*ContainerStatistic
	Containers     map[string]*ContainerStatistic
}

func (s *PodStatistic) Initialize(pod *corev1.Pod) {
	if s.Initialized {
		return
	}
	s.Initialized = true
	s.TimeSource = RealTimeSource{}
	s.Name = pod.Name
	s.Namespace = pod.Namespace
	s.CreationTimestamp = pod.CreationTimestamp.Time
	s.InitContainers = make(map[string]*ContainerStatistic)
	s.Containers = make(map[string]*ContainerStatistic)

	for _, container := range pod.Spec.InitContainers {
		s.InitContainers[container.Name] = NewContainerStatistic(s, true, container)
	}
	for _, container := range pod.Spec.Containers {
		s.Containers[container.Name] = NewContainerStatistic(s, false, container)
	}
}

func (s PodStatistic) logger() zerolog.Logger {
	return log.With().
		Str("kube_namespace", s.Namespace).
		Str("pod_name", s.Name).
		Logger()
}

func (s PodStatistic) event() *zerolog.Event {
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

func (s PodStatistic) report() {
	logger := s.logger()

	event_logger := logger.With().
		Str("kube_transition_metric_type", "pod").
		Dict("kube_transition_metrics", s.event()).Logger()
	event_logger.Info().Msg("")

	for _, container_statistics := range s.InitContainers {
		container_statistics.report()
	}
	for _, container_statistics := range s.Containers {
		container_statistics.report()
	}
}

func (s *PodStatistic) update(pod *corev1.Pod) {
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

func (s *PodStatistic) updateContainers(pod *corev1.Pod) {
	now := s.TimeSource.Now()

	logger := s.logger()

	for _, container_status := range pod.Status.InitContainerStatuses {
		container_statistic, ok := s.InitContainers[container_status.Name]
		if !ok {
			logger.Error().Msgf(
				"Init container statistic does not exist for %s", container_status.Name,
			)

			continue
		}
		container_statistic.update(now, container_status)
	}

	for _, container_status := range pod.Status.ContainerStatuses {
		container_statistic, ok := s.Containers[container_status.Name]
		if !ok {
			logger.Error().Msgf(
				"Container statistic does not exist for %s", container_status.Name,
			)

			continue
		}
		container_statistic.update(now, container_status)
	}

	s.report()
}
