package utils

import (
	"github.com/spf13/cast"
	"reflect"
)

// StructSet - 动态给结构体赋值
func StructSet(obj any, key string, val any) {
	// 获取结构体的值
	value := reflect.ValueOf(obj).Elem()
	// 获取结构体的类型
	typ := value.Type()
	// 遍历结构体的字段
	for i := 0; i < value.NumField(); i++ {
		// 获取结构体的字段
		field := typ.Field(i)
		// 获取字段的tag
		tag := field.Tag.Get("json")
		// 判断字段的tag是否等于传入的key
		if tag == key {
			// 类型断言
			switch field.Type.Kind() {
			case reflect.String:
				value.Field(i).SetString(cast.ToString(val))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value.Field(i).SetInt(cast.ToInt64(val))
			case reflect.Bool:
				value.Field(i).SetBool(cast.ToBool(val))
			case reflect.Float32, reflect.Float64:
				value.Field(i).SetFloat(cast.ToFloat64(val))
			default:
				value.Field(i).Set(reflect.ValueOf(val))
			}
			return
		}
	}
}

// StructGet - 动态获取结构体的值
func StructGet(obj any, key string) (result any) {
	// 获取结构体的值
	value := reflect.ValueOf(obj).Elem()
	// 获取结构体的类型
	typ := value.Type()
	// 遍历结构体的字段
	for i := 0; i < value.NumField(); i++ {
		// 获取结构体的字段
		field := typ.Field(i)
		// 获取字段的tag
		tag := field.Tag.Get("json")
		// 判断字段的tag是否等于传入的key
		if tag == key {
			// 获取字段的值
			return value.Field(i).Interface()
		}
	}
	return nil
}

// StructDel - 删除结构体的字段
func StructDel(obj any, key string) {
	// 获取结构体的值
	value := reflect.ValueOf(obj).Elem()
	// 获取结构体的类型
	typ := value.Type()
	// 遍历结构体的字段
	for i := 0; i < value.NumField(); i++ {
		// 获取结构体的字段
		field := typ.Field(i)
		// 获取字段的tag
		tag := field.Tag.Get("json")
		// 判断字段的tag是否等于传入的key
		if tag == key {
			// 获取字段的值
			value.Field(i).Set(reflect.Zero(field.Type))
		}
	}
}

// StructHas - 判断结构体是否存在某个字段
func StructHas(obj any, key string) (ok bool) {
	// 获取结构体的值
	value := reflect.ValueOf(obj).Elem()
	// 获取结构体的类型
	typ := value.Type()
	// 遍历结构体的字段
	for i := 0; i < value.NumField(); i++ {
		// 获取结构体的字段
		field := typ.Field(i)
		// 获取字段的tag
		tag := field.Tag.Get("json")
		// 判断字段的tag是否等于传入的key
		if tag == key {
			return true
		}
	}
	return false
}

// StructKeys - 获取结构体的字段
func StructKeys(obj any) (slice []string) {
	// 获取结构体的值
	value := reflect.ValueOf(obj).Elem()
	// 获取结构体的类型
	typ := value.Type()
	// 定义一个切片
	keys := make([]string, 0)
	// 遍历结构体的字段
	for i := 0; i < value.NumField(); i++ {
		// 获取结构体的字段
		field := typ.Field(i)
		// 获取字段的tag
		tag := field.Tag.Get("json")
		// 判断字段的tag是否等于传入的key
		keys = append(keys, tag)
	}
	return keys
}

// StructValues - 获取结构体的值
func StructValues(obj any) (slice []any) {
	// 获取结构体的值
	value := reflect.ValueOf(obj).Elem()
	// 定义一个切片
	keys := make([]any, 0)
	// 遍历结构体的字段
	for i := 0; i < value.NumField(); i++ {
		// 获取字段的值
		keys = append(keys, value.Field(i).Interface())
	}
	return keys
}

// StructLen - 获取结构体的长度
func StructLen(obj any) (length int) {
	// 获取结构体的值
	value := reflect.ValueOf(obj).Elem()
	return value.NumField()
}

// StructMap - 将结构体转换为map
func StructMap(obj any) (result map[string]any) {
	// 获取结构体的值
	value := reflect.ValueOf(obj).Elem()
	// 获取结构体的类型
	typ := value.Type()
	// 定义一个map
	// result := make(map[string]any)
	// 遍历结构体的字段
	for i := 0; i < value.NumField(); i++ {
		// 获取结构体的字段
		field := typ.Field(i)
		// 获取字段的tag
		tag := field.Tag.Get("json")
		// 获取字段的值
		result[tag] = value.Field(i).Interface()
	}
	return result
}

// StructSlice - 将结构体转换为切片
func StructSlice(obj any) (slice []any) {
	// 获取结构体的值
	value := reflect.ValueOf(obj).Elem()
	// 定义一个切片
	s := make([]any, 0)
	// 遍历结构体的字段
	for i := 0; i < value.NumField(); i++ {
		// 获取字段的值
		s = append(s, value.Field(i).Interface())
	}
	return s
}