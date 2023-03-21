package logs

import "fmt"

// Extract extracts all the key-value pairs from a nested map into dot notation.
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

func ExtractKeys(kv map[string]string) []string {
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	return keys
}
