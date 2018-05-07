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
