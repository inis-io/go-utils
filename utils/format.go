package utils

import (
	"github.com/spf13/cast"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"sort"
	"strings"
)

var Format *FormatClass

type FormatClass struct{}

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

// LowerSnakeCase 小蛇式命名法
func (this *FormatClass) LowerSnakeCase(name any) string {

	// 处理驼峰命名法（包括大驼峰和小驼峰）
	item := regexp.MustCompile("([a-z])([A-Z])")
	name  = item.ReplaceAllString(cast.ToString(name), "${1}_${2}")

	// 将空格替换为下划线 - 转换为小写
	return strings.ReplaceAll(cases.Lower(language.English).String(cast.ToString(name)), " ", "_")
}

// UpperSnakeCase 大蛇式命名法
func (this *FormatClass) UpperSnakeCase(name any) string {

	// 处理驼峰命名法（包括大驼峰和小驼峰）
	item := regexp.MustCompile("([a-z])([A-Z])")
	name  = item.ReplaceAllString(cast.ToString(name), "${1}_${2}")

	// 将空格替换为下划线 - 转换为大写
	return strings.ReplaceAll(cases.Upper(language.English).String(cast.ToString(name)), " ", "_")
}

// LowerKebabCase 小烤肉串式命名法
func (this *FormatClass) LowerKebabCase(name any) string {

	// 处理驼峰命名法（包括大驼峰和小驼峰）
	item := regexp.MustCompile("([a-z])([A-Z])")
	name  = item.ReplaceAllString(cast.ToString(name), "${1}-${2}")

	// 将空格替换为下划线 - 转换为小写
	return strings.ReplaceAll(cases.Lower(language.English).String(cast.ToString(name)), " ", "-")
}

// UpperKebabCase 大烤肉串式命名法
func (this *FormatClass) UpperKebabCase(name any) string {

	// 处理驼峰命名法（包括大驼峰和小驼峰）
	item := regexp.MustCompile("([a-z])([A-Z])")
	name  = item.ReplaceAllString(cast.ToString(name), "${1}-${2}")

	// 将空格替换为下划线 - 转换为大写
	return strings.ReplaceAll(cases.Upper(language.English).String(cast.ToString(name)), " ", "-")
}

// LowerSpaceCase - 小空格命名法
func (this *FormatClass) LowerSpaceCase(name any) string {

	// 处理驼峰命名法（包括大驼峰和小驼峰）
	item := regexp.MustCompile("([a-z])([A-Z])")
	name  = item.ReplaceAllString(cast.ToString(name), "${1} ${2}")

	// 将空格替换为下划线 - 转换为小写
	return strings.ReplaceAll(cases.Lower(language.English).String(cast.ToString(name)), " ", " ")
}

// UpperSpaceCase - 大空格命名法
func (this *FormatClass) UpperSpaceCase(name any) string {

	// 处理驼峰命名法（包括大驼峰和小驼峰）
	item := regexp.MustCompile("([a-z])([A-Z])")
	name  = item.ReplaceAllString(cast.ToString(name), "${1} ${2}")

	// 将空格替换为下划线 - 转换为大写
	return strings.ReplaceAll(cases.Upper(language.English).String(cast.ToString(name)), " ", " ")
}

// LowerCamelCase - 小驼峰命名法
func (this *FormatClass) LowerCamelCase(name any) string {

	// 先将字符串按分隔符（这里假设是空格或下划线）分割成单词切片
	words := regexp.MustCompile(`[\s_]+`).Split(cast.ToString(name), -1)

	// 转换每个单词为首字母小写（除了第一个单词）
	var result string
	for i, word := range words {
		if i == 0 {
			result += word
		} else {
			result += cases.Title(language.English).String(word)
		}
	}

	// 将结果的第一个字符转换为小写
	return cases.Lower(language.English).String(result[:1]) + result[1:]
}

// UpperCamelCase - 大驼峰命名法
func (this *FormatClass) UpperCamelCase(name any) string {

	// 按分隔符（空格或下划线）分割字符串为单词切片
	words := regexp.MustCompile(`[\s_]+`).Split(cast.ToString(name), -1)

	// 转换每个单词为首字母大写
	var result string
	for _, word := range words {
		result += cases.Title(language.English).String(word)
	}

	return result
}
