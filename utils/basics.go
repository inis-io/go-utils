package utils

// ForMap - 遍历数组，返回新数组
func ForMap[T any](slice []T, fun func(item T) (result T)) (newSlice []T) {
	for key, val := range slice {
		slice[key] = fun(val)
	}
	return slice
}