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

type Try interface {
	IsReturn() bool
	IsError() bool

	Get() interface{}
	Error() error
}

func NewTry(i interface{}, err error) Try {
	if err != nil {
		return &errorT{err}
	} else {
		return &returnT{i}
	}
}

func NewReturn(i interface{}) Try {
	return &returnT{i}
}

func NewError(err error) Try {
	return &errorT{err}
}

type Return interface {
	Try
}

type Error interface {
	Try
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
