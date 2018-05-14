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
	"testing"
	"time"

	tbntime "github.com/turbinelabs/nonstdlib/time"
	"github.com/turbinelabs/test/assert"
)

func triggerNextContext(cs tbntime.ControlledSource) {
	for !cs.TriggerNextContext() {
		time.Sleep(1 * time.Millisecond)
	}
}

func triggerTimers(cs tbntime.ControlledSource) {
	for cs.TriggerAllTimers() == 0 {
		time.Sleep(1 * time.Millisecond)
	}
}

func testExecWithGlobalTimeoutSucceeds(t *testing.T, mk mkExecutor) {
	e := mk(
		WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
		WithMaxAttempts(3),
		WithTimeout(10*time.Second),
	)
	defer e.Stop()

	c := make(chan Try, 10)
	defer close(c)

	invocations := 0

	e.Exec(
		func(_ context.Context) (interface{}, error) {
			invocations++
			if invocations == 3 {
				return "ok", nil
			}
			return nil, errors.New("not yet")
		},
		func(try Try) {
			c <- try
		},
	)

	try := <-c

	assert.Equal(t, invocations, 3)
	if assert.True(t, try.IsReturn()) {
		assert.Equal(t, try.Get(), "ok")
	}
}

func testExecWithGlobalTimeoutTimesOut(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		e := mk(
			WithTimeout(10*time.Millisecond),
			WithTimeSource(cs),
		)
		defer e.Stop()

		c := make(chan Try, 10)
		defer close(c)

		e.Exec(
			func(ctxt context.Context) (interface{}, error) {
				<-ctxt.Done()
				return nil, errors.New("ctxt done")
			},
			func(try Try) {
				c <- try
			},
		)

		triggerNextContext(cs)

		try := <-c

		if assert.True(t, try.IsError()) {
			assert.ErrorContains(t, try.Error(), "action exceeded timeout (10ms)")
		}
	})
}

func testExecWithGlobalTimeoutTimesOutBeforeRetry(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		e := mk(
			WithRetryDelayFunc(NewConstantDelayFunc(10*time.Millisecond)),
			WithMaxAttempts(3),
			WithTimeout(15*time.Millisecond),
			WithTimeSource(cs),
		)
		defer e.Stop()

		c := make(chan Try, 10)
		defer close(c)

		sync := make(chan *struct{}, 10)
		defer close(sync)

		e.Exec(
			func(_ context.Context) (interface{}, error) {
				sync <- nil
				return nil, errors.New("error message")
			},
			func(try Try) {
				c <- try
			},
		)

		<-sync
		cs.Advance(0 * time.Millisecond)

		triggerTimers(cs)

		<-sync

		try := <-c

		if assert.True(t, try.IsError()) {
			assert.ErrorContains(
				t,
				try.Error(),
				"failed action would timeout before next retry: error message",
			)
		}
	})
}

func testExecManyWithGlobalTimeoutSucceeds(t *testing.T, mk mkExecutor) {
	e := mk(
		WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
		WithMaxAttempts(3),
		WithTimeout(10*time.Second),
		WithParallelism(3),
	)
	defer e.Stop()

	c := make(chan Try, 10)
	defer close(c)

	invocations := []int{0, 0, 0}

	mkFunc := func(i int) Func {
		return func(_ context.Context) (interface{}, error) {
			invocations[i]++
			if invocations[i] == 3 {
				return fmt.Sprintf("ok %d", i+1), nil
			}
			return nil, fmt.Errorf("not yet %d", i+1)
		}
	}

	e.ExecMany([]Func{mkFunc(0), mkFunc(1), mkFunc(2)}, func(_ int, try Try) { c <- try })

	try1 := <-c
	try2 := <-c
	try3 := <-c

	assert.DeepEqual(t, invocations, []int{3, 3, 3})
	ok1 := assert.True(t, try1.IsReturn())
	ok2 := assert.True(t, try2.IsReturn())
	ok3 := assert.True(t, try3.IsReturn())

	if ok1 && ok2 && ok3 {
		results := []string{
			try1.Get().(string),
			try2.Get().(string),
			try3.Get().(string),
		}
		assert.HasSameElements(t, results, []string{"ok 1", "ok 2", "ok 3"})
	}
}

func testExecManyWithGlobalTimeoutTimesOut(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		e := mk(
			WithTimeout(10*time.Millisecond),
			WithParallelism(3),
			WithTimeSource(cs),
		)
		defer e.Stop()

		rets := make(chan string, 10)
		defer close(rets)

		errs := make(chan error, 10)
		defer close(errs)

		f := func(_ context.Context) (interface{}, error) {
			return "ok", nil
		}

		timeoutF := func(ctxt context.Context) (interface{}, error) {
			<-ctxt.Done()
			return nil, errors.New("ctxt done")
		}

		e.ExecMany(
			[]Func{f, timeoutF, f},
			func(_ int, try Try) {
				if try.IsError() {
					errs <- try.Error()
				} else {
					rets <- try.Get().(string)
				}
			},
		)

		ret1 := <-rets
		ret2 := <-rets

		cs.Advance(10 * time.Millisecond)
		err := <-errs

		assert.Equal(t, ret1, "ok")
		assert.Equal(t, ret2, "ok")
		assert.ErrorContains(t, err, "action exceeded timeout (10ms)")
	})
}

func testExecGatheredWithGlobalTimeoutSucceeds(t *testing.T, mk mkExecutor) {
	e := mk(
		WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
		WithMaxAttempts(3),
		WithTimeout(10*time.Second),
		WithParallelism(3),
	)
	defer e.Stop()

	c := make(chan Try, 10)
	defer close(c)

	invocations := []int{0, 0, 0}

	mkFunc := func(i int) Func {
		return func(_ context.Context) (interface{}, error) {
			invocations[i]++
			if invocations[i] == 3 {
				return fmt.Sprintf("ok %d", i+1), nil
			}
			return nil, fmt.Errorf("not yet %d", i+1)
		}
	}

	e.ExecGathered([]Func{mkFunc(0), mkFunc(1), mkFunc(2)}, func(try Try) { c <- try })

	try := <-c

	if assert.True(t, try.IsReturn()) {
		assert.HasSameElements(t, try.Get(), []interface{}{"ok 1", "ok 2", "ok 3"})
	}
}

func testExecGatheredWithGlobalTimeoutTimesOut(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		e := mk(
			WithTimeout(10*time.Millisecond),
			WithParallelism(3),
			WithTimeSource(cs),
		)
		defer e.Stop()

		c := make(chan Try, 10)
		defer close(c)

		f := func(_ context.Context) (interface{}, error) {
			return "ok", nil
		}

		timeoutF := func(ctxt context.Context) (interface{}, error) {
			<-ctxt.Done()
			return nil, errors.New("ctxt done")
		}

		e.ExecGathered([]Func{f, timeoutF, f}, func(try Try) { c <- try })

		triggerNextContext(cs)

		try := <-c

		if assert.True(t, try.IsError()) {
			assert.ErrorContains(t, try.Error(), "action exceeded timeout (10ms)")
		}
	})
}

func testExecWithAttemptTimeoutSucceeds(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		e := mk(
			WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
			WithMaxAttempts(3),
			WithAttemptTimeout(10*time.Millisecond),
			WithTimeSource(cs),
		)
		defer e.Stop()

		c := make(chan Try, 10)
		defer close(c)

		invocations := 0

		e.Exec(
			func(ctxt context.Context) (interface{}, error) {
				invocations++
				if invocations == 3 {
					return "ok", nil
				}
				<-ctxt.Done()
				return nil, errors.New("not yet")
			},
			func(try Try) {
				c <- try
			},
		)

		triggerNextContext(cs) // timeout the attempt
		triggerTimers(cs)      // trigger 2nd attempt

		triggerNextContext(cs) // timeout the 2nd attempt
		triggerTimers(cs)      // trigger 3rd attempt

		try := <-c

		assert.Equal(t, invocations, 3)
		if assert.True(t, try.IsReturn()) {
			assert.Equal(t, try.Get(), "ok")
		}
	})
}

func testExecWithAttemptTimeoutTimesOut(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		e := mk(
			WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
			WithMaxAttempts(3),
			WithAttemptTimeout(10*time.Millisecond),
			WithTimeSource(cs),
		)
		defer e.Stop()

		c := make(chan Try, 10)
		defer close(c)

		invocations := 0

		e.Exec(
			func(ctxt context.Context) (interface{}, error) {
				invocations++
				<-ctxt.Done()
				return nil, errors.New("ctxt done")
			},
			func(try Try) {
				c <- try
			},
		)

		triggerNextContext(cs) // timeout the attempt
		triggerTimers(cs)      // trigger 2nd attempt

		triggerNextContext(cs) // timeout the 2nd attempt
		triggerTimers(cs)      // trigger 3rd attempt

		triggerNextContext(cs) // timeout the 3rd attempt
		try := <-c

		assert.Equal(t, invocations, 3)
		if assert.True(t, try.IsError()) {
			assert.ErrorContains(t, try.Error(), "action exceeded attempt timeout (10ms)")
		}
	})
}

func testExecWithGlobalAndAttemptTimeoutsTimesOut(t *testing.T, mk mkExecutor) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		e := mk(
			WithRetryDelayFunc(NewConstantDelayFunc(10*time.Millisecond)),
			WithMaxAttempts(10),
			WithAttemptTimeout(20*time.Millisecond),
			WithTimeout(50*time.Millisecond),
			WithTimeSource(cs),
		)
		defer e.Stop()

		c := make(chan Try, 10)
		defer close(c)

		sync := make(chan *struct{}, 10)
		defer close(sync)

		invocations := 0

		e.Exec(
			func(ctxt context.Context) (interface{}, error) {
				sync <- nil
				invocations++
				<-ctxt.Done()
				return nil, errors.New("ctxt done")
			},
			func(try Try) {
				c <- try
			},
		)

		<-sync
		cs.Advance(20 * time.Millisecond)

		triggerTimers(cs)

		<-sync
		cs.Advance(25 * time.Millisecond)

		try := <-c

		assert.GreaterThan(t, invocations, 1)
		if assert.True(t, try.IsError()) {
			assert.ErrorContains(t, try.Error(), "action exceeded timeout (50ms)")
		}
	})
}
