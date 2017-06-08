package executor

import (
	"time"

	tbntime "github.com/turbinelabs/nonstdlib/time"
)

type semRequest struct{}

type semaphore chan semRequest

func (s semaphore) Acquire() {
	s <- semRequest{}
}

func (s semaphore) Release() {
	<-s
}

type goroutineExecImpl struct {
	sem semaphore
}

// NewGoroutineExecutor constructs a new Executor. Each task attempt
// is executed in a new goroutine, but only a fixed number (the
// parallelism) are allowed to execute at once. By default, the
// Executor never retries, has parallelism of 1, and a maximum queue
// depth of 10.
func NewGoroutineExecutor(options ...Option) Executor {
	impl := &goroutineExecImpl{}

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

	impl.sem = make(semaphore, e.parallelism)

	if e.log != nil {
		e.log.Printf(
			"goroutine executor: max parallelism %d, max attempts %d, global timeout %s, attempt timeout %s",
			e.parallelism,
			e.maxAttempts,
			e.timeout,
			e.attemptTimeout,
		)
	}

	return e
}

func (g *goroutineExecImpl) stop(c *commonExec) {
	for i := 0; i < c.parallelism; i++ {
		g.sem.Acquire()
	}
}

func (g *goroutineExecImpl) add(c *commonExec, r *retry) {
	go g.run(c, r)
}

func (g *goroutineExecImpl) retry(c *commonExec, delay time.Duration, rx *retry) bool {
	if rx.attempts >= c.maxAttempts {
		return false
	}

	c.time.AfterFunc(delay, func() { g.run(c, rx) })
	return true
}

func (g *goroutineExecImpl) run(c *commonExec, r *retry) {
	g.sem.Acquire()
	defer g.sem.Release()

	c.attempt(r)
}
