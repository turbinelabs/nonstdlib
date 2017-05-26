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
	"sync"
	"time"

	tbntime "github.com/turbinelabs/nonstdlib/time"
)

const (
	defaultMaxAttempts   = 1
	defaultMaxQueueDepth = 10
	defaultParallelism   = 1

	execManyWidth     = "exec_many_width"
	execGatheredWidth = "exec_gathered_width"

	actionDelaytime = "action_delaytime"
	actionAttempts  = "action_attempts"

	actionFailure          = "action_failure.error"
	actionGlobalTimeout    = "action_failure.timeout"
	actionRetryTimeout     = "action_failure.attempt_timeout"
	actionCanceled         = "action_failure.canceled"
	actionRetry            = "action_retry"
	actionRetriesExhausted = "action_retries_exhausted"
)

var (
	defaultDelayFunc = NewConstantDelayFunc(1 * time.Second)

	exitSignal = time.Time{}

	gatherError = NewError(errors.New("underflowed results for ExecGathered"))
)

// Retry encapsulates a nextAttempt (when to retry the action), the
// data necessary to perform the retry, and how many attempts have
// already been made.
type retry struct {
	f           Func
	cb          CallbackFunc
	start       time.Time
	nextAttempt time.Time
	ctxt        context.Context
	ctxtCancel  context.CancelFunc
	attempts    int
}

type retryingExecImpl struct {
	sync.Mutex

	nextAttemptChan chan time.Time
	execChan        chan *retry
	q               []*retry

	maxQueueDepth int
}

var _ execImpl = &retryingExecImpl{}

// NewRetryingExecutor constructs a new Executor. Tasks are run by a
// number of long-lived goroutines (equal to the parallelism of the
// Executor). Tasks are scheduled by an additional long-lived
// goroutine. By default, the Executor never retries, has parallelism
// of 1, and a maximum queue depth of 10.
func NewRetryingExecutor(options ...Option) Executor {
	impl := &retryingExecImpl{
		nextAttemptChan: make(chan time.Time, 2),
		q:               make([]*retry, 0, 10),
		maxQueueDepth:   defaultMaxQueueDepth,
	}

	e := &commonExec{
		time:           tbntime.NewSource(),
		parallelism:    defaultParallelism,
		maxAttempts:    defaultMaxAttempts,
		delay:          defaultDelayFunc,
		timeout:        noTimeout,
		attemptTimeout: noTimeout,
		diag:           NewNoopDiagnosticsCallback(),
		impl:           impl,
	}

	for _, apply := range options {
		apply(e)
	}

	impl.execChan = make(chan *retry, impl.maxQueueDepth)

	heap.Init(impl)
	go impl.handleRetries(e.time)

	for i := 0; i < e.parallelism; i++ {
		go impl.handleExec(e.attempt)
	}

	if e.log != nil {
		e.log.Printf(
			"retrying executor: %d workers, max queue %d, max attempts %d, global timeout %s, attempt timeout %s",
			e.parallelism,
			impl.maxQueueDepth,
			e.maxAttempts,
			e.timeout,
			e.attemptTimeout,
		)
	}

	return e
}

// Implements heap.Interface
func (q *retryingExecImpl) Len() int      { return len(q.q) }
func (q *retryingExecImpl) Swap(i, j int) { q.q[i], q.q[j] = q.q[j], q.q[i] }
func (q *retryingExecImpl) Less(i, j int) bool {
	return q.q[i].nextAttempt.Before(q.q[j].nextAttempt)
}

// Implements heap.Interface
func (q *retryingExecImpl) Push(x interface{}) {
	item := x.(*retry)
	q.q = append(q.q, item)
}

// Implements heap.Interface
func (q *retryingExecImpl) Pop() interface{} {
	old := q.q
	n := len(old)
	item := old[n-1]
	q.q = old[0 : n-1]
	return item
}

var _ heap.Interface = &retryingExecImpl{}

func (q *retryingExecImpl) stop(_ *commonExec) {
	// handleRetries() closes q.execChan as it returns to avoid
	// writes to execChan after close.
	q.nextAttemptChan <- exitSignal
}

func (q *retryingExecImpl) add(c *commonExec, r *retry) {
	q.retry(c, 0, r)
}

// Adds an entry into the queue if it has not exhausted its
// attempts. Enqueues the nextAttempt on the queue's channel.
func (q *retryingExecImpl) retry(c *commonExec, _ time.Duration, r *retry) bool {
	if r.attempts >= c.maxAttempts {
		return false
	}

	if r.nextAttempt.IsZero() {
		return false
	}

	q.Lock()
	defer func() {
		q.Unlock()

		q.nextAttemptChan <- r.nextAttempt
	}()

	heap.Push(q, r)

	return true
}

// Removes the first entry of the queue if it's nextAttempt has
// expired. If there are no items or the head of the queue is not
// expired, returns nil.
func (q *retryingExecImpl) removeIfPast(now time.Time) *retry {
	q.Lock()
	defer q.Unlock()

	if len(q.q) > 0 && !q.q[0].nextAttempt.After(now) {
		r := heap.Pop(q)
		return r.(*retry)
	}

	return nil
}

// Peeks at the first entry in the queue. Returns nil if there are no
// entries.
func (q *retryingExecImpl) peek() *retry {
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
func (q *retryingExecImpl) handleRetries(source tbntime.Source) {
	// Close on exit to avoid writes to a closed channel.
	defer close(q.execChan)

	var timer tbntime.Timer

OUTER:
	for {
		// presume empty queue; wait for an entry
		earliestNextAttempt := <-q.nextAttemptChan
		if earliestNextAttempt.IsZero() {
			return
		}

		delay := earliestNextAttempt.Sub(source.Now())
		timer = source.NewTimer(delay)

		for {
			select {
			case newNextAttempt := <-q.nextAttemptChan:
				if newNextAttempt.IsZero() {
					return
				}
				if newNextAttempt.Before(earliestNextAttempt) {
					earliestNextAttempt = newNextAttempt

					delay := newNextAttempt.Sub(source.Now())
					timer.Reset(delay)
				}
			case <-timer.C():
				// issue retries
				for {
					r := q.removeIfPast(source.Now())
					if r == nil {
						break
					}
					q.execChan <- r
				}

				if r := q.peek(); r != nil {
					// reset timer for next known nextAttempt
					earliestNextAttempt = r.nextAttempt

					delay := earliestNextAttempt.Sub(source.Now())
					timer.Reset(delay)
				} else {
					// empty queue, continue outer loop
					continue OUTER
				}
			}
		}
	}
}

func (q *retryingExecImpl) handleExec(attempt func(*retry)) {
	for {
		r, ok := <-q.execChan
		if !ok {
			return
		}

		attempt(r)
	}
}
