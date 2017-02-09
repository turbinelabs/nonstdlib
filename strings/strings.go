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

// Package strings introduces additional utilities for dealing with strings and
// slices of strings.
package strings

// Set represets set-like operations for strings.
type Set map[string]bool

// Transform will apply some transformation to each item in the set.
func (ss Set) Transform(fn func(string) string) {
	for k, v := range ss {
		nk := fn(k)
		if nk != k {
			delete(ss, k)
			ss[nk] = v
		}
	}
}

// Equals returns true if two Sets contain the same items.
func (ss Set) Equals(o Set) bool {
	if len(ss) != len(o) {
		return false
	}
	for k := range ss {
		if !o.Contains(k) {
			return false
		}
	}

	return true
}

// Contains returns true if the set contains some key k.
func (ss Set) Contains(k string) bool {
	return ss[k]
}

// Remove removes an item from the set, no change is made if it is not present.
func (ss Set) Remove(k string) {
	delete(ss, k)
}

// Put adds a new key to the set.
func (ss Set) Put(k string) {
	ss[k] = true
}

// Slice returns the Set as a slice of values.
func (ss Set) Slice() []string {
	slc := make([]string, 0, len(ss))
	for k := range ss {
		slc = append(slc, k)
	}

	return slc
}

// NewSet ccreates a new Set containing the provided strings as an
// initial values stored in the set.
func NewSet(ss ...string) Set {
	sset := map[string]bool{}
	for _, s := range ss {
		sset[s] = true
	}

	return sset
}
