package utils

import (
	"fmt"
	"github.com/spf13/cast"
)

type LangModel struct {
	Directory string	// 语言包目录
	Lang 	  string	// 当前语言
	Mode 	  string	// 文件类型
}

// Lang 实例化
func Lang(model ...LangModel) *LangModel {

	item := new(LangModel)

	// 合并参数
	if len(model) > 0 {
		item = &model[0]
	}

	// 设置默认值
	item.Lang      = Ternary[string](!IsEmpty(item.Lang), item.Lang, "zh-cn")
	item.Mode 	   = Ternary[string](!IsEmpty(item.Mode), item.Mode, "json")
	item.Directory = Ternary[string](!IsEmpty(item.Directory), item.Directory, "config/i18n/")

	return item
}

func (this *LangModel) Value(key any, args ...any) (result any) {

	// 读取语言包
	bytes  := File().Byte(this.Directory + this.Lang + "." + this.Mode)
	text   := cast.ToString(key)

	// 解析语言包
	lang := cast.ToStringMap(JsonDecode(bytes.Result))

	// 获取语言
	result = lang[text]

	// 如果没有找到语言，通过javascript风格获取
	if IsEmpty(result) {
		item, err := JsonGet(bytes.Result, text)
		if err == nil {
			result = item
		}
	}

	// 如果没有找到语言，则返回原文
	if IsEmpty(result) {
		return text
	}

	// 如果有参数，则格式化
	if len(args) > 0 {
		return fmt.Sprintf(cast.ToString(result), args...)
	}

	return
}