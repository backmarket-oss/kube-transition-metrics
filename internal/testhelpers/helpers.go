package testhelpers

import (
	"testing"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/logging"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
)

// ConfigureLogging sets up the global logging configuration for the tests.
// It forbids parallel execution of tests that use this function.
func ConfigureLogging(t *testing.T, opts *options.Options) {
	t.Helper()

	t.Cleanup(func() {
		logging.Unconfigure()
	})
	logging.Configure()
	logging.SetOptions(opts)
}
