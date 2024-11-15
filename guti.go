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
//	type MyStruct struct {
//	}
//
//	var myInt int
//	var myPtr *MyStruct
//	fmt.Println(GetTypeName(myInt)) // prints "int"
//	fmt.Println(GetTypeName(myPtr)) // prints "*MyStruct"
//
//	Playground: https://go.dev/play/p/XUAjQoGwilU
func GetTypeName(myvar interface{}) string {
	t := reflect.TypeOf(myvar)

	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}

	return t.Name()
}

// CompareStructs is a function that takes two input parameters of type interface{},
// and returns a bool indicating whether the two structs are equal or not.
//
// The function uses reflection to determine the type of each input parameter,
// and then compares their values recursively, field by field, until either a mismatch
// is found, or all fields have been compared successfully.
//
// If the input parameters are maps, the function compares each key-value pair in
// the maps recursively. If they are slices, the function compares each element in
// the slices recursively. If they are any other type, the function compares their
// values directly.
//
// The method is designed to work with any kind of struct, as long as it is represented
// as a map or a slice of interfaces. It can be used for testing, data validation,
// or any other use case where you need to compare two structs for equality.
//
// Example usage:
//
//	type People struct {
//		Age int
//	}
//	type ExampleStruct struct {
//		Name   string
//		People People
//	}
//
//	struct1 := ExampleStruct{Name: "John Doe", People: People{Age: 30}}
//	struct2 := ExampleStruct{Name: "John Doe", People: People{Age: 40}}
//	struct3 := ExampleStruct{Name: "John Doe", People: People{Age: 30}}
//
//	fmt.Println(guti.CompareStructs(struct1, struct3)) // should return true
//	fmt.Println(guti.CompareStructs(struct1, struct2)) // should return false
//
// Playground: https://go.dev/play/p/GT_7bK_BRro
func CompareStructs(s1 interface{}, s2 interface{}) bool {
	if reflect.TypeOf(s1) != reflect.TypeOf(s2) {
		return false
	}

	switch s1 := s1.(type) {
	case map[string]interface{}:
		s2 := s2.(map[string]interface{})
		if len(s1) != len(s2) {
			return false
		}
		for key := range s1 {
			if !CompareStructs(s1[key], s2[key]) {
				return false
			}
		}
		return true

	case []interface{}:
		s2 := s2.([]interface{})
		if len(s1) != len(s2) {
			return false
		}
		for i := range s1 {
			if !CompareStructs(s1[i], s2[i]) {
				return false
			}
		}
		return true

	default:
		return s1 == s2
	}
}
