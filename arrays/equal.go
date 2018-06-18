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

package arrays

func EqualInt(a, b []int) bool {
	return CompareIntSlices(a, b) == 0
}

// Treats smaller slices as being less than bigger ones and then does a per
// index comparison betweeen the two slices.
func CompareIntSlices(a, b []int) int {
	if len(a) < len(b) {
		return -1
	}

	if len(a) > len(b) {
		return 1
	}

	for i, av := range a {
		if av < b[i] {
			return -1
		}

		if av > b[i] {
			return 1
		}
	}

	return 0
}

func EqualInt64(a, b []int64) bool {
	return CompareInt64Slices(a, b) == 0
}

// Treats smaller slices as being less than bigger ones and then does a per
// index comparison betweeen the two slices.
func CompareInt64Slices(a, b []int64) int {
	if len(a) < len(b) {
		return -1
	}

	if len(a) > len(b) {
		return 1
	}

	for i, av := range a {
		if av < b[i] {
			return -1
		}

		if av > b[i] {
			return 1
		}
	}

	return 0
}

func EqualFloat64(a, b []float64) bool {
	return CompareFloat64Slices(a, b) == 0
}

// Treats smaller slices as being less than bigger ones and then does a per
// index comparison betweeen the two slices.
func CompareFloat64Slices(a, b []float64) int {
	if len(a) < len(b) {
		return -1
	}

	if len(a) > len(b) {
		return 1
	}

	for i, av := range a {
		if av < b[i] {
			return -1
		}

		if av > b[i] {
			return 1
		}
	}

	return 0
}

func EqualString(a, b []string) bool {
	return CompareStringSlices(a, b) == 0
}

// Treats smaller slices as being less than bigger ones and then does a per
// index comparison betweeen the two slices.
func CompareStringSlices(a, b []string) int {
	if len(a) < len(b) {
		return -1
	}

	if len(a) > len(b) {
		return 1
	}

	for i, av := range a {
		if av < b[i] {
			return -1
		}

		if av > b[i] {
			return 1
		}
	}

	return 0
}
