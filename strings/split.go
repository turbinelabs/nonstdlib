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

package strings

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// SplitFirstEqual invokes Split2 to split a string of the form "A=B".
func SplitFirstEqual(s string) (string, string) {
	return Split2(s, "=")
}

// SplitFirstColon invokes Split2 to split a string of the form "A:B".
func SplitFirstColon(s string) (string, string) {
	return Split2(s, ":")
}

// SplitHostPort splits a string of the form "host:port". After
// splitting the string via net.SplitHostPort, the port is converted
// to an integer. If there is an error splitting the host and port, if
// port cannot be converted into an integer, or if port is not in the
// range [0, 65535], an error is returned.
func SplitHostPort(s string) (string, int, error) {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return "", 0, err
	}

	if host == "" {
		return "", 0, fmt.Errorf("address %s: missing host", s)
	}

	if port == "" {
		return "", 0, fmt.Errorf("address %s: missing port", s)
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, fmt.Errorf("address %s: cannot convert port to integer", s)
	}

	if portNum < 0 || portNum > 65535 {
		return "", 0, fmt.Errorf("address %s: port out of range", s)
	}

	return host, portNum, nil
}

// Split2 splits a string on the first occurrence of the given
// delimiter and returns the portions of the original string before
// and after the delimiter.  If no delimiter appears in the given
// string, the entire string is returned as the first result and the
// second result is empty. If multiple delimiters appear in the
// string, the strings are split at the first occurrence. The
// delimiter may appear at the start of end of the string, resulting
// in either the first or second result being the empty string.
func Split2(s, delim string) (string, string) {
	kv := strings.SplitN(s, delim, 2)
	if len(kv) == 1 {
		return kv[0], ""
	}

	return kv[0], kv[1]
}
