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

// Package time provides utility functions for go time.Time instances.
package time

import (
	"time"
)

const (
	// TurbineFormat is the canonical format we want to display
	// timestamps in. All timestamps within Turbine code should be
	// UTC.
	TurbineFormat = "2006-01-02 15:04:05.000"

	microsPerSecond = int64(1000000)
	microsPerNano   = int64(time.Microsecond)
	millisPerSecond = int64(1000)
	millisPerNano   = int64(time.Millisecond)
)

// Equal compares pointers to two times for equality. If both pointers are
// nil or if they both point to a time that is equivalent then they are equal.
func Equal(t, o *time.Time) bool {
	switch {
	case t == nil && o == nil:
		return true
	case t != nil && o != nil:
		return t.Equal(*o)
	default:
		return false
	}
}

// Format uses the canonical Turbine time format and converts the *time.Time
// to a string using it. Returns empty string if t is nil. If t is not nil we
// also will force the timezone to be UTC before rendering to string.
func Format(t *time.Time) string {
	if t == nil {
		return ""
	}

	return t.UTC().Format(TurbineFormat)
}

// Parse takes a timestamp in the Turbine canonical format and returns
// a pointer to the represented time. An empty string produces a nil
// pointer; an error parsing the input produces a nil pointer and an
// error.
func Parse(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}

	t, err := time.Parse(TurbineFormat, s)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// ToUnixMilli returns a Unix timestamp in milliseconds from "January 1, 1970 UTC".
// The result is undefined if the Unix time cannot be represented by an int64.
// For example, calling ToUnixMilli on a zero Time is undefined.
//
// This utility is useful for service API's such as AWS CloudWatch Logs which require
// their unix time values to be in milliseconds.
//
// See Go stdlib https://golang.org/pkg/time/#Time.UnixNano for more information.
func ToUnixMilli(t time.Time) int64 {
	return t.UnixNano() / int64(millisPerNano)
}

// FromUnixMilli returns a time.Time from the given milliseconds from
// "January 1, 1970 UTC". The timezone of the returns Time is UTC.
func FromUnixMilli(millis int64) time.Time {
	return time.Unix(millis/millisPerSecond, millis%millisPerSecond*millisPerNano).In(time.UTC)
}

// TruncUnixMilli truncates the given time to millisecond resolution.
func TruncUnixMilli(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	return FromUnixMilli(ToUnixMilli(t))
}

// ToUnixMicro returns a Unix timestamp in microseconds from "January 1, 1970 UTC".
// The result is undefined if the Unix time cannot be represented by an int64.
// For example, calling ToUnixMicro on a zero Time is undefined.
//
// See Go stdlib https://golang.org/pkg/time/#Time.UnixNano for more information.
func ToUnixMicro(t time.Time) int64 {
	return t.UnixNano() / int64(microsPerNano)
}

// FromUnixMicro returns a time.Time from the given microseconds from
// "January 1, 1970 UTC". The timezone of the returns Time is UTC.
func FromUnixMicro(micros int64) time.Time {
	return time.Unix(micros/microsPerSecond, micros%microsPerSecond*microsPerNano).In(time.UTC)
}

// TruncUnixMicro truncates the given time to microsecond resolution.
func TruncUnixMicro(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return FromUnixMicro(ToUnixMicro(t))
}

// Min returns the earliest of several time.Time instances.
func Min(a time.Time, bs ...time.Time) time.Time {
	ts := append([]time.Time{a}, bs...)
	ans := time.Time{}

	for _, t := range ts {
		if ans.IsZero() {
			ans = t
		} else if !ans.Before(t) && !t.IsZero() {
			ans = t
		}
	}

	return ans
}

// Max returns the latest of several time.Time instances.
func Max(a time.Time, bs ...time.Time) time.Time {
	ans := a

	for _, t := range bs {
		if ans.Before(t) {
			ans = t
		}
	}

	return ans
}
