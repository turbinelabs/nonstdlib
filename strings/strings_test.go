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

package strings

import (
	"strings"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func ss() Set {
	return Set{
		"has-key":      true,
		"also-has-key": true,
	}
}

func TestSetContains(t *testing.T) {
	assert.True(t, ss().Contains("has-key"))
	assert.False(t, ss().Contains("missing-key"))
}

func TestSetPut(t *testing.T) {
	s := ss()
	s.Put("new-key")
	assert.True(t, s.Contains("new-key"))
}

func TestSetRemove(t *testing.T) {
	s := ss()
	s.Remove("has-key")
	assert.True(t, s.Contains("also-has-key"))
	assert.False(t, s.Contains("has-key"))
}

func TestSetSlice(t *testing.T) {
	assert.HasSameElements(
		t,
		ss().Slice(),
		[]string{"has-key", "also-has-key"},
	)
}

func TestNewSet(t *testing.T) {
	assert.DeepEqual(
		t,
		NewSet(ss().Slice()...),
		ss(),
	)
}

func TestSetEquals(t *testing.T) {
	assert.True(t, ss().Equals(ss()))
}

func TestSetTransform(t *testing.T) {
	s := ss()
	s.Transform(strings.ToUpper)
	assert.DeepEqual(t, s, NewSet("HAS-KEY", "ALSO-HAS-KEY"))
}

func TestSetEqualsFailure(t *testing.T) {
	ss1 := ss()
	ss2 := ss()
	ss2.Put("new-key")
	assert.False(t, ss1.Equals(ss2))
	assert.False(t, ss2.Equals(ss1))
}
