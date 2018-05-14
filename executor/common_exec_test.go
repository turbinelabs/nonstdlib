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

package executor

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	tbntime "github.com/turbinelabs/nonstdlib/time"
	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/log"
)

type testRun struct {
	attemptedRetries chan string
}

type testData struct {
	id       string
	fails    int
	failFunc func()
}

func (d *testData) apply(r *testRun) (interface{}, error) {
	if d.fails > 0 {
		if r.attemptedRetries != nil {
			r.attemptedRetries <- d.id + " fail"
		}
		if d.failFunc != nil {
			defer d.failFunc()
		}
		d.fails--
		return nil, errors.New("failed")
	}
	if r.attemptedRetries != nil {
		r.attemptedRetries <- d.id + " ok"
	}
	return d.id, nil
}

func (d *testData) mkFunc(run *testRun) Func {
	return func(_ context.Context) (interface{}, error) {
		return d.apply(run)
	}
}

type stringer string

func (s stringer) String() string {
	return string(s)
}

var _ fmt.Stringer = stringer("")

type panicStruct struct {
	s string
}

type mkExecutor func(options ...Option) Executor

func triggerNTimers(cs tbntime.ControlledSource, n int) {
	for triggered := 0; triggered < n; triggered += cs.TriggerAllTimers() {
		time.Sleep(1 * time.Millisecond)
	}
}

func testRetriesWithNoCallback(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		e := mk(
			WithTimeSource(cs),
			WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
			WithMaxAttempts(3),
		)
		defer e.Stop()

		run1 := &testRun{
			attemptedRetries: make(chan string, 10),
		}
		defer close(run1.attemptedRetries)

		run2 := &testRun{
			attemptedRetries: make(chan string, 10),
		}
		defer close(run2.attemptedRetries)

		run3 := &testRun{
			attemptedRetries: make(chan string, 10),
		}
		defer close(run3.attemptedRetries)

		p1 := &testData{"p1", 1, nil}
		p2 := &testData{"p2", 0, nil}
		p3 := &testData{"p3", 3, nil}

		e.ExecAndForget(p1.mkFunc(run1))
		e.ExecAndForget(p2.mkFunc(run2))
		e.ExecAndForget(p3.mkFunc(run3))

		// Wait for p1, p2, and p3 to complete their first attempts.
		triggerNTimers(cs, 3)

		assert.ArrayEqual(
			t,
			[]string{<-run1.attemptedRetries, <-run1.attemptedRetries},
			[]string{"p1 fail", "p1 ok"},
		)

		assert.ArrayEqual(
			t,
			[]string{<-run2.attemptedRetries},
			[]string{"p2 ok"},
		)

		assert.ArrayEqual(
			t,
			[]string{<-run3.attemptedRetries, <-run3.attemptedRetries, <-run3.attemptedRetries},
			[]string{"p3 fail", "p3 fail", "p3 fail"},
		)
	})
}

func testEarlierNextRetry(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		e := mk(
			WithTimeSource(cs),
			WithRetryDelayFunc(NewConstantDelayFunc(100*time.Millisecond)),
			WithMaxAttempts(4),
		)
		defer e.Stop()

		var wg sync.WaitGroup
		wg.Add(1)

		run := &testRun{
			attemptedRetries: make(chan string, 10),
		}
		defer close(run.attemptedRetries)

		p1 := &testData{"p1", 1, wg.Done}
		p2 := &testData{"p2", 0, nil}

		// Run p1 and wait for it to fail once.
		e.ExecAndForget(p1.mkFunc(run))
		wg.Wait()
		assert.Equal(t, <-run.attemptedRetries, "p1 fail")

		// Halfway through it's retry delay, start a new task and wait for it to complete.
		cs.Advance(50 * time.Millisecond)
		e.ExecAndForget(
			func(_ context.Context) (interface{}, error) {
				return p2.apply(run)
			},
		)
		assert.Equal(t, <-run.attemptedRetries, "p2 ok")

		// Complete p1's retry delay and wait for it to complete.
		cs.Advance(50 * time.Millisecond)
		assert.Equal(t, <-run.attemptedRetries, "p1 ok")
	})
}

func testExecInvokesCallback(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		e := mk(
			WithTimeSource(cs),
			WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
			WithMaxAttempts(3),
		)
		defer e.Stop()

		callbacks := make(chan string, 10)
		defer close(callbacks)

		run := &testRun{
			attemptedRetries: make(chan string, 10),
		}

		p1 := &testData{"p1", 1, nil}
		p2 := &testData{"p2", 0, nil}
		p3 := &testData{"p3", 3, nil}

		mkCallback := func(s string) CallbackFunc {
			return func(t Try) {
				if t.IsError() {
					callbacks <- s + " fail"
				} else {
					callbacks <- s + " ok"
				}
			}
		}

		// Start tasks, wait for each to complete its first attempt and check for p2's callback.
		e.Exec(p1.mkFunc(run), mkCallback("p1"))
		e.Exec(p2.mkFunc(run), mkCallback("p2"))
		e.Exec(p3.mkFunc(run), mkCallback("p3"))
		<-run.attemptedRetries
		<-run.attemptedRetries
		<-run.attemptedRetries
		assert.Equal(t, <-callbacks, "p2 ok")

		// Trigger retries, wait for attempts, and check for p1's callback.
		cs.Advance(50 * time.Millisecond)
		<-run.attemptedRetries
		<-run.attemptedRetries
		assert.Equal(t, <-callbacks, "p1 ok")

		// Trigger retries, wait for last attempt, and check for p3's callback.
		cs.Advance(50 * time.Millisecond)
		<-run.attemptedRetries
		assert.Equal(t, <-callbacks, "p3 fail")
	})
}

func testExecExecutesInParallel(t *testing.T, mk mkExecutor) {
	e := mk(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(3),
		WithParallelism(3),
	)
	defer e.Stop()

	sequence := make(chan string, 10)
	defer close(sequence)

	run := &testRun{}

	p1 := &testData{"p1", 0, nil}
	p2 := &testData{"p2", 0, nil}

	var in, s1, s2, d1, d2 sync.WaitGroup
	in.Add(2)
	s1.Add(1)
	s2.Add(1)
	d1.Add(1)
	d2.Add(1)

	e.Exec(
		func(_ context.Context) (interface{}, error) {
			sequence <- "p1 enter"
			in.Done()
			s1.Wait()
			sequence <- "p1 run"
			return p1.apply(run)
		},
		func(t Try) {
			sequence <- "p1 done"
			d1.Done()
		},
	)
	e.Exec(
		func(_ context.Context) (interface{}, error) {
			sequence <- "p2 enter"
			in.Done()
			s2.Wait()
			sequence <- "p2 run"
			return p2.apply(run)
		},
		func(t Try) {
			sequence <- "p2 done"
			d2.Done()
		},
	)

	in.Wait()
	s1.Done()
	d1.Wait()
	s2.Done()
	s2.Wait()

	messages := [6]string{}
	for i := 0; i < 6; i++ {
		messages[i] = <-sequence
	}

	assert.HasSameElements(t, messages[0:2], []string{"p1 enter", "p2 enter"})
	assert.ArrayEqual(t, messages[2:6], []string{"p1 run", "p1 done", "p2 run", "p2 done"})
}

func testExecPanicsBecomeErrors(t *testing.T, mk mkExecutor) {
	log, logBuffer := log.NewBufferLogger()

	e := mk(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(1),
		WithLogger(log),
	)
	defer e.Stop()

	tries := make(chan Try, 10)
	defer close(tries)

	p1panic := "p1 panic"
	p2panic := errors.New("p2 panic")
	p3panic := stringer("p3 panic")
	p4panic := panicStruct{"p4 panic"}

	mkPanicFunc := func(i interface{}) Func {
		return func(_ context.Context) (interface{}, error) {
			panic(i)
		}
	}

	e.Exec(mkPanicFunc(p1panic), func(t Try) { tries <- t })
	assert.Equal(t, (<-tries).Error().Error(), "p1 panic")
	assert.NotEqual(t, logBuffer.String(), "")
	logBuffer.Reset()

	e.Exec(mkPanicFunc(p2panic), func(t Try) { tries <- t })
	assert.Equal(t, (<-tries).Error().Error(), "p2 panic")
	assert.NotEqual(t, logBuffer.String(), "")
	logBuffer.Reset()

	e.Exec(mkPanicFunc(p3panic), func(t Try) { tries <- t })
	assert.Equal(t, (<-tries).Error().Error(), "p3 panic")
	assert.NotEqual(t, logBuffer.String(), "")
	logBuffer.Reset()

	e.Exec(mkPanicFunc(p4panic), func(t Try) { tries <- t })
	assert.Equal(t, (<-tries).Error().Error(), fmt.Sprintf("%#v", p4panic))
	assert.NotEqual(t, logBuffer.String(), "")
	logBuffer.Reset()
}

func testExecStopsWithInFlightRetries(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		e := mk(
			WithTimeSource(cs),
			WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
			WithMaxAttempts(2),
		)

		c := make(chan Try, 10)

		wg := &sync.WaitGroup{}
		wg.Add(1)

		e.Exec(
			func(_ context.Context) (interface{}, error) {
				defer wg.Done()
				return nil, errors.New("nope")
			},
			func(t Try) {
				c <- t
			},
		)

		// Await the task's first attempt.
		wg.Wait()

		// Stop the executor.
		e.Stop()

		// Advance the timer and expect that the task's final retry was not triggered.
		cs.Advance(50 * time.Millisecond)
		assert.ChannelEmpty(t, c)
	})
}
