package member

func roleField(item interface{}) string {
	m, ok := item.(map[string]interface{})
	if !ok {
		return ""
	}

	roles, ok := m["roles"].([]interface{})
	if !ok || len(roles) == 0 {
		return ""
	}

	first, ok := roles[0].(map[string]interface{})
	if !ok {
		return ""
	}

	name, _ := first["name"].(string)

	return name
}
