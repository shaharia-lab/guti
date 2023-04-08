// Package guti contains packages
package guti

import (
	"fmt"
	"sort"
	"strconv"
)

// SortStrings A function that can sort a slice of strings in ascending or descending order.
func SortStrings(slice []string, ascending bool) []string {
	if ascending {
		sort.Strings(slice)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(slice)))
	}
	return slice
}

// StringInSlice A function that can check if a given string exists in a slice of strings.
func StringInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// ToString A function that can convert any data type to a string.
func ToString(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}
