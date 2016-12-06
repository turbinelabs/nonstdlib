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
