package main

import "time"

type TimeSource interface {
	Now() time.Time
}

type RealTimeSource struct{}

func (rts RealTimeSource) Now() time.Time {
	return time.Now()
}
