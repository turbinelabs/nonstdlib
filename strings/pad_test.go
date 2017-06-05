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
