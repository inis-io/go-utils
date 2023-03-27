package utils

import (
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

// PasswordCreate - 创建密码
func PasswordCreate(password any) (result string) {

	item, _ := bcrypt.GenerateFromPassword([]byte(password.(string)), bcrypt.MinCost)

	return string(item)
}

// PasswordVerify - 验证密码
func PasswordVerify(encode any, password any) (ok bool) {
	if err := bcrypt.CompareHashAndPassword([]byte(cast.ToString(encode)), []byte(cast.ToString(password))); err != nil {
		return false
	}
	return true
}