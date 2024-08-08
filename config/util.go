package config

func isUnstructuredMap(v any) bool {
	switch v.(type) {
	case map[string]interface{}, *map[string]interface{}:
		return true
	}
	return false
}
