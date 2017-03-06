package flag

import (
	"flag"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestWrap(t *testing.T) {
	fs := &flag.FlagSet{}
	wrapped := Wrap(fs)
	assert.NonNil(t, wrapped)
	assert.SameInstance(t, wrapped.Unwrap(), fs)
}

func TestScopedUnwrap(t *testing.T) {
	fs := &flag.FlagSet{}
	wrapped := Wrap(fs)

	scoped := wrapped.Scope("a", "desc").Scope("b", "desc").Scope("c", "desc")
	assert.NonNil(t, scoped)
	assert.SameInstance(t, scoped.Unwrap(), fs)
}

func TestTestFlagSet(t *testing.T) {
	tfs := NewTestFlagSet()

	assert.Equal(t, tfs.Unwrap().NFlag(), 0)

	tfs.Int("x", 0, "usage")

	var recoveredPanic interface{}
	func() {
		defer func() { recoveredPanic = recover() }()

		tfs.Parse([]string{"-non-existent-flag"})
	}()

	assert.NonNil(t, recoveredPanic)
}
