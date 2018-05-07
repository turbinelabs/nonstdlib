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

package tabwriter

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestFormat(t *testing.T) {
	in := "A\tB\tC\nsome a\tsome b\tsome c"
	out := Format(in)

	assert.Equal(
		t,
		out,
		`A       B       C
some a  some b  some c
`,
	)
}

func TestFormatWithHeaderAndConfig(t *testing.T) {
	cfg := Config{PadChar: '.'}
	tw := New(cfg)

	hdr := "A\tB\tC"
	in := "some a\tsome b\tsome c"
	out := tw.FormatWithHeader(hdr, in)

	assert.Equal(
		t,
		out,
		`A.....B.....C
some asome bsome c
`,
	)
}

func TestFormatWithHeader(t *testing.T) {
	hdr := "A\tB\tC"
	in := "some a\tsome b\tsome c"
	out := FormatWithHeader(hdr, in)

	assert.Equal(
		t,
		out,
		`A       B       C
some a  some b  some c
`,
	)
}
