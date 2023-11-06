package jsonutil

import (
	"encoding/json"
	"regexp"
)

func UpdateJSON(inputJSON []byte) ([]byte, error) {
	var data map[string]interface{}

	if err := json.Unmarshal(inputJSON, &data); err != nil {
		return nil, err
	}

	modifyJSON(data)

	outputJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return outputJSON, nil
}

func modifyJSON(data map[string]interface{}) {
	for key, value := range data {
		if isVowelStart(key) {
			delete(data, key)
		} else {
			switch v := value.(type) {
			case map[string]interface{}:
				modifyJSON(v)
			case []interface{}:
				modifyJSONArray(v)
			case float64:
				if v == float64(int(v)) && int(v)%2 == 0 {
					data[key] = v + 1000
				}
			}
		}
	}
}

func modifyJSONArray(arr []interface{}) {
	for i, value := range arr {
		switch v := value.(type) {
		case map[string]interface{}:
			modifyJSON(v)
		case []interface{}:
			modifyJSONArray(v)
		case float64:
			if v == float64(int(v)) && int(v)%2 == 0 {
				arr[i] = v + 1000
			}
		}
	}
}

func isVowelStart(s string) bool {
	vowelStartPattern := regexp.MustCompile("^[aeiouAEIOU]")
	return vowelStartPattern.MatchString(s)
}
