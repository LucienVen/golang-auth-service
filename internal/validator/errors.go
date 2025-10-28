package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError 验证错误结构
type ValidationError struct {
	Field   string `json:"field"`   // 字段名
	Message string `json:"message"` // 错误消息
	Tag     string `json:"tag"`     // 验证标签
	Value   string `json:"value"`   // 输入值
}

// Error 实现error接口
func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Message)
}

// ValidationResult 验证结果
type ValidationResult struct {
	Valid  bool              `json:"valid"`  // 是否通过验证
	Errors []ValidationError `json:"errors"` // 错误列表
}

// IsEmpty 检查验证结果是否为空
func (vr ValidationResult) IsEmpty() bool {
	return len(vr.Errors) == 0
}

// FirstError 获取第一个错误
func (vr ValidationResult) FirstError() *ValidationError {
	if len(vr.Errors) > 0 {
		return &vr.Errors[0]
	}
	return nil
}

// ToMap 转换为map格式
func (vr ValidationResult) ToMap() map[string]string {
	result := make(map[string]string)
	for _, err := range vr.Errors {
		result[err.Field] = err.Message
	}
	return result
}

// ToSlice 转换为切片格式
func (vr ValidationResult) ToSlice() []string {
	var result []string
	for _, err := range vr.Errors {
		result = append(result, err.Error())
	}
	return result
}

// CustomValidationMessages 自定义验证消息映射
var CustomValidationMessages = map[string]string{
	"required":                "%s不能为空",
	"min":                     "%s长度不能少于%s",
	"max":                     "%s长度不能超过%s",
	"email":                   "邮箱格式不正确",
	"username_valid":          "用户名只能包含字母、数字、下划线，且以字母开头",
	"phone_valid":             "手机号格式不正确",
	"strong_password":         "密码必须包含大小写字母和数字",
	"account_type":            "账号格式不正确（请输入邮箱、手机号或用户名）",
	"duplicate_field":         "%s不能与其他字段相同",
	"same_as_account":         "密码不能与账号相同",
	"same_as_old":             "新密码不能与旧密码相同",
	"required_account":        "用户名、邮箱、手机号至少需要提供一个",
	"required_field":          "至少需要提供一个要更新的字段",
	"account_exists":          "账号已存在",
	"account_not_found":       "账号不存在",
	"account_disabled":        "账号已被禁用",
	"account_not_activated":   "账号未激活",
	"password_incorrect":      "密码错误",
	"weak_password":           "密码强度不足",
	"invalid_email":           "邮箱格式不正确",
	"invalid_phone":           "手机号格式不正确",
	"invalid_username":        "用户名格式不正确",
	"field_too_long":          "%s长度不能超过%s",
	"field_too_short":         "%s长度不能少于%s",
}

// GetCustomValidationMessage 获取自定义验证消息
func GetCustomValidationMessage(tag string) string {
	if msg, exists := CustomValidationMessages[tag]; exists {
		return msg
	}
	return "字段验证失败"
}

// FormatValidationMessage 格式化验证消息
func FormatValidationMessage(field, tag, param string) string {
	template := GetCustomValidationMessage(tag)
	message := strings.ReplaceAll(template, "%s", field)
	if param != "" {
		message = strings.ReplaceAll(message, "%s", param)
	}
	return message
}

// ParseValidationError 解析验证错误
func ParseValidationError(err error) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		result.Valid = false
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()
			value := fmt.Sprintf("%v", e.Value())

			// 生成错误消息
			message := FormatValidationMessage(field, tag, param)

			result.Errors = append(result.Errors, ValidationError{
				Field:   field,
				Message: message,
				Tag:     tag,
				Value:   value,
			})
		}
	} else {
		// 处理非验证错误
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "general",
			Message: err.Error(),
			Tag:     "general",
			Value:   "",
		})
	}

	return result
}

// NewValidationError 创建验证错误
func NewValidationError(field, message, tag string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
		Tag:     tag,
		Value:   "",
	}
}

// NewValidationResult 创建验证结果
func NewValidationResult(valid bool, errors ...ValidationError) ValidationResult {
	return ValidationResult{
		Valid:  valid,
		Errors: errors,
	}
}

// ValidateAndFormat 验证并格式化错误
func ValidateAndFormat(req interface{}) ValidationResult {
	err := ValidateReq(req)
	if err != nil {
		return ParseValidationError(err)
	}
	return ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}
}

// ValidateBusiness 验证业务逻辑
func ValidateBusiness(errors ...ValidationError) ValidationResult {
	if len(errors) > 0 {
		return ValidationResult{
			Valid:  false,
			Errors: errors,
		}
	}
	return ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}
}