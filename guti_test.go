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
