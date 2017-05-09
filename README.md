
[//]: # ( Copyright 2017 Turbine Labs, Inc.                                   )
[//]: # ( you may not use this file except in compliance with the License.    )
[//]: # ( You may obtain a copy of the License at                             )
[//]: # (                                                                     )
[//]: # (     http://www.apache.org/licenses/LICENSE-2.0                      )
[//]: # (                                                                     )
[//]: # ( Unless required by applicable law or agreed to in writing, software )
[//]: # ( distributed under the License is distributed on an "AS IS" BASIS,   )
[//]: # ( WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or     )
[//]: # ( implied. See the License for the specific language governing        )
[//]: # ( permissions and limitations under the License.                      )

# turbinelabs/nonstdlib

[![Apache 2.0](https://img.shields.io/hexpm/l/plug.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/turbinelabs/nonstdlib?status.svg)](https://godoc.org/github.com/turbinelabs/nonstdlib)
[![CircleCI](https://circleci.com/gh/turbinelabs/nonstdlib.svg?style=shield)](https://circleci.com/gh/turbinelabs/nonstdlib)
[![Go Report Card](https://goreportcard.com/badge/github.com/turbinelabs/nonstdlib)](https://goreportcard.com/report/github.com/turbinelabs/nonstdlib)

The nonstdlib project comprises extensions to the Go stdlib, either to increase
feature set or testability, and other utility code. The nonstdlib project has
no external dependencies beyond the go standard library; the tests depend on
our [test package](https://github.com/turbinelabs/test) and
[gomock](https://github.com/golang/mock).

Where possible, we mirror the stdlib package naming, though in practice, we
commonly import packages with a tbn prefix for clarity, e.g.:

```go
import (
  "os"

  tbnos "github.com/turbinelabs/nonstdlib/os"
)
```

## Requirements

- Go 1.8 or later (previous versions may work, but we don't build or test against them)

## Install

```
go get -u github.com/turbinelabs/nonstdlib/...
```

## Clone/Test

```
mkdir -p $GOPATH/src/turbinelabs
git clone https://github.com/turbinelabs/nonstdlib.git > $GOPATH/src/turbinelabs/nonstdlib
go test github.com/turbinelabs/nonstdlib/...
```

## Packages

Each package is best described in its respective Godoc:

- [`arrays`](https://godoc.org/github.com/turbinelabs/nonstdlib/arrays):
  includes several sub-packages allowing type-safe execution of tasks commonly
  applied to slices
- [`editor`](https://godoc.org/github.com/turbinelabs/nonstdlib/editor):
  provides simple wrappers for interacting with an environment configured text
  editor
- [`executor`](https://godoc.org/github.com/turbinelabs/nonstdlib/executor):
  provides a mechanism for asyncronous execution of tasks, using callbacks to
  indicate success or failure
- [`flag`](https://godoc.org/github.com/turbinelabs/nonstdlib/flag):
  provides convenience methods for dealing with golang flag.FlagSets
- [`log`](https://godoc.org/github.com/turbinelabs/nonstdlib/log):
  provides infrastructure for topic based logging to files
- [`math`](https://godoc.org/github.com/turbinelabs/nonstdlib/math):
  provides mathematical utilities
- [`must`](https://godoc.org/github.com/turbinelabs/nonstdlib/must):
  provides extraction of useful information out of (data, error) tuples
- [`net`](https://godoc.org/github.com/turbinelabs/nonstdlib/net):
  provides convenience methods for dealing with the net package of the stdlib
- [`os`](https://godoc.org/github.com/turbinelabs/nonstdlib/os):
  provides an OS interface mirroring a subset of commonly used functions and
  variables from the golang os package
- [`proc`](https://godoc.org/github.com/turbinelabs/nonstdlib/proc):
  provides a mechanism for running processes under management
- [`ptr`](https://godoc.org/github.com/turbinelabs/nonstdlib/ptr):
  provides convenience and conversion methods for working with pointer types
- [`stats`](https://godoc.org/github.com/turbinelabs/nonstdlib/stats):
  provides an interface for reporting simple statistics
- [`strings`](https://godoc.org/github.com/turbinelabs/nonstdlib/strings):
  provides convenience methods for working with strings and string slices
- [`text/tabwriter`](https://godoc.org/github.com/turbinelabs/nonstdlib/text/tabwriter):
   provides a set of sane defaults for converting tab separated values into a
   pretty column formatted output.
- [`time`](https://godoc.org/github.com/turbinelabs/nonstdlib/time):
  provides utility functions for go time.Time instances

## Versioning

Please see [Versioning of Turbine Labs Open Source Projects](http://github.com/turbinelabs/developer/blob/master/README.md#versioning).

## Pull Requests

Patches accepted! Please see [Contributing to Turbine Labs Open Source Projects](http://github.com/turbinelabs/developer/blob/master/README.md#contributing).

## Code of Conduct

All Turbine Labs open-sourced projects are released with a
[Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in our
projects you agree to abide by its terms, which will be carefully enforced.
