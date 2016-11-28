package executor

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"context"

	"github.com/turbinelabs/nonstdlib/stats"
)

// Invoked to execute an action. The given Context should be used to
// make HTTP requests. The function should return as soon as possible
// if the context's Done channel is closed. Must return a nil error if
// the action succeeded. Return an error to try again later.
type Func func(context.Context) (interface{}, error)

// Invoked at most once to return the result of Func.
type CallbackFunc func(Try)

// Invoked at most once each for functions invoked via a single
// ExecMany call. Each invocation includes the index of the function
// in ExecMany's array argument.
type ManyCallbackFunc func(int, Try)

// Executor invokes functions asynchronously with callbacks on
// completion and automatic retries, if configured.
type Executor interface {
	// Invoke the Func, possibly in parallel with other
	// invocations. The function's result is ignored.
	ExecAndForget(Func)

	// Invoke the Func, possibly in parallel with other
	// invocations. Calls back with the result of the call at some
	// future point.
	Exec(Func, CallbackFunc)

	// Invoke the given Funcs, possibly in parallel with other
	// invocations. Calls back with the result of each invocation
	// at some future point.
	ExecMany([]Func, ManyCallbackFunc)

	// Invoke the given Funcs, as in ExecMany. Calls back with a
	// Try containing an []interface{} of the successful results
	// or the first error encountered.
	ExecGathered([]Func, CallbackFunc)

	// Stop executor activity and release related resources. In
	// progress actions will complete their current
	// attempt. Pending actions and retries are dropped and
	// callbacks are not invoked.
	Stop()

	// Submits diagnostic information about executor behavior to
	// the given stats.Stats.
	SetStats(stats.Stats)
}
