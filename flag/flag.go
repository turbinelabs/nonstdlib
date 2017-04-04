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

// Package flag provides convenience methods for dealing with golang
// flag.FlagSets
package flag

import (
	"flag"
)

const requiredPrefix = "[REQUIRED] "

// ConstrainedValue is a flag.Value with constraints on it assigned
// value. Typically the constraints limit a flag to a given set of
// strings.
type ConstrainedValue interface {
	flag.Value

	// ValidValuesDescription provides a description of the value
	// values for this ContrainedValue suiteable for use in flag
	// usage text.
	ValidValuesDescription() string
}

// Enumerate returns a slice containing all Flags in the Flagset
func Enumerate(flagset *flag.FlagSet) []*flag.Flag {
	if flagset == nil {
		return []*flag.Flag{}
	}
	flags := make([]*flag.Flag, 0)
	flagset.VisitAll(func(f *flag.Flag) {
		flags = append(flags, f)
	})
	return flags
}

// IsSet indicates whether a given Flag in a FlagSet has been set or
// not. The name should be the same value (case-sensitive) as that
// passed to the FlagSet methods for constructing flags.
func IsSet(flagset *flag.FlagSet, name string) bool {
	if flagset == nil {
		return false
	}

	found := false
	flagset.Visit(func(flag *flag.Flag) {
		if flag.Name == name {
			found = true
		}
	})
	return found
}
