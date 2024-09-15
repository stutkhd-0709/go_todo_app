package clock

import "time"

type Clocker interface {
	Now() time.Time
}

type RealClock struct{}

func (r RealClock) Now() time.Time {
	return time.Now()
}

type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	return time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
}
