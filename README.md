
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
[![GoDoc](https://https://godoc.org/github.com/turbinelabs/nonstdlib?status.svg)](https://https://godoc.org/github.com/turbinelabs/nonstdlib)
[![CircleCI](https://circleci.com/gh/turbinelabs/nonstdlib.svg?style=shield)](https://circleci.com/gh/turbinelabs/nonstdlib)

The nonstdlib project comprises extensions to the Go stdlib, either to increase
feature set or testability, and other utility code. The nonstdlib project has
no external dependencies beyond the go standard library; the tests depend on
our [test package](https://github.com/turbinelabs/test) and
[gomock](https://github.com/golang/mock).

Where possible, we mirror the stdlib package naming, though in practice, we
commonly import packages with a tbn prefix for clarity, eg:

```go
import (
  "os"

  tbnos "github.com/turbinelabs/nonstdlib/os"
)
```

## Requirements

- Go 1.7.4 or later (previous versions may work, but we don't build or test against them)

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
- [`os`](https://godoc.org/github.com/turbinelabs/nonstdlib/os):
  provides an OS interface mirroring a subset of commonly used functions and
  variables from the golang os package
- [`proc`](https://godoc.org/github.com/turbinelabs/nonstdlib/proc):
  provides a mechanism for running processes under management
- [`ptr`](https://godoc.org/github.com/turbinelabs/nonstdlib/ptr):
  provides convenience and conversion methods for working with pointer types
- [`stats`](https://godoc.org/github.com/turbinelabs/nonstdlib/stats):
  provides an interface for reporting simple statistics
- [`time`](https://godoc.org/github.com/turbinelabs/nonstdlib/time):
  provides utility functions for go time.Time instances

## Versioning

Please see [Versioning of Turbine Labs Open Source Projects](http://github.com/turbinelabs/developer/blob/master/README.md#versioning).

## Pull Requests

Patches accepted! Please see [Contributing to Turbine Labs Open Source Projects](http://github.com/turbinelabs/developer/blob/master/README.md#contributing).

## Code of Conduct

All Turbine Labs open-sourced projects are released with a
[Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in our
projects you agree to abide by its terms, which will be vigorously enforced.
