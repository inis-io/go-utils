package utils

import (
	"github.com/spf13/cast"
	"regexp"
	"strings"
)

// Unity - 统一规范化
var Unity *UnityClass

type UnityClass struct{}

// Ids 参数归一化
func (this *UnityClass) Ids(param ...any) (ids []any) {

	fn := func(param any) (ids []any) {

		types := []string{"string", "int", "int64", "float", "float64"}

		if InArray(Get.Type(param), types) {
			// 正则提取数字部分
			item := regexp.MustCompile(`\d+`).FindAllString(cast.ToString(param), -1)
			for _, val := range item {
				ids = append(ids, cast.ToInt(val))
			}

		}
		if Get.Type(param) == "slice" {
			item := cast.ToStringSlice(param)
			for _, val := range item {
				ids = append(ids, cast.ToInt(val))
			}
		}
		return ids
	}

	for _, val := range param {
		ids = append(ids, fn(val)...)
	}

	return ArrayUnique(ArrayEmpty(ids))
}

// Keys 参数归一化
func (this *UnityClass) Keys(param any, reg ...any) (keys []any) {

	// 正则表达式
	var regex string
	if len(reg) > 0 {
		regex = cast.ToString(reg[0])
	} else {
		regex = `[^,]+`
	}

	if Get.Type(param) == "string" {

		item := regexp.MustCompile(regex).FindAllString(cast.ToString(param), -1)

		for _, val := range item {
			keys = append(keys, val)
		}
	}
	if Get.Type(param) == "slice" {
		item := cast.ToStringSlice(param)
		for _, val := range item {
			keys = append(keys, val)
		}
	}

	// 去重 - 去空
	keys = ArrayUnique(ArrayEmpty(keys))

	// 遍历每一项，去重空格
	for key, val := range keys {
		keys[key] = regexp.MustCompile(`\s+`).ReplaceAllString(cast.ToString(val), "")
	}

	return keys
}

// Int 参数归一化
func (this *UnityClass) Int(value ...any) (array []int) {

	fn := func(value any) (resp []int) {

		types := []string{"string", "int", "int64", "float", "float64"}

		if InArray(Get.Type(value), types) {
			// 正则提取数字部分，包含0
			item := regexp.MustCompile(`-?\d+`).FindAllString(cast.ToString(value), -1)
			for _, val := range item {
				resp = append(resp, cast.ToInt(val))
			}
		}

		if Get.Type(value) == "slice" {
			item := cast.ToStringSlice(value)
			for _, val := range item {
				resp = append(resp, cast.ToInt(val))
			}
		}
		return resp
	}

	for _, val := range value {
		array = append(array, fn(val)...)
	}

	return ArrayUnique(array)
}

// Float 参数归一化
func (this *UnityClass) Float(value ...any) (array []float64) {

	fn := func(value any) (resp []float64) {

		types := []string{"string", "int", "int64", "float", "float64"}

		if InArray(Get.Type(value), types) {
			// 正则提取数字部分，包含正数，负数，0和小数
			item := regexp.MustCompile(`-?\d+(\.\d+)?`).FindAllString(cast.ToString(value), -1)

			for _, val := range item {
				resp = append(resp, cast.ToFloat64(val))
			}
		}
		if Get.Type(value) == "slice" {
			item := cast.ToStringSlice(value)
			for _, val := range item {
				resp = append(resp, cast.ToFloat64(val))
			}
		}
		return resp
	}

	for _, val := range value {
		array = append(array, fn(val)...)
	}

	return ArrayUnique(array)
}

// Join - 数组转字符串
func (this *UnityClass) Join(elems any, unique bool, sep ...string) string {

	// 默认分隔符
	if len(sep) == 0 { sep = []string{","} }

	rows := cast.ToStringSlice(elems)

	// 去重，去空
	if unique {
		rows = ArrayUnique(ArrayEmpty(rows))
	}

	if len(rows) == 0 { return "" }

	return sep[0] + strings.Join(rows, sep[0]) + sep[0]
}