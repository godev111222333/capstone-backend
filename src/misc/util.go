package misc

func MapGetString(m interface{}, fieldName string) string {
	data, ok := m.(map[string]interface{})
	if ok {
		value, ok := data[fieldName].(string)
		if ok {
			return value
		}
	}

	return ""
}
