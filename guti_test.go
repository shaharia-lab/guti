package guti

import "testing"

func TestGetTypeName(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	testCases := []struct {
		input    interface{}
		expected string
	}{
		{1, "int"},
		{3.14, "float64"},
		{"hello", "string"},
		{true, "bool"},
		{&Person{Name: "Alice", Age: 30}, "*Person"},
	}

	for _, tc := range testCases {
		actual := GetTypeName(tc.input)
		if actual != tc.expected {
			t.Errorf("GetTypeName(%v) = %s; expected %s", tc.input, actual, tc.expected)
		}
	}
}

func TestCompareStructs(t *testing.T) {
	testCases := []struct {
		name     string
		s1       interface{}
		s2       interface{}
		expected bool
	}{
		{
			"equal maps",
			map[string]interface{}{"a": 1, "b": "two", "c": true},
			map[string]interface{}{"a": 1, "b": "two", "c": true},
			true,
		},
		{
			"Unequal maps",
			map[string]interface{}{"a": 1, "b": "two", "c": true},
			map[string]interface{}{"a": 1, "b": "two", "c": false},
			false,
		},
		{
			"Equal slices",
			[]interface{}{1, "two", true},
			[]interface{}{1, "two", true},
			true,
		},
		{
			"Unequal slices",
			[]interface{}{1, "two", true},
			[]interface{}{1, "two", false},
			false,
		},
		{
			"Different types",
			map[string]interface{}{"a": 1, "b": "two", "c": true},
			[]interface{}{1, "two", true},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CompareStructs(tc.s1, tc.s2)
			if result != tc.expected {
				t.Errorf("CompareStructs(%v, %v) = %v, expected %v",
					tc.s1, tc.s2, result, tc.expected)
			}
		})
	}
}
