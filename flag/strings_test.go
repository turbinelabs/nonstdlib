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

import (
	"flag"
	"fmt"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestNewStrings(t *testing.T) {
	s := NewStrings()
	assert.Nil(t, s.Strings)
	assert.Nil(t, s.AllowedValues)
	assert.Equal(t, s.Delimiter, ",")
}

func TestNewStringsWithConstraint(t *testing.T) {
	s := NewStringsWithConstraint("a", "b", "c")
	assert.Nil(t, s.Strings)
	assert.DeepEqual(t, s.AllowedValues, []string{"a", "b", "c"})
	assert.Equal(t, s.Delimiter, ",")
}

func TestStringsString(t *testing.T) {
	s := &Strings{Strings: []string{"v1", "v2"}, Delimiter: "-"}
	assert.Equal(t, s.String(), "v1-v2")
}

func TestStringGet(t *testing.T) {
	s := &Strings{Strings: []string{"v1", "v2"}, Delimiter: ","}
	assert.DeepEqual(t, s.Get(), s.Strings)
}

type setTestCase struct {
	delimiter string
	input     string
	expected  []string
}

func (t *setTestCase) name(i int) string {
	return fmt.Sprintf(
		"testcase %d (%s split on '%s')",
		i,
		t.input,
		t.delimiter,
	)
}

func TestStringsSet(t *testing.T) {
	testcases := []setTestCase{
		{",", "a,b,c", []string{"a", "b", "c"}},
		{",", "a", []string{"a"}},
		{",", ",,,,,,", []string{}},
		{",", ",,a,,c,,", []string{"a", "c"}},
		{",", " , , c , , a , , ", []string{"c", "a"}},
		{",", "", []string{}},
		{"-", "a-b-c", []string{"a", "b", "c"}},
		{"-", "a,b,c", []string{"a,b,c"}},
	}

	for i, testcase := range testcases {
		assert.Group(
			testcase.name(i),
			t,
			func(g *assert.G) {
				s := &Strings{Delimiter: testcase.delimiter}
				err := s.Set(testcase.input)
				assert.DeepEqual(g, s.Strings, testcase.expected)
				assert.Nil(g, err)
			},
		)
	}
}

func TestStringsSetWithConstraint(t *testing.T) {
	s := &Strings{AllowedValues: []string{"a", "b", "c"}, Delimiter: ","}
	err := s.Set("a,b,c")
	assert.DeepEqual(t, s.Strings, []string{"a", "b", "c"})
	assert.Nil(t, err)

	s.ResetDefault()
	err = s.Set("a,b,c,d")
	assert.DeepEqual(t, s.Strings, []string{})
	assert.ErrorContains(t, err, "invalid flag value(s): d")

	s.ResetDefault()
	err = s.Set("x,a,y,b,z,c")
	assert.DeepEqual(t, s.Strings, []string{})
	assert.ErrorContains(t, err, "invalid flag value(s): x, y, z")
}

func TestStringsSetMulti(t *testing.T) {
	s := &Strings{Delimiter: ","}
	err := s.Set("a,b")
	assert.Nil(t, err)
	err = s.Set("c,d")
	assert.Nil(t, err)
	assert.ArrayEqual(t, s.Strings, []string{"a", "b", "c", "d"})

	s = &Strings{Delimiter: ","}
	s.ResetDefault("default-value-semantics")

	err = s.Set("a")
	assert.Nil(t, err)
	err = s.Set("b")
	assert.Nil(t, err)
	assert.ArrayEqual(t, s.Strings, []string{"a", "b"})
}

func TestStringsFlagSetIntegration(t *testing.T) {
	strings1 := NewStrings()
	strings2 := NewStrings()
	strings3 := NewStrings()

	strings2.Strings = []string{"some", "default", "values"}

	fs := flag.NewFlagSet("stuff", flag.PanicOnError)
	fs.Var(&strings1, "x", "Flag help")
	fs.Var(&strings2, "y", "Flag help")
	fs.Var(&strings3, "z", "Flag help")

	fs.Parse([]string{
		"-x=a",
		"-z=g,h",
		"-x=b",
		"-y=d,e,f",
		"-x=c",
		"-z=i",
	})

	assert.ArrayEqual(t, strings1.Strings, []string{"a", "b", "c"})
	assert.ArrayEqual(t, strings2.Strings, []string{"d", "e", "f"})
	assert.ArrayEqual(t, strings3.Strings, []string{"g", "h", "i"})
}

func TestStringsValidValuesDescription(t *testing.T) {
	c := &Strings{AllowedValues: []string{"a", "b", "c"}}
	assert.Equal(t, c.ValidValuesDescription(), `"a", "b", or "c"`)

	c = &Strings{AllowedValues: []string{"a", "b"}}
	assert.Equal(t, c.ValidValuesDescription(), `"a" or "b"`)

	c = &Strings{AllowedValues: []string{"a"}}
	assert.Equal(t, c.ValidValuesDescription(), `"a"`)

	c = &Strings{}
	assert.Equal(t, c.ValidValuesDescription(), "")
}
