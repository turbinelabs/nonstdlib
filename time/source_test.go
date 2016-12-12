package time

import (
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func TestDefaultSource(t *testing.T) {
	source := NewSource()

	for i := 0; i < 10; i++ {
		before := time.Now()
		now := source.Now()
		after := time.Now()

		assert.True(t, before.Before(now) || before.Equal(now))
		assert.True(t, after.After(now) || after.Equal(now))
	}
}

func TestControlledSource(t *testing.T) {
	original := time.Now()
	source := &controlledTimeSource{now: original}

	assert.DeepEqual(t, source.Now(), original)
	assert.DeepEqual(t, source.Now(), original)

	source.Advance(5 * time.Minute)
	assert.DeepEqual(t, source.Now(), original.Add(5*time.Minute))

	source.Set(original)
	assert.DeepEqual(t, source.Now(), original)
}

func TestWithTimeAt(t *testing.T) {
	original := time.Now()

	called := false
	WithTimeAt(original, func(ts ControlledSource) {
		called = true

		assert.Equal(t, ts.Now(), original)
	})

	assert.True(t, called)
}

func TestWithCurrentTimeFrozen(t *testing.T) {
	before := time.Now()

	called := false
	frozenTime := time.Time{}
	WithCurrentTimeFrozen(func(ts ControlledSource) {
		called = true
		frozenTime = ts.Now()
	})
	after := time.Now()

	assert.True(t, called)
	assert.True(t, before.Before(frozenTime) || before.Equal(frozenTime))
	assert.True(t, after.After(frozenTime) || after.Equal(frozenTime))
}
