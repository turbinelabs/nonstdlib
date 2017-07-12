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

package flag

import (
	"flag"
)

// Wrap converts an existing *flag.FlagSet into a FlagSet with no
// scope.
func Wrap(fs *flag.FlagSet) FlagSet {
	return &flagSet{fs}
}

// NewTestFlagSet creates a new FlagSet suitable for tests. It has no
// prefix, contains no flags, and the unwrapped flag.FlagSet will
// panic on parse errors.
func NewTestFlagSet() TestFlagSet {
	return &testFlagSet{&flagSet{flag.NewFlagSet("test flags", flag.PanicOnError)}}
}

type flagSet struct {
	*flag.FlagSet
}

func (fs *flagSet) Scope(prefix, description string) FlagSet {
	return newPrefixedFlagSet(fs.FlagSet, prefix, description)
}

func (fs *flagSet) GetScope() string {
	return ""
}

func (fs *flagSet) Unwrap() *flag.FlagSet {
	return fs.FlagSet
}

// FlagSet represents an optionally scoped *flag.FlagSet for tests. It
// differs from FlagSet only in that methods not normally needed by
// consumers of FlagSet are directly available.
type TestFlagSet interface {
	FlagSet

	// Parse invokes the Parse function of the underlying
	// flag.FlagSet as a convenience for tests.
	Parse(args []string) error
}

type testFlagSet struct {
	*flagSet
}
