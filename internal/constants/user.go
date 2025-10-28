package constants

// UserStatus 用户状态常量
const (
	UserStatusInactive int8 = 0 // 未激活
	UserStatusActive   int8 = 1 // 正常
	UserStatusDisabled int8 = 2 // 禁用
	UserStatusLogout   int8 = 3 // 注销
	UserStatusDeleted  int8 = 9 // 已删除
)