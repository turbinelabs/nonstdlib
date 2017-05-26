package executor

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestAttemptResultString(t *testing.T) {
	for i := 0; i < 5; i++ {
		r := AttemptResult(i)
		assert.NotEqual(t, r.String(), "")
		assert.NotEqual(t, r.String(), "AttemptUnknown")
	}

	assert.Equal(t, AttemptResult(-1).String(), "AttemptUnknown")
	assert.Equal(t, AttemptResult(5).String(), "AttemptUnknown")
}
