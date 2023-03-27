package utils

import (
	"github.com/spf13/cast"
	"sort"
	"strings"
)

// FormatQuery 转 Query 格式
func FormatQuery(data any) (result string) {

	body := cast.ToStringMap(data)

	// ========== 此处解决 map 无序问题 - 开始 ==========
	keys := make([]string, 0, len(body))
	for key := range body {
		keys = append(keys, key)
	}
	// 排序 keys
	sort.Strings(keys)
	// ========== 此处解决 map 无序问题 - 开始 ==========

	for key := range keys {
		result += keys[key] + "=" + cast.ToString(body[keys[key]]) + "&"
	}

	return strings.TrimRight(result, "&")
}
