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

package flag

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE --write_package_comment=false

import (
	"flag"
	"regexp"
	"strings"

	"github.com/turbinelabs/nonstdlib/flag/usage"
	tbnos "github.com/turbinelabs/nonstdlib/os"
)

var (
	notAlphaNum         = regexp.MustCompile("[^A-Za-z0-9_]+")
	multipleUnderscores = regexp.MustCompile("_+")
)

// NewFromEnv produces a FromEnv, using the provided FlagSet and scopes.
// The scopes are used to produce the environment key prefix, by uppercasing,
// replacing non-alphanumeric+underscore characters with underscores, and
// concatenating with underscores.
//
// For example:
//  {"foo-foo", "bar.bar", "baz"} -> "FOO_FOO_BAR_BAR_BAZ"
func NewFromEnv(fs *flag.FlagSet, scopes ...string) FromEnv {
	return fromEnv{
		prefix:        EnvKey(scopes...),
		fs:            fs,
		os:            tbnos.New(),
		filledFromEnv: map[string]string{},
	}
}

// FromEnv supports operations on a FlagSet based on environment variables.
// In particular, FromEnv allows one to fill a FlagSet from the environment
// and then inspect the results.
type FromEnv interface {
	// Prefix returns the environment key prefix, eg "SOME_PREFIX_"
	Prefix() string

	// Fill parses all registered flags in the FlagSet, and if they are not already
	// set it attempts to set their values from environment variables. Environment
	// variables take the name of the flag but are UPPERCASE, have the given prefix,
	// and non-alphanumeric+underscore chars are replaced by underscores.
	//
	// For example:
	//  some-flag -> SOME_PREFIX_SOME_FLAG
	//
	// the provided map[string]string is also populated with the keys and values
	// added to the FlagSet.
	Fill() error

	// Filled returns a map of the environment keys and values for flags currently
	// filled from the environment. Values for flags marked sensitive will be
	// redacted
	Filled() map[string]string

	// AllFlags returns a slice containing all Flags in the underlying Flagset
	AllFlags() []*flag.Flag
}

type fromEnv struct {
	prefix        string
	fs            *flag.FlagSet
	os            tbnos.OS
	filledFromEnv map[string]string
}

func (fe fromEnv) Prefix() string {
	return EnvKey(fe.prefix, "")
}

func (fe fromEnv) Fill() error {
	var firstErr error
	alreadySet := map[string]bool{}
	fe.fs.Visit(func(f *flag.Flag) {
		alreadySet[f.Name] = true
	})
	fe.fs.VisitAll(func(f *flag.Flag) {
		if !alreadySet[f.Name] {
			key := EnvKey(fe.prefix, f.Name)
			val, found := fe.os.LookupEnv(key)
			if found {
				if usage.IsSensitive(f) {
					fe.filledFromEnv[key] = "<redacted>"
				} else {
					fe.filledFromEnv[key] = val
				}
				if err := fe.fs.Set(f.Name, val); err != nil {
					if firstErr == nil {
						firstErr = err
					}
					return
				}
			}
		}
	})
	return firstErr
}

func (fe fromEnv) Filled() map[string]string {
	return fe.filledFromEnv
}

func (fe fromEnv) AllFlags() []*flag.Flag {
	return Enumerate(fe.fs)
}

// EnvKey produces a namespaced environment variable key, concatenates a prefix
// and key with an infix underscore, replacing all non-alphanumeric,
// non-underscore characters with underscores, and upper-casing the entire
// string
func EnvKey(parts ...string) string {
	for i, part := range parts {
		parts[i] = notAlphaNum.ReplaceAllString(part, "_")
	}
	joined := strings.ToUpper(strings.Join(parts, "_"))
	return multipleUnderscores.ReplaceAllLiteralString(joined, "_")
}
