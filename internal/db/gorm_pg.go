package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/LucienVen/golang-auth-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormPGDB GORM PostgreSQL 数据库管理器
// GormPGDB 实现了 db.DB 接口
type GormPGDB struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewGormPGDB 创建新的 GORM PostgreSQL 连接
func NewGormPGDB(cfg *config.Config) *GormPGDB {
	return &GormPGDB{
		cfg: cfg,
	}
}

// Connect 连接 PostgreSQL 数据库
func (g *GormPGDB) Connect() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		g.cfg.PGHost,
		g.cfg.PGPort,
		g.cfg.PGUser,
		g.cfg.PGPass,
		g.cfg.PGName,
	)

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接 PostgreSQL 数据库失败: %w", err)
	}

	// 获取通用数据库对象 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	g.db = db
	return nil
}

// Close 关闭数据库连接
func (g *GormPGDB) Close() error {
	if g.db != nil {
		sqlDB, err := g.db.DB()
		if err != nil {
			return fmt.Errorf("获取数据库实例失败: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("关闭数据库连接失败: %w", err)
		}
		log.Println("PostgreSQL 数据库连接已关闭")
	}
	return nil
}

// GetDB 获取数据库连接
func (g *GormPGDB) GetDB() *gorm.DB {
	return g.db
}

// Ping 检查数据库连接
func (g *GormPGDB) Ping() error {
	if g.db == nil {
		return fmt.Errorf("数据库未连接")
	}
	sqlDB, err := g.db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}
	return sqlDB.Ping()
}

// AutoMigrate 自动迁移数据库结构
func (g *GormPGDB) AutoMigrate(models ...interface{}) error {
	if g.db == nil {
		return fmt.Errorf("数据库未连接")
	}
	return g.db.AutoMigrate(models...)
}

// Transaction 事务处理
func (g *GormPGDB) Transaction(fc func(tx *gorm.DB) error) error {
	if g.db == nil {
		return fmt.Errorf("数据库未连接")
	}
	return g.db.Transaction(fc)
}

// Shutdown 实现 Shutdownable 接口
func (g *GormPGDB) Shutdown(ctx context.Context) error {
	return g.Close()
}

// GetType 返回数据库类型
func (g *GormPGDB) GetType() string {
	return DBTypePG
}
