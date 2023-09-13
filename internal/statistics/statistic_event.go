package statistics

import (
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
	EventChan     chan statisticEvent
	blacklistUIDs []types.UID
	statistics    map[types.UID]*podStatistic
}

// NewStatisticEventHandler creates a new StatisticEventHandler which filters
// out events for the provided initial_sync_blacklist Pod UIDs.
func NewStatisticEventHandler(
	initial_sync_blacklist []types.UID,
) *StatisticEventHandler {
	return &StatisticEventHandler{
		EventChan:     make(chan statisticEvent),
		blacklistUIDs: initial_sync_blacklist,
		statistics:    map[types.UID]*podStatistic{},
	}
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
	for event := range eh.EventChan {
		uid := event.PodUID()
		if eh.isBlacklisted(uid) {
			continue
		}

		statistic := eh.getPodStatistic(uid)
		if event.Handle(statistic) {
			delete(eh.statistics, uid)
		}

		prommetrics.PODS_TRACKED.Set(float64(len(eh.statistics)))
		prommetrics.EVENTS_HANDLED.Inc()
	}
}
