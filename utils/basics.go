package utils

import "runtime"

// ForMap - 遍历数组，返回新数组
func ForMap[T any](slice []T, fun func(item T) (result T)) (newSlice []T) {
	for key, val := range slice {
		slice[key] = fun(val)
	}
	return slice
}

// Ternary
// 三元运算符
func Ternary[T any](IF bool, TRUE T, FALSE T) T {
	if IF {
		return TRUE
	}
	return FALSE
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