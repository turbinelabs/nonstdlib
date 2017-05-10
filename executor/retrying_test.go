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
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func TestRetryingExecHeapInterface(t *testing.T) {
	start := time.Now()
	items := []*retry{
		{nextAttempt: start, attempts: 1},
		{nextAttempt: start, attempts: 1},
		{nextAttempt: start, attempts: 1},
		{nextAttempt: start, attempts: 1},
		{nextAttempt: start, attempts: 1},
	}

	q := &retryingExecImpl{
		nextAttemptChan: make(chan time.Time, 10),
		q:               make([]*retry, 0, 10),
	}
	defer close(q.nextAttemptChan)

	e := &commonExec{
		delay:       func(_ int) time.Duration { return 0 * time.Second },
		maxAttempts: 2,
		impl:        q,
	}

	heap.Init(q)

	assert.Nil(t, q.peek())

	for _, item := range items {
		assert.True(t, q.retry(e, 0, item))
	}

	peeked := q.peek()

	r := q.removeIfPast(time.Now())
	assert.NonNil(t, r)
	assert.SameInstance(t, r, peeked)
	assert.Equal(t, r.nextAttempt, start)

	r = q.removeIfPast(time.Now())
	assert.NonNil(t, r)
	assert.Equal(t, r.nextAttempt, start)

	r = q.removeIfPast(time.Now())
	assert.NonNil(t, r)
	assert.Equal(t, r.nextAttempt, start)

	r = q.removeIfPast(time.Now())
	assert.NonNil(t, r)
	assert.Equal(t, r.nextAttempt, start)

	r = q.removeIfPast(time.Now())
	assert.NonNil(t, r)
	assert.Equal(t, r.nextAttempt, start)

	r = q.removeIfPast(time.Now())
	assert.Nil(t, r)

	for _, r := range items {
		nextAttempt := <-q.nextAttemptChan
		assert.Equal(t, nextAttempt, r.nextAttempt)
	}

	assert.False(t, q.retry(e, 0, &retry{nextAttempt: start, attempts: 2}))
	assert.Nil(t, q.peek())
}

func TestRetryingExecHandleRetriesWithNoCallback(t *testing.T) {
	testRetriesWithNoCallback(t, NewRetryingExecutor)
}

func TestRetryingExecHandleRetriesEarlierNextRetry(t *testing.T) {
	testEarlierNextRetry(t, NewRetryingExecutor)
}

func TestRetryingExecInvokesCallback(t *testing.T) {
	testExecInvokesCallback(t, NewRetryingExecutor)
}

func TestRetryingExecExecutesInParallel(t *testing.T) {
	testExecExecutesInParallel(t, NewRetryingExecutor)
}

func TestRetryingExecPanicsBecomeErrors(t *testing.T) {
	testExecPanicsBecomeErrors(t, NewRetryingExecutor)
}

func TestRetryingExecStopsWithInFlightRetries(t *testing.T) {
	testExecStopsWithInFlightRetries(t, NewRetryingExecutor)
}
