package statistics

import (
	"time"

	"github.com/rs/zerolog"
)

type imagePullStatistic struct {
	container *containerStatistic

	startedAt  time.Time
	finishedAt time.Time
}

func (s imagePullStatistic) log(message string) {
	metrics := zerolog.Dict()
	metrics.Bool("init_container", s.container.initContainer)
	if !s.finishedAt.IsZero() && !s.startedAt.IsZero() {
		metrics.Float64(
			"image_pull_duration",
			s.finishedAt.Sub(s.startedAt).Seconds(),
		)
	}

	logger :=
		s.container.logger().
			With().
			Str("kube_transition_metric_type", "image_pull").
			Dict("kube_transition_metrics", metrics).
			Logger()
	logger.Info().Msg(message)
}
