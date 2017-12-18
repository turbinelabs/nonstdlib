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

//go:generate ./prefixed_flagset_gen.sh bool time.Duration float64 int int64 string uint uint64

import (
	"flag"
	"strings"
)

// prefixedFlagSet extends a flag.FlagSet to allow arbitrary-depth scoping
// of flags, using "." as a delemiter.
type prefixedFlagSet struct {
	FlagSet

	prefix     string
	descriptor string
}

// newPrefixedFlagSet produces a new prefixedFlagSet with the given
// FlagSet, prefix and descriptor. The descriptor is included will be
// used to replace the string "{{NAME}}" in usage strings used when
// declaring Flags.
func newPrefixedFlagSet(fs FlagSet, prefix, descriptor string) *prefixedFlagSet {
	if prefix != "" && !strings.HasSuffix(prefix, ".") {
		prefix = prefix + "."
	}

	return &prefixedFlagSet{
		FlagSet:    fs,
		prefix:     prefix,
		descriptor: descriptor,
	}
}

func (f *prefixedFlagSet) mkUsage(usage string) string {
	usage = strings.Replace(usage, "{{NAME}}", f.descriptor, -1)
	return strings.Replace(usage, "{{PREFIX}}", f.prefix, -1)
}

// Var wraps the underlying FlagSet's Var function.
func (f *prefixedFlagSet) Var(value flag.Value, name string, usage string) {
	f.FlagSet.Var(value, f.prefix+name, f.mkUsage(usage))
}

func (f *prefixedFlagSet) HostPortVar(hp *HostPort, name string, value HostPort, usage string) {
	f.FlagSet.HostPortVar(hp, f.prefix+name, value, f.mkUsage(usage))
}

func (f *prefixedFlagSet) HostPort(name string, value HostPort, usage string) *HostPort {
	return f.FlagSet.HostPort(f.prefix+name, value, f.mkUsage(usage))
}

// Scope scopes the target prefixedFlagSet to produce a new FlagSet,
// with the given scope an descriptor.
func (f *prefixedFlagSet) Scope(prefix, descriptor string) FlagSet {
	// apply {{NAME}} to descriptor
	descriptor = f.mkUsage(descriptor)

	return newPrefixedFlagSet(f.FlagSet, f.prefix+prefix, descriptor)
}

func (f *prefixedFlagSet) GetScope() string {
	return f.prefix
}

func (f *prefixedFlagSet) Unwrap() *flag.FlagSet {
	return f.FlagSet.Unwrap()
}
