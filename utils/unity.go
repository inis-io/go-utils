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