package executor

import (
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/log"
)

func TestWithLogger(t *testing.T) {
	exec := &retryingExec{}
	log := log.NewNoopLogger()

	WithLogger(log)(exec)

	assert.SameInstance(t, exec.log, log)
}

func TestWithRetryDelayFunc(t *testing.T) {
	exec := &retryingExec{}
	d := NewConstantDelayFunc(0)

	WithRetryDelayFunc(d)(exec)

	assert.SameInstance(t, exec.delay, d)
}

func TestWithMaxAttempts(t *testing.T) {
	exec := &retryingExec{}

	WithMaxAttempts(0)(exec)
	assert.Equal(t, exec.maxAttempts, 1)

	WithMaxAttempts(100)(exec)
	assert.Equal(t, exec.maxAttempts, 100)

	WithMaxAttempts(-100)(exec)
	assert.Equal(t, exec.maxAttempts, 1)
}

func TestWithParallelism(t *testing.T) {
	exec := &retryingExec{}

	WithParallelism(0)(exec)
	assert.Equal(t, exec.parallelism, 1)

	WithParallelism(100)(exec)
	assert.Equal(t, exec.parallelism, 100)

	WithParallelism(-100)(exec)
	assert.Equal(t, exec.parallelism, 1)
}

func TestWithMaxQueueDepth(t *testing.T) {
	exec := &retryingExec{}

	WithMaxQueueDepth(0)(exec)
	assert.Equal(t, exec.maxQueueDepth, 1)

	WithMaxQueueDepth(100)(exec)
	assert.Equal(t, exec.maxQueueDepth, 100)

	WithMaxQueueDepth(-100)(exec)
	assert.Equal(t, exec.maxQueueDepth, 1)
}

func TestWithTimeout(t *testing.T) {
	exec := &retryingExec{}

	WithTimeout(0 * time.Second)(exec)
	assert.Equal(t, exec.timeout, 0*time.Second)

	WithTimeout(100 * time.Millisecond)(exec)
	assert.Equal(t, exec.timeout, 100*time.Millisecond)

	WithTimeout(-100 * time.Millisecond)(exec)
	assert.Equal(t, exec.timeout, 0*time.Second)
}
