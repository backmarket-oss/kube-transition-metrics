package prommetrics

import "github.com/prometheus/client_golang/prometheus"

//nolint:gochecknoglobals
var (
	PodCollectorErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pod_collector_errors_total",
			Help: "Total number of pod collector errors since the last restart",
		},
	)
	PodCollectorRestarts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pod_collector_restarts_total",
			Help: "Total number of times the pod collector Watch was restarted since " +
				"the process started",
		},
	)
	PodsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pods_watch_processed_total",
			Help: "Total number of Pod Watch messages since the last restart",
		},
		[]string{"event_type"},
	)
	ImagePullCollectorRoutines = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "image_pull_collector_routines",
			Help: "Current number of running image pull collector routines",
		},
	)
	ImagePullCollectorErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_pull_collector_errors_total",
			Help: "Total number of image pull collector errors since the last restart",
		},
	)
	ImagePullCollectorRestarts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_pull_collector_restarts_total",
			Help: "Total number of times the image pull collector Watch was restarted " +
				"since the process started",
		},
	)
	EventsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_watch_processed_total",
			Help: "Total number of Event Watch messages since the last restart",
		},
		[]string{"event_type"},
	)
	PodsTracked = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "pod_statistics_tracked",
			Help: "Current number of pods tracked",
		},
	)
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
