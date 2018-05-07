package regexp

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestGolangIdentifierRegexp(t *testing.T) {
	rx := GolangIdentifierRegexp()
	assert.False(t, rx.MatchString("_"))
	assert.False(t, rx.MatchString(""))
	assert.False(t, rx.MatchString("1a"))
	assert.False(t, rx.MatchString("$a"))
	assert.False(t, rx.MatchString("a@"))
	assert.True(t, rx.MatchString("Aa23B_asd"))
	assert.True(t, rx.MatchString("aA23B_asd2"))
	assert.True(t, rx.MatchString("_a_aA23B_asd2"))
	assert.True(t, rx.MatchString("_1_aA23B_asd2"))
	// Some day my prince will come
	// assert.True(t, rx.MatchString("ƒ"))
}
