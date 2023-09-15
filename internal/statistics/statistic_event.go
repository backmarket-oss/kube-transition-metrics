package statistics

import (
	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/prommetrics"
	"k8s.io/apimachinery/pkg/types"
)

type statisticEvent interface {
	Handle(statistic *podStatistic) bool
	PodUID() types.UID
}

// StatisticEventHandler loops over statistic events sent by collectors to track
// and update metrics for Pod lifecycle events.
type StatisticEventHandler struct {
	options       *options.Options
	eventChan     prommetrics.MonitoredChannel[statisticEvent]
	blacklistUIDs []types.UID
	statistics    map[types.UID]*podStatistic
}

// NewStatisticEventHandler creates a new StatisticEventHandler which filters
// out events for the provided initialSyncBlacklist Pod UIDs.
func NewStatisticEventHandler(
	options *options.Options,
	initialSyncBlacklist []types.UID,
) *StatisticEventHandler {
	return &StatisticEventHandler{
		options: options,
		eventChan: prommetrics.NewMonitoredChannel[statisticEvent](
			"statistic_events", options.StatisticEventQueueLength),
		blacklistUIDs: initialSyncBlacklist,
		statistics:    map[types.UID]*podStatistic{},
	}
}

// Publish sends an event to the StatisticEventHandler loop.
func (eh StatisticEventHandler) Publish(ev statisticEvent) {
	eh.eventChan.Publish(ev)
}

func (eh StatisticEventHandler) isBlacklisted(uid types.UID) bool {
	for _, blacklistedUID := range eh.blacklistUIDs {
		if blacklistedUID == uid {
			return true
		}
	}

	return false
}

func (eh *StatisticEventHandler) getPodStatistic(uid types.UID) *podStatistic {
	if statistic, ok := eh.statistics[uid]; ok {
		return statistic
	} else {
		eh.statistics[uid] = &podStatistic{}

		return eh.statistics[uid]
	}
}

// Run launches the statistic event handling loop. It is blocking and should be
// run in another goroutine to each of the collectors. It provides synchronous
// and ordered execution of statistic events.
func (eh *StatisticEventHandler) Run() {
	for {
		event, ok := eh.eventChan.Read()
		if !ok {
			break
		}

		uid := event.PodUID()
		if eh.isBlacklisted(uid) {
			continue
		}

		statistic := eh.getPodStatistic(uid)
		if event.Handle(statistic) {
			delete(eh.statistics, uid)
		}

		prommetrics.PodsTracked.Set(float64(len(eh.statistics)))
		prommetrics.EventsHandled.Inc()
	}
}
