package utils

import (
	"github.com/spf13/cast"
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
func ArrayUnique[T any](array []T) (slice []T) {
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
func ArrayEmpty[T any](array []T) (slice []T) {
	for _, item := range array {
		if !Is.Empty(item) {
			slice = append(slice, item)
		}
	}
	return slice
}

// ArrayMapWithField - 数组Map保留指定字段
func ArrayMapWithField(array []map[string]any, field any) (slice []any) {

	// 获取字段
	keys := cast.ToStringSlice(Unity.Keys(field))

	if Is.Empty(keys) {
		return cast.ToSlice(array)
	}

	for _, item := range array {
		val := Map.WithField(cast.ToStringMap(item), keys)
		slice = append(cast.ToSlice(slice), val)
	}

	return slice
}

// ArrayReverse - 数组反转
func ArrayReverse[T any](array []T) (slice []T) {
	for i, j := 0, len(array)-1; i < j; i, j = i+1, j-1 {
		array[i], array[j] = array[j], array[i]
	}
	return array
}

// ArrayDiff - 数组差集
func ArrayDiff[T any](array1, array2 []T) (slice []T) {
	for _, item := range array1 {
		if !InArray(item, array2) {
			slice = append(slice, item)
		}
	}
	return slice
}

// ArrayIntersect - 交集
func ArrayIntersect[T any](array1, array2 []T) (slice []T) {
	for _, item := range array1 {
		if InArray(item, array2) {
			slice = append(slice, item)
		}
	}
	return slice
}

// ArrayPaging - 分页处理
func ArrayPaging[T any](rows []T, page, limit int) (code int, count int, data []T) {

	// 获取数据的总长度
	count = len(rows)

	// 计算总页数
	if count == 0 || limit <= 0 {
		code = 0
	} else {
		code = (count + limit - 1) / limit
	}

	// 检查页码是否有效
	if page < 1 {
		page = 1
	}
	if page > code {
		// 如果页码超过总页数，返回空数据
		return code, count, []T{}
	}

	// 计算当前页数据的起始和结束索引
	start := (page - 1) * limit
	end := start + limit
	if end > count { end = count }

	// 截取当前页的数据
	data = rows[start:end]

	return code, count, data
}