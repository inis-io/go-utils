package utils

import (
	"strings"
)

// ArrayFilter - 数组过滤
func ArrayFilter(array []string) (slice []string) {
	for key, val := range array {
		if (key > 0 && array[key-1] == val) || len(val) == 0 {
			continue
		}
		slice = append(slice, val)
	}
	return
}

// ArrayRemove - 数组删除
func ArrayRemove(array []string, args ...string) []string {
	for _, value := range args {
		for key, val := range array {
			// 去除空格
			val = strings.TrimSpace(val)
			// 根据索引删除
			if val == value {
				array = append(array[:key], array[key+1:]...)
			}
		}
	}
	return array
}
