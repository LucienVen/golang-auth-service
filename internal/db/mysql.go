package db

import (
	"fmt"
	"log"

	"github.com/LucienVen/golang-auth-service/config"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQL MySQL 数据库管理器
type MySQL struct {
	db  *sqlx.DB
	cfg *config.Config
}

// NewMySQL 创建新的 MySQL 连接
func NewMySQL(cfg *config.Config) *MySQL {
	return &MySQL{
		cfg: cfg,
	}
}

// Connect 连接数据库
func (m *MySQL) Connect() error {
	cfg := mysql.Config{
		User:                 m.cfg.DBUser,
		Passwd:               m.cfg.DBPass,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", m.cfg.DBHost, m.cfg.DBPort),
		DBName:               m.cfg.DBName,
		AllowNativePasswords: true,
	}

	db, err := sqlx.Connect("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 设置连接池
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	m.db = db
	return nil
}

// Close 关闭数据库连接
func (m *MySQL) Close() error {
	if m.db != nil {
		if err := m.db.Close(); err != nil {
			return fmt.Errorf("关闭数据库连接失败: %w", err)
		}
		log.Println("数据库连接已关闭")
	}
	return nil
}

// GetDB 获取数据库连接
func (m *MySQL) GetDB() *sqlx.DB {
	return m.db
}

// Ping 检查数据库连接
func (m *MySQL) Ping() error {
	if m.db == nil {
		return fmt.Errorf("数据库未连接")
	}
	return m.db.Ping()
}
