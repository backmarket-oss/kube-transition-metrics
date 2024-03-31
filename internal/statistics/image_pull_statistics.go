package statistics

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type imagePullStatistic struct {
	container *containerStatistic

	alreadyPresent    bool
	startedTimestamp  time.Time
	finishedTimestamp time.Time
}

func (s imagePullStatistic) log(message string) {
	imagePullMetrics := zerolog.Dict()

	imagePullMetrics.Str("container_name", s.container.name)
	imagePullMetrics.Bool("already_present", s.alreadyPresent)
	if !s.startedTimestamp.IsZero() {
		imagePullMetrics.Time("started_timestamp", s.startedTimestamp)
	}
	if !s.finishedTimestamp.IsZero() {
		imagePullMetrics.Time("finished_timestamp", s.finishedTimestamp)
		if !s.startedTimestamp.IsZero() {
			imagePullMetrics.Dur("duration_seconds", s.finishedTimestamp.Sub(s.startedTimestamp))
		}
	}

	metrics := zerolog.Dict()
	metrics.Str("type", "image_pull")
	metrics.Dict("image_pull", imagePullMetrics)
	metrics.Str("kube_namespace", s.container.pod.namespace)
	metrics.Str("pod_name", s.container.pod.name)

	logger :=
		log.
			Output(metricOutput).
			With().
			Dict("kube_transition_metrics", metrics).
			Logger()
	logger.Log().Msg(message)
}
