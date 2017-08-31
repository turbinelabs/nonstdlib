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
