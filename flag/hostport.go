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
	"fmt"
	"net"
	"strings"
)

// NewHostPort constructs a HostPort from the given string. If the
// given argument is not a valid HostPort, a HostPort matching ":0" is
// returned.
func NewHostPort(s string) HostPort {
	hp := HostPort{}
	hp.Set(s)
	return hp
}

// NewHostPortWithDefaultPort constructs a HostPort from the given
// host:port and a port lookup function. This makes the port portion
// of the host:port string optional both in this constructor and in
// command line flags. The port lookup function is used to resolve the
// port during a call Addr(), String(), or ParsedHostPort() unless
// command line flags have specified an explicit port name or number.
func NewHostPortWithDefaultPort(s string, getPort func() int) HostPort {
	hp := HostPort{portLookupFunc: getPort}
	hp.Set(s)
	hp.portLookupFunc = getPort
	return hp
}

// HostPort represents the value of a TCP host:port flag.
type HostPort struct {
	host           string
	port           string
	numericPort    int
	hadBrackets    bool
	portLookupFunc func() int
}

var _ flag.Getter = &HostPort{}

func (hp *HostPort) handlePortLookup() {
	if hp.portLookupFunc != nil {
		p := hp.portLookupFunc()
		hp.port = fmt.Sprintf("%d", p)
		hp.numericPort = p
	}
}

func (hp *HostPort) String() string {
	if hp == nil {
		return ":0"
	}

	hp.handlePortLookup()

	if hp.host == "" && hp.port == "" {
		return ":0"
	}

	host := hp.host
	if hp.hadBrackets {
		host = fmt.Sprintf("[%s]", hp.host)
	}

	if hp.port == "" {
		return fmt.Sprintf("%s:0", host)
	}

	return fmt.Sprintf("%s:%s", host, hp.port)
}

// Addr returns a "host:port" string suitable for use by net.Dial or net.Listen.
func (hp *HostPort) Addr() string { return hp.String() }

// ParsedHostPort returns the host string (including brackets if
// they were originally present), and a numeric port.
func (hp *HostPort) ParsedHostPort() (string, int) {
	if hp == nil {
		return "", 0
	}
	hp.handlePortLookup()
	return hp.host, hp.numericPort
}

// Get implements flag.Getter.
func (hp *HostPort) Get() interface{} { return hp }

// Set implements flag.Value.
func (hp *HostPort) Set(s string) error {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") && hp.portLookupFunc != nil {
			h, hadBrackets := stripBrackets(s)
			*hp = HostPort{host: h, hadBrackets: hadBrackets, portLookupFunc: hp.portLookupFunc}
			return nil
		}

		return err
	}

	numericPort, err := net.LookupPort("tcp", port)
	if err != nil {
		return err
	}

	*hp = HostPort{host: host, port: port, numericPort: numericPort, hadBrackets: s[0] == '['}
	return nil
}

func stripBrackets(host string) (string, bool) {
	if len(host) > 0 && host[0] == '[' && strings.HasSuffix(host, "]") {
		return host[1 : len(host)-1], true
	}
	return host, false
}
