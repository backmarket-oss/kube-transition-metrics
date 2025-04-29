package prommetrics

import "github.com/prometheus/client_golang/prometheus"

//nolint:gochecknoglobals
var (
	// PodCollectorErrors tracks the total number of pod collector errors since the
	// last restart.
	PodCollectorErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pod_collector_errors_total",
			Help: "Total number of pod collector errors since the last restart",
		},
	)
	// PodCollectorRestarts tracks the total number of pod collector restarts since
	// the process started.
	PodCollectorRestarts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pod_collector_restarts_total",
			Help: "Total number of times the pod collector Watch was restarted since " +
				"the process started",
		},
	)
	// PodsProcessed tracks the total number of pod watch messages since the last
	// restart.
	PodsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pods_watch_processed_total",
			Help: "Total number of Pod Watch messages since the last restart",
		},
		[]string{"event_type"},
	)
	// ImagePullCollectorRoutines tracks the current number of running image pull
	// collector routines.
	ImagePullCollectorRoutines = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "image_pull_collector_routines",
			Help: "Current number of running image pull collector routines",
		},
	)
	// ImagePullCollectorErrors tracks the total number of image pull collector
	// errors since the last restart.
	ImagePullCollectorErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_pull_collector_errors_total",
			Help: "Total number of image pull collector errors since the last restart",
		},
	)
	// ImagePullCollectorRestarts tracks the total number of image pull collector
	// restarts since the process started.
	ImagePullCollectorRestarts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_pull_collector_restarts_total",
			Help: "Total number of times the image pull collector Watch was restarted " +
				"since the process started",
		},
	)
	// EventsProcessed tracks the total number of event watch messages since the
	// last restart.
	EventsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_watch_processed_total",
			Help: "Total number of Event Watch messages since the last restart",
		},
		[]string{"event_type"},
	)
	// PodsTracked tracks the current number of pods tracked.
	PodsTracked = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "pod_statistics_tracked",
			Help: "Current number of pods tracked",
		},
	)
	// EventsHandled tracks the total number of statistic events handled since the
	// last restart.
	EventsHandled = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "statistic_events_handled_total",
			Help: "Total number of statistic events handled since the last restart",
		},
	)
)

// Register registers the prometheus Collectors (metrics) exported by this
// package.
func Register() {
	prometheus.MustRegister(
		PodCollectorErrors,
		PodCollectorRestarts,
		PodsProcessed,
		ImagePullCollectorRoutines,
		ImagePullCollectorErrors,
		ImagePullCollectorRestarts,
		EventsProcessed,
		PodsTracked,
		EventsHandled,
		MonitoredChannelQueueDepth,
		MonitoredChannelPublishWaitDuration,
		ChannelMonitors,
	)
}
