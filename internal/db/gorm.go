package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/LucienVen/golang-auth-service/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormDB GORM 数据库管理器
type GormDB struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewGormDB 创建新的 GORM 连接
func NewGormDB(cfg *config.Config) *GormDB {
	return &GormDB{
		cfg: cfg,
	}
}

// Connect 连接数据库
func (g *GormDB) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		g.cfg.DBUser,
		g.cfg.DBPass,
		g.cfg.DBHost,
		g.cfg.DBPort,
		g.cfg.DBName,
	)

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取通用数据库对象 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间

	g.db = db
	return nil
}

// Close 关闭数据库连接
func (g *GormDB) Close() error {
	if g.db != nil {
		sqlDB, err := g.db.DB()
		if err != nil {
			return fmt.Errorf("获取数据库实例失败: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("关闭数据库连接失败: %w", err)
		}
		log.Println("数据库连接已关闭")
	}
	return nil
}

// GetDB 获取数据库连接
func (g *GormDB) GetDB() *gorm.DB {
	return g.db
}

// Ping 检查数据库连接
func (g *GormDB) Ping() error {
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
func (g *GormDB) AutoMigrate(models ...interface{}) error {
	if g.db == nil {
		return fmt.Errorf("数据库未连接")
	}
	return g.db.AutoMigrate(models...)
}

// Transaction 事务处理
func (g *GormDB) Transaction(fc func(tx *gorm.DB) error) error {
	if g.db == nil {
		return fmt.Errorf("数据库未连接")
	}
	return g.db.Transaction(fc)
}

// Shutdown 实现 Shutdownable 接口
func (g *GormDB) Shutdown(ctx context.Context) error {
	return g.Close()
}
