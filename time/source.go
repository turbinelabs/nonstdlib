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
	"time"
)

var (
	defaultTimeSourceInstance = &defaultTimeSource{}
)

// Source is a source of time.Time values.
type Source interface {
	// Returns the current time, as with time.Time.Now().
	Now() time.Time
}

// ControlledSource is a source of time.Time values that returns a
// fixed time unless modified with the Set or Advance
// methods. ControlledSource should be used for testing only.
type ControlledSource interface {
	Source

	Set(time.Time)
	Advance(time.Duration)
}

// NewSource creates a new Source for normal use. The Source's Now
// method returns the current time.
func NewSource() Source {
	return defaultTimeSourceInstance
}

// WithTimeAt creates a new ControlledSource with the given time and
// passes it to the given function for testing.
func WithTimeAt(t time.Time, f func(ControlledSource)) {
	s := &controlledTimeSource{now: t}
	f(s)
}

// WithCurrentTimeFrozen Creates a new ControlledSource with the
// current time and passes it to the given function for testing.
func WithCurrentTimeFrozen(f func(ControlledSource)) {
	WithTimeAt(time.Now(), f)
}

type defaultTimeSource struct{}

func (s *defaultTimeSource) Now() time.Time {
	return time.Now()
}

type controlledTimeSource struct {
	now time.Time
}

func (s *controlledTimeSource) Now() time.Time {
	return s.now
}

func (s *controlledTimeSource) Set(t time.Time) {
	s.now = t
}

func (s *controlledTimeSource) Advance(delta time.Duration) {
	s.now = s.now.Add(delta)
}

type incrementingTimeSource struct {
	*controlledTimeSource
	inc time.Duration
}

func (i *incrementingTimeSource) Now() time.Time {
	t := i.now
	i.now = i.now.Add(i.inc)
	return t
}

// NewIncrementingControlledSource returns a new ControlledSource that
// increments the controlled time by some delta every time Now() is called.
func NewIncrementingControlledSource(at time.Time, delta time.Duration) ControlledSource {
	return &incrementingTimeSource{&controlledTimeSource{at}, delta}
}
