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
	"testing"

	"github.com/golang/mock/gomock"

	tbnos "github.com/turbinelabs/nonstdlib/os"
	"github.com/turbinelabs/test/assert"
)

func TestFromEnvPrefix(t *testing.T) {
	want := "A_B_CD_A_B_CD_A_B_CD"
	values := []string{
		"A__B_CD",
		"A-.B*CD",
		"A.b-##cd",
		"a()&b-cd",
		"aö@bë^cD",
		"a\t    b\ncD",
	}
	for _, v := range values {
		assert.Equal(t, EnvKey(v, v, v), want)
		assert.Equal(t, NewFromEnv(nil, v, v, v).Prefix(), want+"_")
	}
}

func TestFromEnvAllFlagsNil(t *testing.T) {
	got := NewFromEnv(nil).AllFlags()
	assert.Equal(t, len(got), 0)
}

func TestFromEnvAllFlagsEmpty(t *testing.T) {
	got := NewFromEnv(&flag.FlagSet{}).AllFlags()
	assert.Equal(t, len(got), 0)
}

func TestFromEnvAllFlags(t *testing.T) {
	fs, _, _, _ := testFlags()
	got := NewFromEnv(fs).AllFlags()
	assert.Equal(t, len(got), 3)
	names := []string{got[0].Name, got[1].Name, got[2].Name}
	assert.HasSameElements(
		t,
		names,
		[]string{"foo-baz", "bar", "qux"},
	)
}

func TestFillFromEnvAllUnset(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	mockOS := tbnos.NewMockOS(ctrl)
	mockOS.EXPECT().LookupEnv("FOO_BAR_FOO_BAZ").Return("", true)
	mockOS.EXPECT().LookupEnv("FOO_BAR_BAR").Return("", false)
	mockOS.EXPECT().LookupEnv("FOO_BAR_QUX").Return("", false)

	fs, fooFlag, barFlag, quxFlag := testFlags()
	*fooFlag = "foo-default"
	*barFlag = "bar-default"
	fs.Parse([]string{})

	fe := NewFromEnv(fs, "foo", "bar").(fromEnv)
	fe.os = mockOS

	assert.Nil(t, fe.Fill())

	assert.Equal(t, *fooFlag, "")
	assert.Equal(t, *barFlag, "bar-default")
	assert.Equal(t, *quxFlag, "qux-default")
	assert.DeepEqual(t, fe.Filled(), map[string]string{
		"FOO_BAR_FOO_BAZ": "",
	})
}

func TestFillFromEnvOneSet(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	mockOS := tbnos.NewMockOS(ctrl)
	mockOS.EXPECT().LookupEnv("FOO_BAR_BAR").Return("", false)
	mockOS.EXPECT().LookupEnv("FOO_BAR_QUX").Return("", false)

	fs, fooFlag, barFlag, _ := testFlags()
	fs.Parse([]string{"--foo-baz=blargo"})

	fe := NewFromEnv(fs, "foo", "bar").(fromEnv)
	fe.os = mockOS

	assert.Nil(t, fe.Fill())

	assert.Equal(t, *fooFlag, "blargo")
	assert.Equal(t, *barFlag, "")
	assert.DeepEqual(t, fe.Filled(), map[string]string{})
}

func TestFillFromEnvIllegalValue(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	mockOS := tbnos.NewMockOS(ctrl)
	mockOS.EXPECT().LookupEnv("FOO_BAR_FOO_BAZ").Return("", false)
	mockOS.EXPECT().LookupEnv("FOO_BAR_BAR").Return("", false)
	mockOS.EXPECT().LookupEnv("FOO_BAR_QUX").Return("", false)
	mockOS.EXPECT().LookupEnv("FOO_BAR_INT").Return("not an int", true)

	fs, _, _, _ := testFlags()
	fs.Int("int", 0, "some int")
	fs.Parse([]string{""})

	fe := NewFromEnv(fs, "foo", "bar").(fromEnv)
	fe.os = mockOS

	assert.ErrorContains(t, fe.Fill(), `strconv.ParseInt: parsing "not an int"`)
	assert.DeepEqual(t, fe.Filled(), map[string]string{"FOO_BAR_INT": "not an int"})
}

func TestFillFromEnvOneSensitive(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	mockOS := tbnos.NewMockOS(ctrl)
	mockOS.EXPECT().LookupEnv("FOO_BAR_BAR").Return("", false)
	mockOS.EXPECT().LookupEnv("FOO_BAR_QUX").Return("it is a duck", true)

	fs, _, _, _ := testFlags()
	fs.Parse([]string{"--foo-baz=blargo"})

	fe := NewFromEnv(fs, "foo", "bar").(fromEnv)
	fe.os = mockOS

	assert.Nil(t, fe.Fill())

	assert.DeepEqual(t, fe.Filled(), map[string]string{"FOO_BAR_QUX": "<redacted>"})
}
