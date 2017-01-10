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

// Package indexof comprises types and functions in support of finding the
// index of a target value within a slice of same-typed values. Implementations
// for some common types are provided.
package indexof

import (
	"strings"
)

const NotFound = -1

// ElementPredicate is a function that applies some predicate to an element of
// a slice at a given index.
type ElementPredicate func(int) bool

// IndexOf takes the length of a sequence and an ElementPredicate that can
// determine if the element at a given index is equal to what we're looking
// for.
//
// It returns NotFound if no element causes test to return true or the index of
// the first element that causes test to return true.
//
// test will be called with the values 0 to (len - 1). Behavior is unspecified
// if len results in a call to test with a value outside the bounds of the
// slice being searched.
func IndexOf(len int, test ElementPredicate) int {
	for i := 0; i < len; i++ {
		if test(i) {
			return i
		}
	}
	return NotFound
}

// String returns the location of a target string within a []string or -1 if
// it is not found.
func String(ss []string, target string) int {
	return IndexOf(len(ss), func(i int) bool { return ss[i] == target })
}

// TrimmedString returns the location of a target string within a
// []string or -1 if it is not found. The comparison ignores leading
// and trailing spaces (see strings.TrimSpace).
func TrimmedString(ss []string, target string) int {
	target = strings.TrimSpace(target)

	return IndexOf(len(ss), func(i int) bool { return strings.TrimSpace(ss[i]) == target })
}

// Int returns the location of a target string within a []int or -1 if
// it is not found.
func Int(ss []int, target int) int {
	return IndexOf(len(ss), func(i int) bool { return ss[i] == target })
}
