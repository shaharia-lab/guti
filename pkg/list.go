// Package pkg contains packages
package pkg

import (
	"math"
	"reflect"
)

const epsilon = 1e-6

// IsExist find item from a slice
func IsExist(what interface{}, in interface{}) bool {
	s := reflect.ValueOf(in)

	if s.Kind() != reflect.Slice {
		panic("IsExist: Second argument must be a slice")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Kind() != reflect.TypeOf(what).Kind() {
			continue
		}

		switch s.Index(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if s.Index(i).Int() == reflect.ValueOf(what).Int() {
				return true
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if s.Index(i).Uint() == reflect.ValueOf(what).Uint() {
				return true
			}
		case reflect.Float32, reflect.Float64:
			if math.Abs(s.Index(i).Float()-reflect.ValueOf(what).Float()) < epsilon {
				return true
			}
		case reflect.String:
			if s.Index(i).String() == reflect.ValueOf(what).String() {
				return true
			}
		case reflect.Bool:
			if s.Index(i).Bool() == reflect.ValueOf(what).Bool() {
				return true
			}
		default:
			if reflect.DeepEqual(what, s.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}
