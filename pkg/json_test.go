package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestJSONToMap(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected map[string]interface{}
		hasError bool
	}{
		{
			input: []byte(`{
				"name": "John Doe",
				"age": 30,
				"city": "New York"
			}`),
			expected: map[string]interface{}{
				"name": "John Doe",
				"age":  30.0,
				"city": "New York",
			},
			hasError: false,
		},
		{
			input:    []byte(`invalid-json`),
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			actual, err := JSONToMap(tc.input)

			if tc.hasError {
				require.Error(t, err)
				require.Nil(t, actual)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}
func TestJSONToString(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "simple object",
			input:    map[string]interface{}{"name": "John", "age": 30},
			expected: `{"age":30,"name":"John"}`,
		},
		{
			name:     "nested object",
			input:    map[string]interface{}{"person": map[string]interface{}{"name": "John", "age": 30}},
			expected: `{"person":{"age":30,"name":"John"}}`,
		},
		{
			name:     "array of objects",
			input:    []interface{}{map[string]interface{}{"name": "John", "age": 30}, map[string]interface{}{"name": "Jane", "age": 25}},
			expected: `[{"age":30,"name":"John"},{"age":25,"name":"Jane"}]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := JSONToString(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Expected '%v', but got '%v'", tc.expected, actual)
			}
		})
	}
}

func TestJSONFileToMap(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		expected map[string]interface{}
	}{
		{
			name:     "Test case 1",
			filename: "testdata/file1.json",
			expected: map[string]interface{}{
				"name": "John Doe",
				"job":  "engineer",
			},
		},
		{
			name:     "Test case 2",
			filename: "testdata/file2.json",
			expected: map[string]interface{}{
				"name": "Jane Doe",
				"job":  "student",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := JSONFileToMap(tc.filename)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDeepMergeJSON(t *testing.T) {
	testCases := []struct {
		name     string
		dst      map[string]interface{}
		src      map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "Test case 1",
			dst: map[string]interface{}{
				"name": "John Doe",
				"age":  30,
				"address": map[string]interface{}{
					"city":    "New York",
					"country": "USA",
				},
			},
			src: map[string]interface{}{
				"age": 35,
				"address": map[string]interface{}{
					"city":  "Boston",
					"state": "MA",
				},
				"phone": "123-456-7890",
			},
			expected: map[string]interface{}{
				"name": "John Doe",
				"age":  35,
				"address": map[string]interface{}{
					"city":    "Boston",
					"country": "USA",
					"state":   "MA",
				},
				"phone": "123-456-7890",
			},
		},
		{
			name: "Test case 2",
			dst: map[string]interface{}{
				"name": "John Doe",
				"age":  30,
			},
			src: map[string]interface{}{
				"address": map[string]interface{}{
					"city": "New York",
				},
			},
			expected: map[string]interface{}{
				"name": "John Doe",
				"age":  30,
				"address": map[string]interface{}{
					"city": "New York",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DeepMergeJSON(tc.dst, tc.src)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("unexpected result: expected=%v, actual=%v", tc.expected, result)
			}
		})
	}
}
