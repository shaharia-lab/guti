// Package guti contains utility functions
package guti

import (
	"encoding/json"
	"io/ioutil"
)

// JSONToMap takes a JSON byte array as input and returns a map containing the parsed JSON data.
// The method uses the json.Unmarshal function to parse the input data into a map with string keys
// and interface{} values. If the input data is not a valid JSON string, the method returns an error.
//
// This method can be useful when you need to convert a JSON string into a map of key-value pairs
// dynamically, such as when you are working with data that is generated or received in JSON format.
//
// Example usage:
//
//	data := []byte(`{"name": "John Doe", "age": 30}`)
//	result, err := JSONToMap(data)
//	fmt.Println(result["name"])  // prints "John Doe"
//	fmt.Println(result["age"])   // prints 30
func JSONToMap(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	return result, err
}

// JSONToString takes a JSON object as input and returns a string representation of the JSON.
// The method uses the json.Marshal function to serialize the input data into a byte array,
// and then converts the byte array to a string.
//
// This method can be useful when you need to convert a JSON object into a string dynamically,
// such as when you are working with data that needs to be serialized or sent over the network
// as a string.
//
// Example usage:
//
//	data := map[string]interface{}{"name": "John", "age": 30}
//	result, err := JSONToString(data)
//	fmt.Println(result)   // prints '{"age":30,"name":"John"}'
func JSONToString(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	return string(b), err
}

// JSONFileToMap reads a JSON file from disk and returns a map containing the parsed JSON data.
// The method reads the contents of the file using the ioutil.ReadFile function, and then uses
// the json.Unmarshal function to parse the contents into a map with string keys and interface{}
// values.
//
// This method can be useful when you need to load a JSON file from disk and convert its contents
// into a map of key-value pairs dynamically, such as when you are working with configuration files
// or other JSON data stored on disk.
//
// Example usage:
//
//	result, err := JSONFileToMap("path/to/file.json")
//	if err != nil {
//		fmt.Println("Error reading JSON file:", err)
//	} else {
//		fmt.Println(result)
//	}
func JSONFileToMap(filename string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// DeepMergeJSON recursively merges two JSON objects together, combining the values of the same keys from both objects.
// The method takes two maps with string keys and interface{} values as input, and returns a new map with the merged values.
// If a key exists in both input maps, the value from the src map will overwrite the value from the dst map.
// If both values are maps, the method will recursively merge the sub-maps. Otherwise, the method will replace the
// value in the dst map with the value from the src map.
//
// This method can be useful when you need to merge two JSON objects together dynamically, such as when you are
// working with configuration files or other JSON data that needs to be merged together from multiple sources.
//
// Example usage:
//
//	dst := map[string]interface{}{
//	  "name": "John Doe",
//	  "age":  30,
//	  "address": map[string]interface{}{
//	    "city":    "New York",
//	    "country": "USA",
//	  },
//	}
//	src := map[string]interface{}{
//	  "age": 35,
//	  "address": map[string]interface{}{
//	    "city":  "Boston",
//	    "state": "MA",
//	  },
//	  "phone": "123-456-7890",
//	}
//	result := DeepMergeJSON(dst, src)
//	fmt.Println(result)
func DeepMergeJSON(dst, src map[string]interface{}) map[string]interface{} {
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			if srcValMap, srcValIsMap := srcVal.(map[string]interface{}); srcValIsMap {
				if dstValMap, dstValIsMap := dstVal.(map[string]interface{}); dstValIsMap {
					dst[key] = DeepMergeJSON(dstValMap, srcValMap)
				}
			} else {
				dst[key] = srcVal
			}
		} else {
			dst[key] = srcVal
		}
	}
	return dst
}
