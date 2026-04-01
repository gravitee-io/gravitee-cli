package member

func roleField(item any) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	roles, ok := m["roles"].([]any)
	if !ok || len(roles) == 0 {
		return ""
	}

	first, ok := roles[0].(map[string]any)
	if !ok {
		return ""
	}

	name, _ := first["name"].(string)

	return name
}
