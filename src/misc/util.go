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

func ArrContainsInt(arr []int, element int) bool {
	for _, ele := range arr {
		if ele == element {
			return true
		}
	}
	return false
}
