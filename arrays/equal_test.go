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
	"fmt"
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

func TestCompareIntSlices(t *testing.T) {
	var z, nz int
	nz = 1

	tcs := []struct {
		left, right []int
		expected    int
	}{
		{
			left:     nil,
			right:    nil,
			expected: 0,
		},
		{
			left:     nil,
			right:    []int{nz},
			expected: -1,
		},
		{
			left:     []int{nz},
			right:    nil,
			expected: 1,
		},
		{
			left:     []int{z, z, nz},
			right:    []int{z, z, z},
			expected: 1,
		},
		{
			left:     []int{z, z, z},
			right:    []int{z, z, nz},
			expected: -1,
		},
		{
			left:     []int{nz, nz},
			right:    []int{z, z, z},
			expected: -1,
		},
		{
			left:     []int{z, z, z},
			right:    []int{nz, nz},
			expected: 1,
		},
		{
			left:     []int{nz, nz, nz},
			right:    []int{nz, nz, nz},
			expected: 0,
		},
	}

	for i, tc := range tcs {
		assert.Group(
			fmt.Sprintf("testCases[%d]: left=[%#v], right=[%#v]", i, tc.left, tc.right),
			t,
			func(g *assert.G) {
				assert.Equal(g, CompareIntSlices(tc.left, tc.right), tc.expected)
			},
		)
	}
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

func TestCompareInt64Slices(t *testing.T) {
	var z, nz int64
	nz = 1

	tcs := []struct {
		left, right []int64
		expected    int
	}{
		{
			left:     nil,
			right:    nil,
			expected: 0,
		},
		{
			left:     nil,
			right:    []int64{nz},
			expected: -1,
		},
		{
			left:     []int64{nz},
			right:    nil,
			expected: 1,
		},
		{
			left:     []int64{z, z, nz},
			right:    []int64{z, z, z},
			expected: 1,
		},
		{
			left:     []int64{z, z, z},
			right:    []int64{z, z, nz},
			expected: -1,
		},
		{
			left:     []int64{nz, nz},
			right:    []int64{z, z, z},
			expected: -1,
		},
		{
			left:     []int64{z, z, z},
			right:    []int64{nz, nz},
			expected: 1,
		},
		{
			left:     []int64{nz, nz, nz},
			right:    []int64{nz, nz, nz},
			expected: 0,
		},
	}

	for i, tc := range tcs {
		assert.Group(
			fmt.Sprintf("testCases[%d]: left=[%#v], right=[%#v]", i, tc.left, tc.right),
			t,
			func(g *assert.G) {
				assert.Equal(g, CompareInt64Slices(tc.left, tc.right), tc.expected)
			},
		)
	}
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

func TestCompareFloat64Slices(t *testing.T) {
	var z, nz float64
	nz = 1

	tcs := []struct {
		left, right []float64
		expected    int
	}{
		{
			left:     nil,
			right:    nil,
			expected: 0,
		},
		{
			left:     nil,
			right:    []float64{nz},
			expected: -1,
		},
		{
			left:     []float64{nz},
			right:    nil,
			expected: 1,
		},
		{
			left:     []float64{z, z, nz},
			right:    []float64{z, z, z},
			expected: 1,
		},
		{
			left:     []float64{z, z, z},
			right:    []float64{z, z, nz},
			expected: -1,
		},
		{
			left:     []float64{nz, nz},
			right:    []float64{z, z, z},
			expected: -1,
		},
		{
			left:     []float64{z, z, z},
			right:    []float64{nz, nz},
			expected: 1,
		},
		{
			left:     []float64{nz, nz, nz},
			right:    []float64{nz, nz, nz},
			expected: 0,
		},
	}

	for i, tc := range tcs {
		assert.Group(
			fmt.Sprintf("testCases[%d]: left=[%#v], right=[%#v]", i, tc.left, tc.right),
			t,
			func(g *assert.G) {
				assert.Equal(g, CompareFloat64Slices(tc.left, tc.right), tc.expected)
			},
		)
	}
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

func TestCompareStringSlices(t *testing.T) {
	var z, nz string
	nz = "X"

	tcs := []struct {
		left, right []string
		expected    int
	}{
		{
			left:     nil,
			right:    nil,
			expected: 0,
		},
		{
			left:     nil,
			right:    []string{nz},
			expected: -1,
		},
		{
			left:     []string{nz},
			right:    nil,
			expected: 1,
		},
		{
			left:     []string{z, z, nz},
			right:    []string{z, z, z},
			expected: 1,
		},
		{
			left:     []string{z, z, z},
			right:    []string{z, z, nz},
			expected: -1,
		},
		{
			left:     []string{nz, nz},
			right:    []string{z, z, z},
			expected: -1,
		},
		{
			left:     []string{z, z, z},
			right:    []string{nz, nz},
			expected: 1,
		},
		{
			left:     []string{nz, nz, nz},
			right:    []string{nz, nz, nz},
			expected: 0,
		},
	}

	for i, tc := range tcs {
		assert.Group(
			fmt.Sprintf("testCases[%d]: left=[%#v], right=[%#v]", i, tc.left, tc.right),
			t,
			func(g *assert.G) {
				assert.Equal(g, CompareStringSlices(tc.left, tc.right), tc.expected)
			},
		)
	}
}
