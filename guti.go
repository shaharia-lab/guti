// Package guti contains utility functions
package guti

import "reflect"

// GetTypeName returns the name of the type of the given object.
// If the object is a pointer, it returns the name of the pointed-to type.
func GetTypeName(myvar interface{}) string {
	t := reflect.TypeOf(myvar)

	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}

	return t.Name()
}
