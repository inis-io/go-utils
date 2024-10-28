package utils

import (
	"fmt"
	"reflect"
	"strings"
)

var Url = &UrlClass{}

type UrlClass struct {
}

// Encoded - 将 map 编码为 URL 查询字符串 - x-www-form-urlencoded
func (this *UrlClass) Encoded(params map[string]any) string {

	var parts []string

	for key, value := range params {
		parts = append(parts, this.build(key, value)...)
	}

	return strings.Join(parts, "&")
}

// EncodedKeys - 获取 URL 查询字符串的键
func (this *UrlClass) EncodedKeys(params map[string]any) []string {

	encoded := this.Encoded(params)
	// & 分隔数组，字符串转数组
	keys := strings.Split(encoded, "&")
	// 去除 = 号之后的内容，包含 = 号
	for i, item := range keys {
		keys[i] = item[:strings.Index(item, "=")]
	}

	return keys
}

// build - 构建 URL 查询字符串
func (this *UrlClass) build(key string, value any) []string {

	var parts []string

	switch item := value.(type) {
	case string:
		parts = append(parts, fmt.Sprintf("%s=%v", key, item))
	case []string:
		for _, sv := range item {
			parts = append(parts, fmt.Sprintf("%s[]=%v", key, sv))
		}
	case []int:
		for _, iv := range item {
			parts = append(parts, fmt.Sprintf("%s[]=%d", key, iv))
		}
	case map[string]any:
		for k, sub := range item {
			parts = append(parts, this.build(fmt.Sprintf("%s[%s]", key, k), sub)...)
		}
	case []any:
		for i, sub := range item {
			parts = append(parts, this.build(fmt.Sprintf("%s[%d]", key, i), sub)...)
		}
	default:
		parts = append(parts, fmt.Sprintf("%s=%v", key, this.stringify(value)))
	}

	return parts
}

// stringify - 将任意类型转换为字符串
func (this *UrlClass) stringify(value any) string {

	item := reflect.ValueOf(value)

	switch item.Kind() {
	case reflect.Array, reflect.Slice:
		var strVals []string
		for i := 0; i < item.Len(); i++ {
			strVals = append(strVals, fmt.Sprintf("%v", item.Index(i)))
		}
		return strings.Join(strVals, ",")
	default:
		return fmt.Sprintf("%v", value)
	}
}
