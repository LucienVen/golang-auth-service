package errors

// 系统级错误码 (1000-1999)
const (
	// 通用错误码
	ErrCodeSuccess         = 0    // 成功
	ErrCodeSystemError     = 1000 // 系统错误
	ErrCodeParamInvalid    = 1001 // 参数无效
	ErrCodeUnauthorized    = 1002 // 未授权
	ErrCodeForbidden       = 1003 // 禁止访问
	ErrCodeNotFound        = 1004 // 资源不存在
	ErrCodeDatabaseError   = 1005 // 数据库错误
	ErrCodeRedisError      = 1006 // Redis错误
	ErrCodeThirdPartyError = 1007 // 第三方服务错误
)

// 用户模块错误码 (2000-2999)
const (
	ErrCodeUserNotFound      = 2000 // 用户不存在
	ErrCodeUserAlreadyExists = 2001 // 用户已存在
	ErrCodePasswordInvalid   = 2002 // 密码无效
	ErrCodeUsernameInvalid   = 2003 // 用户名无效
	ErrCodeEmailInvalid      = 2004 // 邮箱无效
	ErrCodePhoneInvalid      = 2005 // 手机号无效
	ErrCodeUserDisabled      = 2006 // 用户已禁用
	ErrCodeUserDeleted       = 2007 // 用户已删除
)

// 认证模块错误码 (3000-3999)
const (
	ErrCodeTokenInvalid        = 3000 // Token无效
	ErrCodeTokenExpired        = 3001 // Token已过期
	ErrCodeTokenMissing        = 3002 // Token缺失
	ErrCodeRefreshTokenInvalid = 3003 // 刷新Token无效
	ErrCodeRefreshTokenExpired = 3004 // 刷新Token已过期
)

// 权限模块错误码 (4000-4999)
const (
	ErrCodePermissionDenied = 4000 // 权限不足
	ErrCodeRoleNotFound     = 4001 // 角色不存在
	ErrCodeRoleDisabled     = 4002 // 角色已禁用
	ErrCodeRoleDeleted      = 4003 // 角色已删除
)

// 错误码映射表
var codeMessages = map[int]string{
	// 系统级错误码
	ErrCodeSuccess:         "成功",
	ErrCodeSystemError:     "系统错误",
	ErrCodeParamInvalid:    "参数无效",
	ErrCodeUnauthorized:    "未授权",
	ErrCodeForbidden:       "禁止访问",
	ErrCodeNotFound:        "资源不存在",
	ErrCodeDatabaseError:   "数据库错误",
	ErrCodeRedisError:      "Redis错误",
	ErrCodeThirdPartyError: "第三方服务错误",

	// 用户模块错误码
	ErrCodeUserNotFound:      "用户不存在",
	ErrCodeUserAlreadyExists: "用户已存在",
	ErrCodePasswordInvalid:   "密码无效",
	ErrCodeUsernameInvalid:   "用户名无效",
	ErrCodeEmailInvalid:      "邮箱无效",
	ErrCodePhoneInvalid:      "手机号无效",
	ErrCodeUserDisabled:      "用户已禁用",
	ErrCodeUserDeleted:       "用户已删除",

	// 认证模块错误码
	ErrCodeTokenInvalid:        "Token无效",
	ErrCodeTokenExpired:        "Token已过期",
	ErrCodeTokenMissing:        "Token缺失",
	ErrCodeRefreshTokenInvalid: "刷新Token无效",
	ErrCodeRefreshTokenExpired: "刷新Token已过期",

	// 权限模块错误码
	ErrCodePermissionDenied: "权限不足",
	ErrCodeRoleNotFound:     "角色不存在",
	ErrCodeRoleDisabled:     "角色已禁用",
	ErrCodeRoleDeleted:      "角色已删除",
}

// GetMessage 获取错误码对应的消息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// IsSuccess 判断是否成功
func IsSuccess(code int) bool {
	return code == ErrCodeSuccess
}

// IsSystemError 判断是否系统错误
func IsSystemError(code int) bool {
	return code >= ErrCodeSystemError && code < 2000
}

// IsUserError 判断是否用户模块错误
func IsUserError(code int) bool {
	return code >= ErrCodeUserNotFound && code < 3000
}

// IsAuthError 判断是否认证模块错误
func IsAuthError(code int) bool {
	return code >= ErrCodeTokenInvalid && code < 4000
}

// IsPermissionError 判断是否权限模块错误
func IsPermissionError(code int) bool {
	return code >= ErrCodePermissionDenied && code < 5000
}
