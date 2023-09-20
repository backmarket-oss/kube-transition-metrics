package statistics

import (
	"io"
	"os"
)

//nolint:gochecknoglobals
var metricOutput io.Writer = os.Stdout
