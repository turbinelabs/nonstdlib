package executor

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func TestCancelCounter(t *testing.T) {
	for _, n := range []int{1, 2, 4} {
		actualCancels := 0
		c := newCancelCounter(n, func() { actualCancels++ })

		for i := 0; i < n-1; i++ {
			c.cancel()
		}
		assert.Equal(t, actualCancels, 0)

		c.cancel()
		assert.Equal(t, actualCancels, 1)

		c.cancel()
		assert.Equal(t, actualCancels, 2)
	}
}

func TestRetryingExecExecMany(t *testing.T) {
	f1 := func(_ context.Context) (interface{}, error) { return "p1", nil }
	f2 := func(_ context.Context) (interface{}, error) { return "p2", nil }
	f3 := func(_ context.Context) (interface{}, error) { return "p3", nil }

	c := make(chan pair, 10)
	defer close(c)

	mcb := func(idx int, try Try) { c <- pair{idx, try} }

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(1),
		WithParallelism(3),
	)
	defer q.Stop()

	q.ExecMany([]Func{f1, f2, f3}, mcb)

	pairs := make([]pair, 0, 3)
	for i := 0; i < 3; i++ {
		p := <-c

		var j int
		for j = 0; j < len(pairs); j++ {
			if p.idx < pairs[j].idx {
				pairs = append(pairs[:j], append([]pair{p}, pairs[j:]...)...)
				break
			}
		}
		if j == len(pairs) {
			pairs = append(pairs, p)
		}
	}

	assert.Equal(t, pairs[0].idx, 0)
	assert.True(t, pairs[0].try.IsReturn())
	assert.Equal(t, pairs[0].try.Get(), "p1")

	assert.Equal(t, pairs[1].idx, 1)
	assert.True(t, pairs[1].try.IsReturn())
	assert.Equal(t, pairs[1].try.Get(), "p2")

	assert.Equal(t, pairs[2].idx, 2)
	assert.True(t, pairs[2].try.IsReturn())
	assert.Equal(t, pairs[2].try.Get(), "p3")
}

func TestRetryingExecExecManyWithFailure(t *testing.T) {
	f1 := func(_ context.Context) (interface{}, error) { return "p1", nil }
	f2 := func(_ context.Context) (interface{}, error) { return nil, errors.New("p2") }
	f3 := func(_ context.Context) (interface{}, error) { return "p3", nil }

	c := make(chan pair, 10)
	defer close(c)

	mcb := func(idx int, try Try) { c <- pair{idx, try} }

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(1),
		WithParallelism(3),
	)
	defer q.Stop()

	q.ExecMany([]Func{f1, f2, f3}, mcb)

	pairs := make([]pair, 0, 3)
	for i := 0; i < 3; i++ {
		p := <-c

		var j int
		for j = 0; j < len(pairs); j++ {
			if p.idx < pairs[j].idx {
				pairs = append(pairs[:j], append([]pair{p}, pairs[j:]...)...)
				break
			}
		}
		if j == len(pairs) {
			pairs = append(pairs, p)
		}
	}

	assert.Equal(t, pairs[0].idx, 0)
	assert.True(t, pairs[0].try.IsReturn())
	assert.Equal(t, pairs[0].try.Get(), "p1")

	assert.Equal(t, pairs[1].idx, 1)
	assert.True(t, pairs[1].try.IsError())
	assert.Equal(t, pairs[1].try.Error().Error(), "p2")

	assert.Equal(t, pairs[2].idx, 2)
	assert.True(t, pairs[2].try.IsReturn())
	assert.Equal(t, pairs[2].try.Get(), "p3")
}

func TestRetryingExecExecManyNoop(t *testing.T) {
	q := NewRetryingExecutor(
		WithMaxAttempts(1),
		WithParallelism(1),
	)
	q.Stop()

	// If we passed functions this would panic.
	q.ExecMany([]Func{}, func(_ int, _ Try) {})
}

func TestRetryingExecExecGathered(t *testing.T) {
	f1 := func(_ context.Context) (interface{}, error) { return "p1", nil }
	f2 := func(_ context.Context) (interface{}, error) { return "p2", nil }
	f3 := func(_ context.Context) (interface{}, error) { return "p3", nil }

	c := make(chan Try, 10)
	defer close(c)

	cb := func(try Try) { c <- try }

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(1),
		WithParallelism(3),
	)
	defer q.Stop()

	q.ExecGathered([]Func{f1, f2, f3}, cb)

	try := <-c
	assert.True(t, try.IsReturn())

	results := try.Get().([]interface{})
	assert.Equal(t, results[0], "p1")
	assert.Equal(t, results[1], "p2")
	assert.Equal(t, results[2], "p3")
}

func TestRetryingExecExecGatheredWithError(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := make(chan Try, 10)
	defer close(c)

	ctxtErr := make(chan string, 10)
	defer close(ctxtErr)

	f1 := func(_ context.Context) (interface{}, error) {
		return "p1", nil
	}
	f2 := func(ctxt context.Context) (interface{}, error) {
		if err := ctxt.Err(); err != nil {
			ctxtErr <- "p2"
		}

		return nil, errors.New("p2")
	}
	f3 := func(ctxt context.Context) (interface{}, error) {
		if err := ctxt.Err(); err != nil {
			ctxtErr <- "p3"
		}

		return nil, errors.New("p3")
	}

	cb := func(try Try) { c <- try }

	q := NewRetryingExecutor(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(1),
		WithParallelism(1),
	)
	defer q.Stop()

	q.ExecGathered([]Func{f1, f2, f3}, cb)

	try := <-c
	assert.True(t, try.IsError())

	errFunc := <-ctxtErr

	// Because parallelism == 1, one thread gets a context error
	// (the context is canceled when the other thread returns a
	// failure)
	if errFunc == "p2" {
		assert.Equal(t, try.Error().Error(), "p3")
	} else {
		assert.Equal(t, errFunc, "p3")
		assert.Equal(t, try.Error().Error(), "p2")
	}
}

func TestRetryingExecExecGatheredNoCallback(t *testing.T) {
	q := NewRetryingExecutor(
		WithMaxAttempts(1),
		WithParallelism(3),
	)
	defer q.Stop()

	c := make(chan string, 10)
	defer close(c)

	q.ExecGathered(
		[]Func{
			func(_ context.Context) (interface{}, error) { c <- "p1"; return nil, nil },
			func(_ context.Context) (interface{}, error) { c <- "p2"; return nil, nil },
			func(_ context.Context) (interface{}, error) { c <- "p3"; return nil, nil },
		},
		nil,
	)

	results := []string{<-c, <-c, <-c}
	assert.HasSameElements(t, results, []string{"p1", "p2", "p3"})
}

func TestRetryingExecExecGatheredNoop(t *testing.T) {
	q := NewRetryingExecutor(
		WithMaxAttempts(1),
		WithParallelism(1),
	)
	q.Stop()

	cbInvoked := false

	q.ExecGathered(
		[]Func{},
		func(try Try) {
			cbInvoked = true

			assert.True(t, try.IsReturn())
			result, ok := try.Get().([]interface{})
			assert.True(t, ok)
			assert.Equal(t, len(result), 0)
		},
	)

	assert.True(t, cbInvoked)

	// doesn't crash:
	q.ExecGathered([]Func{}, nil)
}
