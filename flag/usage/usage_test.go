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

package usage

import (
	"flag"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestRequired(t *testing.T) {
	json := `{"is_required":true,"is_sensitive":false,"is_deprecated":false,"usage":"foo"}`
	assert.Equal(t, Required("foo"), json)

	usage := New("foo")
	assert.False(t, usage.IsRequired())
	usage.SetRequired()
	assert.True(t, usage.IsRequired())
	assert.Equal(t, usage.Usage(), "foo")
	assert.Equal(t, usage.Pretty(), "[REQUIRED] foo")
	assert.Equal(t, usage.String(), json)

	usage = New(json)
	assert.True(t, usage.IsRequired())
	assert.Equal(t, usage.Usage(), "foo")
	assert.Equal(t, usage.Pretty(), "[REQUIRED] foo")
	assert.Equal(t, usage.String(), json)
}

func TestIsRequired(t *testing.T) {
	assert.True(t, IsRequired(&flag.Flag{Usage: New("foo").SetSensitive().SetRequired().String()}))
	assert.True(t, IsRequired(&flag.Flag{Usage: New("foo").SetRequired().String()}))
	assert.True(t, IsRequired(&flag.Flag{Usage: Required("foo")}))
	assert.False(t, IsRequired(&flag.Flag{Usage: "foo"}))
}

func TestSensitive(t *testing.T) {
	json := `{"is_required":false,"is_sensitive":true,"is_deprecated":false,"usage":"foo"}`
	assert.Equal(t, Sensitive("foo"), json)

	usage := New("foo")
	assert.False(t, usage.IsSensitive())
	usage.SetSensitive()
	assert.True(t, usage.IsSensitive())
	assert.Equal(t, usage.Usage(), "foo")
	assert.Equal(t, usage.Pretty(), "[SENSITIVE] foo")
	assert.Equal(t, usage.String(), json)

	usage = New(json)
	assert.True(t, usage.IsSensitive())
	assert.Equal(t, usage.Usage(), "foo")
	assert.Equal(t, usage.Pretty(), "[SENSITIVE] foo")
	assert.Equal(t, usage.String(), json)
}

func TestIsSensitive(t *testing.T) {
	assert.True(t, IsSensitive(&flag.Flag{Usage: New("foo").SetSensitive().SetRequired().String()}))
	assert.True(t, IsSensitive(&flag.Flag{Usage: New("foo").SetSensitive().String()}))
	assert.True(t, IsSensitive(&flag.Flag{Usage: Sensitive("foo")}))
	assert.False(t, IsSensitive(&flag.Flag{Usage: "foo"}))
}

func TestDeprecated(t *testing.T) {
	json := `{"is_required":false,"is_sensitive":false,"is_deprecated":true,"usage":"foo"}`
	assert.Equal(t, Deprecated("foo"), json)

	usage := New("foo")
	assert.False(t, usage.IsDeprecated())
	usage.SetDeprecated()
	assert.True(t, usage.IsDeprecated())
	assert.Equal(t, usage.Usage(), "foo")
	assert.Equal(t, usage.Pretty(), "[DEPRECATED] foo")
	assert.Equal(t, usage.String(), json)

	usage = New(json)
	assert.True(t, usage.IsDeprecated())
	assert.Equal(t, usage.Usage(), "foo")
	assert.Equal(t, usage.Pretty(), "[DEPRECATED] foo")
	assert.Equal(t, usage.String(), json)
}

func TestIsDeprecated(t *testing.T) {
	assert.True(t, IsDeprecated(&flag.Flag{Usage: New("foo").SetDeprecated().SetSensitive().SetRequired().String()}))
	assert.True(t, IsDeprecated(&flag.Flag{Usage: New("foo").SetDeprecated().String()}))
	assert.True(t, IsDeprecated(&flag.Flag{Usage: Deprecated("foo")}))
	assert.False(t, IsDeprecated(&flag.Flag{Usage: "foo"}))
}

func TestRequiredAndSensitiveAndDeprecated(t *testing.T) {
	json := `{"is_required":true,"is_sensitive":true,"is_deprecated":true,"usage":"foo"}`
	assert.Equal(t, Deprecated(Required(Sensitive("foo"))), json)
	assert.Equal(t, Deprecated(Sensitive(Required("foo"))), json)
	assert.Equal(t, Sensitive(Required(Deprecated("foo"))), json)
	assert.Equal(t, Sensitive(Deprecated(Required("foo"))), json)
	assert.Equal(t, Required(Sensitive(Deprecated("foo"))), json)
	assert.Equal(t, Required(Deprecated(Sensitive("foo"))), json)

	usage := New("foo")
	assert.False(t, usage.IsSensitive())
	assert.False(t, usage.IsRequired())
	assert.False(t, usage.IsDeprecated())
	usage.SetSensitive().SetRequired().SetDeprecated()
	assert.True(t, usage.IsRequired())
	assert.True(t, usage.IsSensitive())
	assert.True(t, usage.IsDeprecated())
	assert.Equal(t, usage.Usage(), "foo")
	assert.Equal(t, usage.Pretty(), "[REQUIRED/SENSITIVE/DEPRECATED] foo")
	assert.Equal(t, usage.String(), json)

	usage = New(json)
	assert.True(t, usage.IsRequired())
	assert.True(t, usage.IsSensitive())
	assert.True(t, usage.IsDeprecated())
	assert.Equal(t, usage.Usage(), "foo")
	assert.Equal(t, usage.Pretty(), "[REQUIRED/SENSITIVE/DEPRECATED] foo")
	assert.Equal(t, usage.String(), json)
}

func TestMissingRequired(t *testing.T) {
	fs := &flag.FlagSet{}
	fs.String("bar", "", Required("harty har to the bar"))
	assert.DeepEqual(t, MissingRequired(fs), []string{"bar"})
	fs.Parse([]string{"--bar=baz"})
	assert.DeepEqual(t, MissingRequired(fs), []string{})
}
