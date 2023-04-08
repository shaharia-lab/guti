package gotil

import (
	"reflect"
	"testing"
)

func TestSortStrings(t *testing.T) {
	testCases := []struct {
		name      string
		input     []string
		ascending bool
		expected  []string
	}{
		{
			name:      "Sort in ascending order",
			input:     []string{"zebra", "apple", "banana", "carrot"},
			ascending: true,
			expected:  []string{"apple", "banana", "carrot", "zebra"},
		},
		{
			name:      "Sort in descending order",
			input:     []string{"zebra", "apple", "banana", "carrot"},
			ascending: false,
			expected:  []string{"zebra", "carrot", "banana", "apple"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := SortStrings(tc.input, tc.ascending)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Unexpected result:\nExpected: %v\nActual:   %v", tc.expected, actual)
			}
		})
	}
}

func TestStringInSlice(t *testing.T) {
	testCases := []struct {
		name      string
		slice     []string
		searchStr string
		expected  bool
	}{
		{
			name:      "string found in slice",
			slice:     []string{"apple", "banana", "cherry"},
			searchStr: "banana",
			expected:  true,
		},
		{
			name:      "string not found in slice",
			slice:     []string{"apple", "banana", "cherry"},
			searchStr: "mango",
			expected:  false,
		},
		{
			name:      "empty slice",
			slice:     []string{},
			searchStr: "banana",
			expected:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := StringInSlice(tc.searchStr, tc.slice)
			if actual != tc.expected {
				t.Errorf("unexpected result: expected=%v, actual=%v", tc.expected, actual)
			}
		})
	}
}

func TestToString(t *testing.T) {
	testCases := []struct {
		input    interface{}
		expected string
	}{
		{"Hello, world!", "Hello, world!"},
		{[]byte("test"), "test"},
		{42, "42"},
		{3.14, "3.14"},
		{true, "true"},
	}

	for _, tc := range testCases {
		actual := ToString(tc.input)
		if actual != tc.expected {
			t.Errorf("ToString(%v) = %s; expected %s", tc.input, actual, tc.expected)
		}
	}
}
