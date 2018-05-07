/*
Copyright 2018 Turbine Labs, Inc.

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

// Package executor provides a mechanism for asynchronous execution of tasks,
// using callbacks to indicate success or failure.
package executor

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE --write_package_comment=false

import "context"

// Func is invoked to execute an action. The given Context should be
// used to make HTTP requests. The function should return as soon as
// possible if the context's Done channel is closed. Must return a nil
// error if the action succeeded. Return an error to try again later.
type Func func(context.Context) (interface{}, error)

// CallbackFunc is invoked at most once to return the result of Func.
type CallbackFunc func(Try)

// ManyCallbackFunc invoked at most once each for functions invoked
// via a single ExecMany call. Each invocation includes the index of
// the function in ExecMany's array argument.
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
	// at some future point. If no Funcs are given, the callback
	// is never invoked.
	ExecMany([]Func, ManyCallbackFunc)

	// Invoke the given Funcs, as in ExecMany. Calls back with a
	// Try containing an []interface{} of the successful results
	// or the first error encountered. If no Funcs are given,
	// the callback is invoked with an empty []interface{}.
	ExecGathered([]Func, CallbackFunc)

	// Stop executor activity and release related resources. In
	// progress actions will complete their current
	// attempt. Pending actions and retries are dropped and
	// callbacks are not invoked.
	Stop()

	// Sets the DiagnosticsCallback for this Executor. Must be
	// called before the first invocation of any Exec
	// function. See also the WithDiagnostics Option to
	// NewRetryingExecutor and NewGoroutineExecutor.
	SetDiagnosticsCallback(DiagnosticsCallback)
}

// NewExecutor constructs a new Executor. The current default
// executor uses goroutines. See NewGoroutineExecutor.
func NewExecutor(options ...Option) Executor { return NewGoroutineExecutor(options...) }
