package time

import (
	"time"
)

var (
	defaultTimeSourceInstance = &defaultTimeSource{}
)

// Source is a source of time.Time values.
type Source interface {
	// Returns the current time, as with time.Time.Now().
	Now() time.Time
}

// ControlledSource is a source of time.Time values that returns a
// fixed time unless modified with the Set or Advance
// methods. ControlledSource should be used for testing only.
type ControlledSource interface {
	Source

	Set(time.Time)
	Advance(time.Duration)
}

// Creates a new Source for normal use. The Source's Now method
// returns the current time.
func NewSource() Source {
	return defaultTimeSourceInstance
}

// Creates a new ControlledSource with the given time and passes it to
// the given function for testing.
func WithTimeAt(t time.Time, f func(ControlledSource)) {
	s := &controlledTimeSource{now: t}
	f(s)
}

// Creates a new ControlledSource with the current time and passes it
// to the given function for testing.
func WithCurrentTimeFrozen(f func(ControlledSource)) {
	WithTimeAt(time.Now(), f)
}

type defaultTimeSource struct{}

func (s *defaultTimeSource) Now() time.Time {
	return time.Now()
}

type controlledTimeSource struct {
	now time.Time
}

func (s *controlledTimeSource) Now() time.Time {
	return s.now
}

func (s *controlledTimeSource) Set(t time.Time) {
	s.now = t
}

func (s *controlledTimeSource) Advance(delta time.Duration) {
	s.now = s.now.Add(delta)
}
