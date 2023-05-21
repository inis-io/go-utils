package utils

import (
	"github.com/spf13/cast"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// VersionGo - 获取当前go版本号
func VersionGo() (result string) {
	return strings.Replace(runtime.Version(), "go", "", -1)
}

// VersionCompare - 版本号比对
/**
 * @param v1 string - 小版本号
 * @param v2 string - 大版本号
 * @return int - 0: 相等，1: v1 < v2，-1: v1 > v2
 * @example：
 * 	utils.VersionCompare("1.2.0", "1.0.0") // 1
 */
func VersionCompare(v1, v2 any) (result int) {

	rule := regexp.MustCompile(`\d+`)
	v1Arr := rule.FindAllString(cast.ToString(v1), -1)
	v2Arr := rule.FindAllString(cast.ToString(v2), -1)

	for i := 0; i < len(v1Arr) || i < len(v2Arr); i++ {
		v1Num := 0
		v2Num := 0

		if i < len(v1Arr) {
			v1Num, _ = strconv.Atoi(v1Arr[i])
		}

		if i < len(v2Arr) {
			v2Num, _ = strconv.Atoi(v2Arr[i])
		}

		if v2Num > v1Num {
			return 1
		} else if v2Num < v1Num {
			return -1
		}
	}

	return 0
}
