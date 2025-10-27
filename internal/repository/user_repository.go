package repository

import (
	"github.com/LucienVen/golang-auth-service/internal/entity"
	"gorm.io/gorm"
)

// UserRepository 用户数据访问层接口
// 定义了用户相关的所有数据库操作
type UserRepository interface {
	// ========== 基础CRUD操作 ==========

	// Create 创建用户记录
	Create(user *entity.User) error

	// GetByID 根据用户ID获取用户信息
	GetByID(id uint) (*entity.User, error)

	// GetByEmail 根据邮箱获取用户信息
	GetByEmail(email string) (*entity.User, error)

	// GetByPhone 根据手机号获取用户信息
	GetByPhone(phone string) (*entity.User, error)

	// GetByUsername 根据用户名获取用户信息
	GetByUsername(username string) (*entity.User, error)

	// Update 更新用户信息
	Update(user *entity.User) error

	// Delete 软删除用户
	Delete(id uint) error

	// HardDelete 硬删除用户（物理删除）
	HardDelete(id uint) error

	// ========== 查询操作 ==========

	// GetByAccount 统一账户查询接口
	// account可以是邮箱、手机号或用户名
	GetByAccount(account string) (*entity.User, error)

	// ListByStatus 根据用户状态批量查询用户
	ListByStatus(status int8, limit, offset int) ([]*entity.User, error)

	// ListActive 查询所有活跃用户
	ListActive(limit, offset int) ([]*entity.User, error)

	// CountByStatus 统计指定状态的用户数量
	CountByStatus(status int8) (int64, error)

	// CountAll 统计所有用户数量
	CountAll() (int64, error)

	// SearchUsers 搜索用户（支持用户名、昵称、邮箱模糊搜索）
	SearchUsers(keyword string, limit, offset int) ([]*entity.User, error)

	// ========== 业务相关查询 ==========

	// ExistsByEmail 检查邮箱是否已存在
	ExistsByEmail(email string) (bool, error)

	// ExistsByPhone 检查手机号是否已存在
	ExistsByPhone(phone string) (bool, error)

	// ExistsByUsername 检查用户名是否已存在
	ExistsByUsername(username string) (bool, error)

	// GetUsersByLastLoginTime 根据最后登录时间查询用户
	GetUsersByLastLoginTime(startTime, endTime int64, limit, offset int) ([]*entity.User, error)

	// GetInactiveUsers 获取长期未登录的用户
	GetInactiveUsers(days int, limit, offset int) ([]*entity.User, error)

	// ========== 批量操作 ==========

	// CreateBatch 批量创建用户
	CreateBatch(users []*entity.User) error

	// UpdateBatch 批量更新用户状态
	UpdateBatch(userIDs []uint, status int8) error

	// DeleteBatch 批量软删除用户
	DeleteBatch(userIDs []uint) error

	// ========== 事务相关 ==========

	// WithTransaction 开始事务
	WithTransaction(tx *gorm.DB) UserRepository

	// ========== 统计和分析 ==========

	// GetRegistrationStats 获取用户注册统计
	GetRegistrationStats(days int) (map[string]int64, error)

	// GetStatusDistribution 获取用户状态分布
	GetStatusDistribution() (map[int8]int64, error)
}