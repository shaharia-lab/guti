// Package guti contains utility functions
package guti

import "reflect"

// GetTypeName returns the name of the type of the given object.
// If the object is a pointer, it returns the name of the pointed-to type with a
// leading asterisk (*).
//
// The method uses reflection to determine the type of the input object, and then
// returns the name of the type as a string. If the input object is a pointer, it
// returns the name of the pointed-to type with a leading asterisk (*).
//
// This method can be useful in situations where you need to get the name of the
// type of an object dynamically, such as in debugging, logging, or error handling.
//
// Example usage:
//
//	var myInt int
//	var myPtr *MyStruct
//	fmt.Println(GetTypeName(myInt))   // prints "int"
//	fmt.Println(GetTypeName(myPtr))   // prints "*MyStruct"
func GetTypeName(myvar interface{}) string {
	t := reflect.TypeOf(myvar)

	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}

	return t.Name()
}
