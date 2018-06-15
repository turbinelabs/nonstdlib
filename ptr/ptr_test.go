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

package ptr

import (
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

var testCasesStringSlice = [][]string{
	{"a", "b", "c", "d", "e"},
	{"a", "b", "", "", "e"},
}

func TestStringSlice(t *testing.T) {
	for _, in := range testCasesStringSlice {
		if in == nil {
			continue
		}
		out := StringSlice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.DeepEqual(t, in[i], *(out[i]))
		}

		out2 := StringValueSlice(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

var testCasesStringValueSlice = [][]*string{
	{String("a"), String("b"), nil, String("c")},
}

func TestStringValueSlice(t *testing.T) {
	for _, in := range testCasesStringValueSlice {
		if in == nil {
			continue
		}
		out := StringValueSlice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			if in[i] == nil {
				assert.Equal(t, "", out[i])
			} else {
				assert.Equal(t, *(in[i]), out[i])
			}
		}

		out2 := StringSlice(out)
		assert.Equal(t, len(out2), len(in))
		for i := range out2 {
			if in[i] == nil {
				assert.Equal(t, *(out2[i]), "")
			} else {
				assert.Equal(t, *in[i], *out2[i])
			}
		}
	}
}

var testCasesStringMap = []map[string]string{
	{"a": "1", "b": "2", "c": "3"},
}

func TestStringMap(t *testing.T) {
	for _, in := range testCasesStringMap {
		if in == nil {
			continue
		}
		out := StringMap(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.Equal(t, in[i], *(out[i]))
		}

		out2 := StringValueMap(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

var testCasesBoolSlice = [][]bool{
	{true, true, false, false},
}

func TestBoolSlice(t *testing.T) {
	for _, in := range testCasesBoolSlice {
		if in == nil {
			continue
		}
		out := BoolSlice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.Equal(t, in[i], *(out[i]))
		}

		out2 := BoolValueSlice(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

func TestBoolValueSlice(t *testing.T) {
	tr := true
	f := false
	testCasesBoolValueSlice := [][]*bool{
		{},
		{&tr, &f},
	}

	for _, in := range testCasesBoolValueSlice {
		if in == nil {
			continue
		}
		out := BoolValueSlice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			if in[i] == nil {
				assert.False(t, out[i])
			} else {
				assert.Equal(t, *(in[i]), out[i])
			}
		}

		out2 := BoolSlice(out)
		assert.Equal(t, len(out2), len(in))
		for i := range out2 {
			if in[i] == nil {
				assert.False(t, *(out2[i]))
			} else {
				assert.DeepEqual(t, in[i], out2[i])
			}
		}
	}
}

var testCasesBoolMap = []map[string]bool{
	{"a": true, "b": false, "c": true},
}

func TestBoolMap(t *testing.T) {
	for _, in := range testCasesBoolMap {
		if in == nil {
			continue
		}
		out := BoolMap(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.Equal(t, in[i], *(out[i]))
		}

		out2 := BoolValueMap(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

var testCasesIntSlice = [][]int{
	{1, 2, 3, 4},
}

func TestIntSlice(t *testing.T) {
	for _, in := range testCasesIntSlice {
		if in == nil {
			continue
		}
		out := IntSlice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.Equal(t, in[i], *(out[i]))
		}

		out2 := IntValueSlice(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

func TestIntValueSlice(t *testing.T) {
	i1 := 1
	i2 := 0
	i3 := 3
	testCasesIntValueSlice := [][]*int{
		{},
		{&i1, &i2, &i3},
	}

	for _, in := range testCasesIntValueSlice {
		if in == nil {
			continue
		}
		out := IntValueSlice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			if in[i] == nil {
				assert.Equal(t, out[i], 0)
			} else {
				assert.Equal(t, *(in[i]), out[i])
			}
		}

		out2 := IntSlice(out)
		assert.Equal(t, len(out2), len(in))
		for i := range out2 {
			if in[i] == nil {
				assert.Equal(t, *(out2[i]), 0)
			} else {
				assert.DeepEqual(t, in[i], out2[i])
			}
		}
	}
}

var testCasesIntMap = []map[string]int{
	{"a": 3, "b": 2, "c": 1},
}

func TestIntMap(t *testing.T) {
	for _, in := range testCasesIntMap {
		if in == nil {
			continue
		}
		out := IntMap(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.Equal(t, in[i], *(out[i]))
		}

		out2 := IntValueMap(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

var testCasesInt64Slice = [][]int64{
	{1, 2, 3, 4},
}

func TestInt64Slice(t *testing.T) {
	for _, in := range testCasesInt64Slice {
		if in == nil {
			continue
		}
		out := Int64Slice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.Equal(t, in[i], *(out[i]))
		}

		out2 := Int64ValueSlice(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

func TestInt64ValueSlice(t *testing.T) {
	a := int64(0)
	b := int64(2)
	c := int64(3)
	testCasesInt64ValueSlice := [][]*int64{
		nil,
		{},
		{&a, &b, &c},
	}

	for _, in := range testCasesInt64ValueSlice {
		if in == nil {
			continue
		}
		out := Int64ValueSlice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			if in[i] == nil {
				assert.Equal(t, out[i], 0)
			} else {
				assert.Equal(t, *(in[i]), out[i])
			}
		}

		out2 := Int64Slice(out)
		assert.Equal(t, len(out2), len(in))
		for i := range out2 {
			if in[i] == nil {
				assert.Equal(t, *(out2[i]), 0)
			} else {
				assert.DeepEqual(t, in[i], out2[i])
			}
		}
	}
}

var testCasesInt64Map = []map[string]int64{
	{"a": 3, "b": 2, "c": 1},
}

func TestInt64Map(t *testing.T) {
	for _, in := range testCasesInt64Map {
		if in == nil {
			continue
		}
		out := Int64Map(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.Equal(t, in[i], *(out[i]))
		}

		out2 := Int64ValueMap(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

var testCasesFloat64Slice = [][]float64{
	{1, 2, 3, 4},
}

func TestFloat64Slice(t *testing.T) {
	for _, in := range testCasesFloat64Slice {
		if in == nil {
			continue
		}
		out := Float64Slice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.Equal(t, in[i], *(out[i]))
		}

		out2 := Float64ValueSlice(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

func TestFloat64ValueSlice(t *testing.T) {
	a := float64(0)
	b := float64(2)
	c := float64(3)
	var testCasesFloat64ValueSlice = [][]*float64{
		nil,
		{},
		{&a, &b, &c},
	}

	for _, in := range testCasesFloat64ValueSlice {
		if in == nil {
			continue
		}
		out := Float64ValueSlice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			if in[i] == nil {
				assert.Equal(t, out[i], 0)
			} else {
				assert.Equal(t, *(in[i]), out[i])
			}
		}

		out2 := Float64Slice(out)
		assert.Equal(t, len(out2), len(in))
		for i := range out2 {
			if in[i] == nil {
				assert.Equal(t, *(out2[i]), 0)
			} else {
				assert.Equal(t, *in[i], *out2[i])
			}
		}
	}
}

var testCasesFloat64Map = []map[string]float64{
	{"a": 3, "b": 2, "c": 1},
}

func TestFloat64Map(t *testing.T) {
	for _, in := range testCasesFloat64Map {
		if in == nil {
			continue
		}
		out := Float64Map(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.Equal(t, in[i], *(out[i]))
		}

		out2 := Float64ValueMap(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

var testCasesTimeSlice = [][]time.Time{
	{time.Now(), time.Now().AddDate(100, 0, 0)},
}

func TestTimeSlice(t *testing.T) {
	for _, in := range testCasesTimeSlice {
		if in == nil {
			continue
		}
		out := TimeSlice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.Equal(t, in[i], *(out[i]))
		}

		out2 := TimeValueSlice(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

func TestTimeValueSlice(t *testing.T) {
	t0 := time.Time{}
	t1 := time.Now()
	testCasesTimeValueSlice := [][]*time.Time{
		{},
		{&t0, &t1},
	}

	for _, in := range testCasesTimeValueSlice {
		if in == nil {
			continue
		}
		out := TimeValueSlice(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			if in[i] == nil {
				assert.Equal(t, out[i], time.Time{})
			} else {
				assert.Equal(t, *(in[i]), out[i])
			}
		}

		out2 := TimeSlice(out)
		assert.Equal(t, len(out2), len(in))
		for i := range out2 {
			if in[i] == nil {
				assert.Equal(t, *(out2[i]), time.Time{})
			} else {
				assert.Equal(t, *in[i], *out2[i])
			}
		}
	}
}

var testCasesTimeMap = []map[string]time.Time{
	{"a": time.Now().AddDate(-100, 0, 0), "b": time.Now()},
}

func TestTimeMap(t *testing.T) {
	for _, in := range testCasesTimeMap {
		if in == nil {
			continue
		}
		out := TimeMap(in)
		assert.Equal(t, len(out), len(in))
		for i := range out {
			assert.Equal(t, in[i], *(out[i]))
		}

		out2 := TimeValueMap(out)
		assert.Equal(t, len(out2), len(in))
		assert.DeepEqual(t, in, out2)
	}
}

func TestBoolEqualWithNils(t *testing.T) {
	assert.True(t, BoolEqual(nil, nil))
	assert.False(t, BoolEqual(Bool(true), nil))
	assert.False(t, BoolEqual(Bool(false), nil))
}

func TestStringEqualWithNils(t *testing.T) {
	assert.True(t, StringEqual(nil, nil))
	assert.False(t, StringEqual(String("hello"), nil))
	assert.False(t, StringEqual(String(""), nil))
}

func TestIntEqualWithNils(t *testing.T) {
	assert.True(t, IntEqual(nil, nil))
	assert.False(t, IntEqual(Int(0), nil))
	assert.False(t, IntEqual(Int(100), nil))
}

func TestBoolEqualWhenEqual(t *testing.T) {
	assert.True(t, BoolEqual(Bool(true), Bool(true)))
	assert.True(t, BoolEqual(Bool(false), Bool(false)))
	bPtr := Bool(true)
	assert.True(t, BoolEqual(bPtr, bPtr))
}

func TestIntEqualWhenEqual(t *testing.T) {
	iPtr := Int(500)
	assert.True(t, IntEqual(Int(100), Int(100)))
	assert.True(t, IntEqual(iPtr, iPtr))
}

func TestStringEqualWhenEqual(t *testing.T) {
	sPtr := String("harro")
	assert.True(t, StringEqual(String("blerp"), String("blerp")))
	assert.True(t, StringEqual(sPtr, sPtr))
}

func TestBoolEqualWhenNotEqual(t *testing.T) {
	assert.False(t, BoolEqual(Bool(true), Bool(false)))
	assert.False(t, BoolEqual(Bool(true), nil))
	assert.False(t, BoolEqual(Bool(false), nil))
}

func TestIntEqualWhenNotEqual(t *testing.T) {
	assert.False(t, IntEqual(Int(5), Int(10)))
	assert.False(t, IntEqual(Int(5), nil))
}

func TestStringEqualWhenNotEqual(t *testing.T) {
	assert.False(t, StringEqual(String("not"), String("equal")))
	assert.False(t, StringEqual(String("nope"), nil))
}

func TestIntValueOkWithNilPtr(t *testing.T) {
	i, b := IntValueOk(nil)
	assert.Equal(t, i, 0)
	assert.False(t, b)
}

func TestIntValueOkWithNonNilPtr(t *testing.T) {
	i, b := IntValueOk(Int(10))
	assert.Equal(t, i, 10)
	assert.True(t, b)
}
