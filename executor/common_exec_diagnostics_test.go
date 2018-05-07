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
	"sync/atomic"
	"testing"
	"time"

	tbntime "github.com/turbinelabs/nonstdlib/time"
	"github.com/turbinelabs/test/assert"
)

type unit struct{}

type testDiag struct {
	taskStarts  chan unit
	taskResults chan AttemptResult

	attemptStarts  chan unit
	attemptResults chan AttemptResult

	callbacks int32
}

func newTestDiag(expectedTasks, expectedAttempts int) *testDiag {
	return &testDiag{
		taskStarts:     make(chan unit, expectedTasks*2),
		taskResults:    make(chan AttemptResult, expectedTasks*2),
		attemptStarts:  make(chan unit, expectedAttempts*2),
		attemptResults: make(chan AttemptResult, expectedAttempts*2),
	}
}

func (t *testDiag) countPendingTaskStarts() int {
	i := 0
	for {
		select {
		case <-t.taskStarts:
			i++
		default:
			return i
		}
	}
}

func (t *testDiag) countPendingAttemptStarts() int {
	i := 0
	for {
		select {
		case <-t.attemptStarts:
			i++
		default:
			return i
		}
	}
}

func (t *testDiag) TaskStarted(i int) {
	for i > 0 {
		select {
		case t.taskStarts <- unit{}:
		default:
			panic("taskStarts channel full")
		}
		i--
	}
}

func (t *testDiag) TaskCompleted(r AttemptResult, _ time.Duration) {
	select {
	case t.taskResults <- r:
	default:
		panic("taskResults channel full")
	}
}

func (t *testDiag) AttemptStarted(_ time.Duration) {
	select {
	case t.attemptStarts <- unit{}:
	default:
		panic("attemptStarts channel full")
	}
}

func (t *testDiag) AttemptCompleted(r AttemptResult, _ time.Duration) {
	select {
	case t.attemptResults <- r:
	default:
		panic("attemptResults channel full")
	}
}

func (t *testDiag) CallbackDuration(_ time.Duration) {
	atomic.AddInt32(&t.callbacks, 1)
}

func testExecTaskDiagnostics(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		diag := newTestDiag(3, 3)

		e := mk(
			WithTimeout(50*time.Millisecond),
			WithRetryDelayFunc(NewConstantDelayFunc(0*time.Millisecond)),
			WithDiagnostics(diag),
			WithParallelism(3),
			WithTimeSource(cs),
		)
		defer e.Stop()

		e.ExecAndForget(
			func(ctxt context.Context) (interface{}, error) {
				return "ok", nil
			},
		)

		e.Exec(
			func(ctxt context.Context) (interface{}, error) {
				return nil, errors.New("i failed")
			},
			func(try Try) {},
		)

		e.Exec(
			func(ctxt context.Context) (interface{}, error) {
				<-ctxt.Done()
				return nil, errors.New("ctxt done")
			},
			func(try Try) {},
		)

		// Await result of successful and failed tasks.
		assert.HasSameElements(
			t,
			[]AttemptResult{<-diag.attemptResults, <-diag.attemptResults},
			[]AttemptResult{AttemptSuccess, AttemptError},
		)

		cs.Advance(50 * time.Millisecond)
		assert.Equal(t, <-diag.attemptResults, AttemptGlobalTimeout)

		assert.HasSameElements(
			t,
			[]AttemptResult{<-diag.taskResults, <-diag.taskResults, <-diag.taskResults},
			[]AttemptResult{AttemptSuccess, AttemptError, AttemptGlobalTimeout},
		)

		assert.Equal(t, diag.countPendingTaskStarts(), 3)
		assert.Equal(t, diag.countPendingAttemptStarts(), 3)
		assert.Equal(t, diag.callbacks, int32(2))
	})
}

func testExecAttemptDiagnostics(t *testing.T, mk mkExecutor) {

	tbntime.WithTimeAt(time.Now().Truncate(time.Hour), func(cs tbntime.ControlledSource) {
		tick := func(s string) { fmt.Println("time:", cs.Now(), "@", s) }

		diag := newTestDiag(3, 5)

		e := mk(
			WithDiagnostics(diag),
			WithAttemptTimeout(50*time.Millisecond),
			WithRetryDelayFunc(NewConstantDelayFunc(10*time.Millisecond)),
			WithParallelism(3),
			WithMaxAttempts(2),
			WithTimeSource(cs),
		)
		defer e.Stop()

		tick("start")

		e.Exec(
			func(ctxt context.Context) (interface{}, error) {
				return "ok", nil
			},
			func(try Try) {},
		)

		e.Exec(
			func(ctxt context.Context) (interface{}, error) {
				return nil, errors.New("i failed")
			},
			func(try Try) {},
		)

		e.Exec(
			func(ctxt context.Context) (interface{}, error) {
				<-ctxt.Done()
				return nil, errors.New("ctxt done")
			},
			func(try Try) {},
		)

		// 3 tasks started
		assert.Equal(t, <-diag.taskStarts, unit{})
		assert.Equal(t, <-diag.taskStarts, unit{})
		assert.Equal(t, <-diag.taskStarts, unit{})

		// 3 attempts started
		assert.Equal(t, <-diag.attemptStarts, unit{})
		assert.Equal(t, <-diag.attemptStarts, unit{})
		assert.Equal(t, <-diag.attemptStarts, unit{})

		// 1 attempt success, 1 fails
		assert.HasSameElements(
			t,
			[]AttemptResult{<-diag.attemptResults, <-diag.attemptResults},
			[]AttemptResult{AttemptSuccess, AttemptError},
		)

		// 1 task succeeds
		assert.Equal(t, <-diag.taskResults, AttemptSuccess)

		// Trigger failed task retry
		tick("trigger retry")
		for !cs.TriggerNextTimer() {
		}
		tick("triggered retry")

		// 1 attempt starts (retry), and fails again
		assert.Equal(t, <-diag.attemptStarts, unit{})
		assert.Equal(t, <-diag.attemptResults, AttemptError)
		assert.Equal(t, <-diag.taskResults, AttemptError)

		// Trigger slow task timeout
		tick("trigger timeout 1")
		for !cs.TriggerNextContext() {
		}
		tick("triggered timeout 1")

		// 1 attempt times outs
		assert.Equal(t, <-diag.attemptResults, AttemptTimeout)

		// Trigger timed out task retry
		tick("trigger retry timeout")
		for !cs.TriggerNextTimer() {
		}
		tick("triggered retry timeout")

		// 1 attempt starts (retry)
		assert.Equal(t, <-diag.attemptStarts, unit{})

		// Trigger slow task timeout #2
		tick("trigger timeout 2")
		for !cs.TriggerNextContext() {
		}
		tick("triggered timeout 2")

		assert.Equal(t, <-diag.attemptResults, AttemptTimeout)
		assert.Equal(t, <-diag.taskResults, AttemptTimeout)

		assert.Equal(t, diag.callbacks, int32(3))
	})
}
