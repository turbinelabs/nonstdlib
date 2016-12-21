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

package dedupe

import (
	"reflect"
	"testing"
)

func doTest(t *testing.T, in, want []int) {
	got := &ints{in}

	Dedupe(got)
	if !reflect.DeepEqual(got.is, want) {
		t.Errorf("got %v; wanted %v", got.is, want)
	}
}

func TestDedupe(t *testing.T) {
	doTest(
		t,
		[]int{0, 1, 2, 3, 4},
		[]int{0, 1, 2, 3, 4},
	)
}

func TestDedupeWithMiddleDupes(t *testing.T) {
	doTest(
		t,
		[]int{0, 2, 2, 2, 1, 3, 3, 3, 5},
		[]int{0, 2, 1, 3, 5},
	)
}

func TestDedupeStartEdgeCase(t *testing.T) {
	doTest(
		t,
		[]int{2, 2, 2, 1, 3, 3, 3, 5},
		[]int{2, 1, 3, 5},
	)
}

func TestDedupeEndEdgeCase(t *testing.T) {
	doTest(
		t,
		[]int{0, 2, 8, 1, 3, 5, 5, 5, 5, 5},
		[]int{0, 2, 8, 1, 3, 5},
	)
}

func TestDedupeStrings(t *testing.T) {
	in := []string{"a", "b", "b", "c", "c", "c"}
	out := Strings(in)
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("got: %v, want: %v", out, want)
	}
}
