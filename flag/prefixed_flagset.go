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

type PrefixedFlagSet struct {
	*flag.FlagSet

	prefix     string
	descriptor string
}

func NewPrefixedFlagSet(fs *flag.FlagSet, prefix, descriptor string) *PrefixedFlagSet {
	if prefix != "" && !strings.HasSuffix(prefix, ".") {
		prefix = prefix + "."
	}

	return &PrefixedFlagSet{
		FlagSet:    fs,
		prefix:     prefix,
		descriptor: descriptor,
	}
}

func (f *PrefixedFlagSet) mkUsage(usage string) string {
	return strings.Replace(usage, "{{NAME}}", f.descriptor, -1)
}

func (f *PrefixedFlagSet) Var(value flag.Value, name string, usage string) {
	f.FlagSet.Var(value, f.prefix+name, f.mkUsage(usage))
}

func (f *PrefixedFlagSet) Scope(prefix, descriptor string) *PrefixedFlagSet {
	return NewPrefixedFlagSet(f.FlagSet, f.prefix+prefix, descriptor)
}

func (f *PrefixedFlagSet) Descriptor() string {
	return f.descriptor
}
