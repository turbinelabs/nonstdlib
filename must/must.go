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

// Package must is a convenience wrapper for extracting useful information out of
// (data, error) tuples when you really just want the data and exiting on errors
// is acceptable.
package must

import (
	"database/sql"

	"github.com/turbinelabs/nonstdlib/log/console"
)

func die(e error) {
	if e != nil {
		console.Error().Fatal(e)
	}
}

// String returns the string s or logs a fatal error e.
func String(s string, e error) string {
	die(e)
	return s
}

// Int returns the int i or logs a fatal error e.
func Int(i int, e error) int {
	die(e)
	return i
}

// Int64 returns the int64 i or logs a fatal error e.
func Int64(i int64, e error) int64 {
	die(e)
	return i
}

// Uint returns the uint u or logs a fatal error e.
func Uint(u uint, e error) uint {
	die(e)
	return u
}

// Uint64 returns the uint u or logs a fatal error e.
func Uint64(u uint64, e error) uint64 {
	die(e)
	return u
}

// Float64 returns the float f or logs a fatal error e.
func Float64(f float64, e error) float64 {
	die(e)
	return f
}

// Any returns the interface{} i or logs a fatal error e.
func Any(i interface{}, e error) interface{} {
	die(e)
	return i
}

// Rows returns the *sql.Rows r or logs a fatal error e.
func Rows(r *sql.Rows, e error) *sql.Rows {
	die(e)
	return r
}

// Result returns the *sql.Result r or logs a fatal error e.
func Result(r *sql.Result, e error) *sql.Result {
	die(e)
	return r
}

// Work will log a fatal error if e != nil.
func Work(e error) {
	die(e)
}
