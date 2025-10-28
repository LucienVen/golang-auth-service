package service

import (
	"context"

	"github.com/LucienVen/golang-auth-service/internal/entity"
)

// UserRepository 用户数据访问接口（在service包中重新定义，避免循环依赖）
type UserRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error

	// 账号查询操作
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByAccount(ctx context.Context, account string) (*entity.User, error) // 统一账号查询

	// 管理操作
	SoftDelete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error

	// 查询操作
	List(ctx context.Context, offset, limit int) ([]*entity.User, error)
	Count(ctx context.Context) (int64, error)

	// 状态查询
	GetByStatus(ctx context.Context, status int, offset, limit int) ([]*entity.User, error)
	CountByStatus(ctx context.Context, status int) (int64, error)
}