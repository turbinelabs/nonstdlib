/*
Copyright 2017 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

func TestIncrementingControlledSource(t *testing.T) {
	before := time.Now()
	delta := 5 * time.Second

	s := NewIncrementingControlledSource(before, delta)
	assert.Equal(t, before, s.Now())
	assert.Equal(t, before.Add(delta), s.Now())
}

func TestIncrementingControlledSourceAdvance(t *testing.T) {
	before := time.Now()
	delta := 5 * time.Second

	s := NewIncrementingControlledSource(before, delta)
	assert.Equal(t, before, s.Now())
	s.Advance(time.Second)

	assert.Equal(t, before.Add(delta+time.Second), s.Now())
}

func TestIncrementingControlledSourceSet(t *testing.T) {
	before := time.Now()
	delta := 5 * time.Second

	s := NewIncrementingControlledSource(before, delta)
	assert.Equal(t, before, s.Now())
	s.Set(before.Add(-1 * time.Hour))

	assert.Equal(t, before.Add(-1*time.Hour), s.Now())
	assert.Equal(t, before.Add(-1*time.Hour).Add(delta), s.Now())
}
