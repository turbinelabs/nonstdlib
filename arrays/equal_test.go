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

package arrays

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestEqualInt(t *testing.T) {
	var z, nz int
	nz = 1
	x := []int{z, nz, z, nz, z, nz}
	assert.True(t, EqualInt(nil, nil))
	assert.True(t, EqualInt(x[0:0], x[1:1]))
	assert.True(t, EqualInt(x[0:0], nil))
	assert.True(t, EqualInt(nil, x[0:0]))
	assert.False(t, EqualInt(x[0:1], x[0:2]))
	assert.False(t, EqualInt(x[0:2], x[0:1]))
	assert.True(t, EqualInt(x, x))
	assert.True(t, EqualInt(x[0:3], x[2:5]))
	assert.False(t, EqualInt(x[0:3], x[1:4]))
}

func TestEqualInt64(t *testing.T) {
	var z, nz int64
	nz = 1
	x := []int64{z, nz, z, nz, z, nz}
	assert.True(t, EqualInt64(nil, nil))
	assert.True(t, EqualInt64(x[0:0], x[1:1]))
	assert.True(t, EqualInt64(x[0:0], nil))
	assert.True(t, EqualInt64(nil, x[0:0]))
	assert.False(t, EqualInt64(x[0:1], x[0:2]))
	assert.False(t, EqualInt64(x[0:2], x[0:1]))
	assert.True(t, EqualInt64(x, x))
	assert.True(t, EqualInt64(x[0:3], x[2:5]))
	assert.False(t, EqualInt64(x[0:3], x[1:4]))
}

func TestEqualFloat64(t *testing.T) {
	var z, nz float64
	nz = 1
	x := []float64{z, nz, z, nz, z, nz}
	assert.True(t, EqualFloat64(nil, nil))
	assert.True(t, EqualFloat64(x[0:0], x[1:1]))
	assert.True(t, EqualFloat64(x[0:0], nil))
	assert.True(t, EqualFloat64(nil, x[0:0]))
	assert.False(t, EqualFloat64(x[0:1], x[0:2]))
	assert.False(t, EqualFloat64(x[0:2], x[0:1]))
	assert.True(t, EqualFloat64(x, x))
	assert.True(t, EqualFloat64(x[0:3], x[2:5]))
	assert.False(t, EqualFloat64(x[0:3], x[1:4]))
}

func TestEqualString(t *testing.T) {
	var z, nz string
	nz = "X"
	x := []string{z, nz, z, nz, z, nz}
	assert.True(t, EqualString(nil, nil))
	assert.True(t, EqualString(x[0:0], x[1:1]))
	assert.True(t, EqualString(x[0:0], nil))
	assert.True(t, EqualString(nil, x[0:0]))
	assert.False(t, EqualString(x[0:1], x[0:2]))
	assert.False(t, EqualString(x[0:2], x[0:1]))
	assert.True(t, EqualString(x, x))
	assert.True(t, EqualString(x[0:3], x[2:5]))
	assert.False(t, EqualString(x[0:3], x[1:4]))
}
