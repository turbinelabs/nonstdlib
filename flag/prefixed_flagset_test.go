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
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

const (
	flagUsage    = "a flag for {{NAME}} with prefix {{PREFIX}}"
	flagUsageFmt = "a flag for %s with prefix %s"
)

type prefixedFlagTestCase struct {
	name     string
	flagType reflect.Type

	addFlag func(f *prefixedFlagSet) interface{}
}

type stringValue struct {
	value string
}

func (sv *stringValue) Set(s string) error {
	sv.value = s
	return nil
}

func (sv *stringValue) String() string {
	return sv.value
}

var _ flag.Value = &stringValue{}

func (tc *prefixedFlagTestCase) run(t testing.TB) {
	fs := flag.NewFlagSet("generated code: "+tc.name, flag.PanicOnError)

	pfs := newPrefixedFlagSet(Wrap(fs), "theprefix", "the-app-name")

	target := tc.addFlag(pfs)

	var value string
	var expectedValue interface{}

	switch tc.flagType.Kind() {
	case reflect.Bool:
		value = "true"
		ev := true
		expectedValue = &ev
	case reflect.Int:
		value = "123"
		ev := 123
		expectedValue = &ev
	case reflect.Int64:
		if tc.flagType == reflect.TypeOf(time.Duration(0)) {
			value = "10s"
			ev := time.Duration(10 * time.Second)
			expectedValue = &ev
		} else {
			value = "123"
			ev := int64(123)
			expectedValue = &ev
		}
	case reflect.Uint:
		value = "123"
		ev := uint(123)
		expectedValue = &ev
	case reflect.Uint64:
		value = "123"
		ev := uint64(123)
		expectedValue = &ev
	case reflect.Float64:
		value = "12.3"
		ev := float64(12.3)
		expectedValue = &ev
	case reflect.String:
		value = "something"
		expectedValue = &value
	default:
		t.Fatalf("unhandled type %s", tc.flagType.String())
		return
	}

	flagName := "theprefix." + tc.name

	fs.Parse([]string{
		"-" + flagName + "=" + value,
	})

	if valueTarget, ok := target.(flag.Value); ok {
		v := valueTarget.String()
		assert.DeepEqual(t, &v, expectedValue)
	} else {
		assert.DeepEqual(t, target, expectedValue)
	}

	f := fs.Lookup(flagName)
	assert.NonNil(t, f)
	assert.Equal(t, f.Name, flagName)
	assert.Equal(t, f.Usage, fmt.Sprintf(flagUsageFmt, "the-app-name", "theprefix."))
}

func TestGeneratedCode(t *testing.T) {
	for _, tc := range generatedTestCases {
		assert.Group(fmt.Sprintf("for flag type %s", tc.name), t, func(g *assert.G) {
			tc.run(g)
		})
	}
}

func TestVar(t *testing.T) {
	testCase := prefixedFlagTestCase{
		name:     "var",
		flagType: reflect.TypeOf(string(0)),
		addFlag: func(f *prefixedFlagSet) interface{} {
			var target stringValue
			f.Var(&target, "var", flagUsage)
			return &target
		},
	}
	testCase.run(t)
}

func TestHostPortVar(t *testing.T) {
	fs := NewTestFlagSet()

	pfs := fs.Scope("theprefix", "the-app-name")

	hp := NewHostPort("default:80")
	pfs.HostPortVar(&hp, "hostport", hp, flagUsage)

	fs.Parse([]string{"-theprefix.hostport=example.com:443"})

	assert.Equal(t, hp.Addr(), "example.com:443")

	f := fs.Unwrap().Lookup("theprefix.hostport")
	assert.NonNil(t, f)
	assert.Equal(t, f.Name, "theprefix.hostport")
	assert.Equal(t, f.Usage, fmt.Sprintf(flagUsageFmt, "the-app-name", "theprefix."))
}

func TestHostPort(t *testing.T) {
	fs := NewTestFlagSet()

	pfs := fs.Scope("theprefix", "the-app-name")

	hp := pfs.HostPort("hostport", NewHostPort("default:80"), flagUsage)

	fs.Parse([]string{"-theprefix.hostport=example.com:443"})

	assert.Equal(t, hp.Addr(), "example.com:443")

	f := fs.Unwrap().Lookup("theprefix.hostport")
	assert.NonNil(t, f)
	assert.Equal(t, f.Name, "theprefix.hostport")
	assert.Equal(t, f.Usage, fmt.Sprintf(flagUsageFmt, "the-app-name", "theprefix."))
}

func TestScope(t *testing.T) {
	fs := flag.NewFlagSet("scoping test", flag.PanicOnError)
	underlying := newPrefixedFlagSet(Wrap(fs), "theprefix", "the-app-name")
	pfs := underlying.Scope("scope", "scope-name")
	assert.SameInstance(t, pfs.Unwrap(), fs)
	pfsImpl := pfs.(*prefixedFlagSet)
	assert.Equal(t, pfsImpl.prefix, "theprefix.scope.")
	assert.Equal(t, pfsImpl.descriptor, "scope-name")

	pfs = underlying.Scope("scope.", "scope-name")
	assert.SameInstance(t, pfs.Unwrap(), fs)
	pfsImpl = pfs.(*prefixedFlagSet)
	assert.Equal(t, pfsImpl.prefix, "theprefix.scope.")
	assert.Equal(t, pfsImpl.descriptor, "scope-name")

	pfs = underlying.Scope("scope.", "{{NAME}}: scope-name")
	assert.SameInstance(t, pfs.Unwrap(), fs)
	pfsImpl = pfs.(*prefixedFlagSet)
	assert.Equal(t, pfsImpl.prefix, "theprefix.scope.")
	assert.Equal(t, pfsImpl.descriptor, "the-app-name: scope-name")
}

func TestGetScope(t *testing.T) {
	fs := flag.NewFlagSet("scoping test", flag.PanicOnError)
	wrapped := Wrap(fs)
	assert.Equal(t, wrapped.GetScope(), "")

	scoped := wrapped.Scope("theprefix", "")
	assert.Equal(t, scoped.GetScope(), "theprefix.")

	scopedAgain := scoped.Scope("more", "")
	assert.Equal(t, scopedAgain.GetScope(), "theprefix.more.")
}
