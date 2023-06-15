package utils

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// MapWithField - 从map中提取指定字段
func MapWithField[T map[string]any](data T, field []string) (result T) {
	result = make(T)
	for _, val := range field {
		result[val] = data[val]
	}
	return
}

// MapWithoutField - 从map中排除指定字段
func MapWithoutField[T map[string]any](data T, field []string) (result T) {
	result = make(T)
	for key, val := range data {
		if !InArray[string](key, field) {
			result[key] = val
		}
	}
	return
}

// MapToURL - 无序 map 转 有序 URL
func MapToURL(params map[string]any) (result string) {

	// 创建一个 URL 结构体实例
	URL := url.URL{}

	// 将 map 的 key 按照字母顺序排序
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 构造有序的 URL 查询参数
	values := url.Values{}
	for _, index := range keys {

		item := params[index]

		switch value := item.(type) {
		case string:
			values.Set(index, value)
		case []any:
			for _, ele := range value {
				values.Add(index, fmt.Sprintf("%v", ele))
			}
		case map[string]any:
			for key, val := range value {
				values.Set(index+"["+key+"]", fmt.Sprintf("%v", val))
			}
		default:
			values.Set(index, fmt.Sprintf("%v", item))
		}
	}

	// 将查询参数附加到 URL 中
	URL.RawQuery = values.Encode()

	// url 解码
	// URL.RawQuery, _ = url.QueryUnescape(URL.RawQuery)
	// 去除末尾的 & 和 前面的 ?
	result = strings.TrimSuffix(URL.String(), "&")

	return strings.TrimPrefix(result, "?")
}
