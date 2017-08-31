package reflect

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

type tcase struct {
	in    interface{}
	isNil bool
}

func (tc tcase) run(t *testing.T) {
	assert.Equal(t, IsNil(tc.in), tc.isNil)
}

func TestIsNil(t *testing.T) {
	var iNil interface{}
	var i interface{} = tcase{}
	var s tcase
	var sptrNil *tcase
	sptr := &tcase{}
	primitive := 1234
	str := "aoseutnh"
	strPtr := &str
	var primPtrNil *int

	cases := []tcase{
		{iNil, true},
		{i, false},
		{s, false},
		{sptrNil, true},
		{sptr, false},
		{primitive, false},
		{str, false},
		{strPtr, false},
		{primPtrNil, true},
	}

	for _, c := range cases {
		c.run(t)
	}
}
