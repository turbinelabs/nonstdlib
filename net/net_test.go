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

package net

import (
	"fmt"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestValidateListenerAddr(t *testing.T) {
	testCases := []struct {
		addr    string
		wantErr bool
	}{
		{addr: ":0"},
		{addr: ":80"},
		{addr: ":http"},
		{addr: ":somederpynet", wantErr: true},
		{addr: "0:80"},
		{addr: "0:http"},
		{addr: "0:somederpynet", wantErr: true},
		{addr: "127.0.0.1:80"},
		{addr: "127.0.0.1:http"},
		{addr: "127.0.0.1:somederpynet", wantErr: true},
		{addr: "localhost:80"},
		{addr: "localhost:http"},
		{addr: "localhost:somederpynet", wantErr: true},
		{addr: "[::1]:80"},
		{addr: "[::1]:http"},
		{addr: "[::1]:somederpynet", wantErr: true},
		{addr: ":999999", wantErr: true},
		{addr: ":-1}", wantErr: true},
		{addr: "", wantErr: true},
		{addr: ":", wantErr: true},
		{addr: "0:", wantErr: true},
		{addr: "localhost", wantErr: true},
	}

	for _, tc := range testCases {
		assert.Group(
			fmt.Sprintf("ValidateListenerAddr(%q)", tc.addr),
			t,
			func(g *assert.G) {
				err := ValidateListenerAddr(tc.addr)
				if tc.wantErr {
					assert.NonNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			},
		)
	}
}
