package utils

import (
	"github.com/goccy/go-json"
)

// JsonEncode 编码
func JsonEncode(data any) (result string) {
	text, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(text)
}

// JsonDecode 解码
func JsonDecode(data string) (result any) {
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil
	}
	return result
}
