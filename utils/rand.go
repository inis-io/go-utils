package utils

import (
	"fmt"
	"github.com/spf13/cast"
	"math/rand"
	"time"
)

// Rand - 随机数
var Rand *RandClass

type RandClass struct {}

// Number - 生成指定长度的随机数
func (this *RandClass) Number(length any) (result string) {

	mac := Hash.Sum32(Get.Mac())
	pid := Get.Pid()
	nano := time.Now().UnixNano()

	// 生成一个随机种子
	seed := fmt.Sprintf("%v%d%d", mac, pid, nano)

	// 如果 种子 超过了 int64 的最大值
	if len(seed) > 19 {
		// 压缩种子
		seed = Hash.Sum32(seed)
	}

	// 种子长度不足 int64 的最大值，补足
	if len(seed) < 19 {
		seed = seed + Hash.Number(19-len(seed))
	}

	rand.NewSource(cast.ToInt64(seed))

	// 生成指定长度的随机数
	for i := 0; i < cast.ToInt(length); i++ {
		result += fmt.Sprintf("%d", rand.Intn(10))
	}

	return result
}

// String - 生成随机字符串
func (this *RandClass) String(length any, chars ...string) (result string) {

	var charset string

	if Is.Empty(chars) {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	} else {
		charset = chars[0]
	}
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	item := make([]byte, cast.ToInt(length))
	for i := range item {
		item[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(item)
}

// Code - 生成随机验证码
// number:数字, letter:字母, mix:混合
func (this *RandClass) Code(length any, mode ...string) (result string) {

	var charset string

	if Is.Empty(mode) {
		charset = "number"
	} else {
		charset = mode[0]
	}

	switch charset {
	case "number":
		return this.Number(length)
	case "letter":
		return this.String(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case "mix":
		return this.String(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	default:
		return this.Number(length)
	}
}

// Int - 生成随机整数
func (this *RandClass) Int(max any, min ...any) (result int) {
	if Is.Empty(min) {
		min = []any{0}
	}
	if cast.ToInt(max) <= cast.ToInt(min[0]) {
		// 交换两个数
		max, min[0] = min[0], cast.ToInt(max)
	}
	if max == min[0] {
		return cast.ToInt(max)
	}
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(cast.ToInt(max)-cast.ToInt(min[0])) + cast.ToInt(min[0])
}

// Slice - 返回随机的指定长度的切片
func (this *RandClass) Slice(slice []any, limit any) (result []any) {

	// 如果切片为空，直接返回
	if len(slice) == 0 {
		return slice
	}

	// 设置随机数种子
	rand.NewSource(time.Now().UnixNano())

	// 创建一个map用于存储选中的元素
	selected := make(map[any]bool)

	// 限制最大长度
	if cast.ToInt(limit) > len(slice) {
		limit = len(slice)
	}

	// 随机选择指定数量的不重复元素
	for len(selected) < cast.ToInt(limit) {
		index := rand.Intn(len(slice))
		selected[slice[index]] = true
	}

	// 将选中的元素存储到切片中
	for key := range selected {
		result = append(result, key)
	}

	return result
}

// MapSlice - 打乱切片顺序
func (this *RandClass) MapSlice(slice []map[string]any) (result []map[string]any) {

	// 如果切片为空，直接返回
	if len(slice) == 0 {
		return slice
	}

	// 设置随机数种子
	rand.NewSource(time.Now().UnixNano())

	// 打乱切片顺序
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	return slice
}
