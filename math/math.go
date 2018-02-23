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

package math

import "math"

// Round computes the nearest integer to f, rounding ties towards
// positive infinity.
func Round(f float64) int64 {
	if f == 0.0 || math.IsNaN(f) {
		return 0
	}

	if f >= float64(math.MaxInt64) || math.IsInf(f, 1) {
		return math.MaxInt64
	}

	if f <= float64(math.MinInt64) || math.IsInf(f, -1) {
		return math.MinInt64
	}

	int, frac := math.Modf(f)
	if frac >= 0.5 {
		return int64(int) + 1
	}

	if frac < -0.5 {
		return int64(int) - 1
	}

	return int64(int)
}

// Sign returns -1 if f is negative, +1 if f is positive, and 0 if f
// is positive or negative zero or NaN.
func Sign(f float64) int {
	if f < 0.0 {
		return -1
	}

	if f > 0.0 {
		return 1
	}

	return 0
}

// SignInt64 returns -1 if f is negative, +1 if f is positive, and 0 if f
// is 0.
func SignInt64(i int64) int {
	if i < 0 {
		return -1
	}

	if i > 0 {
		return 1
	}

	return 0
}
