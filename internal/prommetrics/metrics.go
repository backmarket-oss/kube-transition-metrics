package prommetrics

import "github.com/prometheus/client_golang/prometheus"

//nolint:gochecknoglobals
var (
	POD_COLLECTOR_ERRORS = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pod_collector_errors_total",
			Help: "Total number of pod collector errors since the last restart",
		},
	)
	POD_COLLECTOR_RESTARTS = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pod_collector_restarts_total",
			Help: "Total number of times the pod collector Watch was restarted since " +
				"the process started",
		},
	)
	PODS_PROCESSED = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pods_watch_processed_total",
			Help: "Total number of Pod Watch messages since the last restart",
		},
		[]string{"event_type"},
	)
	IMAGE_PULL_COLLECTOR_ROUTINES = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "image_pull_collector_routines",
			Help: "Current number of running image pull collector routines",
		},
	)
	IMAGE_PULL_COLLECTOR_ERRORS = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_pull_collector_errors_total",
			Help: "Total number of image pull collector errors since the last restart",
		},
	)
	IMAGE_PULL_COLLECTOR_RESTARTS = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_pull_collector_restarts_total",
			Help: "Total number of times the image pull collector Watch was restarted " +
				"since the process started",
		},
	)
	EVENTS_PROCESSED = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_watch_processed_total",
			Help: "Total number of Event Watch messages since the last restart",
		},
		[]string{"event_type"},
	)
	PODS_TRACKED = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "pod_statistics_tracked",
			Help: "Current number of pods tracked",
		},
	)
	EVENTS_HANDLED = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "statistic_events_handled_total",
			Help: "Total number of statistic events handled since the last restart",
		},
	)
	EVENT_PUBLISH_WAIT_DURATION = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "statistic_event_publish_wait_total",
			Help: "Total amount of time in seconds waiting to publish a statistic " +
				"event from a collector",
		},
	)
)

// Register registers the prometheus Collectors (metrics) exported by this
// package.
func Register() {
	prometheus.MustRegister(
		POD_COLLECTOR_ERRORS,
		POD_COLLECTOR_RESTARTS,
		PODS_PROCESSED,
		IMAGE_PULL_COLLECTOR_ROUTINES,
		IMAGE_PULL_COLLECTOR_ERRORS,
		IMAGE_PULL_COLLECTOR_RESTARTS,
		EVENTS_PROCESSED,
		PODS_TRACKED,
		EVENTS_HANDLED,
		EVENT_PUBLISH_WAIT_DURATION,
	)
}
