package gotil

import (
	"reflect"
	"testing"
)

func TestIsExist(t *testing.T) {
	tests := []struct {
		name     string
		in       interface{}
		what     interface{}
		expected bool
	}{
		{
			name:     "int_in_int_slice",
			in:       []int{1, 2, 3, 4, 5},
			what:     3,
			expected: true,
		},
		{
			name:     "string_in_string_slice",
			in:       []string{"foo", "bar", "baz"},
			what:     "baz",
			expected: true,
		},
		{
			name:     "float_in_float_slice",
			in:       []float64{1.5, 2.5, 3.5},
			what:     2.5,
			expected: true,
		},
		{
			name: "object_in_object_slice",
			in: []struct {
				Name string
				Age  int
			}{
				{Name: "Alice", Age: 25},
				{Name: "Bob", Age: 30},
				{Name: "Charlie", Age: 35},
			},
			what: struct {
				Name string
				Age  int
			}{Name: "Bob", Age: 30},
			expected: true,
		},
		{
			name:     "int_not_in_int_slice",
			in:       []int{1, 2, 3, 4, 5},
			what:     6,
			expected: false,
		},
		{
			name:     "uint in uint slice",
			what:     uint(2),
			in:       []uint{1, 2, 3},
			expected: true,
		},
		{
			name:     "uint not in uint slice",
			what:     uint(4),
			in:       []uint{1, 2, 3},
			expected: false,
		},
		{
			name:     "bool in bool slice",
			what:     true,
			in:       []bool{true, false},
			expected: true,
		},
		{
			name:     "bool not in bool slice",
			what:     true,
			in:       []bool{false},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExist(tt.what, tt.in); got != tt.expected {
				t.Errorf("IsExist(%v, %v) = %v, want %v", tt.what, tt.in, got, tt.expected)
			}
		})
	}
}

func TestIsExistSecondArgNotSlice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("IsExist did not panic")
		}
	}()

	IsExist("a", 123)
}

func TestIsExist_TypeMismatch(t *testing.T) {
	var inputSlice = []interface{}{1, 2, 3, 4, 5, "6"}

	// Search for a string in a slice of integers
	searchItem := "6"
	expected := false

	result := IsExist(searchItem, inputSlice)
	if result != false {
		t.Errorf("IsExist(%v, %v) = %v; want %v", searchItem, inputSlice, result, expected)
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []interface{}
		predicate func(interface{}) bool
		want      []interface{}
	}{
		{
			name:      "filter ints",
			input:     []interface{}{1, 2, 3, 4, 5},
			predicate: func(x interface{}) bool { return x.(int)%2 == 0 },
			want:      []interface{}{2, 4},
		},
		{
			name:      "filter strings",
			input:     []interface{}{"hello", "world", "foo", "bar"},
			predicate: func(x interface{}) bool { return len(x.(string)) > 3 },
			want:      []interface{}{"hello", "world"},
		},
		{
			name:      "filter empty list",
			input:     []interface{}{},
			predicate: func(x interface{}) bool { return true },
			want:      []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.input, tt.predicate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		predicate func(interface{}) bool
		expected  bool
	}{
		{
			name:      "any item matches predicate",
			input:     []int{1, 2, 3},
			predicate: func(item interface{}) bool { return item.(int) == 2 },
			expected:  true,
		},
		{
			name:      "no item matches predicate",
			input:     []int{1, 2, 3},
			predicate: func(item interface{}) bool { return item.(int) == 4 },
			expected:  false,
		},
		{
			name:      "empty input",
			input:     []int{},
			predicate: func(item interface{}) bool { return item.(int) == 2 },
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Any(toInterfaceSlice(tt.input), tt.predicate)
			if result != tt.expected {
				t.Errorf("Any(%v, predicate) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func toInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("toInterfaceSlice: not a slice")
	}
	result := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		result[i] = s.Index(i).Interface()
	}
	return result
}
