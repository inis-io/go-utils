package utils

import (
	"github.com/spf13/cast"
	"runtime"
	"strings"
)

// ForMap - 遍历数组，返回新数组
func ForMap[T any](slice []T, fun func(item T) (result T)) (newSlice []T) {
	for key, val := range slice {
		slice[key] = fun(val)
	}
	return slice
}

// Ternary - 三元运算符
func Ternary[T any](IF bool, TRUE T, FALSE T) T {
	if IF {
		return TRUE
	}
	return FALSE
}

// Replace - 字符串替换
func Replace(value any, params map[string]any) (result string) {
	result = cast.ToString(value)
	for key, val := range params {
		result = strings.Replace(result, key, cast.ToString(val), -1)
	}
	return result
}

func GetCaller() (funcName string, fileName string, line int) {
	pc, fileName, line, ok := runtime.Caller(2)
	if !ok {
		return "", "", 0
	}
	funcName = runtime.FuncForPC(pc).Name()
	return funcName, fileName, line
}

type caller struct {
	// 文件名
	FileName string
	// 函数名
	FuncName string
	// 行号
	Line     int
}

// Caller 获取代码调用者
func Caller() *caller {
	funcName, fileName, line := GetCaller()
	return &caller{
		FileName: fileName,
		FuncName: funcName,
		Line:     line,
	}
}