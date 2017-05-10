/*
Copyright 2017 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package time

import (
	"context"
	"sync"
	"time"

	"github.com/turbinelabs/nonstdlib/ptr"
)

// ControlledSource is a source of time.Time values that returns a
// fixed time unless modified with the Set or Advance
// methods. ControlledSource should be used for testing only.
type ControlledSource interface {
	Source

	// Set sets the current time returned by this Source. Timers
	// and context.Contexts created by this Source are
	// triggered/canceled if the new time exceeds their deadline.
	Set(time.Time)

	// Set advances the current time returned by this Source.
	// Timers and context.Contexts created by this Source are
	// triggered/canceled if the new time exceeds their deadline.
	Advance(time.Duration)

	// TriggerAllTimers set the current time to the latest
	// deadline of all timers and returns the number of timers
	// that were triggered. If no timers are live, it returns 0
	// and the time is not advanced. Any Contexts whose deadlines
	// are exceeded will also be triggered.
	TriggerAllTimers() int
}

// WithTimeAt creates a new ControlledSource with the given time and
// passes it to the given function for testing.
func WithTimeAt(t time.Time, f func(ControlledSource)) {
	s := &controlledTimeSource{
		now:   t,
		mutex: &sync.Mutex{},
	}
	f(s)
}

// WithCurrentTimeFrozen Creates a new ControlledSource with the
// current time and passes it to the given function for testing.
func WithCurrentTimeFrozen(f func(ControlledSource)) {
	WithTimeAt(time.Now(), f)
}

// NewIncrementingControlledSource returns a new ControlledSource that
// increments the controlled time by some delta every time Now() is called.
func NewIncrementingControlledSource(at time.Time, delta time.Duration) ControlledSource {
	return &incrementingTimeSource{
		&controlledTimeSource{
			now:   at,
			mutex: &sync.Mutex{},
		},
		delta,
	}
}

// controlledTimeSource is a deterministic Source of time values. It
// provides Timer and context.Context implementations that are tied to
// the current time reported by Now().
type controlledTimeSource struct {
	now      time.Time
	timers   []*controlledTimer
	contexts []*controlledContext
	mutex    *sync.Mutex
}

func (s *controlledTimeSource) Now() time.Time {
	return s.now
}

func (s *controlledTimeSource) NewTimer(d time.Duration) Timer {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timer := newControlledTimer(s, d)
	s.timers = append(s.timers, timer)
	s.checkDeadlines()
	return timer
}

func (s *controlledTimeSource) AfterFunc(d time.Duration, f func()) Timer {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timer := controlledAfterFunc(s, d, f)
	s.timers = append(s.timers, timer)
	s.checkDeadlines()
	return timer
}

func (s *controlledTimeSource) NewContextWithDeadline(
	parent context.Context,
	deadline time.Time,
) (context.Context, context.CancelFunc) {
	if parentDeadline, ok := parent.Deadline(); ok && parentDeadline.Before(deadline) {
		return context.WithCancel(parent)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	ctxt := newControlledContext(parent, s, deadline)
	s.contexts = append(s.contexts, ctxt)
	s.checkDeadlines()

	return ctxt, func() { ctxt.cancel(context.Canceled) }
}

func (s *controlledTimeSource) NewContextWithTimeout(
	parent context.Context,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {
	return s.NewContextWithDeadline(parent, s.now.Add(timeout))
}

func (s *controlledTimeSource) Set(t time.Time) {
	s.now = t

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.checkDeadlines()
}

func (s *controlledTimeSource) Advance(delta time.Duration) {
	s.now = s.now.Add(delta)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.checkDeadlines()
}

func (s *controlledTimeSource) TriggerAllTimers() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var maxDeadline time.Time

	for _, timer := range s.timers {
		if timer != nil && timer.deadline != nil && timer.deadline.After(maxDeadline) {
			maxDeadline = *timer.deadline
		}
	}

	if maxDeadline.IsZero() {
		return 0
	}

	s.now = maxDeadline

	n, _ := s.checkDeadlines()
	return n
}

func (s *controlledTimeSource) checkDeadlines() (int, int) {
	numTimers := 0
	numCtxts := 0
	for _, timer := range s.timers {
		if timer != nil && timer.check(s.now) {
			numTimers++
		}
	}

	for _, ctxt := range s.contexts {
		if ctxt != nil && ctxt.check(s.now) {
			numCtxts++
		}
	}

	return numTimers, numCtxts
}

// incrementingTimeSource wraps a controlledTimeSource and updates
// the current time by inc each time Now() is called, returning the
// previous value.
type incrementingTimeSource struct {
	*controlledTimeSource
	inc time.Duration
}

func (i *incrementingTimeSource) Now() time.Time {
	t := i.now
	i.Advance(i.inc)
	return t
}

// controlledTimer implements Timer. It requires a ControlledSource to
// invoke check when the current time is updated.
type controlledTimer struct {
	source   ControlledSource
	deadline *time.Time
	c        chan time.Time
	f        func()
}

func newControlledTimer(source ControlledSource, delay time.Duration) *controlledTimer {
	return &controlledTimer{
		source:   source,
		deadline: ptr.Time(source.Now().Add(delay)),
		c:        make(chan time.Time, 1),
	}
}

func controlledAfterFunc(source ControlledSource, delay time.Duration, f func()) *controlledTimer {
	return &controlledTimer{
		source:   source,
		deadline: ptr.Time(source.Now().Add(delay)),
		f:        f,
	}
}

func (t *controlledTimer) C() <-chan time.Time {
	return t.c
}

func (t *controlledTimer) Reset(d time.Duration) bool {
	set := t.deadline != nil
	t.deadline = ptr.Time(t.source.Now().Add(d))
	return set
}

func (t *controlledTimer) Stop() bool {
	stopping := t.deadline != nil
	t.deadline = nil
	return stopping
}

func (t *controlledTimer) check(now time.Time) bool {
	if t.deadline != nil && !t.deadline.After(now) {
		t.deadline = nil
		if t.c != nil {
			// Send time if room in channel, otherwise skip it.
			select {
			case t.c <- now:
			default:
			}
		}

		if t.f != nil {
			go t.f()
		}

		return true
	}

	return false
}

// controlledTimer implements context.Context. It requires a
// ControlledSource to invoke check when the current time is updated.
type controlledContext struct {
	context.Context
	source   Source
	deadline time.Time
	done     chan struct{}
	err      error
}

func newControlledContext(
	parent context.Context,
	source Source,
	deadline time.Time,
) *controlledContext {
	ctxt := &controlledContext{
		Context:  parent,
		source:   source,
		deadline: deadline,
		done:     make(chan struct{}),
	}

	// propagate cancellation from parent
	go func() {
		select {
		case <-parent.Done():
			ctxt.cancel(parent.Err())
		case <-ctxt.Done():
		}
	}()

	return ctxt
}

func (c *controlledContext) Deadline() (time.Time, bool) {
	return c.deadline, true
}

func (c *controlledContext) Done() <-chan struct{} {
	return c.done
}

func (c *controlledContext) Err() error {
	return c.err
}

func (c *controlledContext) cancel(err error) {
	if c.err != nil {
		return
	}

	c.err = err
	close(c.done)
}

func (c *controlledContext) check(now time.Time) bool {
	if c.err == nil && !c.deadline.After(now) {
		c.cancel(context.DeadlineExceeded)
		return true
	}

	return false
}
