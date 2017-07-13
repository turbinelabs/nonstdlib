package executor

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"
	"time"

	tbntime "github.com/turbinelabs/nonstdlib/time"
)

// AttemptResult represents whether or not a task attempt (e.g. an
// individual invocation of a Func passed to Exec) succeeded or
// failed. Failures may be retried depending on the Executor
// configuration.
type AttemptResult int

const (
	// AttemptSuccess indicates the attempt succeeded.
	AttemptSuccess AttemptResult = iota

	// AttemptTimeout indicates the attempt timed out. It may be
	// retried, depending on Executor configuration.
	AttemptTimeout

	// AttemptGlobalTimeout indicates the attempt timed out
	// because the overall time out for the task's execution
	// expired.
	AttemptGlobalTimeout

	// AttemptCancellation indicates the task was canceled
	// (typically because another task within an ExecGathered call
	// failed).
	AttemptCancellation

	// AttemptError indicates that the attempt failed because it
	// returned an error.
	AttemptError

	// Internal use only. Must come last.
	attemptUnknown
)

const numAttemptResults = int(attemptUnknown)

// DiagnosticsCallback provides information about tasks and attempts
// within an Executor. Typically, this interface is used to record
// statistics about the Executor.
type DiagnosticsCallback interface {
	// A task was accepted for execution. The value is the "width"
	// of the task. The width is one for calls to Exec or
	// ExecAndForget, and the number of Funcs passed to ExecMany
	// or ExecGathered.
	TaskStarted(int)

	// A task completed. The duration is the total time taken
	// executing the task, including delays between retries. Each
	// Func passed to ExecMany or ExecGathered will trigger a call
	// to this function.
	TaskCompleted(AttemptResult, time.Duration)

	// An attempt was started. The duration is the delay between
	// time an attempt was scheduled to start and when it actually
	// started. The attempt may time out before the actual Func is
	// invoked.
	AttemptStarted(time.Duration)

	// An attempt completed. The duration is the amount of time
	// spent executing the attempt.
	AttemptCompleted(AttemptResult, time.Duration)

	// The amount of time spent executing a task's callback.
	CallbackDuration(time.Duration)
}

// NewNoopDiagnosticsCallback creates an implementation of
// DiagnosticsCallback that does nothing.
func NewNoopDiagnosticsCallback() DiagnosticsCallback {
	return &noopDiagnosticsCallback{}
}

// Valid returns true if the AttemptResult is a valid value.
func (r AttemptResult) Valid() bool {
	return r >= AttemptSuccess && r <= AttemptError
}

// String returns a string representation of the AttemptResult.
func (r AttemptResult) String() string {
	switch r {
	case AttemptSuccess:
		return "AttemptSuccess"
	case AttemptTimeout:
		return "AttemptTimeout"
	case AttemptGlobalTimeout:
		return "AttemptGlobalTimeout"
	case AttemptCancellation:
		return "AttemptCancellation"
	case AttemptError:
		return "AttemptError"
	default:
		return "AttemptUnknown"
	}
}

func ForEachAttemptResult(f func(AttemptResult)) {
	for r := AttemptSuccess; r < attemptUnknown; r++ {
		f(r)
	}
}

type noopDiagnosticsCallback struct{}

func (ndc *noopDiagnosticsCallback) TaskStarted(_ int)                                 {}
func (ndc *noopDiagnosticsCallback) TaskCompleted(_ AttemptResult, _ time.Duration)    {}
func (ndc *noopDiagnosticsCallback) AttemptStarted(_ time.Duration)                    {}
func (ndc *noopDiagnosticsCallback) AttemptCompleted(_ AttemptResult, _ time.Duration) {}
func (ndc *noopDiagnosticsCallback) CallbackDuration(_ time.Duration)                  {}

// NewLoggingDiagnosticsCallback creates an implementation of
// DiagnosticsCallback that logs diagnostics information periodically.
func NewLoggingDiagnosticsCallback(logger *log.Logger, period time.Duration) DiagnosticsCallback {
	return newLoggingDiagnosticsCallback(logger, period, tbntime.NewSource())
}

func newLoggingDiagnosticsCallback(
	logger *log.Logger,
	period time.Duration,
	source tbntime.Source,
) DiagnosticsCallback {
	ldc := &loggingDiagnosticsCallback{
		logger: logger,
		period: period,
		data:   newDiagnosticsData(),
		time:   source,
		quit:   make(chan struct{}),
	}

	go ldc.logPeriodically()

	return ldc
}

type countedDuration struct {
	count         int64
	totalDuration int64
	maxDuration   int64
}

func (cd *countedDuration) add(d time.Duration) {
	atomic.AddInt64(&cd.count, 1)
	atomic.AddInt64(&cd.totalDuration, int64(d))

	// Loop to handle the case where another goroutine
	// concurrently updates the maximum. Loop ends when d is no
	// longer the maximum or if we successfully update the
	// maximum.
	for {
		currentMax := atomic.LoadInt64(&cd.maxDuration)
		if d <= time.Duration(currentMax) {
			break
		}

		if atomic.CompareAndSwapInt64(&cd.maxDuration, currentMax, int64(d)) {
			break
		}
	}
}

func (cd *countedDuration) format(prefix string) (string, bool) {
	if cd.count == 0 {
		return fmt.Sprintf("%s: 0", prefix), false
	}

	return fmt.Sprintf(
		"%s: %d (avg %s; max %s)",
		prefix,
		cd.count,
		time.Duration(cd.totalDuration/cd.count).String(),
		time.Duration(cd.maxDuration).String(),
	), true
}

type countedDurationsByResult map[AttemptResult]*countedDuration

func newCountedDurationsByResult() countedDurationsByResult {
	cdbr := make(countedDurationsByResult, numAttemptResults+1)

	ForEachAttemptResult(func(r AttemptResult) {
		cdbr[r] = &countedDuration{}
	})

	cdbr[attemptUnknown] = &countedDuration{}

	return cdbr
}

func (cdbr countedDurationsByResult) add(result AttemptResult, d time.Duration) {
	cdbr[result].add(d)
}

func (cdbr countedDurationsByResult) format(prefix string) ([]string, bool) {
	s := make([]string, 0, numAttemptResults+1)
	ForEachAttemptResult(
		func(r AttemptResult) {
			if row, nonZero := cdbr[r].format(prefix + r.String()); nonZero {
				s = append(s, row)
			}
		},
	)
	if row, nonZero := cdbr[attemptUnknown].format(prefix + attemptUnknown.String()); nonZero {
		s = append(s, row)
	}

	return s, len(s) != 0
}

type diagnosticsData struct {
	tasksStarted      int64
	tasksCompleted    countedDurationsByResult
	attemptsStarted   *countedDuration
	attemptsCompleted countedDurationsByResult
	callbacks         *countedDuration
}

type loggingDiagnosticsCallback struct {
	logger *log.Logger
	period time.Duration
	time   tbntime.Source
	lock   sync.RWMutex
	quit   chan struct{}
	data   *diagnosticsData
}

func (ldc *loggingDiagnosticsCallback) Close() error {
	close(ldc.quit)
	return nil
}

func (ldc *loggingDiagnosticsCallback) logPeriodically() {
	timer := ldc.time.NewTimer(ldc.period)
	for {
		select {
		case <-timer.C():
			ldc.log()
			timer.Reset(ldc.period)
		case <-ldc.quit:
			timer.Stop()
			return
		}
	}
}

func newDiagnosticsData() *diagnosticsData {
	return &diagnosticsData{
		tasksCompleted:    newCountedDurationsByResult(),
		attemptsStarted:   &countedDuration{},
		attemptsCompleted: newCountedDurationsByResult(),
		callbacks:         &countedDuration{},
	}
}

func (ldc *loggingDiagnosticsCallback) resetData() *diagnosticsData {
	var (
		newData = newDiagnosticsData()
		oldData *diagnosticsData
	)

	ldc.lock.Lock()
	defer ldc.lock.Unlock()

	oldData, ldc.data = ldc.data, newData

	return oldData
}

func (ldc *loggingDiagnosticsCallback) log() {
	data := ldc.resetData()

	l := ldc.logger
	l.Printf("tasks started: %d", data.tasksStarted)
	if rows, any := data.tasksCompleted.format("tasks completed, "); any {
		for _, s := range rows {
			l.Println(s)
		}
	}
	if attemptsStarted, any := data.attemptsStarted.format("attempts started"); any {
		l.Println(attemptsStarted)
	}
	if rows, any := data.attemptsCompleted.format("attempts completed, "); any {
		for _, s := range rows {
			l.Println(s)
		}
	}
	if callbacks, any := data.callbacks.format("callbacks"); any {
		l.Println(callbacks)
	}
}

func (ldc *loggingDiagnosticsCallback) TaskStarted(n int) {
	ldc.lock.RLock()
	defer ldc.lock.RUnlock()

	atomic.AddInt64(&ldc.data.tasksStarted, int64(n))
}

func (ldc *loggingDiagnosticsCallback) TaskCompleted(result AttemptResult, d time.Duration) {
	if !result.Valid() {
		result = attemptUnknown
	}

	ldc.lock.RLock()
	defer ldc.lock.RUnlock()

	ldc.data.tasksCompleted.add(result, d)
}

func (ldc *loggingDiagnosticsCallback) AttemptStarted(d time.Duration) {
	ldc.lock.RLock()
	defer ldc.lock.RUnlock()

	ldc.data.attemptsStarted.add(d)
}

func (ldc *loggingDiagnosticsCallback) AttemptCompleted(result AttemptResult, d time.Duration) {
	if !result.Valid() {
		result = attemptUnknown
	}

	ldc.lock.RLock()
	defer ldc.lock.RUnlock()

	ldc.data.attemptsCompleted.add(result, d)
}

func (ldc *loggingDiagnosticsCallback) CallbackDuration(d time.Duration) {
	ldc.lock.RLock()
	defer ldc.lock.RUnlock()

	ldc.data.callbacks.add(d)
}

var _ io.Closer = &loggingDiagnosticsCallback{}
