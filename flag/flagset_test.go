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

package flag

import (
	"flag"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestWrap(t *testing.T) {
	fs := &flag.FlagSet{}
	wrapped := Wrap(fs)
	assert.NonNil(t, wrapped)
	assert.SameInstance(t, wrapped.Unwrap(), fs)
}

func TestScopedUnwrap(t *testing.T) {
	fs := &flag.FlagSet{}
	wrapped := Wrap(fs)

	scoped := wrapped.Scope("a", "desc").Scope("b", "desc").Scope("c", "desc")
	assert.NonNil(t, scoped)
	assert.SameInstance(t, scoped.Unwrap(), fs)
}

func TestTestFlagSet(t *testing.T) {
	tfs := NewTestFlagSet()

	assert.Equal(t, tfs.Unwrap().NFlag(), 0)

	tfs.Int("x", 0, "usage")

	var recoveredPanic interface{}
	func() {
		defer func() { recoveredPanic = recover() }()

		tfs.Parse([]string{"-non-existent-flag"})
	}()

	assert.NonNil(t, recoveredPanic)
}
