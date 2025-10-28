package errors

import (
	"fmt"
)

// AppError 应用错误接口
type AppError interface {
	error
	Code() int
	Message() string
}

// BaseError 基础错误结构
type BaseError struct {
	code    int
	message string
}

// Error 实现error接口
func (e *BaseError) Error() string {
	return e.message
}

// Code 返回错误码
func (e *BaseError) Code() int {
	return e.code
}

// Message 返回错误消息
func (e *BaseError) Message() string {
	return e.message
}

// NewError 创建新的应用错误
func NewError(code int, message string) AppError {
	return &BaseError{
		code:    code,
		message: message,
	}
}

// ValidationError 验证错误
type ValidationError struct {
	BaseError
	Field   string // 字段名
	Value   string // 字段值
	Tag     string // 验证标签
}

// NewValidationError 创建验证错误
func NewValidationError(field, message, tag string) *ValidationError {
	return &ValidationError{
		BaseError: BaseError{
			code:    ErrCodeParamInvalid,
			message: message,
		},
		Field: field,
		Value: "",
		Tag:   tag,
	}
}

// UserNotFoundError 用户不存在错误
type UserNotFoundError struct {
	BaseError
}

// NewUserNotFoundError 创建用户不存在错误
func NewUserNotFoundError(message string) *UserNotFoundError {
	return &UserNotFoundError{
		BaseError: BaseError{
			code:    ErrCodeUserNotFound,
			message: message,
		},
	}
}

// UserExistsError 用户已存在错误
type UserExistsError struct {
	BaseError
}

// NewUserExistsError 创建用户已存在错误
func NewUserExistsError(message string) *UserExistsError {
	return &UserExistsError{
		BaseError: BaseError{
			code:    ErrCodeUserAlreadyExists,
			message: message,
		},
	}
}

// UserDisabledError 用户已禁用错误
type UserDisabledError struct {
	BaseError
}

// NewUserDisabledError 创建用户已禁用错误
func NewUserDisabledError(message string) *UserDisabledError {
	return &UserDisabledError{
		BaseError: BaseError{
			code:    ErrCodeUserDisabled,
			message: message,
		},
	}
}

// UserCreationError 用户创建错误
type UserCreationError struct {
	BaseError
}

// NewUserCreationError 创建用户创建错误
func NewUserCreationError(message string) *UserCreationError {
	return &UserCreationError{
		BaseError: BaseError{
			code:    ErrCodeSystemError,
			message: message,
		},
	}
}

// InvalidCredentialsError 无效凭据错误
type InvalidCredentialsError struct {
	BaseError
}

// NewInvalidCredentialsError 创建无效凭据错误
func NewInvalidCredentialsError(message string) *InvalidCredentialsError {
	return &InvalidCredentialsError{
		BaseError: BaseError{
			code:    ErrCodePasswordInvalid,
			message: message,
		},
	}
}

// TokenGenerationError 令牌生成错误
type TokenGenerationError struct {
	BaseError
}

// NewTokenGenerationError 创建令牌生成错误
func NewTokenGenerationError(message string) *TokenGenerationError {
	return &TokenGenerationError{
		BaseError: BaseError{
			code:    ErrCodeSystemError,
			message: message,
		},
	}
}

// TokenRefreshError 令牌刷新错误
type TokenRefreshError struct {
	BaseError
}

// NewTokenRefreshError 创建令牌刷新错误
func NewTokenRefreshError(message string) *TokenRefreshError {
	return &TokenRefreshError{
		BaseError: BaseError{
			code:    ErrCodeRefreshTokenInvalid,
			message: message,
		},
	}
}

// TokenRevocationError 令牌撤销错误
type TokenRevocationError struct {
	BaseError
}

// NewTokenRevocationError 创建令牌撤销错误
func NewTokenRevocationError(message string) *TokenRevocationError {
	return &TokenRevocationError{
		BaseError: BaseError{
			code:    ErrCodeSystemError,
			message: message,
		},
	}
}

// InvalidTokenError 无效令牌错误
type InvalidTokenError struct {
	BaseError
}

// NewInvalidTokenError 创建无效令牌错误
func NewInvalidTokenError(message string) *InvalidTokenError {
	return &InvalidTokenError{
		BaseError: BaseError{
			code:    ErrCodeTokenInvalid,
			message: message,
		},
	}
}

// SessionError 会话错误
type SessionError struct {
	BaseError
}

// NewSessionError 创建会话错误
func NewSessionError(message string) *SessionError {
	return &SessionError{
		BaseError: BaseError{
			code:    ErrCodeSystemError,
			message: message,
		},
	}
}

// InternalError 内部错误
type InternalError struct {
	BaseError
}

// NewInternalError 创建内部错误
func NewInternalError(message string) *InternalError {
	return &InternalError{
		BaseError: BaseError{
			code:    ErrCodeSystemError,
			message: message,
		},
	}
}

// IsError 检查错误是否为指定类型
func IsError(err error, target AppError) bool {
	if err == nil {
		return false
	}
	if appErr, ok := err.(AppError); ok {
		return appErr.Code() == target.Code()
	}
	return false
}

// WrapError 包装错误
func WrapError(err error, code int, message string) AppError {
	if err == nil {
		return NewError(code, message)
	}
	return &BaseError{
		code:    code,
		message: fmt.Sprintf("%s: %v", message, err),
	}
}