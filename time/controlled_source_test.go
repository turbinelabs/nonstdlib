/*
Copyright 2018 Turbine Labs, Inc.

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
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func checkTimer(timer Timer) (time.Time, bool) {
	select {
	case tm := <-timer.C():
		return tm, true
	default:
		return time.Time{}, false
	}
}

func checkContext(ctxt context.Context) (bool, error) {
	select {
	case <-ctxt.Done():
		return true, ctxt.Err()
	default:
		return false, nil
	}
}

func TestControlledSource(t *testing.T) {
	original := time.Now()
	source := &controlledTimeSource{now: original, mutex: &sync.Mutex{}}

	assert.DeepEqual(t, source.Now(), original)
	assert.DeepEqual(t, source.Now(), original)

	source.Advance(5 * time.Minute)
	assert.DeepEqual(t, source.Now(), original.Add(5*time.Minute))

	source.Set(original)
	assert.DeepEqual(t, source.Now(), original)
}

func TestControlledSourceNewTimer(t *testing.T) {
	original := time.Now()
	source := &controlledTimeSource{now: original, mutex: &sync.Mutex{}}

	timer := source.NewTimer(1 * time.Second)

	source.Advance(999 * time.Millisecond)
	_, ok := checkTimer(timer)
	assert.False(t, ok)

	source.Advance(1 * time.Millisecond)
	tm, ok := checkTimer(timer)
	assert.Equal(t, tm, source.Now())
	assert.True(t, ok)

	source.Advance(1 * time.Second)
	_, ok = checkTimer(timer)
	assert.False(t, ok)

	immediate := source.NewTimer(0 * time.Second)
	tm, ok = checkTimer(immediate)
	assert.Equal(t, tm, source.Now())
	assert.True(t, ok)
}

func TestControlledSourceAfterFunc(t *testing.T) {
	original := time.Now()
	source := &controlledTimeSource{now: original, mutex: &sync.Mutex{}}

	calls := make(chan time.Time, 10)
	expectCall := func(expectedTime time.Time) {
		timer := time.NewTimer(time.Second)
		select {
		case tm := <-calls:
			assert.DeepEqual(t, tm, expectedTime)
		case <-timer.C:
			assert.Failed(t, "no call in 1 second")
		}
	}

	timer := source.AfterFunc(1*time.Second, func() { calls <- source.Now() })

	source.Advance(999 * time.Millisecond)
	assert.ChannelEmpty(t, calls)

	source.Advance(1 * time.Millisecond)
	expectCall(original.Add(1 * time.Second))

	source.Advance(1 * time.Second)
	assert.ChannelEmpty(t, calls)

	timer = source.AfterFunc(1*time.Second, func() { calls <- source.Now() })
	timer.Stop()
	source.Advance(2 * time.Second)
	assert.ChannelEmpty(t, calls)

	source.AfterFunc(0*time.Second, func() { calls <- source.Now() })
	expectCall(original.Add(4 * time.Second))
}

func TestControlledSourceNewContextWithTimeout(t *testing.T) {
	original := time.Now()
	source := &controlledTimeSource{now: original, mutex: &sync.Mutex{}}

	ctxt, cancel := source.NewContextWithTimeout(context.TODO(), 1*time.Second)
	defer cancel()

	source.Advance(999 * time.Millisecond)
	ok, _ := checkContext(ctxt)
	assert.False(t, ok)

	source.Advance(1 * time.Millisecond)
	ok, err := checkContext(ctxt)
	assert.Equal(t, err, context.DeadlineExceeded)
	assert.True(t, ok)

	immediate, immediateCancel := source.NewContextWithTimeout(context.TODO(), 0*time.Second)
	defer immediateCancel()
	ok, err = checkContext(immediate)
	assert.Equal(t, err, context.DeadlineExceeded)
	assert.True(t, ok)
}

func TestControlledSourceTriggerAllTimers(t *testing.T) {
	original := time.Now()
	source := &controlledTimeSource{now: original, mutex: &sync.Mutex{}}

	timer1 := source.NewTimer(1 * time.Second)
	timer2 := source.NewTimer(2 * time.Second)
	timer3 := source.NewTimer(3 * time.Second)

	ctxt1, cancel1 := source.NewContextWithTimeout(context.TODO(), 3*time.Second)
	defer cancel1()
	ctxt2, cancel2 := source.NewContextWithTimeout(context.TODO(), 10*time.Second)
	defer cancel2()

	source.Advance(1 * time.Second)

	n := source.TriggerAllTimers()
	assert.Equal(t, n, 2)
	assert.Equal(t, source.Now(), original.Add(3*time.Second))

	tm, ok := checkTimer(timer1)
	assert.Equal(t, tm, original.Add(1*time.Second))

	tm, ok = checkTimer(timer2)
	assert.Equal(t, tm, original.Add(3*time.Second))

	tm, ok = checkTimer(timer3)
	assert.Equal(t, tm, original.Add(3*time.Second))

	ok, err := checkContext(ctxt1)
	assert.Equal(t, err, context.DeadlineExceeded)
	assert.True(t, ok)

	ok, _ = checkContext(ctxt2)
	assert.False(t, ok)

	n = source.TriggerAllTimers()
	assert.Equal(t, n, 0)
	assert.Equal(t, source.Now(), original.Add(3*time.Second))

	ok, _ = checkContext(ctxt2)
	assert.False(t, ok)
}

func TestControlledSourceTriggerNextTimer(t *testing.T) {
	original := time.Now()
	source := &controlledTimeSource{now: original, mutex: &sync.Mutex{}}

	timer1 := source.NewTimer(1 * time.Second)
	timer2 := source.NewTimer(2 * time.Second)
	timer3 := source.NewTimer(3 * time.Second)

	ctxt1, cancel1 := source.NewContextWithTimeout(context.TODO(), 3*time.Second)
	defer cancel1()
	ctxt2, cancel2 := source.NewContextWithTimeout(context.TODO(), 10*time.Second)
	defer cancel2()

	assert.True(t, source.TriggerNextTimer())
	assert.Equal(t, source.Now(), original.Add(1*time.Second))
	tm, ok := checkTimer(timer1)
	assert.Equal(t, tm, original.Add(1*time.Second))

	_, ok = checkTimer(timer2)
	assert.False(t, ok)

	assert.True(t, source.TriggerNextTimer())
	assert.Equal(t, source.Now(), original.Add(2*time.Second))
	tm, ok = checkTimer(timer2)
	assert.Equal(t, tm, original.Add(2*time.Second))

	_, ok = checkTimer(timer3)
	assert.False(t, ok)

	assert.True(t, source.TriggerNextTimer())
	assert.Equal(t, source.Now(), original.Add(3*time.Second))
	tm, ok = checkTimer(timer3)
	assert.Equal(t, tm, original.Add(3*time.Second))
	ok, err := checkContext(ctxt1)
	assert.Equal(t, err, context.DeadlineExceeded)
	assert.True(t, ok)

	assert.False(t, source.TriggerNextTimer())

	ok, _ = checkContext(ctxt2)
	assert.False(t, ok)
}

func TestControlledSourceTriggerNextContext(t *testing.T) {
	original := time.Now()
	source := &controlledTimeSource{now: original, mutex: &sync.Mutex{}}

	timer1 := source.NewTimer(1 * time.Second)
	timer2 := source.NewTimer(2 * time.Second)
	timer3 := source.NewTimer(3 * time.Second)

	ctxt1, cancel1 := source.NewContextWithTimeout(context.TODO(), 3*time.Second)
	defer cancel1()
	ctxt2, cancel2 := source.NewContextWithTimeout(context.TODO(), 10*time.Second)
	defer cancel2()

	assert.True(t, source.TriggerNextContext())
	assert.Equal(t, source.Now(), original.Add(3*time.Second))
	tm, ok := checkTimer(timer1)
	assert.Equal(t, tm, original.Add(3*time.Second))
	tm, ok = checkTimer(timer2)
	assert.Equal(t, tm, original.Add(3*time.Second))
	tm, ok = checkTimer(timer3)
	assert.Equal(t, tm, original.Add(3*time.Second))
	ok, err := checkContext(ctxt1)
	assert.Equal(t, err, context.DeadlineExceeded)
	assert.True(t, ok)

	ok, _ = checkContext(ctxt2)
	assert.False(t, ok)

	assert.True(t, source.TriggerNextContext())
	assert.Equal(t, source.Now(), original.Add(10*time.Second))
	ok, err = checkContext(ctxt2)
	assert.Equal(t, err, context.DeadlineExceeded)
	assert.True(t, ok)

	assert.False(t, source.TriggerNextContext())
}

func TestWithTimeAt(t *testing.T) {
	original := time.Now()

	called := false
	WithTimeAt(original, func(ts ControlledSource) {
		called = true

		assert.Equal(t, ts.Now(), original)
	})

	assert.True(t, called)
}

func TestWithCurrentTimeFrozen(t *testing.T) {
	before := time.Now()

	called := false
	frozenTime := time.Time{}
	WithCurrentTimeFrozen(func(ts ControlledSource) {
		called = true
		frozenTime = ts.Now()
	})
	after := time.Now()

	assert.True(t, called)
	assert.True(t, before.Before(frozenTime) || before.Equal(frozenTime))
	assert.True(t, after.After(frozenTime) || after.Equal(frozenTime))
}

func TestIncrementingControlledSource(t *testing.T) {
	before := time.Now()
	delta := 5 * time.Second

	s := NewIncrementingControlledSource(before, delta)
	assert.Equal(t, before, s.Now())
	assert.Equal(t, before.Add(delta), s.Now())
}

func TestIncrementingControlledSourceAdvance(t *testing.T) {
	before := time.Now()
	delta := 5 * time.Second

	s := NewIncrementingControlledSource(before, delta)
	assert.Equal(t, before, s.Now())
	s.Advance(time.Second)

	assert.Equal(t, before.Add(delta+time.Second), s.Now())
}

func TestIncrementingControlledSourceSet(t *testing.T) {
	before := time.Now()
	delta := 5 * time.Second

	s := NewIncrementingControlledSource(before, delta)
	assert.Equal(t, before, s.Now())
	s.Set(before.Add(-1 * time.Hour))

	assert.Equal(t, before.Add(-1*time.Hour), s.Now())
	assert.Equal(t, before.Add(-1*time.Hour).Add(delta), s.Now())
}

func TestIncrementingControlledSourceNewTimer(t *testing.T) {
	before := time.Now()
	delta := 5 * time.Second
	timerDelay := 6 * time.Second

	s := NewIncrementingControlledSource(before, delta)

	timer := s.NewTimer(timerDelay)

	assert.Equal(t, s.Now(), before)
	_, ok := checkTimer(timer)
	assert.False(t, ok)

	assert.Equal(t, s.Now(), before.Add(delta))
	tm, ok := checkTimer(timer)
	assert.True(t, !tm.Before(before.Add(timerDelay)))
	assert.True(t, ok)
}

func TestControlledTimerReset(t *testing.T) {
	now := time.Now()
	s := &controlledTimeSource{now: now, mutex: &sync.Mutex{}}

	timer := s.NewTimer(10 * time.Second)

	s.Advance(5 * time.Second)

	assert.True(t, timer.Reset(10*time.Second))

	s.Advance(5 * time.Second)

	// original deadline ignored
	_, ok := checkTimer(timer)
	assert.False(t, ok)

	s.Advance(5 * time.Second)

	// reset deadline heeded
	tm, ok := checkTimer(timer)
	assert.Equal(t, tm, now.Add(15*time.Second))
	assert.True(t, ok)

	assert.False(t, timer.Reset(1*time.Second))
	s.Advance(1 * time.Second)
	tm, ok = checkTimer(timer)
	assert.Equal(t, tm, now.Add(16*time.Second))
	assert.True(t, ok)
}

func TestControlledTimerStop(t *testing.T) {
	now := time.Now()
	s := &controlledTimeSource{now: now, mutex: &sync.Mutex{}}

	timer := s.NewTimer(10 * time.Second)

	s.Advance(5 * time.Second)

	assert.True(t, timer.Stop())

	s.Advance(5 * time.Second)

	// original deadline ignored
	_, ok := checkTimer(timer)
	assert.False(t, ok)

	s.Advance(5 * time.Second)
	assert.False(t, timer.Stop())

	_, ok = checkTimer(timer)
	assert.False(t, ok)
}

func TestControlledContextDeadline(t *testing.T) {
	now := time.Now()
	c := &controlledContext{
		deadline: now,
	}

	deadline, ok := c.Deadline()
	assert.Equal(t, deadline, now)
	assert.True(t, ok)
}

func TestControlledContextCancelsOnce(t *testing.T) {
	c := &controlledContext{
		done: make(chan struct{}),
	}

	err := errors.New("this is my cancellation")
	otherErr := errors.New("this is some other cancellation")

	c.cancel(err)
	assert.Equal(t, c.Err(), err)

	c.cancel(otherErr)
	assert.Equal(t, c.Err(), err)
}

func TestControlledContextPropagation(t *testing.T) {
	original := time.Now()
	source := &controlledTimeSource{now: original, mutex: &sync.Mutex{}}

	parent, parentCancel := context.WithCancel(context.TODO())

	child := newControlledContext(parent, source, original.Add(1*time.Hour))

	assert.Nil(t, child.Err())

	parentCancel()

	for child.Err() == nil {
		time.Sleep(1 * time.Millisecond)
	}

	assert.Equal(t, child.Err(), context.Canceled)
}

func TestControlledContextParentExpiresSooner(t *testing.T) {
	original := time.Now()
	source := &controlledTimeSource{now: original, mutex: &sync.Mutex{}}

	parent, parentCancel := source.NewContextWithTimeout(context.TODO(), 1*time.Second)
	defer parentCancel()

	child, childCancel := source.NewContextWithTimeout(parent, 2*time.Second)
	defer childCancel()

	// Prove we didn't make our own context for the child.
	_, wrongType := child.(*controlledContext)
	assert.False(t, wrongType)
}
