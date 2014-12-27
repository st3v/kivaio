package kivaio

import "time"

type Clock interface {
	Now() time.Time
}

var clock Clock = &realClock{}

type realClock struct{}

func (r *realClock) Now() time.Time {
	return time.Now()
}
