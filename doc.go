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

// Package nonstdlib comprises extensions to the Go stdlib, either to increase
// featureset or testability, and other utility code. Where possible, we mirror
// the stdlib package naming, though in practice, we commonly import packages
// with a tbn prefix for clarity, eg:
//
//     import (
//       "os"
//
//       tbnos "github.com/turbinelabs/nonstdlib/os"
//     )
package nonstdlib
