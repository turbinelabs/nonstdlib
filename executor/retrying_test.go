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

package executor

import (
	"container/heap"
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

func TestRetryingExecHeapInterface(t *testing.T) {
	start := time.Now()
	items := []*retry{
		{deadline: start, attempts: 1},
		{deadline: start, attempts: 1},
		{deadline: start, attempts: 1},
		{deadline: start, attempts: 1},
		{deadline: start, attempts: 1},
	}

	q := &retryingExec{
		deadlineChan: make(chan time.Time, 10),
		q:            make([]*retry, 0, 10),
		delay:        func(_ int) time.Duration { return 0 * time.Second },
		maxAttempts:  2,
		time:         tbntime.NewSource(),
	}
	defer close(q.deadlineChan)

	heap.Init(q)

	assert.Nil(t, q.peek())

	for _, item := range items {
		assert.True(t, q.add(item))
	}

	peeked := q.peek()

	r := q.removeIfPast()
	assert.NonNil(t, r)
	assert.SameInstance(t, r, peeked)
	assert.Equal(t, r.deadline, start)

	r = q.removeIfPast()
	assert.NonNil(t, r)
	assert.Equal(t, r.deadline, start)

	r = q.removeIfPast()
	assert.NonNil(t, r)
	assert.Equal(t, r.deadline, start)

	r = q.removeIfPast()
	assert.NonNil(t, r)
	assert.Equal(t, r.deadline, start)

	r = q.removeIfPast()
	assert.NonNil(t, r)
	assert.Equal(t, r.deadline, start)

	r = q.removeIfPast()
	assert.Nil(t, r)

	for _, r := range items {
		deadline := <-q.deadlineChan
		assert.Equal(t, deadline, r.deadline)
	}

	assert.False(t, q.add(&retry{deadline: start, attempts: 2}))
	assert.Nil(t, q.peek())
}

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
			d.failFunc()
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

func TestRetryingExecHandleRetriesWithNoCallback(t *testing.T) {
	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
		WithMaxAttempts(3),
	)
	defer q.Stop()

	run := &testRun{
		attemptedRetries: make(chan string, 10),
	}
	defer close(run.attemptedRetries)

	p1 := &testData{"p1", 1, nil}
	p2 := &testData{"p2", 0, nil}
	p3 := &testData{"p3", 3, nil}

	q.ExecAndForget(p1.mkFunc(run))
	q.ExecAndForget(p2.mkFunc(run))
	q.ExecAndForget(p3.mkFunc(run))

	messages := [6]string{}
	for i := 0; i < 6; i++ {
		messages[i] = <-run.attemptedRetries
	}

	assert.HasSameElements(t, messages[0:3], []string{"p1 fail", "p2 ok", "p3 fail"})
	assert.HasSameElements(t, messages[3:5], []string{"p1 ok", "p3 fail"})
	assert.HasSameElements(t, messages[5:6], []string{"p3 fail"})
}

func TestRetryingExecHandleRetriesEarlierDeadline(t *testing.T) {
	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewConstantDelayFunc(100*time.Millisecond)),
		WithMaxAttempts(4),
	)
	defer q.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	run := &testRun{
		attemptedRetries: make(chan string, 10),
	}
	defer close(run.attemptedRetries)

	p1 := &testData{"p1", 1, wg.Done}
	p2 := &testData{"p2", 0, nil}

	go func() {
		wg.Wait()
		time.Sleep(50 * time.Millisecond)
		q.ExecAndForget(
			func(_ context.Context) (interface{}, error) {
				return p2.apply(run)
			},
		)
	}()

	q.ExecAndForget(p1.mkFunc(run))

	messages := make([]string, 3)
	for i := 0; i < 3; i++ {
		messages[i] = <-run.attemptedRetries
	}

	assert.HasSameElements(t, messages, []string{"p1 fail", "p2 ok", "p1 ok"})
}

func TestRetryingExecInvokesCallback(t *testing.T) {
	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
		WithMaxAttempts(2),
	)
	defer q.Stop()

	attemptedRetries := make(chan string, 10)
	defer close(attemptedRetries)

	run := &testRun{}

	p1 := &testData{"p1", 1, nil}
	p2 := &testData{"p2", 0, nil}
	p3 := &testData{"p3", 3, nil}

	mkCallback := func(s string) CallbackFunc {
		return func(t Try) {
			if t.IsError() {
				attemptedRetries <- s + " fail"
			} else {
				attemptedRetries <- s + " ok"
			}
		}
	}

	q.Exec(p1.mkFunc(run), mkCallback("p1"))
	q.Exec(p2.mkFunc(run), mkCallback("p2"))
	q.Exec(p3.mkFunc(run), mkCallback("p3"))

	messages := [3]string{}
	for i := 0; i < 3; i++ {
		messages[i] = <-attemptedRetries
	}

	assert.ArrayEqual(t, messages, []string{"p2 ok", "p1 ok", "p3 fail"})
}

func TestRetryingExecExecutesInParallel(t *testing.T) {
	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(3),
		WithParallelism(3),
	)
	defer q.Stop()

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

	q.Exec(
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
	q.Exec(
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

type stringer string

func (s stringer) String() string {
	return string(s)
}

var _ fmt.Stringer = stringer("")

type panicStruct struct {
	s string
}

func TestRetryingExecPanicsBecomeErrors(t *testing.T) {
	log, logBuffer := log.NewBufferLogger()

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(1),
		WithLogger(log),
	)
	defer q.Stop()

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

	q.Exec(mkPanicFunc(p1panic), func(t Try) { tries <- t })
	assert.Equal(t, (<-tries).Error().Error(), "p1 panic")
	assert.NotEqual(t, logBuffer.String(), "")
	logBuffer.Reset()

	q.Exec(mkPanicFunc(p2panic), func(t Try) { tries <- t })
	assert.Equal(t, (<-tries).Error().Error(), "p2 panic")
	assert.NotEqual(t, logBuffer.String(), "")
	logBuffer.Reset()

	q.Exec(mkPanicFunc(p3panic), func(t Try) { tries <- t })
	assert.Equal(t, (<-tries).Error().Error(), "p3 panic")
	assert.NotEqual(t, logBuffer.String(), "")
	logBuffer.Reset()

	q.Exec(mkPanicFunc(p4panic), func(t Try) { tries <- t })
	assert.Equal(t, (<-tries).Error().Error(), fmt.Sprintf("%#v", p4panic))
	assert.NotEqual(t, logBuffer.String(), "")
	logBuffer.Reset()
}

func TestRetryingExecStopsWithInFlightRetries(t *testing.T) {
	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewConstantDelayFunc(time.Second)),
		WithMaxAttempts(2),
	)

	c := make(chan Try, 10)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	q.Exec(
		func(_ context.Context) (interface{}, error) {
			defer wg.Done()
			return nil, errors.New("nope")
		},
		func(t Try) {
			c <- t
		},
	)

	wg.Wait()
	q.Stop()

	select {
	case <-c:
		assert.Failed(t, "unexpected invocation of pending retry's callback")
	default:
		// expected nothing
	}
}
