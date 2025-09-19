package log

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitLogger(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("APP_NAME", "test_app")
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_NAME")
	defer os.Unsetenv("APP_ENV")

	// 初始化日志
	InitLogger()

	// 验证 logger 和 sugarLogger 是否已初始化
	assert.NotNil(t, logger)
	assert.NotNil(t, sugarLogger)

	// 测试日志写入
	Info("test info message", zap.String("key", "value"))
	Warn("test warn message", zap.String("key", "value"))
	Error("test error message", zap.String("key", "value"))
	Debug("test debug message", zap.String("key", "value"))

	// 测试格式化日志
	Infof("test info message: %s", "value")
	Warnf("test warn message: %s", "value")
	Errorf("test error message: %s", "value")
	Debugf("test debug message: %s", "value")

	// 验证日志文件是否创建
	execPath, err := os.Executable()
	assert.NoError(t, err)

	logDir := filepath.Join(filepath.Dir(execPath), "logs")
	logFile := filepath.Join(logDir, "test_app.log")

	_, err = os.Stat(logFile)
	assert.NoError(t, err)

	// 清理
	Sync()
}

func TestLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		expected zapcore.Level
	}{
		{
			name:     "Production environment",
			env:      "production",
			expected: zap.InfoLevel,
		},
		{
			name:     "Development environment",
			env:      "development",
			expected: zap.DebugLevel,
		},
		{
			name:     "Test environment",
			env:      "test",
			expected: zap.DebugLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := getLogLevel(tt.env)
			assert.Equal(t, tt.expected, level)
		})
	}
}
