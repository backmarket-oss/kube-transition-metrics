package testhelpers

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/logging"
	"github.com/stretchr/testify/require"
)

// MetricWriter is an [io.Writer] that validates each written JSON document against the
// kube_transition_metrics JSON schema and captures the data for inspection.
type MetricWriter struct {
	t        *testing.T
	buf      *bytes.Buffer
	validate io.Writer
}

// NewMetricWriter creates a new MetricWriter for the given test.
func NewMetricWriter(t *testing.T) *MetricWriter {
	t.Helper()

	return &MetricWriter{
		t:        t,
		buf:      &bytes.Buffer{},
		validate: logging.NewValidationWriter(),
	}
}

// Write validates the data against the schema and captures it in the buffer.
// If validation fails, the test is marked as failed via [testing.T.Errorf].
func (w *MetricWriter) Write(data []byte) (int, error) {
	w.t.Helper()

	_, err := w.validate.Write(data)
	if err != nil {
		w.t.Errorf("metric output failed schema validation: %v", err)
	}

	n, err := w.buf.Write(data)

	//nolint:wrapcheck
	return n, err
}

// DecodeMetricOutput parses the captured metric output and returns the kube_transition_metrics
// objects from each JSON line.
func DecodeMetricOutput(t *testing.T, writer *MetricWriter) []map[string]any {
	t.Helper()

	var metrics []map[string]any

	for line := range strings.SplitSeq(writer.buf.String(), "\n") {
		if line == "" {
			continue
		}

		var record map[string]any
		require.NoError(t, json.Unmarshal([]byte(line), &record), "failed to decode metric output line")

		metric, ok := record["kube_transition_metrics"].(map[string]any)
		require.True(t, ok, "expected kube_transition_metrics to be an object in: %s", line)

		metrics = append(metrics, metric)
	}

	return metrics
}
