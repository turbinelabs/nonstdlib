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

package strings

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestPadLeft(t *testing.T) {
	in := "foo\nbar\nbaz\n"
	assert.Equal(t, PadLeft(in, 3), "   foo\n   bar\n   baz\n   ")
}

func TestPadLeftWith(t *testing.T) {
	in := "foo\nbar"
	assert.Equal(t, PadLeftWith(in, 2, "-*"), "-*-*foo\n-*-*bar")
}

func TestPadLeftEmpty(t *testing.T) {
	assert.Equal(t, PadLeft("", 3), "   ")
}
