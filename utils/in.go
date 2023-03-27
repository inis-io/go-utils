package utils

// InArray - 判断某个值是否在数组中
func InArray(value any, array []any) (ok bool) {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
