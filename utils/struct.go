package utils

import (
	"github.com/spf13/cast"
	"reflect"
)

// Struct - 操作结构体
var Struct *StructClass

type StructClass struct {}

// Set - 动态给结构体赋值
func (this *StructClass) Set(obj any, key string, val any) {
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
				if val == nil {
					val = reflect.Zero(field.Type).Interface()
				} else {
					value.Field(i).Set(reflect.ValueOf(val))
				}
			}
			return
		}
	}
}

// Get - 动态获取结构体的值
func (this *StructClass) Get(obj any, key string) (result any) {
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

// Del - 删除结构体的字段
func (this *StructClass) Del(obj any, key string) {
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

// Has - 判断结构体是否存在某个字段
func (this *StructClass) Has(obj any, key string) (ok bool) {
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

// Keys - 获取结构体的字段
func (this *StructClass) Keys(obj any) (slice []string) {
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

// Values - 获取结构体的值
func (this *StructClass) Values(obj any) (slice []any) {
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

// Len - 获取结构体的长度
func (this *StructClass) Len(obj any) (length int) {
	// 获取结构体的值
	value := reflect.ValueOf(obj).Elem()
	return value.NumField()
}

// Map - 将结构体转换为map
func (this *StructClass) Map(obj any) (result map[string]any) {
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

// Slice - 将结构体转换为切片
func (this *StructClass) Slice(obj any) (slice []any) {
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

// ToStringMap - 将结构体转换为map[string]any
func (this *StructClass) ToStringMap(obj any) (result map[string]any) {
	return cast.ToStringMap(Json.Encode(obj))
}

// ToAsciiString - 将结构体转换为ASCII字符串
func (this *StructClass) ToAsciiString(obj any) (result string) {
	return Ascii.ToString(this.ToStringMap(obj), true)
}

// Fields - 获取结构体的字段名
func (this *StructClass) Fields(dest any) []string {

	// 获取结构体的反射值
	item := reflect.ValueOf(dest)
	// 确保传递的是一个指针
	if item.Kind() == reflect.Ptr {
		item = item.Elem()
	}
	// 确保反射值的类型是结构体
	if item.Kind()!= reflect.Struct { return nil }

	_type := item.Type()
	field := make([]string, item.NumField())

	for i := 0; i < item.NumField(); i++ {
		field[i] = _type.Field(i).Name
	}

	return field
}