package utils

import (
	"fmt"
	"math/rand"
	"time"
	
	"github.com/spf13/cast"
)

// Gen - 生成
var Gen *GenClass

type GenClass struct {}

// SerialNo 生成指定前缀和长度的序列号
/**
 * prefix: 前缀字符串
 * len: 序列号总长度
 * 格式: 前缀 + 日期(8位) + 时间(6位) + 随机数(动态长度)
 * 当指定长度小于前缀+日期+时间的长度时，会截断到指定长度
 */
func (this *GenClass) SerialNo(prefix any, length int) string {
	
	// 种子
	seed   := Hash.Sum32(fmt.Sprintf("%s-%s-%d-%d", cast.ToString(prefix), Get.Mac(), Get.Pid(), time.Now().UnixNano()))
	// 使用当前时间戳创建随机数生成器
	source := rand.New(rand.NewSource(cast.ToInt64(seed)))
	
	// 获取当前日期时间
	now := time.Now()
	datePart := now.Format("20060102") // 8位日期
	timePart := now.Format("150405")   // 6位时间
	
	// 计算固定部分的总长度
	fixedPart   := cast.ToString(prefix) + datePart + timePart
	fixedLength := len(fixedPart)
	
	var serialNo string
	
	// 如果指定长度小于等于固定部分长度，直接截断到指定长度
	if length <= fixedLength {
		
		serialNo = fixedPart[:length]
		
	} else {
		
		// 计算需要的随机数长度
		randomLength := length - fixedLength
		
		// 生成指定长度的随机数
		// 计算10的 randomLength 次方，作为随机数的上限
		maxLimit := 1
		for i := 0; i < randomLength; i++ { maxLimit *= 10 }
		
		// 生成随机数并格式化到指定长度
		randomPart := fmt.Sprintf("%0" + fmt.Sprintf("%dd", randomLength), source.Intn(maxLimit))
		
		// 组合所有部分
		serialNo = fixedPart + randomPart
	}
	
	return serialNo
}