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

package executor

import (
	"math"
	"time"
)

// DelayFunc is invoked to compute a new deadline. The value passed is
// the number of times the action has been previously attempted. It is
// always greater than or equal to 1.
type DelayFunc func(int) time.Duration

// NewExponentialDelayFunc creates a new DelayFunc where the first
// retry occurs after a duration of delay and each subsequent retry is
// delayed by twice the previous delay. E.g., given a delay of 1s, the
// delays for retries are 1s, 2s, 4s, 8s, ... The return value is
// capped at the specified maximum delay. Delays of less than 0 are
// treated as 0.
func NewExponentialDelayFunc(delay time.Duration, maxDelay time.Duration) DelayFunc {
	if delay <= 0 {
		return NewConstantDelayFunc(0)
	}

	if maxDelay < delay {
		maxDelay = delay
	}

	// after this attempt, we would exceed maxDelay
	maxAttempts := int(math.Floor(math.Log2(float64(maxDelay/delay)))) + 1
	if maxAttempts >= 64 {
		maxAttempts = 63
	}

	return func(attempt int) time.Duration {
		if attempt <= 0 {
			// guard against bad input
			return delay
		} else if attempt > maxAttempts {
			return maxDelay
		}

		exp := int64(1) << uint(attempt-1)

		d := time.Duration(exp) * delay
		if d > maxDelay {
			return maxDelay
		}

		return d
	}
}

// NewConstantDelayFunc creates a DelayFunc where all retries occur
// after a fixed delay.
func NewConstantDelayFunc(delay time.Duration) DelayFunc {
	if delay < 0 {
		delay = 0
	}

	return func(attempt int) time.Duration {
		return delay
	}
}
