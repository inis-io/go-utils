package utils

import (
	"sort"
	"strings"
	
	"github.com/spf13/cast"
)

var Format *FormatClass

type FormatClass struct {}

// Query 转 Query 格式
func (this *FormatClass) Query(data any) (result string) {

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