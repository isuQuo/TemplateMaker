package logs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocopper/copper/cerrors"
)

// Import a JSON file and return the JSON object.
// TODO: Prompt user for file with GUI.
func ImportJSONFile(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, cerrors.New(err, "failed to open file", map[string]interface{}{
			"path": path,
		})
	}
	defer file.Close()

	var jsonObject map[string]interface{}
	if err := json.NewDecoder(file).Decode(&jsonObject); err != nil {
		return nil, cerrors.New(err, "failed to decode json", map[string]interface{}{
			"path": path,
		})
	}

	return jsonObject, nil
}

// ExtractKeyValues extracts all the key-value pairs from a nested map into dot notation.
func ExtractKeyValues(prefix string, data map[string]interface{}, kv *map[string]string) error {
	for key, value := range data {
		newPrefix := prefix
		if newPrefix != "" {
			newPrefix += "."
		}
		newPrefix += key

		if nestedMap, ok := value.(map[string]interface{}); ok {
			ExtractKeyValues(newPrefix, nestedMap, kv)
		} else if nestedSlice, ok := value.([]interface{}); ok {
			for i, item := range nestedSlice {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					ExtractKeyValues(fmt.Sprintf("%s.%d", newPrefix, i), nestedMap, kv)
				}
			}
		} else {
			(*kv)[newPrefix] = fmt.Sprintf("%v", value)
		}
	}

	return nil
}

// ExtractKeys extracts all the keys from a map.
func ExtractKeys(kv map[string]string) []string {
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	return keys
}
