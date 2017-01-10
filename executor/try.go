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

// A Try represents the result of a computation, which may return a value
// or an error. The following represents the possible return values for
// IsReturn, IsError, Get and Error:
//
//    Function | Success | Failure
//    ---------|---------|--------
//    IsReturn | true    | false
//    IsError  | false   | true
//    Get      | result  | panic
//    Error    | panic   | error
type Try interface {
	// if true, the computation produced a return value
	IsReturn() bool

	// if true, the computation resulted in failure
	IsError() bool

	// Get returns the successul result of the computation.
	// All calls to Get should be guarded by IsReturn; if the computation
	// produced an error, calls to Get will panic.
	Get() interface{}

	// Error returns the error that caused the compuation to fail.
	// All calls to Error should be guarded by IsError; if the computation
	// succeeded, calls to Error will panic.
	Error() error
}

// NewTry produces a Try based oen the given interface{} and error. If the
// error is non-nil, a Try is returned for which IsError returns true.
// Otherwise, a Try is returned for which IsReturn returns true.
func NewTry(i interface{}, err error) Try {
	if err != nil {
		return &errorT{err}
	} else {
		return &returnT{i}
	}
}

// NewReturn produces a Try representing a successful computation
func NewReturn(i interface{}) Try {
	return &returnT{i}
}

// NewError produces a Try representing a failed computation
func NewError(err error) Try {
	return &errorT{err}
}

type returnT struct {
	r interface{}
}

func (r *returnT) IsReturn() bool {
	return true
}

func (r *returnT) IsError() bool {
	return false
}

func (r *returnT) Get() interface{} {
	return r.r
}

func (r *returnT) Error() error {
	panic("this Try is not an Error")
}

type errorT struct {
	e error
}

func (e *errorT) IsReturn() bool {
	return false
}

func (e *errorT) IsError() bool {
	return true
}

func (e *errorT) Get() interface{} {
	panic("this Try is not a Return")
}

func (e *errorT) Error() error {
	return e.e
}
