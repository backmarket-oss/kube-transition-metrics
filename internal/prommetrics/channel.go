package prommetrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//nolint:gochecknoglobals
var (
	MonitoredChannelPublishWaitDuration = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "channel_publish_wait_total",
			Help: "Total amount of time in seconds waiting to publish a to the channel",
		},
		[]string{"channel_name"},
	)
	MonitoredChannelQueueDepth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "channel_queue_depth",
			Help: "Current queue depth of the channel",
		},
		[]string{"channel_name"},
	)
	ChannelMonitors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "channel_monitors",
			Help: "Current number of channel monitor goroutines",
		},
	)
)

// MonitoredChannel wraps a channel and publishes prometheus metrics for the
// time spent waiting to publish an item to it and the queue depth (len) of the
// channel updated each second.
type MonitoredChannel[T interface{}] struct {
	name   string
	c      chan T
	cancel chan interface{}
}

// NewMonitoredChannel creates a new MonitoredChannel, initializing the channel
// with the provided maximum length. The provided name will be used in
// the prometheus label `channel_name`. It lauches a monitor goroutine in the
// background, which is shutdown when MonitoredChannel.Close() is called.
func NewMonitoredChannel[T interface{}](
	name string,
	length int,
) MonitoredChannel[T] {
	monitoredChannel := MonitoredChannel[T]{
		name:   name,
		c:      make(chan T, length),
		cancel: make(chan interface{}),
	}

	go monitoredChannel.monitor()

	return monitoredChannel
}

// Close closes the underlying channel and stops the monitoring goroutine.
func (mc MonitoredChannel[T]) Close() {
	mc.cancel <- nil
	close(mc.cancel)
	close(mc.c)
}

// Publish sends an item to the channel, and updates the prometheus metrics
// tracking the duration waiting to publish this event.
func (mc MonitoredChannel[T]) Publish(item T) {
	start := time.Now()

	mc.c <- item
	end := time.Now()

	waitDuration := end.Sub(start)
	MonitoredChannelPublishWaitDuration.
		With(mc.prometheusLabels()).
		Add(waitDuration.Seconds())
}

// Channel returns the underlying channel, only use for chan receive!
func (mc MonitoredChannel[T]) Channel() chan T {
	return mc.c
}

func (mc MonitoredChannel[T]) monitor() {
	logger := mc.logger()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	ChannelMonitors.Inc()
	defer ChannelMonitors.Dec()

	for {
		select {
		case _, ok := <-mc.cancel:
			if !ok {
				logger.Panic().Msg("Cancel channel closed prematurely.")
			}

			return
		case _, ok := <-ticker.C:
			if !ok {
				logger.Panic().Msg("Ticker closed prematurely.")
			}

			MonitoredChannelQueueDepth.
				With(mc.prometheusLabels()).
				Set(float64(len(mc.c)))
		}
	}
}

func (mc MonitoredChannel[T]) logger() *zerolog.Logger {
	logger := log.With().
		Str("subsystem", "monitored_channel").
		Str("channel_name", mc.name).
		Logger()

	return &logger
}

func (mc MonitoredChannel[T]) prometheusLabels() prometheus.Labels {
	return prometheus.Labels{"channel_name": mc.name}
}
