package math

import (
	"fmt"
	"math"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestRound(t *testing.T) {
	testCases := []struct {
		input    float64
		expected int64
	}{
		{
			input:    0.0,
			expected: 0,
		},
		{
			input:    math.NaN(),
			expected: 0,
		},
		{
			input:    math.Inf(1),
			expected: math.MaxInt64,
		},
		{
			input:    math.Inf(-1),
			expected: math.MinInt64,
		},
		{
			input:    math.Nextafter(float64(math.MaxInt64), math.Inf(1)),
			expected: math.MaxInt64,
		},
		{
			// float64's mantissa is 53 bits (of which 1
			// is implicit), so this rounds to an int64
			// with the high 53 bits set.
			input:    math.Nextafter(float64(math.MaxInt64), 0),
			expected: 0x7FFFFFFFFFFFFC00,
		},
		{
			input:    math.Nextafter(float64(math.MinInt64), math.Inf(-1)),
			expected: math.MinInt64,
		},
		{
			input:    math.Nextafter(float64(math.MinInt64), 0),
			expected: -0x7FFFFFFFFFFFFC00,
		},
		{
			input:    math.Nextafter(0.5, 0.0),
			expected: 0,
		},
		{
			input:    0.5,
			expected: 1,
		},
		{
			input:    math.Nextafter(0.5, 1.0),
			expected: 1,
		},
		{
			input:    math.Nextafter(1.0, 0.0),
			expected: 1,
		},
		{
			input:    math.Nextafter(1.0, 2.0),
			expected: 1,
		},
		{
			input:    math.Nextafter(-0.5, 0),
			expected: 0,
		},
		{
			input:    -0.5,
			expected: 0,
		},
		{
			input:    math.Nextafter(-0.5, -1.0),
			expected: -1,
		},
	}

	for i, tc := range testCases {
		assert.Group(
			fmt.Sprintf("test %d of %d", i+1, len(testCases)),
			t,
			func(g *assert.G) {
				assert.Equal(g, Round(tc.input), tc.expected)
			},
		)
	}
}
