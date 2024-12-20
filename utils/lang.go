package utils

import (
	"fmt"
	"github.com/spf13/cast"
)

type LangClass struct {
	Directory string // 语言包目录
	Lang      string // 当前语言
	Mode      string // 文件类型
}

// Lang 实例化
func Lang(model ...LangClass) *LangClass {

	item := new(LangClass)

	// 合并参数
	if len(model) > 0 {
		item = &model[0]
	}

	// 设置默认值
	item.Lang = Ternary[string](!Is.Empty(item.Lang), item.Lang, "zh-cn")
	item.Mode = Ternary[string](!Is.Empty(item.Mode), item.Mode, "json")
	item.Directory = Ternary[string](!Is.Empty(item.Directory), item.Directory, "config/i18n/")

	return item
}

func (this *LangClass) Value(key any, args ...any) (result any) {

	// 读取语言包
	bytes := File().Byte(this.Directory + this.Lang + "." + this.Mode)

	if bytes.Error != nil {
		return
	}

	text := cast.ToString(key)

	// 解析语言包
	lang := cast.ToStringMap(Json.Decode(bytes.Text))

	// 获取语言
	result = lang[text]

	// 如果没有找到语言，通过javascript风格获取
	if Is.Empty(result) {
		item, err := Json.Get(bytes.Text, text)
		if err == nil {
			result = item
		}
	}

	// 如果没有找到语言，则返回原文
	if Is.Empty(result) {
		return fmt.Sprintf(text, args...)
	}

	// 如果有参数，则格式化
	if len(args) > 0 {
		return fmt.Sprintf(cast.ToString(result), args...)
	}

	return
}
