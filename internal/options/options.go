package options

import flag "github.com/spf13/pflag"

// Options contains the options for the controller.
type Options struct {
	ListenAddress             string
	KubeconfigPath            string
	ImagePullCancelDelay      float64
	KubeWatchTimeout          int64
	KubeWatchMaxEvents        int64
	StatisticEventQueueLength int
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
	//nolint:gomnd
	flag.Float64Var(
		&options.ImagePullCancelDelay,
		"image-pull-cancel-delay",
		3,
		"The delay before canceling an image pull routine to ensure events are "+
			"flushed (ADVANCED)")
	//nolint:gomnd
	flag.Int64Var(
		&options.KubeWatchTimeout,
		"kube-watch-timeout",
		60,
		"The Kubernetes Watch API timeout (ADVANCED)")
	//nolint:gomnd
	flag.Int64Var(
		&options.KubeWatchMaxEvents,
		"kube-watch-max-events",
		100,
		"The Kubernetes Watch maximum events per response (ADVANCED)")
	//nolint:gomnd
	flag.IntVar(
		&options.StatisticEventQueueLength,
		"statistic-event-queue-length",
		1000,
		"The maximum number of queued statistic events (ADVANCED)")

	flag.Parse()

	return &options
}
