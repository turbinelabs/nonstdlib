package executor

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/turbinelabs/nonstdlib/ptr"
	"github.com/turbinelabs/nonstdlib/stats"
)

const (
	defaultMaxAttempts   = 1
	defaultMaxQueueDepth = 10
	defaultParallelism   = 1
	noTimeout            = time.Duration(0)

	execHandletime         = "exec_handletime"
	execManyHandletime     = "exec_many_handletime"
	execGatheredHandletime = "exec_gathered_handletime"
	execManyWidth          = "exec_many_width"
	execGatheredWidth      = "exec_gathered_width"

	actionDelaytime = "action_delaytime"
	actionAttempts  = "action_attempts"

	actionFailure          = "action_failure.error"
	actionGlobalTimeout    = "action_failure.timeout"
	actionRetryTimeout     = "action_failure.attempt_timeout"
	actionCanceled         = "action_failure.canceled"
	actionRetry            = "action_retry"
	actionRetriesExhausted = "action_retries_exhausted"
)

type contextErrorType int

const (
	noError contextErrorType = iota
	cancellation
	attemptTimeout
	globalTimeout
)

var (
	defaultDelayFunc = NewConstantDelayFunc(1 * time.Second)

	noLimit    = time.Time{}
	exitSignal = time.Time{}

	gatherError = NewError(errors.New("underflowed results for ExecGathered"))
)

// Retry encapsulates a deadline (when to retry the action), a limit
// (when to time out), the data necessary to perform the retry, and
// how many attempts have already been made.
type retry struct {
	f          Func
	cb         CallbackFunc
	start      time.Time
	deadline   time.Time
	ctxt       context.Context
	ctxtCancel context.CancelFunc
	attempts   int
}

type retryingExec struct {
	sync.Mutex

	deadlineChan chan time.Time
	execChan     chan *retry
	q            []*retry

	parallelism    int
	maxQueueDepth  int
	maxAttempts    int
	delay          DelayFunc
	timeout        time.Duration
	attemptTimeout time.Duration

	log   *log.Logger
	stats stats.Stats
}

func newCancelCounter(i int, cancelFunc context.CancelFunc) *cancelCounter {
	return &cancelCounter{
		n:          ptr.Int32(int32(i)),
		cancelFunc: cancelFunc,
	}
}

type cancelCounter struct {
	n          *int32
	cancelFunc context.CancelFunc
}

func (c *cancelCounter) cancel() {
	if atomic.AddInt32(c.n, -1) <= 0 {
		c.cancelFunc()
	}
}

var _ context.CancelFunc = (&cancelCounter{}).cancel

type RetryingExecutorOption func(*retryingExec)

// Sets a Logger for panics recovered while executing actions.
func WithLogger(log *log.Logger) RetryingExecutorOption {
	return func(q *retryingExec) {
		q.log = log
	}
}

// Sets the DelayFunc used when retrying actions.
func WithRetryDelayFunc(d DelayFunc) RetryingExecutorOption {
	return func(q *retryingExec) {
		q.delay = d
	}
}

// Sets the absolute maximum number of attempts made to complete an
// action (including the initial attempt). Values less than 1 act as
// if 1 had been passed.
func WithMaxAttempts(maxAttempts int) RetryingExecutorOption {
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	return func(q *retryingExec) {
		q.maxAttempts = maxAttempts
	}
}

// Sets the number of goroutines used to execute actions. No more than
// this many actions can be executing at once. Values less than 1 act
// as if 1 has been passed.
func WithParallelism(parallelism int) RetryingExecutorOption {
	if parallelism < 1 {
		parallelism = 1
	}

	return func(q *retryingExec) {
		q.parallelism = parallelism
	}
}

// Sets the maximum number of actions pending immediate execution. If
// all worker goroutines are processing actions, the number of items
// that can be pending execution (initial or retry) before blocking
// occurs.
func WithMaxQueueDepth(maxQueueDepth int) RetryingExecutorOption {
	if maxQueueDepth < 1 {
		maxQueueDepth = 1
	}

	return func(q *retryingExec) {
		q.maxQueueDepth = maxQueueDepth
	}
}

// Sets the timeout for completion of actions. If the action has not
// completed (including retries) within the given duration, it is
// canceled. Timeouts less than or equal to zero are treated as "no
// time out."
func WithTimeout(timeout time.Duration) RetryingExecutorOption {
	if timeout <= noTimeout {
		timeout = noTimeout
	}

	return func(q *retryingExec) {
		q.timeout = timeout
	}
}

// Sets the timeout for completion individual attempts of an
// action. If the attempt has not completed within the given duration,
// it is canceled (and potentially retried). Timeouts less than or
// equal to zero are treated as "no time out."
func WithAttemptTimeout(timeout time.Duration) RetryingExecutorOption {
	if timeout <= noTimeout {
		timeout = noTimeout
	}

	return func(q *retryingExec) {
		q.attemptTimeout = timeout
	}
}

// Constructs a new Executor. By default, it never retries, has
// parallelism of 1, and a maximum queue depth of 10.
func NewRetryingExecutor(options ...RetryingExecutorOption) Executor {
	q := &retryingExec{
		deadlineChan:   make(chan time.Time, 2),
		q:              make([]*retry, 0, 10),
		parallelism:    defaultParallelism,
		maxQueueDepth:  defaultMaxQueueDepth,
		maxAttempts:    defaultMaxAttempts,
		delay:          defaultDelayFunc,
		timeout:        noTimeout,
		attemptTimeout: noTimeout,
		stats:          stats.NewNoopStats(),
	}

	for _, apply := range options {
		apply(q)
	}

	q.execChan = make(chan *retry, q.maxQueueDepth)

	heap.Init(q)
	go q.handleRetries()

	for i := 0; i < q.parallelism; i++ {
		go q.handleExec()
	}

	if q.log != nil {
		q.log.Printf(
			"retrying executor: %d workers, max queue %d, max attempts %d, global timeout %s, attempt timeout %s",
			q.parallelism,
			q.maxQueueDepth,
			q.maxAttempts,
			q.timeout,
			q.attemptTimeout,
		)
	}
	return q
}

// Implements heap.Interface
func (q *retryingExec) Len() int           { return len(q.q) }
func (q *retryingExec) Less(i, j int) bool { return q.q[i].deadline.Before(q.q[j].deadline) }
func (q *retryingExec) Swap(i, j int)      { q.q[i], q.q[j] = q.q[j], q.q[i] }

// Implements heap.Interface
func (q *retryingExec) Push(x interface{}) {
	item := x.(*retry)
	q.q = append(q.q, item)
}

// Implements heap.Interface
func (q *retryingExec) Pop() interface{} {
	old := q.q
	n := len(old)
	item := old[n-1]
	q.q = old[0 : n-1]
	return item
}

var _ heap.Interface = &retryingExec{}

func (q *retryingExec) ExecAndForget(f Func) {
	q.Exec(f, nil)
}

func (q *retryingExec) mkContext(
	t time.Time,
	mayCancel bool,
) (context.Context, context.CancelFunc) {
	if q.timeout > noTimeout {
		return context.WithDeadline(context.Background(), t.Add(q.timeout))
	}

	if mayCancel {
		return context.WithCancel(context.Background())
	}

	return context.Background(), func() {}
}

func (q *retryingExec) Exec(f Func, cb CallbackFunc) {
	start := time.Now()
	defer func() {
		q.stats.TimingDuration(execHandletime, time.Now().Sub(start))
	}()

	ctxt, ctxtCancel := q.mkContext(start, false)

	r := &retry{
		f:          f,
		cb:         cb,
		start:      start,
		deadline:   start,
		ctxt:       ctxt,
		ctxtCancel: ctxtCancel,
		attempts:   0,
	}

	q.add(r)
}

func (q *retryingExec) execMany(
	start time.Time,
	ctxt context.Context,
	ctxtCancel context.CancelFunc,
	fs []Func,
	cb ManyCallbackFunc,
) {
	canceler := newCancelCounter(len(fs), ctxtCancel)

	for i, f := range fs {
		var indexedCb CallbackFunc
		if cb != nil {
			idx := i
			indexedCb = func(t Try) { cb(idx, t) }
		}

		r := &retry{
			f:          f,
			cb:         indexedCb,
			start:      start,
			deadline:   start,
			ctxt:       ctxt,
			ctxtCancel: canceler.cancel,
			attempts:   0,
		}

		q.add(r)
	}
}

func (q *retryingExec) ExecMany(fs []Func, cb ManyCallbackFunc) {
	start := time.Now()
	defer func() {
		q.stats.TimingDuration(execManyHandletime, time.Now().Sub(start))
		q.stats.Inc(execManyWidth, int64(len(fs)))
	}()

	if len(fs) == 0 {
		return
	}

	ctxt, ctxtCancel := q.mkContext(start, false)

	q.execMany(start, ctxt, ctxtCancel, fs, cb)
}

type pair struct {
	idx int
	try Try
}

func (q *retryingExec) ExecGathered(fs []Func, cb CallbackFunc) {
	n := len(fs)
	start := time.Now()
	defer func() {
		q.stats.TimingDuration(execGatheredHandletime, time.Now().Sub(start))
		q.stats.Inc(execGatheredWidth, int64(n))
	}()

	if n == 0 {
		if cb != nil {
			cb(NewReturn([]interface{}{}))
		}
		return
	}

	if cb == nil {
		// Don't bother tracking results if the caller doesn't
		// want a call back.
		ctxt, ctxtCancel := q.mkContext(start, false)
		q.execMany(start, ctxt, ctxtCancel, fs, nil)
		return
	}

	completed := make(chan pair, n)
	ctxt, ctxtCancel := q.mkContext(start, true)

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

	q.execMany(start, ctxt, ctxtCancel, fs, mcb)
}

func (q *retryingExec) Stop() {
	// handleRetries() closes q.execChan as it returns to avoid
	// writes to execChan after close.
	q.deadlineChan <- exitSignal
}

// Adds an entry into the queue if it has not exhausted its
// attempts. Enqueues the deadline on the queue's channel.
func (q *retryingExec) add(r *retry) bool {
	if r.attempts >= q.maxAttempts {
		return false
	}

	if r.deadline.IsZero() {
		return false
	}

	q.Lock()
	defer func() {
		q.Unlock()

		q.deadlineChan <- r.deadline
	}()

	heap.Push(q, r)

	return true
}

// Removes the first entry of the queue if it's deadline has
// expired. If there are no items or the head of the queue is not
// expired, returns nil.
func (q *retryingExec) removeIfPast() *retry {
	q.Lock()
	defer q.Unlock()

	if len(q.q) > 0 && !q.q[0].deadline.After(time.Now()) {
		r := heap.Pop(q)
		return r.(*retry)
	}

	return nil
}

// Peeks at the first entry in the queue. Returns nil if there are no
// entries.
func (q *retryingExec) peek() *retry {
	q.Lock()
	defer q.Unlock()

	if len(q.q) > 0 {
		return q.q[0]
	}

	return nil
}

// Loops over the queue, executing retries. Exits if the queue's
// channel is closed. The channel is used to track the earliest
// scheduled retry. Retries are executed in another goroutine to avoid
// blocking the goroutine running this loop.
func (q *retryingExec) handleRetries() {
	// Close on exit to avoid writes to a closed channel.
	defer close(q.execChan)

	var timer *time.Timer

OUTER:
	for {
		// presume empty queue; wait for an entry
		earliestDeadline := <-q.deadlineChan
		if earliestDeadline.IsZero() {
			return
		}

		delay := earliestDeadline.Sub(time.Now())
		timer = time.NewTimer(delay)

		for {
			select {
			case newDeadline := <-q.deadlineChan:
				if newDeadline.IsZero() {
					return
				}
				if newDeadline.Before(earliestDeadline) {
					earliestDeadline = newDeadline

					delay := newDeadline.Sub(time.Now())
					timer.Reset(delay)
				}
			case <-timer.C:
				// issue retries
				for r := q.removeIfPast(); r != nil; r = q.removeIfPast() {
					q.execChan <- r
				}

				if r := q.peek(); r != nil {
					// reset timer for next known deadline
					earliestDeadline = r.deadline

					delay := earliestDeadline.Sub(time.Now())
					timer.Reset(delay)
				} else {
					// empty queue, continue outer loop
					continue OUTER
				}
			}
		}
	}
}

func (q *retryingExec) rescuedCall(f Func, ctxt context.Context) (t Try) {
	defer func() {
		if p := recover(); p != nil {
			stack := make([]byte, 2048)
			runtime.Stack(stack, false)

			if q.log != nil {
				q.log.Printf(
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
				t = NewError(errors.New(fmt.Sprintf("%#v", p)))
			}

		}
	}()

	t = NewTry(f(ctxt))
	return
}

func checkCtxtError(parent context.Context, child context.Context) contextErrorType {
	switch err := parent.Err(); err {
	case context.DeadlineExceeded:
		return globalTimeout

	case nil:
		// check child

	default:
		return cancellation
	}

	switch err := child.Err(); err {
	case context.DeadlineExceeded:
		return attemptTimeout

	case nil:
		return noError

	default:
		return cancellation
	}

}

func (q *retryingExec) mkRetryContext(c context.Context) (context.Context, context.CancelFunc) {
	if q.attemptTimeout > noTimeout {
		return context.WithTimeout(c, q.attemptTimeout)
	}

	return context.WithCancel(c)
}

func (q *retryingExec) handleExec() {
	for {
		r, ok := <-q.execChan
		if !ok {
			return
		}

		r.attempts++

		ctxt, localCancel := q.mkRetryContext(r.ctxt)

		q.stats.TimingDuration(actionDelaytime, time.Now().Sub(r.deadline))

		t := q.rescuedCall(r.f, ctxt)

		ctxtErrType := checkCtxtError(r.ctxt, ctxt)
		localCancel()

		if t.IsError() {
			if ctxtErrType == globalTimeout {
				// global timeout expired
				q.stats.Inc(actionGlobalTimeout, 1)
				t = NewError(
					fmt.Errorf(
						"action exceeded timeout (%s)",
						q.timeout,
					),
				)
			} else if ctxtErrType == cancellation {
				// canceled
				q.stats.Inc(actionCanceled, 1)
				t = NewError(errors.New("action canceled"))
			} else {
				if ctxtErrType == attemptTimeout {
					// retry timeout expired, just count it
					q.stats.Inc(actionRetryTimeout, 1)
					t = NewError(
						fmt.Errorf(
							"action exceeded attempt timeout (%s)",
							q.attemptTimeout,
						),
					)
				} else {
					q.stats.Inc(actionFailure, 1)
				}

				// TODO: check if error is something want actually want to
				// retry (see #1686)
				r.deadline = time.Now().Add(q.delay(r.attempts))

				if limit, ok := r.ctxt.Deadline(); ok && limit.Before(r.deadline) {
					// context will timeout before retry
					q.stats.Inc(actionGlobalTimeout, 1)
					t = NewError(fmt.Errorf(
						"failed action would timeout before next retry: %s",
						t.Error().Error(),
					))
				} else if q.add(r) {
					// retrying
					q.stats.Inc(actionRetry, 1)
					continue
				} else {
					q.stats.Inc(actionRetriesExhausted, 1)
				}
			}
		}

		complete := time.Now()
		q.stats.TimingDuration("action_time", complete.Sub(r.start))

		if r.cb != nil {
			r.cb(t)
			q.stats.TimingDuration("action_callback_time", time.Now().Sub(complete))
		}

		if r.ctxtCancel != nil {
			r.ctxtCancel()
		}
	}
}

func (q *retryingExec) SetStats(s stats.Stats) {
	q.stats = stats.NewAsyncStats(s)
}
