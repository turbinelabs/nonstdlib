package executor

import (
	"math"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func TestNewExponentialDelayFunc(t *testing.T) {
	delayFunc := NewExponentialDelayFunc(100*time.Millisecond, 1*time.Second)

	assert.Equal(t, delayFunc(-1), 100*time.Millisecond)
	assert.Equal(t, delayFunc(0), 100*time.Millisecond)
	assert.Equal(t, delayFunc(1), 100*time.Millisecond)
	assert.Equal(t, delayFunc(2), 200*time.Millisecond)
	assert.Equal(t, delayFunc(3), 400*time.Millisecond)
	assert.Equal(t, delayFunc(4), 800*time.Millisecond)
	assert.Equal(t, delayFunc(5), 1*time.Second)
	assert.Equal(t, delayFunc(6), 1*time.Second)
	assert.Equal(t, delayFunc(100), 1*time.Second)

	delayFunc = NewExponentialDelayFunc(1*time.Nanosecond, time.Duration(math.MaxInt64))
	assert.Equal(t, delayFunc(1), 1*time.Nanosecond)
	assert.Equal(t, delayFunc(64), time.Duration(math.MaxInt64))
	assert.Equal(t, delayFunc(65), time.Duration(math.MaxInt64))
}

func TestNewExponentialDelayFuncBadMaxDelay(t *testing.T) {
	delayFunc := NewExponentialDelayFunc(100*time.Millisecond, 10*time.Millisecond)
	assert.Equal(t, delayFunc(1), 100*time.Millisecond)
	assert.Equal(t, delayFunc(2), 100*time.Millisecond)
}

func TestNewExponentialDelayFuncInvalidDelayIsZero(t *testing.T) {
	delayFunc := NewExponentialDelayFunc(-1*time.Second, 1*time.Hour)
	assert.Equal(t, delayFunc(0), time.Duration(0))

	delayFunc = NewExponentialDelayFunc(0*time.Second, 1*time.Hour)
	assert.Equal(t, delayFunc(0), time.Duration(0))
}

func TestNewConstantDelayFunc(t *testing.T) {
	delay := 100 * time.Millisecond
	delayFunc := NewConstantDelayFunc(delay)
	for i := 0; i < 100; i++ {
		assert.Equal(t, delayFunc(i), delay)
	}
}

func TestNewConstantDelayFuncInvalidDelayIsZero(t *testing.T) {
	delayFunc := NewConstantDelayFunc(-1 * time.Second)
	assert.Equal(t, delayFunc(0), time.Duration(0))

	delayFunc = NewConstantDelayFunc(0 * time.Second)
	assert.Equal(t, delayFunc(0), time.Duration(0))
}
