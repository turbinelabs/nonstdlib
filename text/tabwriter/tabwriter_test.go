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
