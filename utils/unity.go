package utils

import (
	"github.com/spf13/cast"
	"regexp"
)

// UnityIds 参数归一化
func UnityIds(param ...any) (ids []any) {

	fn := func(param any) (ids []any) {
		if GetType(param) == "string" {
			// 正则提取数字部分
			item := regexp.MustCompile(`\d+`).FindAllString(cast.ToString(param), -1)
			for _, val := range item {
				ids = append(ids, cast.ToInt(val))
			}
		}
		if GetType(param) == "slice" {
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

// UnityKeys 参数归一化
func UnityKeys(param any, reg ...any) (keys []any) {

	// 正则表达式
	var regex string
	if len(reg) > 0 {
		regex = cast.ToString(reg[0])
	} else {
		regex = `[^,]+`
	}

	if GetType(param) == "string" {

		item := regexp.MustCompile(regex).FindAllString(cast.ToString(param), -1)

		for _, val := range item {
			keys = append(keys, val)
		}
	}
	if GetType(param) == "slice" {
		item := cast.ToStringSlice(param)
		for _, val := range item {
			keys = append(keys, val)
		}
	}

	return ArrayUnique(ArrayEmpty(keys))
}
