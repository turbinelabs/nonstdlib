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

// Package net provides convenience methods for dealing with the net package
// of the stdlib.
package net

import (
	"fmt"
	"net"
)

// ValidateListenerAddr validates the given listener address. The
// address must be in the form host:port.
//
// The host may be omitted to indicate all local addresses. The host
// may be enclosed in square brackets to represent an IPv6 address
// literal or host name. Otherwise it may be an IPv4 address literal
// or host name. Host names are not checked for resolvability.
//
// The port is required. It may be an integer in the range [0, 65535]
// or a well-known TCP/IP service name, such as "http".
func ValidateListenerAddr(addr string) error {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	if port == "" {
		return fmt.Errorf("address %s: no port in listener address", addr)
	}

	_, err = net.LookupPort("tcp", port)
	if err != nil {
		return fmt.Errorf("address %s: could not resolve port", addr)
	}
	return nil
}
