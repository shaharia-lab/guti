// Package gotil contains packages
package gotil

import (
	"encoding/json"
	"io/ioutil"
)

// JSONToMap takes a JSON byte array as input and returns a map containing the parsed JSON data.
func JSONToMap(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	return result, err
}

// JSONToString takes a JSON object as input and returns a string representation of the JSON
func JSONToString(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	return string(b), err
}

// JSONFileToMap reads a JSON file from disk and returns a map containing the parsed JSON data
func JSONFileToMap(filename string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// DeepMergeJSON recursively merges two JSON objects together, combining the values of the same keys from both objects
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
