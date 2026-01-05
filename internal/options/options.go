package options

import (
	"log"

	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
)

// Options contains the options for the controller.
type Options struct {
	// ListenAddress is the host and port for the HTTP server delivering prometheus metrics and pprof profiling.
	ListenAddress string
	// KubeconfigPath is the path to the kube configuration file.
	KubeconfigPath string
	// ImagePullCancelDelay is the delay before canceling an image pull routine to ensure all events related to the pod
	// have been processed.
	ImagePullCancelDelay float64
	// KubeWatchTimeout is the timeout for the Kubernetes Watch API.
	KubeWatchTimeout int64
	// KubeWatchMaxEvents is the maximum number of events to receive from the Kubernetes Watch API per response.
	KubeWatchMaxEvents int64
	// StatisticEventQueueLength is the maximum number of queued statistic events.
	//
	// TODO(Izzette): consider splitting into PodStatisticEventQueueLength and ImagePullStatisticEventQueueLength
	StatisticEventQueueLength int
	// EmitPartialStatistics enables emitting statistics for pods that have not yet become Ready and image pulls that have
	// not yet completed.
	EmitPartialStatistics bool
	// LogLevel is the global logging level.
	LogLevel zerolog.Level
}

// Parse parses the options and returns them as a pointer to an Options struct.
//
//nolint:funlen
func Parse() *Options {
	options := Options{}
	flag.StringVar(
		&options.ListenAddress,
		"listen-address",
		"127.0.0.1:8080",
		"The host and port for HTTP server delivering prometheus metrics over "+
			"`/metrics` and pprof profiling over `/debug/pprof` endpoints.")
	flag.StringVar(
		&options.KubeconfigPath,
		"kubeconfig-path",
		"",
		"The path to the kube configuration file, if it's not set the value of "+
			"`$KUBECONFIG` will be used, if that's not set `$HOME/.kube/config` will "+
			"be used.")
	flag.Float64Var(
		&options.ImagePullCancelDelay,
		"image-pull-cancel-delay",
		3,
		"The delay (in seconds) before canceling an image pull collector routine to ensure all events related to the pod "+
			"have been processed. (ADVANCED)")
	flag.Int64Var(
		&options.KubeWatchTimeout,
		"kube-watch-timeout",
		60,
		"The Kubernetes Watch API timeout (ADVANCED)")
	flag.Int64Var(
		&options.KubeWatchMaxEvents,
		"kube-watch-max-events",
		100,
		"The Kubernetes Watch maximum events per response (ADVANCED)")
	flag.IntVar(
		&options.StatisticEventQueueLength,
		"statistic-event-queue-length",
		1000,
		"The maximum number of queued statistic events (ADVANCED)")
	flag.BoolVar(
		&options.EmitPartialStatistics,
		"emit-partial",
		false,
		"Emit partial statistics for pods that have not yet become Ready and image pulls that have not yet completed. When "+
			"set to false, pods that never become Ready and image pulls that never complete will not be included in the "+
			"statistics. Partial statistics will always be emitted for pods that are deleted before they become Ready. When "+
			"set to true, multiple statistics will be emitted for the same pod/image pull. (ADVANCED)")

	logLevel := flag.String(
		"log-level",
		"INFO",
		`The global logging level, one of "trace", "debug", "info", "warn", `+
			`"error", "fatal", "panic", "disabled", or "" (empty string). This option's`+
			`values are case-insensitive. Setting a value of "disabled" will result in`+
			`no metrics being emitted.`)

	flag.Parse()

	logLevelParsed, err := zerolog.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("Invalid value for --log-level (%s): %q\n", *logLevel, err)
	} else {
		options.LogLevel = logLevelParsed
	}

	return &options
}
