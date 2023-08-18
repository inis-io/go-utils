package utils

import (
	"github.com/spf13/cast"
	"math/rand"
	"time"
)

// RandString - 生成随机字符串
func RandString(length int, chars ...string) (result string) {

	var charset string

	if IsEmpty(chars) {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	} else {
		charset = chars[0]
	}
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	item := make([]byte, length)
	for i := range item {
		item[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(item)
}

// RandInt - 生成随机整数
func RandInt(max int, min ...int) (result int) {
	if IsEmpty(min) {
		min = []int{0}
	}
	if max <= min[0] {
		// 交换两个数
		max, min[0] = min[0], max
	}
	if max == min[0] {
		return max
	}
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(max-min[0]) + min[0]
}

// RandSlice - 返回随机的指定长度的切片
func RandSlice(slice []any, limit any) (result []any) {

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

// RandMapSlice - 打乱切片顺序
func RandMapSlice(slice []map[string]any) (result []map[string]any) {

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
