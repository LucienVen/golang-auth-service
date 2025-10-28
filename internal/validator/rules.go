package validator

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// usernameValid 验证用户名格式
// 规则：3-50位，只能包含字母、数字、下划线，且以字母开头
func usernameValid(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 50 {
		return false
	}

	// 必须以字母开头
	if !unicode.IsLetter(rune(username[0])) {
		return false
	}

	// 只能包含字母、数字、下划线
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, username)
	return matched
}

// phoneValid 验证手机号格式
// 规则：11位数字，以1开头
func phoneValid(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if len(phone) != 11 {
		return false
	}

	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

// strongPassword 验证密码强度
// 规则：8-72位，必须包含大小写字母和数字
func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 || len(password) > 72 {
		return false
	}

	var (
		hasUpper bool
		hasLower bool
		hasDigit bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

// accountType 验证账号类型（自动检测邮箱/手机号/用户名）
func accountType(fl validator.FieldLevel) bool {
	account := fl.Field().String()
	if account == "" {
		return false
	}

	// 检查是否为邮箱
	if isEmail(account) {
		return true
	}

	// 检查是否为手机号
	if isPhone(account) {
		return true
	}

	// 检查是否为用户名
	if isUsername(account) {
		return true
	}

	return false
}

// isEmail 检查是否为邮箱格式
func isEmail(email string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return matched
}

// isPhone 检查是否为手机号格式
func isPhone(phone string) bool {
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

// isUsername 检查是否为用户名格式
func isUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}

	// 必须以字母开头
	if !unicode.IsLetter(rune(username[0])) {
		return false
	}

	// 只能包含字母、数字、下划线
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, username)
	return matched
}