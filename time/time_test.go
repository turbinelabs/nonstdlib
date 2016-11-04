package time

import (
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func TestTimeEqualNilNil(t *testing.T) {
	assert.True(t, Equal(nil, nil))
}

func TestTimeEqualNilNonNil(t *testing.T) {
	now := time.Now()
	assert.False(t, Equal(&now, nil))
}

func TestTimeEqualNonNilNil(t *testing.T) {
	now := time.Now()
	assert.False(t, Equal(nil, &now))
}

func TestTimeEqualNotEqual(t *testing.T) {
	now := time.Now()
	then := now.Add(-time.Second)
	assert.False(t, Equal(&now, &then))
}

func TestTimeEqualEqual(t *testing.T) {
	now := time.Now()
	assert.True(t, Equal(&now, &now))
}

func TestTimeParseBad(t *testing.T) {
	ts := "2016-01-01 whee 11:30:45.850 UTC"
	tm, err := Parse(ts)
	assert.Nil(t, tm)
	assert.NonNil(t, err)
}

func TestTimeParseNoTZ(t *testing.T) {
	ts := "2016-01-01 11:30:45.850"
	tm, err := Parse(ts)

	l, err := time.LoadLocation("UTC")
	assert.Nil(t, err)

	want := time.Date(2016, 1, 1, 11, 30, 45, 850000000, l)
	assert.Nil(t, err)
	assert.DeepEqual(t, *tm, want)
}

func TestTimeParseWithTZ(t *testing.T) {
	ts := "2016-01-01 11:30:45.850 UTC"
	tm, err := Parse(ts)
	assert.Nil(t, tm)
	assert.NonNil(t, err)
}

func TestTimeFormat(t *testing.T) {
	tm := time.Unix(1451647845, 850000000)
	ts := Format(&tm)
	assert.Equal(t, ts, "2016-01-01 11:30:45.850")
}

func TestTimeFormatNil(t *testing.T) {
	assert.Equal(t, Format(nil), "")
}

func TestTimeFromUnixMilli(t *testing.T) {
	millis := int64(1468970983150)
	ts := FromUnixMilli(millis)
	expected, err := time.Parse(time.RFC3339Nano, "2016-07-19T23:29:43.150000000Z")
	assert.Nil(t, err)

	assert.Equal(t, ts, expected)
}

func TestTimeToUnixMilli(t *testing.T) {
	ts, err := time.Parse(time.RFC3339Nano, "2016-07-19T23:29:43.150000000Z")
	assert.Nil(t, err)

	expected := int64(1468970983150)

	millis := ToUnixMilli(ts)

	assert.Equal(t, millis, expected)
}

func TestTimeTruncUnixMilli(t *testing.T) {
	tsOriginal, err := time.Parse(time.RFC3339Nano, "2016-07-19T23:29:43.150001000Z")
	tsTrunc, err := time.Parse(time.RFC3339Nano, "2016-07-19T23:29:43.150000000Z")
	assert.Nil(t, err)

	assert.Equal(t, TruncUnixMilli(tsOriginal), tsTrunc)
}

func TestTimeTruncUnixMilliZero(t *testing.T) {
	tsOriginal := time.Time{}
	assert.Equal(t, TruncUnixMilli(tsOriginal), tsOriginal)
}

func TestTimeFromUnixMicro(t *testing.T) {
	micros := int64(1468970983123456)
	ts := FromUnixMicro(micros)
	expected, err := time.Parse(time.RFC3339Nano, "2016-07-19T23:29:43.123456000Z")
	assert.Nil(t, err)

	assert.Equal(t, ts, expected)
}

func TestTimeToUnixMicro(t *testing.T) {
	ts, err := time.Parse(time.RFC3339Nano, "2016-07-19T23:29:43.123456000Z")
	assert.Nil(t, err)

	expected := int64(1468970983123456)

	micros := ToUnixMicro(ts)

	assert.Equal(t, micros, expected)
}

func TestTimeTruncUnixMicro(t *testing.T) {
	tsOriginal, err := time.Parse(time.RFC3339Nano, "2016-07-19T23:29:43.150000001Z")
	tsTrunc, err := time.Parse(time.RFC3339Nano, "2016-07-19T23:29:43.150000000Z")
	assert.Nil(t, err)

	assert.Equal(t, TruncUnixMicro(tsOriginal), tsTrunc)
}

func TestTimeTruncUnixMicroZero(t *testing.T) {
	tsOriginal := time.Time{}
	assert.Equal(t, TruncUnixMicro(tsOriginal), tsOriginal)
}

type times []time.Time
type testcase struct {
	name     string
	inputs   []time.Time
	expected time.Time
	fn       func(time.Time, ...time.Time) time.Time
}

func (tc testcase) run(t *assert.G) {
	assert.Equal(
		t,
		tc.fn(tc.inputs[0], tc.inputs[1:]...),
		tc.expected)
}

func TestTimeSelectors(t *testing.T) {
	hr := -1 * time.Hour
	now := time.Now()
	then1 := now.Add(hr)
	then2 := now.Add(hr)
	then3 := now.Add(hr)
	then4 := now.Add(hr)
	then5 := now.Add(hr)
	unset := time.Time{}

	sortedT := times{now, then1, unset, then2, then3, then4, then5}
	revSortedT := times{then5, then4, unset, then3, then2, then1, now}
	mixedT := times{then1, then5, now, unset, then4, then3}

	cases := []testcase{
		{"Trivial Min", times{now}, now, Min},
		{"Trivial Max", times{now}, now, Max},
		{"Simple Min", times{now, then1}, then1, Min},
		{"Simple Min 2", times{then1, now}, then1, Min},
		{"Simple Unset Min", times{unset, now}, now, Min},
		{"Simple Unset Min 2", times{now, unset}, now, Min},
		{"Simple Max", times{now, then1}, now, Max},
		{"Simple Max 2", times{then1, now}, now, Max},
		{"Simple Unset Max", times{unset, now}, now, Max},
		{"Simple Unset Max 2", times{now, unset}, now, Max},
		{"Many Times Min", sortedT, then5, Min},
		{"Many Times Min 2", revSortedT, then5, Min},
		{"Many Times Min 3", mixedT, then5, Min},
		{"Many Times Max", sortedT, now, Max},
		{"Many Times Max 2", revSortedT, now, Max},
		{"Many Times Max 3", mixedT, now, Max},
	}

	// TODO: include unset
	for _, c := range cases {
		assert.Group(c.name, t, func(tg *assert.G) {
			c.run(tg)
		})
	}
}
