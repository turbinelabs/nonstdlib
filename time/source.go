/*
Copyright 2018 Turbine Labs, Inc.

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
	"context"
	"time"
)

var (
	defaultTimeSourceInstance = &defaultTimeSource{}
)

// Source is a source of time.Time values and Timer instances.
type Source interface {
	// Returns the current time, as with time.Time.Now().
	Now() time.Time

	// NewTimer creates a new Timer that will send the current
	// time on its channel after the given duration.
	NewTimer(time.Duration) Timer

	// AfterFunc creates a new Timer that will invoke the given
	// function in its own goroutine after the duration has
	// elapsed.
	AfterFunc(time.Duration, func()) Timer

	// NewContextWithDeadline creates a new context.Context and
	// associated context.CancelFunc with the given deadline.
	NewContextWithDeadline(
		parent context.Context,
		deadline time.Time,
	) (context.Context, context.CancelFunc)

	// NewContextWithTimeout creates a new context.Context and
	// associated context.CancelFunc with the given timeout.
	NewContextWithTimeout(
		parent context.Context,
		timeout time.Duration,
	) (context.Context, context.CancelFunc)
}

// NewSource creates a new Source for normal use. The Source's Now
// method returns the current time.
func NewSource() Source {
	return defaultTimeSourceInstance
}

type defaultTimeSource struct{}

func (s *defaultTimeSource) Now() time.Time {
	return time.Now()
}

func (s *defaultTimeSource) NewTimer(d time.Duration) Timer {
	return NewTimer(d)
}

func (s *defaultTimeSource) AfterFunc(d time.Duration, f func()) Timer {
	return AfterFunc(d, f)
}

func (s *defaultTimeSource) NewContextWithDeadline(
	parent context.Context,
	deadline time.Time,
) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, deadline)
}

func (s *defaultTimeSource) NewContextWithTimeout(
	parent context.Context,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}
