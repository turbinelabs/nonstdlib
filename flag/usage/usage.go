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

// Package usages provides a mechanism to insert and recover richer usage
// information (eg whether a flag is required, whether it contains sensitive
// information) into and from a flag.Flag usage string, respectively.
package usage

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
)

// Usage represents richer usage information for a flag.Flag
type Usage interface {
	// Usage returns the original usage string, without any decoration
	Usage() string

	// IsRequired returns true if the flag has been marked as required
	IsRequired() bool

	// IsSensitive returns true if the flag has been marked as containing
	// sensitive information that should't be dispayed.
	IsSensitive() bool

	// IsDeprecated returns true if the flag has been marked as deprecated
	IsDeprecated() bool

	// SetRequired marks the flag as required
	SetRequired() Usage

	// SetSensitive marks the flag as containing sensitive information that
	// should't be dispayed.
	SetSensitive() Usage

	// SetDeprecated marks the flag as deprecated
	SetDeprecated() Usage

	// Pretty() returns a human-friendly usage string, decorated with an
	// indication of whether the flag has been marked as required or sensitive
	Pretty() string

	// String() returns an encoded string that can be passed to New() to recover
	// the full state of the Usage later.
	String() string
}

// New produces a new Usage from a string. The string can be either a conventual
// usage descrition, or a usage-encoded string.
func New(str string) Usage {
	u := &usage{}
	if err := json.Unmarshal([]byte(str), u); err != nil {
		return &usage{UsageStr: str}
	}
	return u
}

type usage struct {
	Required   bool   `json:"is_required"`
	Sensitive  bool   `json:"is_sensitive"`
	Deprecated bool   `json:"is_deprecated"`
	UsageStr   string `json:"usage"`
}

func (u *usage) Usage() string {
	return u.UsageStr
}

func (u *usage) IsRequired() bool {
	return u.Required
}

func (u *usage) IsSensitive() bool {
	return u.Sensitive
}

func (u *usage) IsDeprecated() bool {
	return u.Deprecated
}

func (u *usage) SetRequired() Usage {
	u.Required = true
	return u
}

func (u *usage) SetSensitive() Usage {
	u.Sensitive = true
	return u
}

func (u *usage) SetDeprecated() Usage {
	u.Deprecated = true
	return u
}

func (u *usage) Pretty() string {
	str := u.Usage()
	features := []string{}
	if u.IsRequired() {
		features = append(features, "REQUIRED")
	}
	if u.IsSensitive() {
		features = append(features, "SENSITIVE")
	}
	if u.IsDeprecated() {
		features = append(features, "DEPRECATED")
	}
	if len(features) > 0 {
		str = fmt.Sprintf("[%s] %s", strings.Join(features, "/"), str)
	}
	return str
}

func (u *usage) String() string {
	b, _ := json.Marshal(u)
	return string(b)
}

// Required produces an encoded usage string for a Flag indiciating that the Flag is
// required. It can be passed through usage.New() to recover the full Usage.
func Required(usage string) string {
	return New(usage).SetRequired().String()
}

// Sensitive produces a usage string for a Flag indiciating that the Flag is
// required. It can be passed through usage.New() to recover the full Usage.
func Sensitive(usage string) string {
	return New(usage).SetSensitive().String()
}

// Deprecated produces an encoded usage string for a Flag indiciating that the Flag is
// deprecated. It can be passed through usage.New() to recover the full Usage.
func Deprecated(usage string) string {
	return New(usage).SetDeprecated().String()
}

// IsRequired checks the usage string of the given Flag to see if it is
// marked as required.
func IsRequired(f *flag.Flag) bool {
	return New(f.Usage).IsRequired()
}

// IsSensitive checks the usage string of the given Flag to see if it is
// marked as sensitive.
func IsSensitive(f *flag.Flag) bool {
	return New(f.Usage).IsSensitive()
}

// IsDeprecated checks the usage string of the given Flag to see if it is
// marked as deprecated.
func IsDeprecated(f *flag.Flag) bool {
	return New(f.Usage).IsDeprecated()
}

// FlagSetFilterFn is a predicate function that takes a flag and whether or not
// the flag is set as its input.
type FlagSetFilterFn func(f *flag.Flag, set bool) bool

// FilterFlagSet produces a slice of the names of flags for which the given
// FilterFn returns true
func FilterFlagSet(fs *flag.FlagSet, fn FlagSetFilterFn) []string {
	seen := map[string]bool{}
	fs.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
	})

	result := []string{}
	fs.VisitAll(func(f *flag.Flag) {
		if fn(f, seen[f.Name]) {
			result = append(result, f.Name)
		}
	})

	return result
}

// MissingRequired produces a slice of the names of all flags for which the
// Usage string indiciates requiredness but no value is set.
func MissingRequired(fs *flag.FlagSet) []string {
	return FilterFlagSet(fs, func(f *flag.Flag, set bool) bool {
		return !set && IsRequired(f)
	})
}

// DeprecatedAndSet produces a slice of the names of all flags for which the
// Usage string indiciates deprecation but are set
func DeprecatedAndSet(fs *flag.FlagSet) []string {
	return FilterFlagSet(fs, func(f *flag.Flag, set bool) bool {
		return set && IsDeprecated(f)
	})
}
