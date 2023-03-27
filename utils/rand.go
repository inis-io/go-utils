package utils

import (
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
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min[0]) + min[0]
}
