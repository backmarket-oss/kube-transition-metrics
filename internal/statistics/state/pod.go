package state

import (
	"io"
	"time"

	"github.com/Izzette/go-safeconcurrency/eventloop/snapshot"
	"github.com/benbjohnson/immutable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// PodStatistic holds the transition statistics for a pod.
// PodStatistic is immutable, all the methods return a new instance of the struct.
// Do not lose track of the returned instance, it should be assigned to the containing structure.
type PodStatistic struct {
	name      string
	namespace string

	// The timestamp for when the pod was created, same as timestamp of when pod
	// was in Pending and containers were Waiting.
	creationTimestamp time.Time

	// The timestamp for when the pod was scheduled.
	scheduledTimestamp time.Time

	// The timestamp for when the pod was initialized.
	initializedTimestamp time.Time

	// The timestamp for when the pod first turned Ready.
	readyTimestamp time.Time

	initContainers *immutable.Map[string, *InitContainerStatistic]
	containers     *immutable.Map[string, *NonInitContainerStatistic]
}

// NewPodStatistic creates a new PodStatistic instance populated with the containers in the pod.
func NewPodStatistic(now time.Time, pod *corev1.Pod) *PodStatistic {
	podStatistic := &PodStatistic{
		name:              pod.Name,
		namespace:         pod.Namespace,
		creationTimestamp: pod.CreationTimestamp.Time,
	}

	initContainers := immutable.NewMapBuilder[string, *InitContainerStatistic](nil)
	for _, container := range pod.Spec.InitContainers {
		initContainers.Set(container.Name, &InitContainerStatistic{
			&ContainerStatistic{
				name: container.Name,
			},
		})
	}
	podStatistic.initContainers = initContainers.Map()

	containers := immutable.NewMapBuilder[string, *NonInitContainerStatistic](nil)
	for _, container := range pod.Spec.Containers {
		containers.Set(container.Name, &NonInitContainerStatistic{
			&ContainerStatistic{
				name: container.Name,
			},
		})
	}
	podStatistic.containers = containers.Map()

	return podStatistic.Update(now, pod)
}

// Report reports the pod statistic to the given output writer.
func (s *PodStatistic) Report(output io.Writer, pod *corev1.Pod) {
	logger := s.logger()

	metrics := zerolog.Dict().
		Func(commonPodLabels(pod)).
		Dict("pod", s.event())
	logMetrics(output, "pod", metrics, "")

	initContainers := s.initContainers.Iterator()
	var previous *ContainerStatistic
	for !initContainers.Done() {
		_, containerStatistics, ok := initContainers.Next()
		if !ok {
			// This should never happen as we're checking `.Done()` on the iterator.
			logger.Panic().Msg("init container statistics not found")
		}

		containerStatistics.Report(output, pod, s, previous)
		previous = containerStatistics.ContainerStatistic
	}

	containers := s.containers.Iterator()
	for !containers.Done() {
		_, containerStatistics, ok := containers.Next()
		if !ok {
			// This should never happen as we're checking `.Done()` on the iterator.
			logger.Panic().Msg("container statistics not found")
		}

		containerStatistics.Report(output, pod, s)
	}
}

// Update updates the pod statistic with the provided pod.
// It returns a new instance of the pod statistic with the updated values.
func (s *PodStatistic) Update(now time.Time, pod *corev1.Pod) *PodStatistic {
	// We will return a copy of the pod statistic, so that we can safely update the pod statistic in the event loop.
	// As this type is immutable, we should shadow the receiver.
	s = s.Copy()

	logger := s.logger()

	for _, condition := range pod.Status.Conditions {
		logger.Debug().Msgf("Pod condition: %+v", condition)
		if condition.Status != corev1.ConditionTrue {
			continue
		}

		// TODO: include core/v1.ContainersReady and core/v1.DisruptionTarget
		switch condition.Type { //nolint:exhaustive
		case corev1.PodScheduled:
			if s.scheduledTimestamp.IsZero() {
				s.scheduledTimestamp = condition.LastTransitionTime.Time
			}
		case corev1.PodInitialized:
			if s.initializedTimestamp.IsZero() {
				s.initializedTimestamp = condition.LastTransitionTime.Time
			}
		case corev1.PodReady:
			if s.readyTimestamp.IsZero() {
				s.readyTimestamp = condition.LastTransitionTime.Time
			}
		}
	}

	s = s.updateContainers(now, pod)

	return s
}

// Copy implements [github.com/Izzette/go-safeconcurrency/types.Copyable.Copy].
func (s *PodStatistic) Copy() *PodStatistic {
	return snapshot.CopyPtr(s)
}

// logger returns a logger for the pod statistic.
//
// TODO(Izzette): replace with [zerolog.Ctx].
func (s *PodStatistic) logger() zerolog.Logger {
	return log.With().
		Str("kube_namespace", s.namespace).
		Str("pod_name", s.name).
		Logger()
}

// event returns the event dictionary for the pod statistic.
func (s *PodStatistic) event() *zerolog.Event {
	event := zerolog.Dict()

	event.Time("creation_timestamp", s.creationTimestamp)
	if !s.scheduledTimestamp.IsZero() {
		event.Time("scheduled_timestamp", s.scheduledTimestamp)
		event.Dur("creation_to_scheduled_seconds", s.scheduledTimestamp.Sub(s.creationTimestamp))
	}
	if !s.initializedTimestamp.IsZero() {
		event.Time("initialized_timestamp", s.initializedTimestamp)
		event.Dur("creation_to_initialized_seconds", s.initializedTimestamp.Sub(s.creationTimestamp))
		if !s.scheduledTimestamp.IsZero() {
			event.Dur("scheduled_to_initialized_seconds", s.initializedTimestamp.Sub(s.scheduledTimestamp))
		}
	}
	if !s.readyTimestamp.IsZero() {
		event.Time("ready_timestamp", s.readyTimestamp)
		event.Dur("creation_to_ready_seconds", s.readyTimestamp.Sub(s.creationTimestamp))
		if !s.initializedTimestamp.IsZero() {
			event.Dur("initialized_to_ready_seconds", s.readyTimestamp.Sub(s.initializedTimestamp))
		}
	}

	return event
}

// updateContainers updates the pod statistic with the provided pod.
// It returns a new instance of the pod statistic with the updated values.
func (s *PodStatistic) updateContainers(now time.Time, pod *corev1.Pod) *PodStatistic {
	// We will return a copy of the pod statistic, so that we can safely update the pod statistic in the event loop.
	// As this type is immutable, we should shadow the receiver.
	s = s.Copy()

	logger := s.logger()

	for _, containerStatus := range pod.Status.InitContainerStatuses {
		containerStatistic, ok := s.initContainers.Get(containerStatus.Name)
		if !ok {
			logger.Error().Msgf(
				"Init container statistic does not exist for %s", containerStatus.Name,
			)

			continue
		}
		containerStatistic = containerStatistic.Update(now, containerStatus, s)
		s.initContainers = s.initContainers.Set(containerStatus.Name, containerStatistic)
	}

	for _, containerStatus := range pod.Status.ContainerStatuses {
		containerStatistic, ok := s.containers.Get(containerStatus.Name)
		if !ok {
			logger.Error().Msgf(
				"Container statistic does not exist for %s", containerStatus.Name,
			)

			continue
		}

		containerStatistic = containerStatistic.Update(now, containerStatus, s)
		s.containers = s.containers.Set(containerStatus.Name, containerStatistic)
	}

	return s
}

// PodStatistics implements [github.com/Izzette/go-safeconcurrency/api/types.Copyable].
type PodStatistics struct {
	blacklistUIDs immutable.Set[types.UID]
	statistics    *immutable.Map[types.UID, *PodStatistic]
}

// NewPodStatistics creates a new PodStatistics with the provided blacklist.
func NewPodStatistics(blacklistUIDs []types.UID) *PodStatistics {
	return &PodStatistics{
		blacklistUIDs: immutable.NewSet(nil, blacklistUIDs...),
		statistics:    &immutable.Map[types.UID, *PodStatistic]{},
	}
}

// Copy implements [github.com/Izzette/go-safeconcurrency/api/types.Copyable].
func (eh *PodStatistics) Copy() *PodStatistics {
	// As the StatisticState is updated Copy-on-Write, we can just return a copy of the underlying struct.
	return snapshot.CopyPtr(eh)
}

// Len returns the number of pod statistics.
func (eh *PodStatistics) Len() int {
	return eh.statistics.Len()
}

// Get returns the pod statistic for the given UID. If the UID is not found,.
func (eh *PodStatistics) Get(uid types.UID) (*PodStatistic, bool) {
	// The pod statistic is immutable, so we can just return it.
	return eh.statistics.Get(uid)
}

// Set sets the pod statistic for the given UID.
func (eh *PodStatistics) Set(uid types.UID, statistic *PodStatistic) *PodStatistics {
	eh = eh.Copy()
	eh.statistics = eh.statistics.Set(uid, statistic)

	return eh
}

// Delete deletes the pod statistic for the given UID, if it exists.
func (eh *PodStatistics) Delete(uid types.UID) *PodStatistics {
	eh = eh.Copy()
	eh.statistics = eh.statistics.Delete(uid)

	return eh
}

// IsBlacklisted checks if the given UID is in the blacklist.
func (eh *PodStatistics) IsBlacklisted(uid types.UID) bool {
	return eh.blacklistUIDs.Has(uid)
}
