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

package indexof

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestIndexOfStringFound(t *testing.T) {
	s := []string{"a", "b", "c"}
	assert.Equal(t, String(s, "a"), 0)
	assert.Equal(t, String(s, "c"), 2)
}

func TestIndexOfStringNotFound(t *testing.T) {
	s := []string{"a", "b", "c"}
	assert.Equal(t, String(s, "d"), -1)
}

func TestTrimmedString(t *testing.T) {
	s := []string{" a", "b ", " c ", "d"}

	assert.Equal(t, TrimmedString(s, "a"), 0)
	assert.Equal(t, TrimmedString(s, " a"), 0)
	assert.Equal(t, TrimmedString(s, " a "), 0)
	assert.Equal(t, TrimmedString(s, "b"), 1)
	assert.Equal(t, TrimmedString(s, " b"), 1)
	assert.Equal(t, TrimmedString(s, " b "), 1)
	assert.Equal(t, TrimmedString(s, "c"), 2)
	assert.Equal(t, TrimmedString(s, " c"), 2)
	assert.Equal(t, TrimmedString(s, " c "), 2)
	assert.Equal(t, TrimmedString(s, "d"), 3)
	assert.Equal(t, TrimmedString(s, " d"), 3)
	assert.Equal(t, TrimmedString(s, " d "), 3)

	assert.Equal(t, TrimmedString(s, "e"), -1)
	assert.Equal(t, TrimmedString(s, " e"), -1)
	assert.Equal(t, TrimmedString(s, " e "), -1)
}

func TestInt(t *testing.T) {
	is := []int{0, 1, 2, 3, 4, 5}

	for i := range is {
		assert.Equal(t, Int(is, i), i)
	}

	assert.Equal(t, Int(is, 7), -1)
}
