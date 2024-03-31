package logging

import (
	"time"

	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
	"github.com/rs/zerolog"
)

// Configure configures zerolog with the required global settings.
func Configure() {
	zerolog.DurationFieldInteger = false
	zerolog.DurationFieldUnit = time.Second
}

// SetOptions configures zerolog global settings based on user-configured options.
func SetOptions(options *options.Options) {
	zerolog.SetGlobalLevel(options.LogLevel)
}
