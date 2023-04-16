// Package guti contains packages
package guti

import (
	"encoding/csv"
	"math"
	"os"
	"reflect"
)

const epsilon = 1e-6

// IsExist searches for an item in a slice and returns true if it is found, and false otherwise.
// It supports searching for items of various types, including integers, floats, strings, booleans,
// and objects. It uses reflection to determine the type of the items in the slice, and to compare them
// to the search item. If the second argument is not a slice, it will panic. If the search item is not
// of the same type as the items in the slice, it will be skipped.
//
// Example usage:
//
//	intSlice := []int{1, 2, 3, 4, 5}
//	fmt.Println(guti.IsExist(3, intSlice)) // prints "true"
//	fmt.Println(guti.IsExist(6, intSlice)) // prints "false"
//
//	strSlice := []string{"foo", "bar", "baz"}
//	fmt.Println(guti.IsExist("qux", strSlice)) // prints "false"
//	fmt.Println(guti.IsExist("foo", strSlice)) // prints "true"
//
//	objectSlice := []struct {
//		Name string
//		Age  int
//	}{
//		{Name: "Alice", Age: 25},
//		{Name: "Bob", Age: 25},
//		{Name: "Charlie", Age: 35},
//	}
//	fmt.Println(guti.IsExist(struct {
//		Name string
//		Age  int
//	}{Name: "Bob", Age: 25}, objectSlice)) // prints "true"
//
//	boolSlice := []bool{true, false}
//	fmt.Println(guti.IsExist(true, boolSlice)) // prints "true"
//
//	emptySlice := []int{}
//	fmt.Println(guti.IsExist(1, emptySlice)) // prints "false"
//
// Playground: https://go.dev/play/p/jHua3iwd6xT
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

// Filter a function that filters a list based on a given predicate function. The function returns a new list with the elements that satisfy the predicate function.
func Filter(data []interface{}, predicate func(interface{}) bool) []interface{} {
	result := []interface{}{}
	for _, d := range data {
		if predicate(d) {
			result = append(result, d)
		}
	}
	return result
}

// Any a function that returns true if at least one element of a list satisfies a given predicate function.
func Any(data []interface{}, predicate func(interface{}) bool) bool {
	for _, d := range data {
		if predicate(d) {
			return true
		}
	}
	return false
}

// Reduce a function that applies a reducing function to a list and returns a
// single value. The reducing function takes two arguments, an accumulator and a
// value, and returns a new accumulator.
func Reduce(data []interface{}, reduce func(interface{}, interface{}) interface{}, initial interface{}) interface{} {
	acc := initial
	for _, d := range data {
		acc = reduce(acc, d)
	}
	return acc
}

// Map a function that applies a transformation function to each element
// of a list and returns a new list with the transformed elements.
func Map(data []interface{}, transform func(interface{}) interface{}) []interface{} {
	result := []interface{}{}
	for _, d := range data {
		result = append(result, transform(d))
	}
	return result
}

// IndexOf returns the index of the first occurrence of a given element in a list. If the element is not found, it returns -1.
func IndexOf(data []interface{}, element interface{}) int {
	for i, d := range data {
		if d == element {
			return i
		}
	}
	return -1
}

// ContainsAll returns true if all elements in the first slice are present in the second slice, otherwise returns false.
func ContainsAll(s1, s2 []interface{}) bool {
	for _, e1 := range s1 {
		found := false
		for _, e2 := range s2 {
			if e1 == e2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Reverse returns a new slice with the elements of the given slice in reverse order.
func Reverse(slice []interface{}) []interface{} {
	result := make([]interface{}, len(slice))
	for i, j := 0, len(slice)-1; i <= j; i, j = i+1, j-1 {
		result[i], result[j] = slice[j], slice[i]
	}
	return result
}

// FilterNil returns a new slice with all nil elements removed from the given slice.
func FilterNil(slice []interface{}) []interface{} {
	result := make([]interface{}, 0, len(slice))
	for _, v := range slice {
		if v != nil {
			result = append(result, v)
		}
	}
	return result
}

// MapReduce takes a slice of items and applies a mapper function to each item to get a slice of results. It then applies a reducer function to the slice of results to get a single result
func MapReduce(items interface{}, mapper func(interface{}) interface{}, reducer func(interface{}, interface{}) interface{}) interface{} {
	mappedItems := make([]interface{}, 0)
	itemsValue := reflect.ValueOf(items)

	for i := 0; i < itemsValue.Len(); i++ {
		mappedItems = append(mappedItems, mapper(itemsValue.Index(i).Interface()))
	}

	reducedResult := mappedItems[0]
	for i := 1; i < len(mappedItems); i++ {
		reducedResult = reducer(reducedResult, mappedItems[i])
	}

	return reducedResult
}

// Batch takes a slice of items and a batch size, and returns a slice of slices, where each inner slice contains at most batchSize items from the input slice
func Batch(items interface{}, batchSize int) [][]interface{} {
	var batches [][]interface{}
	itemsValue := reflect.ValueOf(items)
	batchSize = int(math.Min(float64(batchSize), float64(itemsValue.Len())))

	for i := 0; i < itemsValue.Len(); i += batchSize {
		end := int(math.Min(float64(i+batchSize), float64(itemsValue.Len())))
		batches = append(batches, ConvertSliceInterfaceToSlice(itemsValue.Slice(i, end)))
	}

	return batches
}

// ConvertSliceInterfaceToSlice takes a reflect.Value of a slice of unknown type and returns a new slice of interface{} type
func ConvertSliceInterfaceToSlice(slice reflect.Value) []interface{} {
	s := make([]interface{}, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		s[i] = slice.Index(i).Interface()
	}
	return s
}

// SaveAsCSV save data to csv
func SaveAsCSV(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// Get the type of the data slice and write the header row
	dataType := reflect.TypeOf(data).Elem()
	headerRow := make([]string, dataType.NumField())
	for i := 0; i < dataType.NumField(); i++ {
		headerRow[i] = dataType.Field(i).Name
	}
	writer.Write(headerRow)

	// Write each row of data to the CSV file
	dataValue := reflect.ValueOf(data)
	for i := 0; i < dataValue.Len(); i++ {
		row := make([]string, dataType.NumField())
		for j := 0; j < dataType.NumField(); j++ {
			fieldValue := dataValue.Index(i).Field(j)
			row[j] = fieldValue.Interface().(string)
		}
		writer.Write(row)
	}

	writer.Flush()

	return nil
}
