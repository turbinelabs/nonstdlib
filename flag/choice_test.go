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
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestNewChoice(t *testing.T) {
	c := NewChoice("a", "b", "c")
	assert.Nil(t, c.Choice)
	assert.DeepEqual(t, c.AllowedValues, []string{"a", "b", "c"})
}

func TestChoiceString(t *testing.T) {
	value := "v1"
	c := &Choice{Choice: &value}
	assert.Equal(t, c.String(), value)
}

func TestChoiceGet(t *testing.T) {
	value := "v1"
	c := &Choice{Choice: &value}
	assert.DeepEqual(t, c.Get(), &value)
}

func TestChoiceSet(t *testing.T) {
	v := "b"
	c := &Choice{AllowedValues: []string{"a", "b", "c"}}
	err := c.Set("b")
	assert.DeepEqual(t, c.Choice, &v)
	assert.Nil(t, err)

	c.Choice = nil
	err = c.Set("nope")
	assert.Nil(t, c.Choice)
	assert.ErrorContains(t, err, "invalid flag value: nope, must be one of a, b, c")
}

func TestChoiceWithDefault(t *testing.T) {
	c := Choice{AllowedValues: []string{"a", "b", "c"}}
	assert.Equal(t, c.String(), "")

	c = c.WithDefault("b")
	assert.Equal(t, c.String(), "b")
}

func TestChoiceValidValuesDescription(t *testing.T) {
	c := &Choice{AllowedValues: []string{"a", "b", "c"}}
	assert.Equal(t, c.ValidValuesDescription(), `"a", "b", or "c"`)

	c = &Choice{AllowedValues: []string{"a", "b"}}
	assert.Equal(t, c.ValidValuesDescription(), `"a" or "b"`)

	c = &Choice{AllowedValues: []string{"a"}}
	assert.Equal(t, c.ValidValuesDescription(), `"a"`)

	c = &Choice{}
	assert.Equal(t, c.ValidValuesDescription(), "")
}
