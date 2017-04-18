package strings

import (
	"fmt"
	"testing"

	"github.com/turbinelabs/nonstdlib/ptr"
	"github.com/turbinelabs/test/assert"
)

func TestSplit2(t *testing.T) {
	testCases := [][4]string{
		{"a=b", "=", "a", "b"},
		{"a=b=c", "=", "a", "b=c"},
		{"a", "=", "a", ""},
		{"a=", "=", "a", ""},
		{"=a", "=", "", "a"},
		{"bananas", "a", "b", "nanas"},
		{"bananas", "x", "bananas", ""},
	}

	for _, tc := range testCases {
		input := tc[0]
		delim := tc[1]
		expLeft := tc[2]
		expRight := tc[3]

		assert.Group(
			fmt.Sprintf("SplitDelimiter(%q, %q)", input, delim),
			t,
			func(g *assert.G) {
				left, right := Split2(input, delim)
				assert.Equal(g, left, expLeft)
				assert.Equal(g, right, expRight)
			},
		)
	}
}

func testSplit(
	t *testing.T,
	name string,
	testFunc func(string) (string, string),
	testCases [][3]string,
) {
	for _, tc := range testCases {
		input := tc[0]
		expLeft := tc[1]
		expRight := tc[2]

		assert.Group(
			fmt.Sprintf("%s(%q)", name, input),
			t,
			func(g *assert.G) {
				left, right := testFunc(input)
				assert.Equal(g, left, expLeft)
				assert.Equal(g, right, expRight)
			},
		)
	}
}

func TestSplitFirstEqual(t *testing.T) {
	testCases := [][3]string{
		{"a=b", "a", "b"},
		{"a=b=c", "a", "b=c"},
		{"a", "a", ""},
		{"a=", "a", ""},
		{"=a", "", "a"},
		{"bananas", "bananas", ""},
		{"=bananas", "", "bananas"},
		{"a:b", "a:b", ""},
	}

	testSplit(t, "SplitFirstEqual", SplitFirstEqual, testCases)
}

func TestSplitFirstColon(t *testing.T) {
	testCases := [][3]string{
		{"a:b", "a", "b"},
		{"a:b:c", "a", "b:c"},
		{"a", "a", ""},
		{"a:", "a", ""},
		{":a", "", "a"},
		{"bananas", "bananas", ""},
		{":bananas", "", "bananas"},
		{"a=b", "a=b", ""},
	}

	testSplit(t, "SplitFirstColon", SplitFirstColon, testCases)
}

func TestSplitHostPort(t *testing.T) {
	testCases := []struct {
		s                     string
		expectedHost          string
		expectedPort          int
		expectedErrorContains *string
	}{
		{"a:1", "a", 1, nil},
		{"localhost:80", "localhost", 80, nil},
		{"10.0.0.1:8000", "10.0.0.1", 8000, nil},
		{"[::1]:99", "::1", 99, nil},
		{"a:", "", 0, ptr.String(`address a:: missing port`)},
		{":1", "", 0, ptr.String(`address :1: missing host`)},
		{"a", "", 0, ptr.String(`address a: missing port in address`)},
		{"a:b", "", 0, ptr.String(`address a:b: cannot convert port to integer`)},
		{"a:b:b", "", 0, ptr.String(`address a:b:b: too many colons in address`)},
		{"a:-1", "", 0, ptr.String(`address a:-1: port out of range`)},
		{"a:65536", "", 0, ptr.String(`address a:65536: port out of range`)},
	}

	for _, tc := range testCases {
		assert.Group(
			fmt.Sprintf("SplitHostPort(%q)", tc.s),
			t,
			func(g *assert.G) {
				host, port, err := SplitHostPort(tc.s)
				assert.Equal(t, host, tc.expectedHost)
				assert.Equal(t, port, tc.expectedPort)
				if tc.expectedErrorContains != nil {
					assert.ErrorContains(g, err, *tc.expectedErrorContains)
				} else {
					assert.Nil(g, err)
				}
			},
		)
	}
}
