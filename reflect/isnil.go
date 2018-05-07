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

package reflect

import "reflect"

// IsNil returns true if the parameter is:
//   - untyped nil (<nil, nil>)
//   - a nil value with some fixed type, e.g., (<io.Reader, nil>)
//
// It differs from '== nil' in that a pointer to some type will still be
// considered nil if it is being passed to an interface.
//
// Because this uses reflection it is much slower than a direct comparison to
// nil but is generally <100 microseconds per check. In a tight loop with a
// potentially very large number of iterations this may be a consideration.
func IsNil(i interface{}) bool {
	v := reflect.ValueOf(i)
	return (v.Kind() == reflect.Ptr && v.IsNil()) || v.Kind() == reflect.Invalid
}
