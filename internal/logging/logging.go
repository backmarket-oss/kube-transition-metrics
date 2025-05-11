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

// Unconfigure restores zerolog to its default settings.
func Unconfigure() {
	zerolog.DurationFieldInteger = true
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.SetGlobalLevel(zerolog.Level(0))
}

// SetOptions configures zerolog global settings based on user-configured options.
func SetOptions(options *options.Options) {
	zerolog.SetGlobalLevel(options.LogLevel)
}
