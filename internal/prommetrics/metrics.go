package prommetrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

//nolint:gochecknoglobals
var (
	// summaryObjectives is a the quantile objectives for the summary metrics.
	//nolint:mnd
	summaryObjectives = map[float64]float64{
		0.5:  0.05,
		0.9:  0.01,
		0.99: 0.001,
	}

	// PodCollectorErrors tracks the total number of pod collector errors since the last restart.
	PodCollectorErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pod_collector_errors_total",
			Help: "Total number of pod collector errors since the last restart",
		},
	)
	// PodCollectorRestarts tracks the total number of pod collector restarts since the process started.
	PodCollectorRestarts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pod_collector_restarts_total",
			Help: "Total number of times the pod collector Watch was restarted since " +
				"the process started",
		},
	)
	// PodWatchEvents tracks the total number of pod watch messages since the last restart.
	PodWatchEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pod_watch_events_total",
			Help: "Total number of Pod Watch messages since the last restart",
		},
		[]string{"event_type"},
	)
	// ImagePullCollectorRoutines tracks the current number of running image pull collector routines.
	ImagePullCollectorRoutines = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "image_pull_collector_routines",
			Help: "Current number of running image pull collector routines",
		},
	)
	// ImagePullCollectorErrors tracks the total number of image pull collector errors since the last restart.
	ImagePullCollectorErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_pull_collector_errors_total",
			Help: "Total number of image pull collector errors since the last restart",
		},
	)
	// ImagePullCollectorRestarts tracks the total number of image pull collector restarts since the process started.
	ImagePullCollectorRestarts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_pull_collector_restarts_total",
			Help: "Total number of times the image pull collector Watch was restarted " +
				"since the process started",
		},
	)
	// ImagePullWatchEvents tracks the total number of event watch messages since the last restart.
	ImagePullWatchEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "image_pull_watch_events_total",
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
	// ImagePullTracked tracks the current number of image pulls tracked.
	ImagePullTracked = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "image_pull_statistics_tracked",
			Help: "Current number of image pulls tracked",
		},
	)
	// StatisticEventPublish tracks the time spent waiting to publish an event and the number of events published.
	StatisticEventPublish = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "statistic_event_publish_seconds",
			Help: fmt.Sprintf("Time spent waiting to publish an event in seconds (quarantiles over %v)", prometheus.DefMaxAge),

			Objectives: summaryObjectives,
		},
		[]string{"event_loop"},
	)
	// StatisticEventQueueDepth tracks the current queue depth of the event queue.
	StatisticEventQueueDepth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "statistic_event_queue_depth",
			Help: "Current queue depth of the event queue",
		},
		[]string{"event_loop"},
	)
	// StatisticEventProcessing tracks the time spent processing events and the number of events processed.
	StatisticEventProcessing = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "statistic_event_processing_seconds",
			Help: "Time spent processing events in seconds (quarantiles over " + prometheus.DefMaxAge.String() + ")",

			Objectives: summaryObjectives,
		},
		[]string{"event_loop"},
	)
)

// Register registers the prometheus Collectors (metrics) exported by this
// package.
func Register() {
	prometheus.MustRegister(
		PodCollectorErrors,
		PodCollectorRestarts,
		PodWatchEvents,
		ImagePullCollectorRoutines,
		ImagePullCollectorErrors,
		ImagePullCollectorRestarts,
		ImagePullWatchEvents,
		PodsTracked,
		ImagePullTracked,
		StatisticEventPublish,
		StatisticEventQueueDepth,
		StatisticEventProcessing,
	)
}
