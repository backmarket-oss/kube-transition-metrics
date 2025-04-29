package statistics

import (
	"github.com/BackMarket-oss/kube-transition-metrics/internal/options"
	"github.com/BackMarket-oss/kube-transition-metrics/internal/prommetrics"
	"k8s.io/apimachinery/pkg/types"
)

type statisticEvent interface {
	handle(statistic *podStatistic) bool
	podUID() types.UID
}

// StatisticEventHandler loops over statistic events sent by collectors to track
// and update metrics for Pod lifecycle events.
type StatisticEventHandler struct {
	options       *options.Options
	eventChan     prommetrics.MonitoredChannel[statisticEvent]
	resyncChan    prommetrics.MonitoredChannel[[]types.UID]
	blacklistUIDs []types.UID
	statistics    map[types.UID]*podStatistic
}

// NewStatisticEventHandler creates a new StatisticEventHandler which filters
// out events for the provided initialSyncBlacklist Pod UIDs.
func NewStatisticEventHandler(options *options.Options) *StatisticEventHandler {
	return &StatisticEventHandler{
		options: options,
		eventChan: prommetrics.NewMonitoredChannel[statisticEvent](
			"statistic_events", options.StatisticEventQueueLength),
		// Must not have queue as it ensures that new PodAdded statistic events aren't
		// generated before resync is processed.
		resyncChan: prommetrics.NewMonitoredChannel[[]types.UID](
			"pod_resync", 0),
		statistics: map[types.UID]*podStatistic{},
	}
}

// Publish sends an event to the StatisticEventHandler loop.
func (eh *StatisticEventHandler) Publish(ev statisticEvent) {
	eh.eventChan.Publish(ev)
}

func (eh *StatisticEventHandler) isBlacklisted(uid types.UID) bool {
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

func (eh *StatisticEventHandler) handleEvent(event statisticEvent) {
	uid := event.podUID()
	if eh.isBlacklisted(uid) {
		return
	}

	statistic := eh.getPodStatistic(uid)
	if event.handle(statistic) {
		delete(eh.statistics, uid)
	}

	prommetrics.PodsTracked.Set(float64(len(eh.statistics)))
	prommetrics.EventsHandled.Inc()
}

func (eh *StatisticEventHandler) handleResync(resyncUIDs []types.UID) {
	resyncUIDsSet := map[types.UID]interface{}{}
	for _, resyncUID := range resyncUIDs {
		resyncUIDsSet[resyncUID] = nil
	}

	for uid := range eh.statistics {
		if _, ok := resyncUIDsSet[uid]; !ok {
			delete(eh.statistics, uid)
		}
	}

	eh.blacklistUIDs = []types.UID{}
	for _, uid := range resyncUIDs {
		if _, ok := eh.statistics[uid]; !ok {
			eh.blacklistUIDs = append(eh.blacklistUIDs, uid)
		}
	}
}

// Run launches the statistic event handling loop. It is blocking and should be
// run in another goroutine to each of the collectors. It provides synchronous
// and ordered execution of statistic events.
func (eh *StatisticEventHandler) Run() {
	for {
		select {
		case event, ok := <-eh.eventChan.Channel():
			if !ok {
				break
			}

			eh.handleEvent(event)
		case resyncUIDs, ok := <-eh.resyncChan.Channel():
			if !ok {
				break
			}

			eh.handleResync(resyncUIDs)
		}
	}
}
