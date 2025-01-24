package statistics

import (
	"io"
	"os"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/logging"
	"github.com/rs/zerolog"
)

//nolint:gochecknoglobals
var metricOutput io.Writer = zerolog.MultiLevelWriter(
	os.Stdout,
	logging.NewValidationWriter(),
)
