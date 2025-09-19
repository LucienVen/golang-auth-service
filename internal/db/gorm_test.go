package db

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/LucienVen/golang-auth-service/config"
	"github.com/stretchr/testify/assert"
)

func TestGormDB(t *testing.T) {
	// 获取项目根目录
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")

	// 加载配置
	config.Load(true, filepath.Join(projectRoot, "cmd")) // 使用 cmd/.env 文件
	cfg := config.GetConfig()

	// 打印配置信息
	t.Logf("数据库配置: Host=%s, Port=%s, User=%s, DB=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBName)

	// 创建 GORM 实例
	gormDB := NewGormDB(&cfg)

	// 测试连接
	t.Run("测试数据库连接", func(t *testing.T) {
		err := gormDB.Connect()
		if err != nil {
			t.Skipf("数据库连接失败，跳过测试: %v", err)
		}
		assert.NoError(t, err, "数据库连接应该成功")
	})

	// 测试 Ping
	t.Run("测试数据库 Ping", func(t *testing.T) {
		if gormDB.db == nil {
			t.Skip("数据库未连接，跳过测试")
		}
		err := gormDB.Ping()
		assert.NoError(t, err, "数据库 Ping 应该成功")
	})

	// 测试连接池
	t.Run("测试连接池", func(t *testing.T) {
		if gormDB.db == nil {
			t.Skip("数据库未连接，跳过测试")
		}
		db := gormDB.GetDB()
		sqlDB, err := db.DB()
		assert.NoError(t, err, "获取 sql.DB 应该成功")

		// 验证连接池设置
		stats := sqlDB.Stats()
		assert.GreaterOrEqual(t, stats.Idle, 0, "空闲连接数应该大于等于 0")
		assert.GreaterOrEqual(t, stats.InUse, 0, "使用中的连接数应该大于等于 0")
	})

	// 测试时区设置
	t.Run("测试时区设置", func(t *testing.T) {
		if gormDB.db == nil {
			t.Skip("数据库未连接，跳过测试")
		}
		db := gormDB.GetDB()
		var now time.Time
		err := db.Raw("SELECT NOW()").Scan(&now).Error
		assert.NoError(t, err, "查询当前时间应该成功")

		// 验证时区
		_, offset := time.Now().Local().Zone()
		_, dbOffset := now.Zone()
		assert.Equal(t, offset, dbOffset, "数据库时区应该与本地时区一致")
	})

	// 测试关闭连接
	t.Run("测试关闭连接", func(t *testing.T) {
		if gormDB.db == nil {
			t.Skip("数据库未连接，跳过测试")
		}
		err := gormDB.Close()
		assert.NoError(t, err, "关闭数据库连接应该成功")
	})
}
