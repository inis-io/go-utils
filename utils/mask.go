package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// Mask - 脱敏
var Mask *MaskClass

type MaskClass struct {}


// Phone - 手机号脱敏
func (this *MaskClass) Phone(phone string) string {

	if Is.Empty(phone) { return "" }

	// 简单的手机号格式正则校验，实际中可以使用更严格的校验方式
	match, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)

	if !match { return phone }

	return fmt.Sprintf("%s****%s", phone[:3], phone[7:])
}

// Email - 邮箱脱敏
func (this *MaskClass) Email(email string) string {

	if Is.Empty(email) { return "" }

	// 简单的邮箱格式正则校验，实际中可以使用更严格的校验方式
	match, _ := regexp.MatchString(`^[\w-]+(\.[\w-]+)*@[\w-]+(\.[\w-]+)+$`, email)

	if !match { return email }

	emailArr    := strings.Split(email, "@")
	emailName   := emailArr[0]
	emailDomain := emailArr[1]

	if len(emailName) > 3 {
		emailName = emailName[:3] + "****"
	}

	return fmt.Sprintf("%s@%s", emailName, emailDomain)
}

// IDCard - 身份证号脱敏
func (this *MaskClass) IDCard(idCard string) string {

	if Is.Empty(idCard) { return "" }

	// 简单的身份证号格式正则校验，实际中可以使用更严格的校验方式
	match, _ := regexp.MatchString(`^\d{17}[\dXx]$`, idCard)

	if !match { return idCard }

	return fmt.Sprintf("%s****%s", idCard[:4], idCard[14:])
}

// BankCard - 银行卡号脱敏
func (this *MaskClass) BankCard(bankCard string) string {

	if Is.Empty(bankCard) { return "" }

	return fmt.Sprintf("%s****%s", bankCard[:4], bankCard[len(bankCard)-4:])
}

// Password - 密码脱敏
func (this *MaskClass) Password(password string) string {

	if Is.Empty(password) { return "" }

	return "******"
}

// Custom - 自定义脱敏
func (this *MaskClass) Custom(str string, start, end int) string {

	if Is.Empty(str) { return "" }

	if start < 0 || end < 0 || start > end || end > len(str) {
		return str
	}

	return fmt.Sprintf("%s****%s", str[:start], str[end:])
}