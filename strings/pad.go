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

package strings

import (
	"strings"
)

// PadLeft prepends string s with i spaces.
func PadLeft(s string, i int) string {
	return PadLeftWith(s, i, " ")
}

// PadLeftWith prepends string s with i instances of padding p.
func PadLeftWith(s string, i int, p string) string {
	if i < 0 {
		return s
	}
	padding := strings.Repeat(p, i)

	ss := strings.Split(s, "\n")
	for i, e := range ss {
		ss[i] = padding + e
	}

	return strings.Join(ss, "\n")
}
