package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// 全局验证器实例
var valid *validator.Validate

// Init 初始化验证器，注册自定义验证规则
func Init() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		valid = v

		// 注册字段级验证器
		if err := registerFieldValidators(); err != nil {
			return fmt.Errorf("注册字段验证器失败: %w", err)
		}

		// 注册结构体验证器
		if err := registerStructValidators(); err != nil {
			return fmt.Errorf("注册结构体验证器失败: %w", err)
		}

		return nil
	}
	return fmt.Errorf("获取gin验证器引擎失败")
}

// ValidateReq 验证请求结构体的统一入口
func ValidateReq(req any) error {
	if valid == nil {
		return fmt.Errorf("验证器未初始化")
	}
	return valid.Struct(req)
}

// GetValidator 获取全局验证器实例
func GetValidator() *validator.Validate {
	return valid
}

// GetValidationError 获取详细的验证错误信息
func GetValidationError(err error) map[string]string {
	errorMap := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()

			switch tag {
			case "required":
				errorMap[field] = fmt.Sprintf("%s不能为空", field)
			case "min":
				errorMap[field] = fmt.Sprintf("%s长度不能少于%s", field, param)
			case "max":
				errorMap[field] = fmt.Sprintf("%s长度不能超过%s", field, param)
			case "email":
				errorMap[field] = "邮箱格式不正确"
			case "username_valid":
				errorMap[field] = "用户名只能包含字母、数字、下划线，且以字母开头"
			case "phone_valid":
				errorMap[field] = "手机号格式不正确"
			case "strong_password":
				errorMap[field] = "密码必须包含大小写字母和数字"
			case "account_type":
				errorMap[field] = "账号格式不正确"
			default:
				errorMap[field] = fmt.Sprintf("%s验证失败", field)
			}
		}
	} else {
		errorMap["general"] = err.Error()
	}

	return errorMap
}

// FieldNames 获取结构体的字段名（小写带下划线）
func FieldNames(obj any) []string {
	t := reflect.TypeOf(obj)
	var fields []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			// 处理 json:"field_name" 格式
			parts := strings.Split(jsonTag, ",")
			if len(parts) > 0 && parts[0] != "" {
				fields = append(fields, parts[0])
			}
		}
	}

	return fields
}

// registerFieldValidators 注册字段级验证器
func registerFieldValidators() error {
	// 用户名验证
	if err := valid.RegisterValidation("username_valid", usernameValid); err != nil {
		return err
	}

	// 手机号验证
	if err := valid.RegisterValidation("phone_valid", phoneValid); err != nil {
		return err
	}

	// 强密码验证
	if err := valid.RegisterValidation("strong_password", strongPassword); err != nil {
		return err
	}

	// 账号类型验证（自动检测邮箱/手机号/用户名）
	if err := valid.RegisterValidation("account_type", accountType); err != nil {
		return err
	}

	return nil
}

// registerStructValidators 注册结构体验证器
func registerStructValidators() error {
	// 这里会注册结构体验证器，但需要避免循环导入
	// 具体注册会在请求结构体更新时完成
	return nil
}