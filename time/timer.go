package time

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"time"
)

// Timer is a trivial wrapper around time.Timer which allows a timer
// to be mocked and/or replaced with an implementation that can be
// triggered deterministically.
type Timer interface {
	C() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

// Creates a new Timer, wrapping a time.Timer that will expire after
// the given duration.
func NewTimer(d time.Duration) Timer {
	return &timer{time.NewTimer(d)}
}

type timer struct {
	*time.Timer
}

var _ Timer = &timer{}

func (t *timer) C() <-chan time.Time {
	return t.Timer.C
}
