package statistics

import "time"

type timeSource interface {
	Now() time.Time
}

type realTimeSource struct{}

func (rts realTimeSource) Now() time.Time {
	return time.Now()
}
