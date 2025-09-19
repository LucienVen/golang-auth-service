package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Config 应用配置结构体
type Config struct {
	// 可以添加更多配置项
	AppEnv   string `mapstructure:"APP_ENV"`
	AppName  string `mapstructure:"APP_NAME"`
	DBHost   string `mapstructure:"DB_HOST"`
	DBPort   string `mapstructure:"DB_PORT"`
	DBUser   string `mapstructure:"DB_USER"`
	DBPass   string `mapstructure:"DB_PASS"`
	DBName   string `mapstructure:"DB_NAME"`
	HTTPPort int64  `mapstructure:"HTTP_PORT"`
}

var (
	conf Config
	v    *viper.Viper
)

// Load 加载配置
// overload为true时，env文件变量覆盖系统变量
// overload为false时，系统变量覆盖env文件变量
// configPath为配置文件所在目录，如果为空则使用当前目录
func Load(overload bool, configPath ...string) {
	// 初始化 viper
	v = viper.New()
	v.SetConfigType("env")
	v.AutomaticEnv()

	// 设置环境变量前缀
	v.SetEnvPrefix("")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置默认值
	setDefaults()

	// 设置配置文件路径
	var envPath string
	if len(configPath) > 0 && configPath[0] != "" {
		// 1. 首先检查环境变量指定的配置文件路径
		if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
			envPath = configFile
		} else {
			// 2. 使用默认配置文件路径（cmd/.env）
			envPath = filepath.Join("cmd", ".env")
		}
	} else {
		envPath = ".env"
	}

	// 加载 .env 文件
	v.SetConfigFile(envPath)
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Warning: Error reading .env file from %s: %v\n", envPath, err)
	}

	// 根据 overload 参数决定优先级
	if !overload {
		// 系统变量覆盖 env 文件变量
		for _, env := range os.Environ() {
			pair := strings.SplitN(env, "=", 2)
			if len(pair) == 2 {
				v.Set(pair[0], pair[1])
			}
		}
	}

	// 绑定环境变量
	bindEnvs()

	// 解析配置到结构体
	if err := v.Unmarshal(&conf); err != nil {
		panic(fmt.Sprintf("unable to decode config: %v", err))
	}

	// 确保 AppEnv 为小写
	conf.AppEnv = strings.ToLower(conf.AppEnv)

	fmt.Printf("load config success from %s\n", envPath)
	fmt.Printf("config env: %s\n", conf.AppEnv)
	fmt.Printf("config port: %d\n", conf.HTTPPort)
}

// setDefaults 设置默认值
func setDefaults() {
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_NAME", "app")
	v.SetDefault("HTTP_PORT", 8080)
}

// bindEnvs 绑定环境变量
func bindEnvs() {
	// 绑定所有配置项到环境变量
	v.BindEnv("APP_ENV")
	v.BindEnv("APP_NAME")
	v.BindEnv("DB_HOST")
	v.BindEnv("DB_PORT")
	v.BindEnv("DB_USER")
	v.BindEnv("DB_PASS")
	v.BindEnv("DB_NAME")
	v.BindEnv("HTTP_PORT")
}

// GetConfig 获取配置
func GetConfig() *Config {
	return &conf
}

// GetGinMode 根据环境变量返回对应的 Gin 模式
func (c *Config) GetGinMode() string {
	switch c.AppEnv {
	case "production":
		return gin.ReleaseMode
	case "testing":
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}
