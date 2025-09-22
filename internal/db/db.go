package db

import (
	"context"

	"gorm.io/gorm"
)

// DB 通用接口
// 只要实现 Connect() error 即可
// 可根据需要扩展更多方法
type DB interface {
	Connect() error
	Shutdown(ctx context.Context) error
	Ping() error
	GetDB() *gorm.DB
	GetType() string
}

// DB 类型
const (
	DBTypeMySQL = "mysql"
	DBTypePG    = "pg"
)
