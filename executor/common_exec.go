package executor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	tbntime "github.com/turbinelabs/nonstdlib/time"
)

const (
	noTimeout = time.Duration(0)
)

var (
	noDeadline = time.Time{}
)

type contextErrorType int

const (
	noError contextErrorType = iota
	cancellationError
	attemptTimeoutError
	globalTimeoutError
)

type execImpl interface {
	add(*commonExec, *retry)
	retry(*commonExec, time.Duration, *retry) bool
	stop(*commonExec)
}

type commonExec struct {
	impl execImpl

	parallelism    int
	maxAttempts    int
	delay          DelayFunc
	timeout        time.Duration
	attemptTimeout time.Duration

	time tbntime.Source
	log  *log.Logger
	diag DiagnosticsCallback
}

type pair struct {
	idx int
	try Try
}

func (c *commonExec) ExecAndForget(f Func) {
	c.Exec(f, nil)
}

func (c *commonExec) Exec(f Func, cb CallbackFunc) {
	defer c.diag.TaskStarted(1)

	start := c.time.Now()
	globalDeadline := mkDeadline(start, c.timeout)
	ctxt, ctxtCancel := c.mkContext(globalDeadline, false)

	r := &retry{
		f:           f,
		cb:          cb,
		start:       start,
		nextAttempt: start,
		ctxt:        ctxt,
		ctxtCancel:  ctxtCancel,
		attempts:    0,
	}

	c.impl.add(c, r)
}

func (c *commonExec) ExecMany(fs []Func, cb ManyCallbackFunc) {
	defer c.diag.TaskStarted(len(fs))

	c.execMany(fs, cb)
}

func (c *commonExec) ExecGathered(fs []Func, cb CallbackFunc) {
	defer c.diag.TaskStarted(len(fs))

	c.execGathered(fs, cb)
}

func (c *commonExec) Stop() {
	c.impl.stop(c)
}

func (c *commonExec) SetDiagnosticsCallback(diag DiagnosticsCallback) {
	c.diag = diag
}

func (c *commonExec) execMany(
	fs []Func,
	cb ManyCallbackFunc,
) {
	if len(fs) == 0 {
		return
	}

	if cb == nil {
		cb = func(_ int, _ Try) {}
	}

	start := c.time.Now()
	globalDeadline := mkDeadline(start, c.timeout)
	ctxt, ctxtCancel := c.mkContext(globalDeadline, true)

	childWaiter := &sync.WaitGroup{}
	childWaiter.Add(len(fs))

	go func() {
		defer ctxtCancel()
		childWaiter.Wait()
	}()

	cancelingCb := func(i int, t Try) {
		defer cb(i, t)
		childWaiter.Done()
	}

	c.execChildren(ctxt, ctxtCancel, start, fs, cancelingCb)
}

func (c *commonExec) execGathered(
	fs []Func,
	cb CallbackFunc,
) {
	if cb == nil {
		c.execMany(fs, nil)
		return
	}

	n := len(fs)
	if n == 0 {
		cb(NewReturn([]interface{}{}))
		return
	}

	completed := make(chan pair, n)
	start := c.time.Now()
	globalDeadline := mkDeadline(start, c.timeout)
	ctxt, ctxtCancel := c.mkContext(globalDeadline, true)

	go func() {
		defer close(completed)
		results := make([]interface{}, n)

		for remaining := n; remaining > 0; remaining-- {
			p := <-completed
			if p.try.IsError() {
				if cb != nil {
					cb(p.try)
					cb = nil
				}
			} else {
				results[p.idx] = p.try.Get()
			}
		}

		if cb != nil {
			cb(NewReturn(results))
		}
	}()

	mcb := func(i int, t Try) {
		completed <- pair{i, t}
		if t.IsError() {
			ctxtCancel()
		}
	}

	c.execChildren(ctxt, ctxtCancel, start, fs, mcb)

}

func (c *commonExec) execChildren(
	ctxt context.Context,
	ctxtCancel context.CancelFunc,
	start time.Time,
	fs []Func,
	cb ManyCallbackFunc,
) {
	for i, f := range fs {
		idx := i
		indexingCb := func(t Try) {
			cb(idx, t)
		}

		childCtxt, childCancelFunc := c.mkChildContext(ctxt, noDeadline)

		r := &retry{
			f:           f,
			cb:          indexingCb,
			start:       start,
			nextAttempt: start,
			ctxt:        childCtxt,
			ctxtCancel:  childCancelFunc,
			attempts:    0,
		}

		c.impl.add(c, r)
	}
}

func (c *commonExec) mkContext(
	deadline time.Time,
	mayCancel bool,
) (context.Context, context.CancelFunc) {
	if !deadline.IsZero() {
		return c.time.NewContextWithDeadline(context.Background(), deadline)
	}

	if mayCancel {
		return context.WithCancel(context.Background())
	}

	return context.Background(), func() {}
}

func (c *commonExec) mkChildContext(
	ctxt context.Context,
	deadline time.Time,
) (context.Context, context.CancelFunc) {
	if !deadline.IsZero() {
		return c.time.NewContextWithDeadline(ctxt, deadline)
	}

	return context.WithCancel(ctxt)
}

func checkCtxtError(parent, child context.Context) contextErrorType {
	switch err := parent.Err(); err {
	case context.DeadlineExceeded:
		return globalTimeoutError

	case nil:
		if child == nil {
			return noError
		}
		// check child

	default:
		return cancellationError
	}

	switch err := child.Err(); err {
	case context.DeadlineExceeded:
		return attemptTimeoutError

	case nil:
		return noError

	default:
		return cancellationError
	}

}

func mkDeadline(start time.Time, d time.Duration) time.Time {
	if d > noTimeout {
		return start.Add(d)
	}

	return noDeadline
}

func (c *commonExec) attempt(r *retry) {
	attemptStart := c.time.Now()
	c.diag.AttemptStarted(attemptStart.Sub(r.nextAttempt))

	var t Try
	ctxtErrType := checkCtxtError(r.ctxt, nil)
	if ctxtErrType == noError {
		r.attempts++

		retryDeadline := mkDeadline(c.time.Now(), c.attemptTimeout)
		ctxt, localCancel := c.mkChildContext(r.ctxt, retryDeadline)

		t = rescuedCall(ctxt, r.f, c.log)
		ctxtErrType = checkCtxtError(r.ctxt, ctxt)
		localCancel()
	} else {
		t = NewError(r.ctxt.Err())
	}

	attemptDuration := c.time.Now().Sub(attemptStart)
	attemptResult := AttemptSuccess
	retry := false

	if t.IsError() {
		if ctxtErrType == globalTimeoutError {
			// global timeout expired
			attemptResult = AttemptGlobalTimeout

			t = NewError(
				fmt.Errorf(
					"action exceeded timeout (%s)",
					c.timeout,
				),
			)
		} else if ctxtErrType == cancellationError {
			// canceled
			attemptResult = AttemptCancellation

			t = NewError(errors.New("action canceled"))
		} else if ctxtErrType == attemptTimeoutError {
			// retry timeout expired, just count it
			attemptResult = AttemptTimeout
			t = NewError(
				fmt.Errorf(
					"action exceeded attempt timeout (%s)",
					c.attemptTimeout,
				),
			)

			retry = true
		} else {
			attemptResult = AttemptError

			// TODO: check if error is something want actually want to
			// retry (see #1686)
			retry = true
		}
	}

	c.diag.AttemptCompleted(attemptResult, attemptDuration)

	if retry {
		delay := c.delay(r.attempts)
		r.nextAttempt = c.time.Now().Add(delay)

		if limit, ok := r.ctxt.Deadline(); ok && limit.Before(r.nextAttempt) {
			// context will timeout before retry
			attemptResult = AttemptGlobalTimeout

			t = NewError(fmt.Errorf(
				"failed action would timeout before next retry: %s",
				t.Error().Error(),
			))
		} else if c.impl.retry(c, delay, r) {
			return
		}
	}

	defer func() { c.diag.TaskCompleted(attemptResult, c.time.Now().Sub(r.start)) }()

	if r.ctxtCancel != nil {
		r.ctxtCancel()
	}

	if r.cb != nil {
		callbackStart := c.time.Now()
		r.cb(t)
		c.diag.CallbackDuration(c.time.Now().Sub(callbackStart))
	}
}

func rescuedCall(ctxt context.Context, f Func, log *log.Logger) (t Try) {
	defer func() {
		if p := recover(); p != nil {
			stack := make([]byte, 2048)
			runtime.Stack(stack, false)

			if log != nil {
				log.Printf(
					"rescued retry queue call\npanic: %v\n\n%s\n",
					p,
					string(stack),
				)
			}

			switch s := p.(type) {
			case string:
				t = NewError(errors.New(s))

			case error:
				t = NewError(s)

			case fmt.Stringer:
				t = NewError(errors.New(s.String()))

			default:
				t = NewError(fmt.Errorf("%#v", p))
			}

		}
	}()

	t = NewTry(f(ctxt))
	return
}
