package options

import (
	"log"

	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
)

// Options contains the options for the controller.
type Options struct {
	ListenAddress             string
	KubeconfigPath            string
	ImagePullCancelDelay      float64
	KubeWatchTimeout          int64
	KubeWatchMaxEvents        int64
	StatisticEventQueueLength int
	LogLevel                  zerolog.Level
}

// Parse parses the options and returns them as a pointer to an Options struct.
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
		"The delay before canceling an image pull routine to ensure events are "+
			"flushed (ADVANCED)")
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
	logLevel := flag.String(
		"log-level",
		"INFO",
		`The global logging level, one of "trace", "debug", "info", "warn", `+
			`"error", "fatal", "panic", "disabled", or "" (empty string). This option's`+
			`values are case-insensitive. Setting a value of "disabled" will result in`+
			`no metrics being emitted.`)

	flag.Parse()
	if logLevelParsed, err := zerolog.ParseLevel(*logLevel); err != nil {
		log.Fatalf("Invalid value for --log-level (%s): %q\n", *logLevel, err)
	} else {
		options.LogLevel = logLevelParsed
	}

	return &options
}
