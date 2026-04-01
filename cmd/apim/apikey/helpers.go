package apikey

import "fmt"

func boolField(item any, key string) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	v, ok := m[key]
	if !ok {
		return ""
	}

	b, ok := v.(bool)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%t", b)
}
