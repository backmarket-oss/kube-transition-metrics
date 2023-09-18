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

func (cs containerStatistic) event() *zerolog.Event {
	event := zerolog.Dict()

	event.Bool("init_container", cs.initContainer)
	if !cs.startedTimestamp.IsZero() {
		event.Float64("started_latency",
			cs.startedTimestamp.Sub(cs.pod.creationTimestamp).Seconds())
	}
	if !cs.readyTimestamp.IsZero() {
		event.Float64("ready_latency",
			cs.readyTimestamp.Sub(cs.pod.creationTimestamp).Seconds())
	}
	if !cs.runningTimestamp.IsZero() {
		event.Float64("running_latency",
			cs.runningTimestamp.Sub(cs.pod.creationTimestamp).Seconds())
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
