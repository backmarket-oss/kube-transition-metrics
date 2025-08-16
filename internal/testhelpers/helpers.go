package testhelpers

import (
	"testing"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/logging"
)

// ConfigureLogging sets up the global logging configuration for the tests.
// It forbids parallel execution of tests that use this function.
func ConfigureLogging(t *testing.T) {
	t.Helper()

	t.Cleanup(func() {
		logging.Unconfigure()
	})
	logging.Configure()
}
