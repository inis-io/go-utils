package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func typeof(args ...any) (typeof string, empty bool) {
	typeof, empty = "string", false
	for _, item := range args {

		// 判断是否为空
		if item == nil {
			empty = true
			continue
		}

		// 先利用反射获取数据类型，再进入不同类型的判空逻辑
		switch reflect.TypeOf(item).Kind().String() {
		case "int":
			typeof = "int"
			if item == 0 {
				empty = true
			}
		case "string":
			typeof = "string"
			if item == "" {
				empty = true
			}
		case "int64":
			typeof = "int64"
			if item == 0 {
				empty = true
			}
		case "uint8":
			typeof = "bool"
			if item == false {
				empty = true
			}
		case "float64":
			typeof = "float"
			if item == 0.0 {
				empty = true
			}
		case "byte":
			typeof = "byte"
			if item == 0 {
				empty = true
			}
		case "ptr":
			typeof = "ptr"
			// 接口状态下，它不认为自己是nil，所以要用反射判空
			if item == nil {
				empty = true
			}
			// 反射判空逻辑
			if reflect.ValueOf(item).IsNil() {
				// 利用反射直接判空
				empty = true
			}
		case "struct":
			typeof = "struct"
			if item == nil {
				empty = true
			}
		case "slice":
			typeof = "slice"
			s := reflect.ValueOf(item)
			if s.Len() == 0 {
				empty = true
			}
		case "array":
			typeof = "array"
			s := reflect.ValueOf(item)
			if s.Len() == 0 {
				empty = true
			}
		case "map":
			typeof = "map"
			s := reflect.ValueOf(item)
			if s.Len() == 0 {
				empty = true
			}
		case "chan":
			typeof = "chan"
			s := reflect.ValueOf(item)
			if s.Len() == 0 {
				empty = true
			}
		case "func":
			typeof = "func"
			if item == nil {
				empty = true
			}
		default:
			if item == "true" || item == "false" || item == true || item == false {
				typeof = "bool"
				if item == false {
					empty = true
				}
			} else {
				typeof = "other"
				empty = true
				fmt.Println("奇怪的数据类型")
			}
		}

	}
	return
}

func CustomProcessApi(url string, api string) (result string) {
	if empty := IsEmpty(api); empty {
		api = "api"
	}
	result = url
	if empty := IsEmpty(url); !empty {
		prefix := "//"
		if check := strings.HasPrefix(url, "https://"); check {
			prefix = "https://"
		} else if check := strings.HasPrefix(url, "http://"); check {
			prefix = "http://"
		}
		// 正则匹配 http(s):// - 并去除
		url = regexp.MustCompile("^((https|http)?:\\/\\/)").ReplaceAllString(url, "")
		array := ArrayFilter(strings.Split(url, `/`))
		if len(array) == 1 {
			result = prefix + array[0] + "/" + api + "/"
		} else if len(array) == 2 {
			result = prefix + array[0] + "/" + array[1] + "/"
		}
	}
	return
}

// InMapKey
// 在 map key 中
func InMapKey(key string, array map[string]any) bool {
	for k := range array {
		if key == k {
			return true
		}
	}
	return false
}

// MapMerge
// map合并
func MapMerge(map1 map[any]any, map2 map[any]any) map[any]any {

	result := make(map[any]any)
	for i, v := range map1 {
		for j, w := range map2 {
			if i == j {
				result[i] = w
			} else {
				if _, ok := result[i]; !ok {
					result[i] = v
				}
				if _, ok := result[j]; !ok {
					result[j] = w
				}
			}
		}
	}

	return result
}

// MapMergeString
// map合并
func MapMergeString(map1 map[string]string, map2 map[string]string) map[string]string {

	result := make(map[string]string)
	for i, v := range map1 {
		for j, w := range map2 {
			if i == j {
				result[i] = w
			} else {
				if _, ok := result[i]; !ok {
					result[i] = v
				}
				if _, ok := result[j]; !ok {
					result[j] = w
				}
			}
		}
	}

	return result
}

// InMapValue
// 在 map value 中
func InMapValue(value any, array map[string]string) bool {
	for _, val := range array {
		if value == val {
			return true
		}
	}
	return false
}

// Ternary
// 三元运算符
func Ternary[T any](IF bool, TRUE T, FALSE T) T {
	if IF {
		return TRUE
	}
	return FALSE
}
