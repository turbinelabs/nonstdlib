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
	"strings"
)

const requiredPrefix = "[REQUIRED] "

// Required prefixes its argument with "[REQUIRED] " which, in addition from
// documenting for the users of the command.Cmd on which the flag is declared
// that the argument is required, will also cause it to be checked when the
// Cmd's Run method is invoked.
func Required(usage string) string {
	return requiredPrefix + usage
}

// IsRequired checks the usage string of the given Flag to see if it is
// prefixed with "[REQUIRED] ".
func IsRequired(f *flag.Flag) bool {
	return strings.HasPrefix(f.Usage, requiredPrefix)
}

// AllRequired produces a slice of the names of all flags for which the Usage
// string is prefxied with "[REQUIRED] ".
func AllRequired(fs *flag.FlagSet) []string {
	result := []string{}
	fs.VisitAll(func(f *flag.Flag) {
		if IsRequired(f) {
			result = append(result, f.Name)
		}
	})
	return result
}

// MissingRequired produces a slice of the names of all flags for which the
// Usage string is prefixed with "[REQUIRED] " but no value has been set.
func MissingRequired(fs *flag.FlagSet) []string {
	seen := map[string]bool{}
	fs.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
	})

	result := []string{}
	fs.VisitAll(func(f *flag.Flag) {
		if !seen[f.Name] && IsRequired(f) {
			result = append(result, f.Name)
		}
	})

	return result
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
