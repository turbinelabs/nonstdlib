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

// Package ptr provides convenience and conversion methods for working
// with pointer types. This package was initially forked/lifted from
// https://github.com/aws/aws-sdk-go/blob/v1.1.35/aws/convert_types.go
// and has since drifted to include a few other conversions.
package ptr

import "time"

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// StringSlice converts a slice of string values into a slice of
// string pointers
func StringSlice(src []string) []*string {
	dst := make([]*string, len(src))
	for i := range src {
		dst[i] = &(src[i])
	}
	return dst
}

// StringValueSlice converts a slice of string pointers into a slice of
// string values
func StringValueSlice(src []*string) []string {
	dst := make([]string, len(src))
	for i := range src {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}
	return dst
}

// StringMap converts a string map of string values into a string
// map of string pointers
func StringMap(src map[string]string) map[string]*string {
	dst := make(map[string]*string)
	for k, val := range src {
		v := val
		dst[k] = &v
	}
	return dst
}

// StringValueMap converts a string map of string pointers into a string
// map of string values
func StringValueMap(src map[string]*string) map[string]string {
	dst := make(map[string]string)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}
	return dst
}

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}

// BoolValue returns the value of the bool pointer passed in or
// false if the pointer is nil.
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

// BoolSlice converts a slice of bool values into a slice of
// bool pointers
func BoolSlice(src []bool) []*bool {
	dst := make([]*bool, len(src))
	for i := range src {
		dst[i] = &(src[i])
	}
	return dst
}

// BoolValueSlice converts a slice of bool pointers into a slice of
// bool values
func BoolValueSlice(src []*bool) []bool {
	dst := make([]bool, len(src))
	for i := range src {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}
	return dst
}

// BoolMap converts a string map of bool values into a string
// map of bool pointers
func BoolMap(src map[string]bool) map[string]*bool {
	dst := make(map[string]*bool)
	for k, val := range src {
		v := val
		dst[k] = &v
	}
	return dst
}

// BoolValueMap converts a string map of bool pointers into a string
// map of bool values
func BoolValueMap(src map[string]*bool) map[string]bool {
	dst := make(map[string]bool)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}
	return dst
}

// Int returns a pointer to the int value passed in.
func Int(v int) *int {
	return &v
}

// IntValue returns the value of the int pointer passed in or
// 0 if the pointer is nil.
func IntValue(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

// IntSlice converts a slice of int values into a slice of
// int pointers
func IntSlice(src []int) []*int {
	dst := make([]*int, len(src))
	for i := range src {
		dst[i] = &(src[i])
	}
	return dst
}

// IntValueSlice converts a slice of int pointers into a slice of
// int values
func IntValueSlice(src []*int) []int {
	dst := make([]int, len(src))
	for i := range src {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}
	return dst
}

// IntMap converts a string map of int values into a string
// map of int pointers
func IntMap(src map[string]int) map[string]*int {
	dst := make(map[string]*int)
	for k, val := range src {
		v := val
		dst[k] = &v
	}
	return dst
}

// IntValueMap converts a string map of int pointers into a string
// map of int values
func IntValueMap(src map[string]*int) map[string]int {
	dst := make(map[string]int)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}
	return dst
}

// Int32 returns a pointer to the int32 value passed in.
func Int32(v int32) *int32 {
	return &v
}

// Int32Value returns the value of the int32 pointer passed in or
// 0 if the pointer is nil.
func Int32Value(v *int32) int32 {
	if v != nil {
		return *v
	}
	return 0
}

// Int64 returns a pointer to the int64 value passed in.
func Int64(v int64) *int64 {
	return &v
}

// Int64Value returns the value of the int64 pointer passed in or
// 0 if the pointer is nil.
func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

// Int64Slice converts a slice of int64 values into a slice of
// int64 pointers
func Int64Slice(src []int64) []*int64 {
	dst := make([]*int64, len(src))
	for i := range src {
		dst[i] = &(src[i])
	}
	return dst
}

// Int64ValueSlice converts a slice of int64 pointers into a slice of
// int64 values
func Int64ValueSlice(src []*int64) []int64 {
	dst := make([]int64, len(src))
	for i := range src {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}
	return dst
}

// Int64Map converts a string map of int64 values into a string
// map of int64 pointers
func Int64Map(src map[string]int64) map[string]*int64 {
	dst := make(map[string]*int64)
	for k, val := range src {
		v := val
		dst[k] = &v
	}
	return dst
}

// Int64ValueMap converts a string map of int64 pointers into a string
// map of int64 values
func Int64ValueMap(src map[string]*int64) map[string]int64 {
	dst := make(map[string]int64)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}
	return dst
}

// Uint returns a pointer to the uint value passed in.
func Uint(u uint) *uint {
	return &u
}

// Uint64 returns a pointer to the uint64 value passed in.
func Uint64(v uint64) *uint64 {
	return &v
}

// Uint64Value returns the value of the uint64 pointer passed in or
// 0 if the pointer is nil.
func Uint64Value(v *uint64) uint64 {
	if v != nil {
		return *v
	}
	return 0
}

// Float64 returns a pointer to the float64 value passed in.
func Float64(v float64) *float64 {
	return &v
}

// Float64Value returns the value of the float64 pointer passed in or
// 0 if the pointer is nil.
func Float64Value(v *float64) float64 {
	if v != nil {
		return *v
	}
	return 0
}

// Float64Slice converts a slice of float64 values into a slice of
// float64 pointers
func Float64Slice(src []float64) []*float64 {
	dst := make([]*float64, len(src))
	for i := range src {
		dst[i] = &(src[i])
	}
	return dst
}

// Float64ValueSlice converts a slice of float64 pointers into a slice of
// float64 values
func Float64ValueSlice(src []*float64) []float64 {
	dst := make([]float64, len(src))
	for i := range src {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}
	return dst
}

// Float64Map converts a string map of float64 values into a string
// map of float64 pointers
func Float64Map(src map[string]float64) map[string]*float64 {
	dst := make(map[string]*float64)
	for k, val := range src {
		v := val
		dst[k] = &v
	}
	return dst
}

// Float64ValueMap converts a string map of float64 pointers into a string
// map of float64 values
func Float64ValueMap(src map[string]*float64) map[string]float64 {
	dst := make(map[string]float64)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}
	return dst
}

// Time returns a pointer to the time.Time value passed in.
func Time(v time.Time) *time.Time {
	return &v
}

// TimeValue returns the value of the time.Time pointer passed in or
// time.Time{} if the pointer is nil.
func TimeValue(v *time.Time) time.Time {
	if v != nil {
		return *v
	}
	return time.Time{}
}

// TimeSlice converts a slice of time.Time values into a slice of
// time.Time pointers
func TimeSlice(src []time.Time) []*time.Time {
	dst := make([]*time.Time, len(src))
	for i := range src {
		dst[i] = &(src[i])
	}
	return dst
}

// TimeValueSlice converts a slice of time.Time pointers into a slice of
// time.Time values
func TimeValueSlice(src []*time.Time) []time.Time {
	dst := make([]time.Time, len(src))
	for i := range src {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}
	return dst
}

// TimeMap converts a string map of time.Time values into a string
// map of time.Time pointers
func TimeMap(src map[string]time.Time) map[string]*time.Time {
	dst := make(map[string]*time.Time)
	for k, val := range src {
		v := val
		dst[k] = &v
	}
	return dst
}

// TimeValueMap converts a string map of time.Time pointers into a string
// map of time.Time values
func TimeValueMap(src map[string]*time.Time) map[string]time.Time {
	dst := make(map[string]time.Time)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}
	return dst
}

// Duration returns a pointer to the time.Duration value passed in.
func Duration(v time.Duration) *time.Duration {
	return &v
}

// DurationValue returns the value of the time.Duration pointer passed in or
// 0 if the pointer is nil.
func DurationValue(v *time.Duration) time.Duration {
	if v != nil {
		return *v
	}
	return 0
}

// BoolEqual compares two bool pointer values
func BoolEqual(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return *a == *b
}

// StringEqual compares two string pointer values
func StringEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return *a == *b
}

// IntEqual compares two int pointer values
func IntEqual(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return *a == *b
}
