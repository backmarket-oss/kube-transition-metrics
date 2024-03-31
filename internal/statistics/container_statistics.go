package statistics

import (
	"time"

	"github.com/rs/zerolog"
	corev1 "k8s.io/api/core/v1"
)

type containerStatistic struct {
	name          string
	initContainer bool
	pod           *podStatistic
	imagePull     imagePullStatistic

	// Previous init container, null if first init container or non-init container
	previous *containerStatistic

	// The timestamp for when the container first turned Running.
	runningTimestamp time.Time

	// The timestamp for when the container first started (postStart hook run and
	// startupProbe passed).
	startedTimestamp time.Time

	// The timestamp for when the container first turned Ready (readinessProbe
	// passed).
	readyTimestamp time.Time
}

func newContainerStatistic(
	statistic *podStatistic,
	initContainer bool,
	container corev1.Container,
) *containerStatistic {
	containerStatistic := &containerStatistic{
		name:          container.Name,
		initContainer: initContainer,
		pod:           statistic,
	}
	containerStatistic.imagePull.container = containerStatistic

	return containerStatistic
}

func (cs containerStatistic) logger() zerolog.Logger {
	return cs.pod.logger().With().
		Str("container_name", cs.name).
		Logger()
}

func (cs containerStatistic) appendInitFields(event *zerolog.Event) {
	if !cs.runningTimestamp.IsZero() && cs.previous != nil && !cs.previous.readyTimestamp.IsZero() {
		event.Dur("previous_to_running_seconds", cs.runningTimestamp.Sub(cs.previous.readyTimestamp))
	}
}

func (cs containerStatistic) appendNonInitFields(event *zerolog.Event) {
	if !cs.runningTimestamp.IsZero() && !cs.pod.scheduledTimestamp.IsZero() {
		event.Dur("initialized_to_running_seconds", cs.runningTimestamp.Sub(cs.pod.scheduledTimestamp))
	}
}

func (cs containerStatistic) event() *zerolog.Event {
	event := zerolog.Dict()

	event.Bool("init_container", cs.initContainer)
	if cs.initContainer {
		cs.appendInitFields(event)
	} else {
		cs.appendNonInitFields(event)
	}

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

	return event
}

func (cs containerStatistic) report() {
	logger := cs.logger()

	eventLogger := logger.Output(metricOutput).With().
		Str("kube_transition_metric_type", "container").
		Dict("kube_transition_metrics", cs.event()).
		Logger()
	eventLogger.Log().Msg("")
}

func (cs containerStatistic) logContainerStatus(status corev1.ContainerStatus) {
	logger := cs.logger()

	switch {
	case status.State.Waiting != nil:
		logger := logger.With().
			Str("container_state", "Waiting").
			Str("waiting_reason", status.State.Waiting.Reason).
			Str("waiting_message", status.State.Waiting.Message).
			Logger()
		logger.Debug().Msg("Container is Waiting.")
	case status.State.Running != nil:
		logger := logger.With().
			Str("container_state", "Running").
			Str("started_at", status.State.Running.StartedAt.String()).
			Logger()
		logger.Debug().Msg("Container is Running.")
	case status.State.Terminated != nil:
		logger := logger.With().
			Str("container_state", "Terminated").
			Str("terminated_reason", status.State.Terminated.Reason).
			Str("terminated_message", status.State.Terminated.Message).
			Int32("exit_code", status.State.Terminated.ExitCode).
			Int32("signal", status.State.Terminated.Signal).
			Logger()
		logger.Debug().Msg("Container is Terminated.")
	}
}

func (cs *containerStatistic) update(
	now time.Time,
	status corev1.ContainerStatus,
) {
	cs.logContainerStatus(status)

	if cs.runningTimestamp.IsZero() && status.State.Running != nil {
		cs.runningTimestamp = now
	}
	if cs.startedTimestamp.IsZero() && status.Started != nil && *status.Started {
		cs.startedTimestamp = now
	}
	if cs.readyTimestamp.IsZero() && status.Ready {
		cs.readyTimestamp = now
	}
}
