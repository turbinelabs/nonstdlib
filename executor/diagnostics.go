package executor

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import "time"

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
)

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

type noopDiagnosticsCallback struct{}

func (ndc *noopDiagnosticsCallback) TaskStarted(_ int)                                 {}
func (ndc *noopDiagnosticsCallback) TaskCompleted(_ AttemptResult, _ time.Duration)    {}
func (ndc *noopDiagnosticsCallback) AttemptStarted(_ time.Duration)                    {}
func (ndc *noopDiagnosticsCallback) AttemptCompleted(_ AttemptResult, _ time.Duration) {}
func (ndc *noopDiagnosticsCallback) CallbackDuration(_ time.Duration)                  {}
