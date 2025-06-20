package utils

import (
	"errors"
	"strings"
	"unicode"
)

// Identify - 识别
var Identify *IdentifyClass

type IdentifyClass struct {}

// EmailOrPhone 识别邮箱或手机号
func (this *IdentifyClass) EmailOrPhone(input string) (mode string, err error) {

	// 去除所有空白字符（包括中间的空格）
	cleanInput := strings.Map(func(value rune) rune {
		// 删除所有空白字符
		if unicode.IsSpace(value) { return -1 }
		return value
	}, input)

	// 检查是否是邮箱（允许原始输入中有空格，但clean后必须符合格式）
	if Is.Email(cleanInput) {
		// 但原始输入中@前后不能有空格
		if this.hasSpaceAroundAt(input) {
			return "", errors.New("邮箱不能在 @ 符号前后包含空格")
		}
		return "email", nil
	}

	// 检查是否是手机号（必须严格11位数字，无任何其他字符）
	if Is.Phone(cleanInput) {
		// 如果原始输入和cleanInput不同，说明有空格，视为错误
		if input != cleanInput {
			return "", errors.New("手机号不能包含空格")
		}
		return "phone", nil
	}

	// 都不是则返回错误
	return "", errors.New("既不是有效的邮箱也不是手机号")
}

// 检查邮箱中@符号周围是否有空格
func (this *IdentifyClass) hasSpaceAroundAt(input string) bool {

	atIndex := strings.Index(input, "@")
	if atIndex == -1 { return false }

	// 检查@前面或后面是否是空格
	if atIndex > 0 && unicode.IsSpace(rune(input[atIndex-1])) {
		return true
	}
	if atIndex < len(input)-1 && unicode.IsSpace(rune(input[atIndex+1])) {
		return true
	}
	return false
}