package executor

import (
	"log"
	"time"

	tbntime "github.com/turbinelabs/nonstdlib/time"
)

// Option is used to supply configuration for an Executor
// implementation
type Option func(*commonExec)

// WithLogger sets a Logger for panics recovered while executing
// actions.
func WithLogger(log *log.Logger) Option {
	return func(e *commonExec) {
		e.log = log
	}
}

// WithRetryDelayFunc sets the DelayFunc used when retrying actions.
func WithRetryDelayFunc(d DelayFunc) Option {
	return func(e *commonExec) {
		e.delay = d
	}
}

// WithMaxAttempts sets the absolute maximum number of attempts made
// to complete an action (including the initial attempt). Values less
// than 1 act as if 1 had been passed.
func WithMaxAttempts(maxAttempts int) Option {
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	return func(e *commonExec) {
		e.maxAttempts = maxAttempts
	}
}

// WithParallelism sets the number of goroutines used to execute
// actions. No more than this many actions can be executing at
// once. Values less than 1 act as if 1 has been passed.
func WithParallelism(parallelism int) Option {
	if parallelism < 1 {
		parallelism = 1
	}

	return func(e *commonExec) {
		e.parallelism = parallelism
	}
}

// WithMaxQueueDepth sets the maximum number of actions pending
// immediate execution. If all worker goroutines are processing
// actions, the number of items that can be pending execution (initial
// or retry) before blocking occurs.
func WithMaxQueueDepth(maxQueueDepth int) Option {
	if maxQueueDepth < 1 {
		maxQueueDepth = 1
	}

	return func(e *commonExec) {
		if r, ok := e.impl.(*retryingExecImpl); ok {
			r.maxQueueDepth = maxQueueDepth
		}
	}
}

// WithTimeout sets the timeout for completion of actions. If the
// action has not completed (including retries) within the given
// duration, it is canceled. Timeouts less than or equal to zero are
// treated as "no time out."
func WithTimeout(timeout time.Duration) Option {
	if timeout <= noTimeout {
		timeout = noTimeout
	}

	return func(e *commonExec) {
		e.timeout = timeout
	}
}

// WithAttemptTimeout sets the timeout for completion individual
// attempts of an action. If the attempt has not completed within the
// given duration, it is canceled (and potentially retried). Timeouts
// less than or equal to zero are treated as "no time out."
func WithAttemptTimeout(timeout time.Duration) Option {
	if timeout <= noTimeout {
		timeout = noTimeout
	}

	return func(e *commonExec) {
		e.attemptTimeout = timeout
	}
}

// WithTimeSource sets the tbntime.Source used for obtaining the
// current time. This option should only be used for testing.
func WithTimeSource(src tbntime.Source) Option {
	return func(e *commonExec) {
		e.time = src
	}
}
