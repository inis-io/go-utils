package utils

import (
	"github.com/spf13/cast"
	"reflect"
	"regexp"
)

// Is - 是否为
var Is *IsStruct

type IsStruct struct{}

// Email - 是否为邮箱
func (this *IsStruct) Email(email any) (ok bool) {
	if email == nil {
		return false
	}
	return regexp.MustCompile(`\w[-\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\.)+[A-Za-z]{2,14}`).MatchString(cast.ToString(email))
}

// Phone - 是否为手机号
func (this *IsStruct) Phone(phone any) (ok bool) {
	if phone == nil {
		return false
	}
	return regexp.MustCompile(`^1[3456789]\d{9}$`).MatchString(cast.ToString(phone))
}

// Mobile - 是否为手机号
func (this *IsStruct) Mobile(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^1[3456789]\d{9}$`).MatchString(cast.ToString(value))
}

// Empty - 是否为空
func (this *IsStruct) Empty(value any) (ok bool) {
	_, empty := typeof(value)
	return empty
}

// Domain - 是否为域名
func (this *IsStruct) Domain(domain any) (ok bool) {
	if domain == nil {
		return false
	}
	return regexp.MustCompile(`^((https|http|ftp|rtsp|mms)?://)\S+`).MatchString(cast.ToString(domain))
}

// True - 是否为真
func (this *IsStruct) True(value any) (ok bool) {
	return cast.ToBool(value)
}

// False - 是否为假
func (this *IsStruct) False(value any) (ok bool) {
	return !cast.ToBool(value)
}

// Number - 是否为数字
func (this *IsStruct) Number(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[0-9]+$`).MatchString(cast.ToString(value))
}

// Float - 是否为浮点数
func (this *IsStruct) Float(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[0-9]+(.[0-9]+)?$`).MatchString(cast.ToString(value))
}

// Bool - 是否为bool
func (this *IsStruct) Bool(value any) (ok bool) {
	return cast.ToBool(value)
}

// Accepted - 验证某个字段是否为为 yes, on, 或是 1
func (this *IsStruct) Accepted(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^(yes|on|1)$`).MatchString(cast.ToString(value))
}

// Date - 是否为日期类型
func (this *IsStruct) Date(date any) (ok bool) {
	if date == nil {
		return false
	}
	return regexp.MustCompile(`^\d{4}-\d{1,2}-\d{1,2}$`).MatchString(cast.ToString(date))
}

// Alpha - 只能包含字母
func (this *IsStruct) Alpha(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(cast.ToString(value))
}

// AlphaNum - 只能包含字母和数字
func (this *IsStruct) AlphaNum(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(cast.ToString(value))
}

// AlphaDash - 只能包含字母、数字和下划线_及破折号-
func (this *IsStruct) AlphaDash(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(cast.ToString(value))
}

// Chs - 是否为汉字
func (this *IsStruct) Chs(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x{4e00}-\x{9fa5}]+$`).MatchString(cast.ToString(value))
}

// ChsAlpha - 只能是汉字、字母
func (this *IsStruct) ChsAlpha(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x{4e00}-\x{9fa5}a-zA-Z]+$`).MatchString(cast.ToString(value))
}

// ChsAlphaNum - 只能是汉字、字母和数字
func (this *IsStruct) ChsAlphaNum(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x{4e00}-\x{9fa5}a-zA-Z0-9]+$`).MatchString(cast.ToString(value))
}

// ChsDash - 只能是汉字、字母、数字和下划线_及破折号-
func (this *IsStruct) ChsDash(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x{4e00}-\x{9fa5}a-zA-Z0-9_-]+$`).MatchString(cast.ToString(value))
}

// Cntrl - 是否为控制字符 - （换行、缩进、空格）
func (this *IsStruct) Cntrl(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x00-\x1F\x7F]+$`).MatchString(cast.ToString(value))
}

// Graph - 是否为可见字符 - （除空格外的所有可打印字符）
func (this *IsStruct) Graph(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x21-\x7E]+$`).MatchString(cast.ToString(value))
}

// Lower - 是否为小写字母
func (this *IsStruct) Lower(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[a-z]+$`).MatchString(cast.ToString(value))
}

// Upper - 是否为大写字母
func (this *IsStruct) Upper(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[A-Z]+$`).MatchString(cast.ToString(value))
}

// Space - 是否为空白字符 - （空格、制表符、换页符等）
func (this *IsStruct) Space(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\s]+$`).MatchString(cast.ToString(value))
}

// Xdigit - 是否为十六进制字符 - （0-9、a-f、A-F）
func (this *IsStruct) Xdigit(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\da-fA-F]+$`).MatchString(cast.ToString(value))
}

// ActiveUrl - 是否为有效的域名或者IP
func (this *IsStruct) ActiveUrl(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^([a-z0-9-]+\.)+[a-z]{2,6}$`).MatchString(cast.ToString(value))
}

// Ip - 是否为IP
func (this *IsStruct) Ip(ip any) (ok bool) {
	if ip == nil {
		return false
	}
	return regexp.MustCompile(`(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`).MatchString(cast.ToString(ip))
}

// Url - 是否为URL
func (this *IsStruct) Url(url any) (ok bool) {
	if url == nil {
		return false
	}
	return regexp.MustCompile(`^((https|http|ftp|rtsp|mms)?://)\S+`).MatchString(cast.ToString(url))
}

// IdCard - 是否为有效的身份证号码
func (this *IsStruct) IdCard(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)`).MatchString(cast.ToString(value))
}

// MacAddr - 是否为有效的MAC地址
func (this *IsStruct) MacAddr(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^([A-Fa-f0-9]{2}:){5}[A-Fa-f0-9]{2}$`).MatchString(cast.ToString(value))
}

// Zip - 是否为有效的邮政编码
func (this *IsStruct) Zip(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[1-9]\d{5}$`).MatchString(cast.ToString(value))
}

// String - 是否为字符串
func (this *IsStruct) String(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.String
}

// Slice - 是否为切片
func (this *IsStruct) Slice(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Slice
}

// Array - 是否为数组
func (this *IsStruct) Array(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Array
}

// JsonString - 是否为json字符串
func (this *IsStruct) JsonString(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[{\[].*[}\]]$`).MatchString(cast.ToString(value))
}

// Map - 是否为map
func (this *IsStruct) Map(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Map
}

// SliceSlice - 是否为二维切片
func (this *IsStruct) SliceSlice(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Slice && reflect.TypeOf(value).Elem().Kind() == reflect.Slice
}

// MapAny - 是否为[]map[string]any
func (this *IsStruct) MapAny(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Map && reflect.TypeOf(value).Elem().Kind() == reflect.Interface
}
