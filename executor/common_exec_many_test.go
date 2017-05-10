package executor

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func testExecMany(t *testing.T, mk mkExecutor) {
	e := mk(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(1),
		WithParallelism(3),
	)
	defer e.Stop()

	f1 := func(_ context.Context) (interface{}, error) { return "p1", nil }
	f2 := func(_ context.Context) (interface{}, error) { return "p2", nil }
	f3 := func(_ context.Context) (interface{}, error) { return "p3", nil }

	c := make(chan pair, 10)
	defer close(c)

	mcb := func(idx int, try Try) { c <- pair{idx, try} }

	e.ExecMany([]Func{f1, f2, f3}, mcb)

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

func testExecManyWithFailure(t *testing.T, mk mkExecutor) {
	e := mk(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(1),
		WithParallelism(3),
	)
	defer e.Stop()

	f1 := func(_ context.Context) (interface{}, error) { return "p1", nil }
	f2 := func(_ context.Context) (interface{}, error) { return nil, errors.New("p2") }
	f3 := func(_ context.Context) (interface{}, error) { return "p3", nil }

	c := make(chan pair, 10)
	defer close(c)

	mcb := func(idx int, try Try) { c <- pair{idx, try} }

	e.ExecMany([]Func{f1, f2, f3}, mcb)

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

func testExecManyNoop(t *testing.T, mk mkExecutor) {
	e := mk(
		WithMaxAttempts(1),
		WithParallelism(1),
	)
	e.Stop()

	e.ExecMany(
		[]Func{},
		func(i int, try Try) {
			assert.Tracing(t).Fatalf("unexpected callback: %d; %+v", i, try)
		},
	)
}

func testExecGathered(t *testing.T, mk mkExecutor) {
	e := mk(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(1),
		WithParallelism(3),
	)
	defer e.Stop()

	f1 := func(_ context.Context) (interface{}, error) { return "p1", nil }
	f2 := func(_ context.Context) (interface{}, error) { return "p2", nil }
	f3 := func(_ context.Context) (interface{}, error) { return "p3", nil }

	c := make(chan Try, 10)
	defer close(c)

	cb := func(try Try) { c <- try }

	e.ExecGathered([]Func{f1, f2, f3}, cb)

	try := <-c
	assert.True(t, try.IsReturn())

	results := try.Get().([]interface{})
	assert.Equal(t, results[0], "p1")
	assert.Equal(t, results[1], "p2")
	assert.Equal(t, results[2], "p3")
}

func testExecGatheredWithError(t *testing.T, mk mkExecutor) {
	e := mk(
		WithRetryDelayFunc(NewExponentialDelayFunc(50*time.Millisecond, time.Second)),
		WithMaxAttempts(1),
		WithParallelism(1),
	)
	defer e.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := make(chan Try, 10)
	defer close(c)

	// Because parallelism == 1, p1 succeeds and one of the other
	// attempts fails. The third is skipped because its context
	// has been canceled.
	f1 := func(_ context.Context) (interface{}, error) {
		return "p1", nil
	}
	f2 := func(ctxt context.Context) (interface{}, error) {
		if err := ctxt.Err(); err != nil {
			assert.Failed(t, "unexpected context error in p2")
		}

		return nil, errors.New("p2")
	}
	f3 := func(ctxt context.Context) (interface{}, error) {
		if err := ctxt.Err(); err != nil {
			assert.Failed(t, "unexpected context error in p3")
		}

		return nil, errors.New("p3")
	}

	cb := func(try Try) {
		c <- try
	}

	e.ExecGathered([]Func{f1, f2, f3}, cb)

	try := <-c
	assert.True(t, try.IsError())
	assert.True(t, try.Error().Error() == "p2" || try.Error().Error() == "p3")
}

func testExecGatheredNoCallback(t *testing.T, mk mkExecutor) {
	e := mk(
		WithMaxAttempts(1),
		WithParallelism(3),
	)
	defer e.Stop()

	c := make(chan string, 10)
	defer close(c)

	e.ExecGathered(
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

func testExecGatheredNoop(t *testing.T, mk mkExecutor) {
	e := mk(
		WithMaxAttempts(1),
		WithParallelism(1),
	)
	defer e.Stop()

	cbInvoked := false

	e.ExecGathered(
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
	e.ExecGathered([]Func{}, nil)
}
