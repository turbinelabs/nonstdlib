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
	"errors"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestNewTry(t *testing.T) {
	try := NewTry(123, nil)
	assert.True(t, try.IsReturn())
	assert.False(t, try.IsError())
	assert.Equal(t, try.Get(), 123)
	assert.Panic(t, try.Error)

	err := errors.New("boom")
	try = NewTry(123, err)
	assert.False(t, try.IsReturn())
	assert.True(t, try.IsError())
	assert.Equal(t, try.Error(), err)
	assert.Panic(t, try.Get)
}

func TestNewReturn(t *testing.T) {
	try := NewTry(123, nil)
	ret := NewReturn(123)
	assert.DeepEqual(t, ret, try)
}

func TestNewError(t *testing.T) {
	e := errors.New("boom")
	try := NewTry(nil, e)
	err := NewError(e)
	assert.DeepEqual(t, err, try)
}
