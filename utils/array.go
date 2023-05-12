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

// ArrayUnique - 数组去重
func ArrayUnique[T any](array []T) (slice []any) {
	list := make(map[any]bool)
	for _, item := range array {
		if !list[item] {
			list[item] = true
			slice = append(slice, item)
		}
	}
	return slice
}

// ArrayEmpty - 数组去空
func ArrayEmpty[T any](array []T) (slice []any) {
	for _, item := range array {
		if !IsEmpty(item) {
			slice = append(slice, item)
		}
	}
	return slice
}
