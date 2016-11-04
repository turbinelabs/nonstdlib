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
