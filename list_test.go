package guti

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
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

func TestReduce(t *testing.T) {
	testCases := []struct {
		name       string
		data       []interface{}
		reduceFunc func(interface{}, interface{}) interface{}
		initial    interface{}
		expected   interface{}
	}{
		{
			name: "Test case 1",
			data: []interface{}{1, 2, 3, 4, 5},
			reduceFunc: func(acc interface{}, value interface{}) interface{} {
				return acc.(int) + value.(int)
			},
			initial:  0,
			expected: 15,
		},
		{
			name: "Test case 2",
			data: []interface{}{"foo", "bar", "baz"},
			reduceFunc: func(acc interface{}, value interface{}) interface{} {
				return acc.(string) + value.(string)
			},
			initial:  "",
			expected: "foobarbaz",
		},
		{
			name: "Test case 3",
			data: []interface{}{1.0, 2.5, 3.75},
			reduceFunc: func(acc interface{}, value interface{}) interface{} {
				return acc.(float64) * value.(float64)
			},
			initial:  1.0,
			expected: 9.375,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Reduce(tc.data, tc.reduceFunc, tc.initial)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("unexpected result: expected=%v, actual=%v", tc.expected, result)
			}
		})
	}
}

func TestMap(t *testing.T) {
	testCases := []struct {
		name      string
		input     []interface{}
		transform func(interface{}) interface{}
		expected  []interface{}
	}{
		{
			name: "Test case 1",
			input: []interface{}{
				1, 2, 3, 4, 5,
			},
			transform: func(d interface{}) interface{} {
				return d.(int) * 2
			},
			expected: []interface{}{
				2, 4, 6, 8, 10,
			},
		},
		{
			name: "Test case 2",
			input: []interface{}{
				"hello", "world",
			},
			transform: func(d interface{}) interface{} {
				return strings.ToUpper(d.(string))
			},
			expected: []interface{}{
				"HELLO", "WORLD",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Map(tc.input, tc.transform)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("unexpected result: expected=%v, actual=%v", tc.expected, result)
			}
		})
	}
}

func TestIndexOf(t *testing.T) {
	testCases := []struct {
		name     string
		data     []interface{}
		element  interface{}
		expected int
	}{
		{
			name:     "Test case 1",
			data:     []interface{}{1, 2, 3, 4},
			element:  2,
			expected: 1,
		},
		{
			name:     "Test case 2",
			data:     []interface{}{"hello", "world", "foo", "bar"},
			element:  "world",
			expected: 1,
		},
		{
			name:     "Test case 3",
			data:     []interface{}{1.5, 2.5, 3.5},
			element:  3.5,
			expected: 2,
		},
		{
			name:     "Test case 4",
			data:     []interface{}{"a", "b", "c", "d"},
			element:  "e",
			expected: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IndexOf(tc.data, tc.element)
			if result != tc.expected {
				t.Errorf("unexpected result: expected=%v, actual=%v", tc.expected, result)
			}
		})
	}
}

func TestContainsAll(t *testing.T) {
	testCases := []struct {
		name     string
		s1       []interface{}
		s2       []interface{}
		expected bool
	}{
		{
			name:     "Test_case_1",
			s1:       []interface{}{"a", "b", "c"},
			s2:       []interface{}{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "Test_case_2",
			s1:       []interface{}{"a", "b", "c"},
			s2:       []interface{}{"a", "b", "d"},
			expected: false,
		},
		{
			name:     "Test_case_3",
			s1:       []interface{}{"a", "b", "c"},
			s2:       []interface{}{"a", "b"},
			expected: false,
		},
		{
			name:     "Test_case_4",
			s1:       []interface{}{"a", "b", "c"},
			s2:       []interface{}{"a", "b", "b", "c"},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ContainsAll(tc.s1, tc.s2)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestReverse(t *testing.T) {
	testCases := []struct {
		name     string
		input    []interface{}
		expected []interface{}
	}{
		{
			name:     "Test case 1",
			input:    []interface{}{1, 2, 3, 4, 5},
			expected: []interface{}{5, 4, 3, 2, 1},
		},
		{
			name:     "Test case 2",
			input:    []interface{}{"a", "b", "c", "d", "e"},
			expected: []interface{}{"e", "d", "c", "b", "a"},
		},
		{
			name:     "Test case 3",
			input:    []interface{}{true, false},
			expected: []interface{}{false, true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Reverse(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("unexpected result: expected=%v, actual=%v", tc.expected, result)
			}
		})
	}
}

func TestFilterNil(t *testing.T) {
	testCases := []struct {
		name     string
		input    []interface{}
		expected []interface{}
	}{
		{
			name:     "Test case 1",
			input:    []interface{}{"a", nil, "b", nil},
			expected: []interface{}{"a", "b"},
		},
		{
			name:     "Test case 2",
			input:    []interface{}{1, 2, nil, 4, nil},
			expected: []interface{}{1, 2, 4},
		},
		{
			name:     "Test case 3",
			input:    []interface{}{nil, nil, nil},
			expected: []interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FilterNil(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("unexpected result: expected=%v, actual=%v", tc.expected, result)
			}
		})
	}
}

func TestMapReduce(t *testing.T) {
	testCases := []struct {
		name     string
		items    interface{}
		mapper   func(interface{}) interface{}
		reducer  func(interface{}, interface{}) interface{}
		expected interface{}
	}{
		{
			name:  "Test case 1",
			items: []int{1, 2, 3, 4, 5},
			mapper: func(item interface{}) interface{} {
				return item.(int) * 2
			},
			reducer: func(acc interface{}, item interface{}) interface{} {
				return acc.(int) + item.(int)
			},
			expected: 30,
		},
		{
			name:  "Test case 2",
			items: []string{"hello", "world"},
			mapper: func(item interface{}) interface{} {
				return len(item.(string))
			},
			reducer: func(acc interface{}, item interface{}) interface{} {
				return acc.(int) + item.(int)
			},
			expected: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MapReduce(tc.items, tc.mapper, tc.reducer)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("unexpected result: expected=%v, actual=%v", tc.expected, result)
			}
		})
	}
}

func TestBatch(t *testing.T) {
	testCases := []struct {
		name      string
		items     interface{}
		batchSize int
		expected  [][]interface{}
	}{
		{
			name:      "Test case 1",
			items:     []int{1, 2, 3, 4, 5, 6},
			batchSize: 2,
			expected: [][]interface{}{
				{1, 2},
				{3, 4},
				{5, 6},
			},
		},
		{
			name:      "Test case 2",
			items:     []string{"a", "b", "c", "d", "e", "f", "g", "h"},
			batchSize: 3,
			expected: [][]interface{}{
				{"a", "b", "c"},
				{"d", "e", "f"},
				{"g", "h"},
			},
		},
		{
			name:      "Test case 3",
			items:     []float64{1.2, 2.4, 3.6, 4.8},
			batchSize: 3,
			expected: [][]interface{}{
				{1.2, 2.4, 3.6},
				{4.8},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Batch(tc.items, tc.batchSize)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("unexpected result: expected=%v, actual=%v", tc.expected, result)
			}
		})
	}
}

func TestConvertSliceInterfaceToSlice(t *testing.T) {
	testCases := []struct {
		name     string
		input    reflect.Value
		expected []interface{}
	}{
		{
			name:     "Test case 1",
			input:    reflect.ValueOf([]string{"foo", "bar", "baz"}),
			expected: []interface{}{"foo", "bar", "baz"},
		},
		{
			name:     "Test case 2",
			input:    reflect.ValueOf([]int{1, 2, 3}),
			expected: []interface{}{1, 2, 3},
		},
		{
			name:     "Test case 3",
			input:    reflect.ValueOf([]interface{}{"foo", 42, true}),
			expected: []interface{}{"foo", 42, true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertSliceInterfaceToSlice(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("unexpected result: expected=%v, actual=%v", tc.expected, result)
			}
		})
	}
}

type TestData struct {
	Name  string
	Email string
	Phone string
}

func TestSaveAsCSV(t *testing.T) {
	testCases := []struct {
		name     string
		data     interface{}
		filename string
	}{
		{
			name:     "Test case 1",
			data:     []TestData{{"John Doe", "john.doe@example.com", "123-456-7890"}, {"Jane Doe", "jane.doe@example.com", "987-654-3210"}},
			filename: "testdata/testdata.csv",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SaveAsCSV(tc.data, tc.filename)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify that the file was created and has the expected contents
			fileContents, err := ioutil.ReadFile(tc.filename)
			if err != nil {
				t.Errorf("unexpected error reading file: %v", err)
			}
			expectedContents := "Name,Email,Phone\nJohn Doe,john.doe@example.com,123-456-7890\nJane Doe,jane.doe@example.com,987-654-3210\n"
			if string(fileContents) != expectedContents {
				t.Errorf("unexpected file contents: expected=%q, actual=%q", expectedContents, string(fileContents))
			}

			// Clean up the test file
			os.Remove(tc.filename)
		})
	}
}
