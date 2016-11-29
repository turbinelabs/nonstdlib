package executor

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func TestRetryingExecExecWithGlobalTimeoutSucceeds(t *testing.T) {
	c := make(chan Try, 10)
	defer close(c)

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
		WithMaxAttempts(3),
		WithTimeout(1*time.Second),
	)
	defer q.Stop()

	invocations := 0

	q.Exec(
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

func TestRetryingExecExecWithGlobalTimeoutTimesOut(t *testing.T) {
	c := make(chan Try, 10)
	defer close(c)

	q := NewRetryingExecutor(
		WithTimeout(10 * time.Millisecond),
	)
	defer q.Stop()

	q.Exec(
		func(ctxt context.Context) (interface{}, error) {
			<-ctxt.Done()
			return nil, errors.New("ctxt done")
		},
		func(try Try) {
			c <- try
		},
	)

	try := <-c

	if assert.True(t, try.IsError()) {
		assert.ErrorContains(t, try.Error(), "action exceeded timeout (10ms)")
	}
}

func TestRetryingExecExecWithGlobalTimeoutTimesOutBeforeRetry(t *testing.T) {
	c := make(chan Try, 10)
	defer close(c)

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
		WithMaxAttempts(3),
		WithTimeout(140*time.Millisecond),
	)
	defer q.Stop()

	q.Exec(
		func(_ context.Context) (interface{}, error) {
			return nil, errors.New("error message")
		},
		func(try Try) {
			c <- try
		},
	)

	try := <-c

	if assert.True(t, try.IsError()) {
		assert.ErrorContains(
			t,
			try.Error(),
			"failed action would timeout before next retry: error message",
		)
	}
}

func TestRetryingExecExecManyWithGlobalTimeoutSucceeds(t *testing.T) {
	c := make(chan Try, 10)
	defer close(c)

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
		WithMaxAttempts(3),
		WithTimeout(time.Second),
		WithParallelism(3),
	)
	defer q.Stop()

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

	q.ExecMany([]Func{mkFunc(0), mkFunc(1), mkFunc(2)}, func(_ int, try Try) { c <- try })

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

func TestRetryingExecExecManyWithGlobalTimeoutTimesOut(t *testing.T) {
	rets := make(chan string, 10)
	defer close(rets)

	errs := make(chan error, 10)
	defer close(errs)

	q := NewRetryingExecutor(
		WithTimeout(10*time.Millisecond),
		WithParallelism(3),
	)
	defer q.Stop()

	f := func(_ context.Context) (interface{}, error) {
		return "ok", nil
	}

	timeoutF := func(ctxt context.Context) (interface{}, error) {
		<-ctxt.Done()
		return nil, errors.New("ctxt done")
	}

	q.ExecMany(
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
	err := <-errs

	assert.Equal(t, ret1, "ok")
	assert.Equal(t, ret2, "ok")
	assert.ErrorContains(t, err, "action exceeded timeout (10ms)")
}

func TestRetryingExecExecGatheredWithGlobalTimeoutSucceeds(t *testing.T) {
	c := make(chan Try, 10)
	defer close(c)

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
		WithMaxAttempts(3),
		WithTimeout(time.Second),
		WithParallelism(3),
	)
	defer q.Stop()

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

	q.ExecGathered([]Func{mkFunc(0), mkFunc(1), mkFunc(2)}, func(try Try) { c <- try })

	try := <-c

	if assert.True(t, try.IsReturn()) {
		assert.HasSameElements(t, try.Get(), []interface{}{"ok 1", "ok 2", "ok 3"})
	}
}

func TestRetryingExecExecGatheredWithGlobalTimeoutTimesOut(t *testing.T) {
	c := make(chan Try, 10)
	defer close(c)

	q := NewRetryingExecutor(
		WithTimeout(10*time.Millisecond),
		WithParallelism(3),
	)
	defer q.Stop()

	f := func(_ context.Context) (interface{}, error) {
		return "ok", nil
	}

	timeoutF := func(ctxt context.Context) (interface{}, error) {
		<-ctxt.Done()
		return nil, errors.New("ctxt done")
	}

	q.ExecGathered([]Func{f, timeoutF, f}, func(try Try) { c <- try })

	try := <-c

	if assert.True(t, try.IsError()) {
		assert.ErrorContains(t, try.Error(), "action exceeded timeout (10ms)")
	}
}

func TestRetryingExecExecWithAttemptTimeoutSucceeds(t *testing.T) {
	c := make(chan Try, 10)
	defer close(c)

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewConstantDelayFunc(50*time.Millisecond)),
		WithMaxAttempts(3),
		WithAttemptTimeout(10*time.Millisecond),
	)
	defer q.Stop()

	invocations := 0

	q.Exec(
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

	try := <-c

	assert.Equal(t, invocations, 3)
	if assert.True(t, try.IsReturn()) {
		assert.Equal(t, try.Get(), "ok")
	}
}

func TestRetryingExecExecWithAttemptTimeoutTimesOut(t *testing.T) {
	c := make(chan Try, 10)
	defer close(c)

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewConstantDelayFunc(1*time.Millisecond)),
		WithMaxAttempts(3),
		WithAttemptTimeout(10*time.Millisecond),
	)
	defer q.Stop()

	invocations := 0

	q.Exec(
		func(ctxt context.Context) (interface{}, error) {
			invocations++
			<-ctxt.Done()
			return nil, errors.New("ctxt done")
		},
		func(try Try) {
			c <- try
		},
	)

	try := <-c

	assert.Equal(t, invocations, 3)
	if assert.True(t, try.IsError()) {
		assert.ErrorContains(t, try.Error(), "action exceeded attempt timeout (10ms)")
	}
}

func TestRetryingExecExecWithGlobalAndAttemptTimeoutsTimesOut(t *testing.T) {
	c := make(chan Try, 10)
	defer close(c)

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewConstantDelayFunc(1*time.Millisecond)),
		WithMaxAttempts(10),
		WithAttemptTimeout(20*time.Millisecond),
		WithTimeout(50*time.Millisecond),
	)
	defer q.Stop()

	invocations := 0

	q.Exec(
		func(ctxt context.Context) (interface{}, error) {
			invocations++
			<-ctxt.Done()
			return nil, errors.New("ctxt done")
		},
		func(try Try) {
			c <- try
		},
	)

	try := <-c

	assert.True(t, invocations > 1)
	if assert.True(t, try.IsError()) {
		assert.ErrorContains(t, try.Error(), "action exceeded timeout (50ms)")
	}
}
