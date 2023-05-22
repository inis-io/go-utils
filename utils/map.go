package utils

// MapWithField - 从map中提取指定字段
func MapWithField[T map[string]any](data T, field []string) (result T) {
	result = make(T)
	for _, val := range field {
		result[val] = data[val]
	}
	return
}

// MapWithoutField - 从map中排除指定字段
func MapWithoutField[T map[string]any](data T, field []string) (result T) {
	result = make(T)
	for key, val := range data {
		if !InArray[string](key, field) {
			result[key] = val
		}
	}
	return
}
