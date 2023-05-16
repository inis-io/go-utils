package utils

import (
	"github.com/spf13/cast"
	"reflect"
	"regexp"
)

// IsEmail - 是否为邮箱
func IsEmail(email any) (ok bool) {
	if email == nil {
		return false
	}
	return regexp.MustCompile(`\w[-\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\.)+[A-Za-z]{2,14}`).MatchString(cast.ToString(email))
}

// IsPhone - 是否为手机号
func IsPhone(phone any) (ok bool) {
	if phone == nil {
		return false
	}
	return regexp.MustCompile(`^1[3456789]\d{9}$`).MatchString(cast.ToString(phone))
}

// IsMobile - 是否为手机号
func IsMobile(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^1[3456789]\d{9}$`).MatchString(cast.ToString(value))
}

// IsEmpty - 是否为空
func IsEmpty(value any) (ok bool) {
	_, empty := typeof(value)
	return empty
}

// IsDomain - 是否为域名
func IsDomain(domain any) (ok bool) {
	if domain == nil {
		return false
	}
	return regexp.MustCompile(`^((https|http|ftp|rtsp|mms)?:\/\/)[^\s]+`).MatchString(cast.ToString(domain))
}

// IsTrue - 是否为真
func IsTrue(value any) (ok bool) {
	return cast.ToBool(value)
}

// IsFalse - 是否为假
func IsFalse(value any) (ok bool) {
	return !cast.ToBool(value)
}

// IsNumber - 是否为数字
func IsNumber(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[0-9]+$`).MatchString(cast.ToString(value))
}

// IsFloat - 是否为浮点数
func IsFloat(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[0-9]+(.[0-9]+)?$`).MatchString(cast.ToString(value))
}

// IsBool - 是否为bool
func IsBool(value any) (ok bool) {
	return cast.ToBool(value)
}

// IsAccepted - 验证某个字段是否为为 yes, on, 或是 1
func IsAccepted(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^(yes|on|1)$`).MatchString(cast.ToString(value))
}

// IsDate - 是否为日期类型
func IsDate(date any) (ok bool) {
	if date == nil {
		return false
	}
	return regexp.MustCompile(`^\d{4}-\d{1,2}-\d{1,2}$`).MatchString(cast.ToString(date))
}

// IsAlpha - 只能包含字母
func IsAlpha(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(cast.ToString(value))
}

// IsAlphaNum - 只能包含字母和数字
func IsAlphaNum(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(cast.ToString(value))
}

// IsAlphaDash - 只能包含字母、数字和下划线_及破折号-
func IsAlphaDash(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(cast.ToString(value))
}

// IsChs - 是否为汉字
func IsChs(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x{4e00}-\x{9fa5}]+$`).MatchString(cast.ToString(value))
}

// IsChsAlpha - 只能是汉字、字母
func IsChsAlpha(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x{4e00}-\x{9fa5}a-zA-Z]+$`).MatchString(cast.ToString(value))
}

// IsChsAlphaNum - 只能是汉字、字母和数字
func IsChsAlphaNum(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x{4e00}-\x{9fa5}a-zA-Z0-9]+$`).MatchString(cast.ToString(value))
}

// IsChsDash - 只能是汉字、字母、数字和下划线_及破折号-
func IsChsDash(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x{4e00}-\x{9fa5}a-zA-Z0-9_-]+$`).MatchString(cast.ToString(value))
}

// IsCntrl - 是否为控制字符 - （换行、缩进、空格）
func IsCntrl(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x00-\x1F\x7F]+$`).MatchString(cast.ToString(value))
}

// IsGraph - 是否为可见字符 - （除空格外的所有可打印字符）
func IsGraph(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\x21-\x7E]+$`).MatchString(cast.ToString(value))
}

// IsLower - 是否为小写字母
func IsLower(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[a-z]+$`).MatchString(cast.ToString(value))
}

// IsUpper - 是否为大写字母
func IsUpper(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[A-Z]+$`).MatchString(cast.ToString(value))
}

// IsSpace - 是否为空白字符 - （空格、制表符、换页符等）
func IsSpace(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\s]+$`).MatchString(cast.ToString(value))
}

// IsXdigit - 是否为十六进制字符 - （0-9、a-f、A-F）
func IsXdigit(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\da-fA-F]+$`).MatchString(cast.ToString(value))
}

// IsActiveUrl - 是否为有效的域名或者IP
func IsActiveUrl(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^([a-z0-9-]+\.)+[a-z]{2,6}$`).MatchString(cast.ToString(value))
}

// IsIp - 是否为IP
func IsIp(ip any) (ok bool) {
	if ip == nil {
		return false
	}
	return regexp.MustCompile(`(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`).MatchString(cast.ToString(ip))
}

// IsUrl - 是否为URL
func IsUrl(url any) (ok bool) {
	if url == nil {
		return false
	}
	return regexp.MustCompile(`^((https|http|ftp|rtsp|mms)?:\/\/)[^\s]+`).MatchString(cast.ToString(url))
}

// IsIdCard - 是否为有效的身份证号码
func IsIdCard(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)`).MatchString(cast.ToString(value))
}

// IsMacAddr - 是否为有效的MAC地址
func IsMacAddr(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^([A-Fa-f0-9]{2}:){5}[A-Fa-f0-9]{2}$`).MatchString(cast.ToString(value))
}

// IsZip - 是否为有效的邮政编码
func IsZip(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[1-9]\d{5}$`).MatchString(cast.ToString(value))
}

// IsString - 是否为字符串
func IsString(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.String
}

// IsSlice - 是否为切片
func IsSlice(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Slice
}

// IsArray - 是否为数组
func IsArray(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Array
}

// IsJsonString - 是否为json字符串
func IsJsonString(value any) (ok bool) {
	if value == nil {
		return false
	}
	return regexp.MustCompile(`^[\{\[].*[\}\]]$`).MatchString(cast.ToString(value))
}

// IsMap - 是否为map
func IsMap(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Map
}

// IsSliceSlice - 是否为二维切片
func IsSliceSlice(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Slice && reflect.TypeOf(value).Elem().Kind() == reflect.Slice
}

// IsMapAny - 是否为[]map[string]any
func IsMapAny(value any) (ok bool) {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Map && reflect.TypeOf(value).Elem().Kind() == reflect.Interface
}
