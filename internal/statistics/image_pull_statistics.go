package statistics

import (
	"time"

	"github.com/rs/zerolog"
)

type imagePullStatistic struct {
	container *containerStatistic

	startedTimestamp  time.Time
	finishedTimestamp time.Time
}

func (s imagePullStatistic) log(message string) {
	metrics := zerolog.Dict()

	if !s.startedTimestamp.IsZero() {
		metrics.Time("started_timestamp", s.startedTimestamp)
	}
	if !s.finishedTimestamp.IsZero() {
		metrics.Time("finished_timestamp", s.finishedTimestamp)
		if !s.startedTimestamp.IsZero() {
			metrics.Dur("duration_seconds", s.finishedTimestamp.Sub(s.startedTimestamp))
		}
	}

	logger :=
		s.container.logger().
			Output(metricOutput).
			With().
			Str("kube_transition_metric_type", "image_pull").
			Dict("kube_transition_metrics", metrics).
			Logger()
	logger.Log().Msg(message)
}
