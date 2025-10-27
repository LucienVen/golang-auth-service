package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// 验证错误类型
var (
	ErrEmailEmpty     = errors.New("邮箱不能为空")
	ErrEmailInvalid   = errors.New("邮箱格式不正确")
	ErrEmailTooLong   = errors.New("邮箱长度不能超过100个字符")
	ErrPhoneEmpty     = errors.New("手机号不能为空")
	ErrPhoneInvalid   = errors.New("手机号格式不正确，请输入有效的中国大陆手机号")
	ErrUsernameEmpty  = errors.New("用户名不能为空")
	ErrUsernameTooShort = errors.New("用户名长度不能少于3个字符")
	ErrUsernameTooLong = errors.New("用户名长度不能超过50个字符")
	ErrUsernameInvalid = errors.New("用户名只能包含字母、数字、下划线，且不能以下划线开头或结尾")
	ErrPasswordEmpty  = errors.New("密码不能为空")
	ErrAccountEmpty   = errors.New("账户不能为空")
	ErrAccountInvalid = errors.New("账户格式不正确")
)

// 预编译正则表达式
var (
	// 邮箱格式正则表达式（RFC 5322 兼容的简化版本）
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// 中国大陆手机号正则表达式：1开头，第二位3-9，共11位数字
	phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

	// 用户名正则表达式：字母数字下划线，不能以下划线开头或结尾
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_]*[a-zA-Z0-9]$`)

	// 账户识别正则表达式（邮箱或手机号）
	accountRegex = regexp.MustCompile(`^([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}|1[3-9]\d{9})$`)
)

// ========== 邮箱验证 ==========

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	// 清理输入
	email = strings.TrimSpace(email)

	// 检查是否为空
	if email == "" {
		return ErrEmailEmpty
	}

	// 检查长度
	if len(email) > 100 {
		return ErrEmailTooLong
	}

	// 检查格式
	if !emailRegex.MatchString(email) {
		return ErrEmailInvalid
	}

	// 额外检查：不能包含连续的点
	if strings.Contains(email, "..") {
		return errors.New("邮箱不能包含连续的点")
	}

	// 检查域名部分是否有效
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ErrEmailInvalid
	}

	localPart := parts[0]
	domainPart := parts[1]

	// 本地部分不能以点开始或结束
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return errors.New("邮箱本地部分不能以点开始或结束")
	}

	// 域名部分必须包含至少一个点
	if !strings.Contains(domainPart, ".") {
		return errors.New("邮箱域名部分无效")
	}

	return nil
}

// IsValidEmail 检查邮箱格式是否有效（不返回错误）
func IsValidEmail(email string) bool {
	return ValidateEmail(email) == nil
}

// ========== 手机号验证 ==========

// ValidatePhone 验证中国大陆手机号格式
func ValidatePhone(phone string) error {
	// 清理输入
	phone = strings.TrimSpace(phone)

	// 检查是否为空
	if phone == "" {
		return ErrPhoneEmpty
	}

	// 检查长度
	if len(phone) != 11 {
		return ErrPhoneInvalid
	}

	// 检查格式
	if !phoneRegex.MatchString(phone) {
		return ErrPhoneInvalid
	}

	// 检查是否全为数字
	for _, r := range phone {
		if r < '0' || r > '9' {
			return ErrPhoneInvalid
		}
	}

	return nil
}

// IsValidPhone 检查手机号格式是否有效（不返回错误）
func IsValidPhone(phone string) bool {
	return ValidatePhone(phone) == nil
}

// ========== 用户名验证 ==========

// ValidateUsername 验证用户名格式
func ValidateUsername(username string) error {
	// 清理输入
	username = strings.TrimSpace(username)

	// 检查是否为空
	if username == "" {
		return ErrUsernameEmpty
	}

	// 检查长度（按UTF-8字符计算）
	length := utf8.RuneCountInString(username)
	if length < 3 {
		return ErrUsernameTooShort
	}
	if length > 50 {
		return ErrUsernameTooLong
	}

	// 检查格式
	if !usernameRegex.MatchString(username) {
		return ErrUsernameInvalid
	}

	// 检查是否全为下划线
	if strings.Trim(username, "_") == "" {
		return errors.New("用户名不能全为下划线")
	}

	return nil
}

// IsValidUsername 检查用户名格式是否有效（不返回错误）
func IsValidUsername(username string) bool {
	return ValidateUsername(username) == nil
}

// ========== 密码验证（复用password工具） ==========

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	return CheckPasswordStrength(password)
}

// IsValidPassword 检查密码强度是否足够（不返回错误）
func IsValidPassword(password string) bool {
	return ValidatePassword(password) == nil
}

// ========== 账户验证（支持邮箱/手机号/用户名） ==========

// ValidateAccount 验证账户格式（支持邮箱、手机号、用户名）
func ValidateAccount(account string) error {
	// 清理输入
	account = strings.TrimSpace(account)

	// 检查是否为空
	if account == "" {
		return ErrAccountEmpty
	}

	// 尝试匹配邮箱
	if IsValidEmail(account) {
		return nil
	}

	// 尝试匹配手机号
	if IsValidPhone(account) {
		return nil
	}

	// 尝试匹配用户名
	if IsValidUsername(account) {
		return nil
	}

	return ErrAccountInvalid
}

// IsValidAccount 检查账户格式是否有效（不返回错误）
func IsValidAccount(account string) bool {
	return ValidateAccount(account) == nil
}

// DetectAccountType 检测账户类型
// 返回: "email", "phone", "username", "unknown"
func DetectAccountType(account string) string {
	account = strings.TrimSpace(account)

	if IsValidEmail(account) {
		return "email"
	}

	if IsValidPhone(account) {
		return "phone"
	}

	if IsValidUsername(account) {
		return "username"
	}

	return "unknown"
}

// ========== 通用验证工具 ==========

// ValidateRequired 验证必填字段
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s不能为空", fieldName)
	}
	return nil
}

// ValidateLength 验证字符串长度
func ValidateLength(value string, minLength, maxLength int) error {
	length := utf8.RuneCountInString(strings.TrimSpace(value))
	if length < minLength {
		return fmt.Errorf("长度不能少于%d个字符", minLength)
	}
	if length > maxLength {
		return fmt.Errorf("长度不能超过%d个字符", maxLength)
	}
	return nil
}

// ValidatePattern 验证字符串格式
func ValidatePattern(value, pattern, errorMessage string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("正则表达式编译失败: %w", err)
	}

	if !regex.MatchString(value) {
		return errors.New(errorMessage)
	}

	return nil
}

// ValidateChoices 验证选择值是否在允许的范围内
func ValidateChoices(value string, allowedChoices []string, fieldName string) error {
	for _, choice := range allowedChoices {
		if value == choice {
			return nil
		}
	}

	return fmt.Errorf("%s的值无效，允许的值为: %v", fieldName, allowedChoices)
}

// ========== 数据清理工具 ==========

// SanitizeString 清理字符串（去除首尾空白并标准化）
func SanitizeString(value string) string {
	return strings.TrimSpace(value)
}

// SanitizeEmail 清理邮箱（转换为小写并去除空白）
func SanitizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

// SanitizeUsername 清理用户名（去除空白并转换为小写）
func SanitizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

// SanitizePhone 清理手机号（只保留数字）
func SanitizePhone(phone string) string {
	// 移除所有非数字字符
	regex := regexp.MustCompile(`[^\d]`)
	return regex.ReplaceAllString(strings.TrimSpace(phone), "")
}

// ========== 批量验证工具 ==========

// ValidationErrors 批量验证错误
type ValidationErrors struct {
	Errors map[string]string
}

// NewValidationErrors 创建批量验证错误实例
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make(map[string]string),
	}
}

// Add 添加验证错误
func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors[field] = message
}

// HasErrors 检查是否有错误
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// Error 实现error接口
func (ve *ValidationErrors) Error() string {
	if !ve.HasErrors() {
		return ""
	}

	var messages []string
	for field, message := range ve.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", field, message))
	}

	return strings.Join(messages, "; ")
}

// ToMap 转换为map格式
func (ve *ValidationErrors) ToMap() map[string]string {
	return ve.Errors
}

// ValidateStruct 验证结构体字段（示例用法）
func ValidateStruct(data interface{}) *ValidationErrors {
	// 这里可以添加基于结构体标签的验证逻辑
	// 目前提供一个基础实现
	return NewValidationErrors()
}

// ========== 安全验证 ==========

// ValidateNoSQLInjection 检查潜在的SQL注入攻击
func ValidateNoSQLInjection(input string) error {
	// 基本的SQL注入检测
	dangerousPatterns := []string{
		"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_",
		"drop ", "delete ", "insert ", "update ", "select ",
		"exec ", "execute ", "union ", "script ",
	}

	inputLower := strings.ToLower(input)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(inputLower, pattern) {
			return fmt.Errorf("输入包含潜在的危险字符: %s", pattern)
		}
	}

	return nil
}

// ValidateNoXSS 检查潜在的XSS攻击
func ValidateNoXSS(input string) error {
	// 基本的XSS检测
	dangerousPatterns := []string{
		"<script", "</script>", "javascript:", "onload=", "onerror=",
		"onclick=", "onmouseover=", "onfocus=", "onblur=", "onchange=",
		"<iframe", "</iframe>", "<object", "</object>", "<embed",
		"vbscript:", "data:text/html", "expression(",
	}

	inputLower := strings.ToLower(input)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(inputLower, pattern) {
			return fmt.Errorf("输入包含潜在的危险脚本: %s", pattern)
		}
	}

	return nil
}