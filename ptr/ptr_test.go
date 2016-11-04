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
