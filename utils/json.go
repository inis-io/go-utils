package utils

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"sort"
	"strings"
)

// JsonEncode 编码
func JsonEncode(data any) (result string) {
	text, err := json.Marshal(data)
	return Ternary(err != nil, "", string(text))
}

// JsonDecode 解码
func JsonDecode(data any) (result any) {
	err := json.Unmarshal([]byte(cast.ToString(data)), &result)
	return Ternary(err != nil, nil, result)
}

// JsonGet 获取json中的值 - 支持多级
func JsonGet(jsonString any, key any) (result any, err error) {

	if err := json.Unmarshal([]byte(cast.ToString(jsonString)), &result); err != nil {
		return nil, err
	}

	keys := strings.Split(cast.ToString(key), ".")

	for _, key := range keys {
		object, ok := result.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid key: %v", key)
		}
		val, ok := object[key]
		if !ok {
			return nil, fmt.Errorf("key not found: %v", key)
		}
		result = val
	}

	return result, nil
}

// JsonString map转json字符串
func JsonString(data any) (result string) {

	item := cast.ToStringMap(data)

	// 将map的key排序
	keys := make([]string, 0, len(item))
	for key := range item {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 构建有序的map
	order := make(map[string]any)
	for _, key := range keys {
		order[key] = item[key]
	}

	// 转换为JSON字符串
	bytes, err := json.Marshal(order)
	if err != nil {
		return ""
	}

	// 判断是否为空
	if string(bytes) == "{}" {
		return ""
	}

	return string(bytes)
}
