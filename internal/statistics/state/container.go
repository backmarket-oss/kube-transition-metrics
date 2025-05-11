package state

import (
	"io"
	"time"

	"github.com/Izzette/go-safeconcurrency/eventloop/snapshot"
	"github.com/rs/zerolog"
	corev1 "k8s.io/api/core/v1"
)

// ContainerStatistic holds the transition statistics for a container in a pod.
// ContainerStatistic is immutable, all the methods return a new instance of the struct.
// Do not lose track of the returned instance, it should be assigned to the containing structure.
type ContainerStatistic struct {
	// name is the name of the container.
	name string

	// runningTimestamp for when the container first turned Running.
	runningTimestamp time.Time

	// startedTimestamp for when the container first started (postStart hook run and startupProbe passed).
	startedTimestamp time.Time

	// readyTimestamp for when the container first turned Ready (readinessProbe passed).
	readyTimestamp time.Time
}

// Copy implements [github.com/Izzette/go-safeconcurrency/types.Copyable.Copy].
func (cs *ContainerStatistic) Copy() *ContainerStatistic {
	return snapshot.CopyPtr(cs)
}

// logger returns a logger scoped to the container statistic.
//
// TODO(Izzette): Replace with [log.Ctx] / [zerolog.Ctx] / [zerolog.Event.Ctx].
func (cs *ContainerStatistic) logger(logger zerolog.Logger) zerolog.Logger {
	return logger.With().
		Str("container_name", cs.name).
		Logger()
}

// event appends the container statistics to the event.
func (cs *ContainerStatistic) event(event *zerolog.Event) {
	if !cs.runningTimestamp.IsZero() {
		event.Time("running_timestamp", cs.runningTimestamp)
	}
	if !cs.startedTimestamp.IsZero() {
		event.Time("started_timestamp", cs.startedTimestamp)
		if !cs.runningTimestamp.IsZero() {
			event.Dur("running_to_started_seconds", cs.startedTimestamp.Sub(cs.runningTimestamp))
		}
	}
	if !cs.readyTimestamp.IsZero() {
		event.Time("ready_timestamp", cs.readyTimestamp)

		// As init containers do not supported startup, liveliness, or readiness probes the Started container status field is
		// not set for init containers.
		// Instead, readiness represents the time an init container has excited successfully,allowing the following containers
		// to proceed.
		// Given this, presenting both running_to_ready_seconds and started_to_ready_seconds is useful to cover the differing
		// meanings for both container types.
		// See: https://github.com/kubernetes/website/blob/b397a8f/content/en/docs/concepts/workloads/pods/init-containers.md
		if !cs.runningTimestamp.IsZero() {
			event.Dur("running_to_ready_seconds", cs.readyTimestamp.Sub(cs.runningTimestamp))
		}
		if !cs.startedTimestamp.IsZero() {
			event.Dur("started_to_ready_seconds", cs.readyTimestamp.Sub(cs.runningTimestamp))
		}
	}
}

// logContainerStatus logs the container status to the logger.
func (cs *ContainerStatistic) logContainerStatus(pod *PodStatistic, status corev1.ContainerStatus) {
	// TODO(Izzette): Replace with [log.Ctx] / [zerolog.Ctx] / [zerolog.Event.Ctx].
	logger := cs.logger(pod.logger())

	switch {
	case status.State.Waiting != nil:
		logger.Debug().
			Str("container_state", "Waiting").
			Str("waiting_reason", status.State.Waiting.Reason).
			Str("waiting_message", status.State.Waiting.Message).
			Msg("Container is Waiting.")
	case status.State.Running != nil:
		logger.Debug().
			Str("container_state", "Running").
			Str("started_at", status.State.Running.StartedAt.String()).
			Msg("Container is Running.")
	case status.State.Terminated != nil:
		logger.Debug().
			Str("container_state", "Terminated").
			Str("terminated_reason", status.State.Terminated.Reason).
			Str("terminated_message", status.State.Terminated.Message).
			Int32("exit_code", status.State.Terminated.ExitCode).
			Int32("signal", status.State.Terminated.Signal).
			Msg("Container is Terminated.")
	}
}

// update updates the containers statistic based on the latest Kubernetes container status.
// update returns a new instance of the container statistic with the updated fields.
func (cs *ContainerStatistic) update(
	now time.Time,
	status corev1.ContainerStatus,
	pod *PodStatistic,
) *ContainerStatistic {
	// We will return a copy of the container statistic, so that we can safely update the container statistic in the event
	// loop.
	// As this type is immutable, we should shadow the receiver.
	cs = cs.Copy()

	cs.logContainerStatus(pod, status)

	if cs.runningTimestamp.IsZero() && status.State.Running != nil {
		cs.runningTimestamp = now
	}
	if cs.startedTimestamp.IsZero() && status.Started != nil && *status.Started {
		cs.startedTimestamp = now
	}
	if cs.readyTimestamp.IsZero() && status.Ready {
		cs.readyTimestamp = now
	}

	return cs
}

// InitContainerStatistic holds the transition statistics for an init container in a pod.
type InitContainerStatistic struct {
	*ContainerStatistic
}

// Report reports the container statistic to the output writer.
func (cs *InitContainerStatistic) Report(
	output io.Writer,
	pod *corev1.Pod,
	podStatistic *PodStatistic,
	previous *ContainerStatistic,
) {
	logger := cs.logger(podStatistic.logger())
	container := findContainer(cs.name, pod.Spec.InitContainers)

	metrics := zerolog.Dict().
		Func(commonPodLabels(pod)).
		Func(commonContainerLabels(&logger, container)).
		Dict("container", cs.event(previous))

	logMetrics(output, "container", metrics, "")
}

// Update updates the init container statistic based on the latest Kubernetes container status.
func (cs *InitContainerStatistic) Update(
	now time.Time,
	status corev1.ContainerStatus,
	pod *PodStatistic,
) *InitContainerStatistic {
	// We will return a copy of the container statistic, so that we can safely update the container statistic in the event
	// op.
	// As this type is immutable, we should shadow the receiver.
	cs = snapshot.CopyPtr(cs)

	cs.ContainerStatistic = cs.update(now, status, pod)

	return cs
}

// event returns the event dictionary for the init container statistic.
func (cs *InitContainerStatistic) event(previous *ContainerStatistic) *zerolog.Event {
	event := zerolog.Dict()
	event.Bool("init_container", true)
	if !cs.runningTimestamp.IsZero() && previous != nil && !previous.readyTimestamp.IsZero() {
		event.Dur("previous_to_running_seconds", cs.runningTimestamp.Sub(previous.readyTimestamp))
	}
	cs.ContainerStatistic.event(event)

	return event
}

// NonInitContainerStatistic holds the transition statistics for a non-init container in a pod.
type NonInitContainerStatistic struct {
	*ContainerStatistic
}

// Report reports the container statistic to the output writer.
func (cs *NonInitContainerStatistic) Report(output io.Writer, pod *corev1.Pod, podStatistic *PodStatistic) {
	logger := cs.logger(podStatistic.logger())
	container := findContainer(cs.name, pod.Spec.Containers)
	if container == nil {
		logger.Panic().Msg("container not found")
	}

	metrics := zerolog.Dict().
		Func(commonPodLabels(pod)).
		Func(commonContainerLabels(&logger, container)).
		Dict("container", cs.event(podStatistic))

	logMetrics(output, "container", metrics, "")
}

// Update updates the non-init container statistic based on the latest Kubernetes container status.
func (cs *NonInitContainerStatistic) Update(
	now time.Time,
	status corev1.ContainerStatus,
	pod *PodStatistic,
) *NonInitContainerStatistic {
	// We will return a copy of the container statistic, so that we can safely update the container statistic in the event
	// loop.
	// As this type is immutable, we should shadow the receiver.
	cs = snapshot.CopyPtr(cs)

	cs.ContainerStatistic = cs.update(now, status, pod)

	return cs
}

// event returns the event dictionary for the non-init container statistic.
func (cs *NonInitContainerStatistic) event(pod *PodStatistic) *zerolog.Event {
	event := zerolog.Dict()
	event.Bool("init_container", false)
	if !cs.runningTimestamp.IsZero() && !pod.scheduledTimestamp.IsZero() {
		event.Dur("initialized_to_running_seconds", cs.runningTimestamp.Sub(pod.scheduledTimestamp))
	}

	cs.ContainerStatistic.event(event)

	return event
}
