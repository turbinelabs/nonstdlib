package executor

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

type testDiag struct {
	sync.Mutex

	taskCompletedCh chan struct{}

	tasksStarted int
	taskResults  []AttemptResult

	attemptsStarted int
	attemptResults  []AttemptResult

	callbacks int
}

func (t *testDiag) awaitTasks(n int) {
	for t.numCompleted() < n {
		t.awaitTask()
	}
}

func (t *testDiag) numCompleted() int {
	t.Lock()
	defer t.Unlock()

	return len(t.taskResults)
}

func (t *testDiag) awaitTask() {
	timer := time.NewTimer(10 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-t.taskCompletedCh:
	case <-timer.C:
	}
}

func (t *testDiag) TaskStarted(i int) {
	t.Lock()
	defer t.Unlock()

	t.tasksStarted++
}

func (t *testDiag) TaskCompleted(r AttemptResult, _ time.Duration) {
	t.Lock()
	defer t.Unlock()

	t.taskResults = append(t.taskResults, r)

	select {
	case t.taskCompletedCh <- struct{}{}:
	default:
	}
}

func (t *testDiag) AttemptStarted(_ time.Duration) {
	t.Lock()
	defer t.Unlock()

	t.attemptsStarted++
}

func (t *testDiag) AttemptCompleted(r AttemptResult, _ time.Duration) {
	t.Lock()
	defer t.Unlock()

	t.attemptResults = append(t.attemptResults, r)
}

func (t *testDiag) CallbackDuration(_ time.Duration) {
	t.Lock()
	defer t.Unlock()

	t.callbacks++
}

func newTestDiag() *testDiag {
	return &testDiag{
		taskCompletedCh: make(chan struct{}, 1),
	}
}

func testExecTaskDiagnostics(t *testing.T, mk mkExecutor) {
	diag := &testDiag{}

	e := mk(
		WithTimeout(50*time.Millisecond),
		WithDiagnostics(diag),
		WithParallelism(3),
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

	diag.awaitTasks(3)

	assert.Equal(t, diag.tasksStarted, 3)
	assert.Equal(t, len(diag.taskResults), 3)
	assert.Equal(t, diag.attemptsStarted, 3)
	assert.HasSameElements(
		t,
		diag.attemptResults,
		[]AttemptResult{AttemptSuccess, AttemptError, AttemptGlobalTimeout},
	)
	assert.Equal(t, diag.callbacks, 2)
}

func testExecAttemptDiagnostics(t *testing.T, mk mkExecutor) {
	diag := &testDiag{}

	e := mk(
		WithDiagnostics(diag),
		WithAttemptTimeout(50*time.Millisecond),
		WithRetryDelayFunc(NewConstantDelayFunc(10*time.Millisecond)),
		WithParallelism(3),
		WithMaxAttempts(2),
	)
	defer e.Stop()

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

	diag.awaitTasks(3)

	assert.Equal(t, diag.tasksStarted, 3)
	assert.HasSameElements(
		t,
		diag.taskResults,
		[]AttemptResult{AttemptSuccess, AttemptError, AttemptTimeout},
	)

	assert.Equal(t, diag.attemptsStarted, 5)
	assert.HasSameElements(
		t,
		diag.attemptResults,
		[]AttemptResult{
			AttemptSuccess,
			AttemptError, AttemptError,
			AttemptTimeout, AttemptTimeout,
		},
	)
	assert.Equal(t, diag.callbacks, 3)
}
