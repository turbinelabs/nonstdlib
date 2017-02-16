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
	"testing"

	"github.com/turbinelabs/test/assert"
)

func testFlags() (*flag.FlagSet, *string, *string, *string) {
	var fs flag.FlagSet
	fooFlag := fs.String("foo-baz", "", "do the foo")
	barFlag := fs.String("bar", "", Required("harty har to the bar"))
	quxFlag := fs.String("qux", "qux-default", "if it qux like a duck...")
	return &fs, fooFlag, barFlag, quxFlag
}

func TestRequired(t *testing.T) {
	assert.Equal(t, Required("foo"), "[REQUIRED] foo")
}

func TestIsRequired(t *testing.T) {
	assert.True(t, IsRequired(&flag.Flag{Usage: Required("foo")}))
	assert.False(t, IsRequired(&flag.Flag{}))
}

func TestMissingRequired(t *testing.T) {
	fs, _, _, _ := testFlags()
	assert.DeepEqual(t, MissingRequired(fs), []string{"bar"})
	fs.Parse([]string{"--bar=baz"})
	assert.DeepEqual(t, MissingRequired(fs), []string{})
}

func TestAllRequired(t *testing.T) {
	fs, _, _, _ := testFlags()
	assert.DeepEqual(t, AllRequired(fs), []string{"bar"})
	fs.Parse([]string{"--bar=baz"})
	assert.DeepEqual(t, AllRequired(fs), []string{"bar"})
}

func TestEnumerateNil(t *testing.T) {
	got := Enumerate(nil)
	assert.Equal(t, len(got), 0)
}

func TestEnumerateEmpty(t *testing.T) {
	got := Enumerate(&flag.FlagSet{})
	assert.Equal(t, len(got), 0)
}

func TestEnumerate(t *testing.T) {
	fs, _, _, _ := testFlags()
	got := Enumerate(fs)
	assert.Equal(t, len(got), 3)

	names := []string{got[0].Name, got[1].Name, got[2].Name}
	assert.HasSameElements(
		t,
		names,
		[]string{"foo-baz", "bar", "qux"},
	)
}

func TestIsSet(t *testing.T) {
	fs, _, _, _ := testFlags()
	assert.False(t, IsSet(fs, "foo-baz"))
	assert.False(t, IsSet(fs, "bar"))
	assert.False(t, IsSet(fs, "qux"))

	fs, _, _, _ = testFlags()
	fs.Parse([]string{"--foo-baz=x"})
	assert.True(t, IsSet(fs, "foo-baz"))
	assert.False(t, IsSet(fs, "bar"))
	assert.False(t, IsSet(fs, "qux"))

	fs, _, _, _ = testFlags()
	fs.Parse([]string{"--foo-baz=x", "--bar=y", "--qux=z"})
	assert.True(t, IsSet(fs, "foo-baz"))
	assert.True(t, IsSet(fs, "bar"))
	assert.True(t, IsSet(fs, "qux"))

	// IsSet == true even if CLI has flag with default value
	fs, _, _, _ = testFlags()
	fs.Parse([]string{"--qux=qux-default"})
	assert.False(t, IsSet(fs, "foo-baz"))
	assert.False(t, IsSet(fs, "bar"))
	assert.True(t, IsSet(fs, "qux"))

	// Simulate flags set by FromEnv, including one set to default
	// value.
	fs, _, _, _ = testFlags()
	fs.VisitAll(func(flag *flag.Flag) {
		if flag.Name == "bar" {
			fs.Set(flag.Name, "x")
		} else if flag.Name == "qux" {
			// Set to the default value
			fs.Set(flag.Name, "qux-default")
		}
	})
	assert.False(t, IsSet(fs, "foo-baz"))
	assert.True(t, IsSet(fs, "bar"))
	assert.True(t, IsSet(fs, "qux"))
}
