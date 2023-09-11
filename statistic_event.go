package main

import (
	"k8s.io/apimachinery/pkg/types"
)

type StatisticEvent interface {
	Handle(statistic *PodStatistic) bool
	PodUID() types.UID
}

type StatisticEventHandler struct {
	EventChan     chan StatisticEvent
	blacklistUIDs []types.UID
	statistics    map[types.UID]*PodStatistic
}

func NewStatisticEventHandler() *StatisticEventHandler {
	return &StatisticEventHandler{
		EventChan:     make(chan StatisticEvent),
		blacklistUIDs: []types.UID{},
		statistics:    map[types.UID]*PodStatistic{},
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

func (eh *StatisticEventHandler) GetPodStatistic(uid types.UID) *PodStatistic {
	if statistic, ok := eh.statistics[uid]; ok {
		return statistic
	} else {
		eh.statistics[uid] = &PodStatistic{}

		return eh.statistics[uid]
	}
}

func (eh *StatisticEventHandler) Run() {
	for event := range eh.EventChan {
		uid := event.PodUID()
		if eh.isBlacklisted(uid) {
			continue
		}

		statistic := eh.GetPodStatistic(uid)
		if event.Handle(statistic) {
			delete(eh.statistics, uid)
		}

		PODS_TRACKED.Set(float64(len(eh.statistics)))
		EVENTS_HANDLED.Inc()
	}
}
