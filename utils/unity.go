package utils

import (
	"github.com/spf13/cast"
	"regexp"
)

// Unity - 统一规范化
var Unity *UnityStruct

type UnityStruct struct{}

// Ids 参数归一化
func (this *UnityStruct) Ids(param ...any) (ids []any) {

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
func (this *UnityStruct) Keys(param any, reg ...any) (keys []any) {

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
