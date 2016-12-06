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
