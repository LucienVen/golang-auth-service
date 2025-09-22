package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/LucienVen/golang-auth-service/internal/appcontext"
	"net/http"
	"time"

	"github.com/LucienVen/golang-auth-service/api"
	"github.com/LucienVen/golang-auth-service/config"
	"github.com/LucienVen/golang-auth-service/internal/db"
	"github.com/LucienVen/golang-auth-service/pkg/log"
	"github.com/gin-gonic/gin"
)

// Application 应用程序核心结构
type Application struct {
	config   *config.Config
	db       db.DB
	router   *api.Router
	health   *db.HealthChecker
	server   *http.Server
	logger   *log.Logger
	ctx      context.Context
	shutdown *ShutdownManager
	redis    *db.RedisClient
	appCtx   *appcontext.AppContext
}

// NewApplication 创建新的应用实例
func NewApplication() *Application {
	return &Application{
		ctx:      context.Background(),
		shutdown: NewShutdownManager(),
	}
}

// Start 启动应用
func (app *Application) Start() error {
	// 1. 配置初始化
	if err := app.initConfig(); err != nil {
		return fmt.Errorf("配置初始化失败: %w", err)
	}

	// 打印数据库配置
	fmt.Printf("DB 配置: host=%s port=%s user=%s pass=%s name=%s\n",
		app.config.DBHost, app.config.DBPort, app.config.DBUser, app.config.DBPass, app.config.DBName)

	// 2. 日志初始化
	if err := app.initLogger(); err != nil {
		return fmt.Errorf("日志初始化失败: %w", err)
	}

	// 3. 数据库初始化
	if err := app.initDatabase(); err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}

	if err := app.initRedis(); err != nil {
		return fmt.Errorf("redis初始化失败: %w", err)
	}

	app.appCtx = &appcontext.AppContext{DB: app.db, Redis: app.redis}

	// 4. 健康检查初始化
	if err := app.initHealthCheck(); err != nil {
		return fmt.Errorf("健康检查初始化失败: %w", err)
	}

	// 5. 路由初始化
	if err := app.initRouter(); err != nil {
		return fmt.Errorf("路由初始化失败: %w", err)
	}

	// 6. 服务器初始化
	if err := app.initServer(); err != nil {
		return fmt.Errorf("服务器初始化失败: %w", err)
	}

	// 7. 启动服务器
	go func() {
		if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("服务器启动失败: %v", err)
		}
	}()

	return nil
}

// Shutdown 关闭应用
func (app *Application) Shutdown(ctx context.Context) error {
	return app.shutdown.Shutdown(ctx)
}

// initConfig 初始化配置
func (app *Application) initConfig() error {
	config.Load(true)
	app.config = config.GetConfig()
	return nil
}

// initLogger 初始化日志
func (app *Application) initLogger() error {
	log.InitLogger()
	app.logger = log.GetLogger()
	return nil
}

// initDatabase 初始化数据库
func (app *Application) initDatabase() error {
	switch app.config.DBType {
	case "mysql":
		app.db = db.NewGormDB(app.config)
	case "pg":
		app.db = db.NewGormPGDB(app.config)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", app.config.DBType)
	}
	if err := app.db.Connect(); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	app.shutdown.Register(app.db)
	return nil
}

// initHealthCheck 初始化健康检查
func (app *Application) initHealthCheck() error {
	app.health = db.NewHealthChecker(app.db)
	app.health.Start(app.ctx, 30*time.Second)
	app.shutdown.Register(app.health)
	return nil
}

// initRouter 初始化路由
func (app *Application) initRouter() error {
	// 设置 Gin 模式
	gin.SetMode(app.config.GetGinMode())

	// 创建路由
	app.router = api.NewRouter(app.db)
	app.router.RegisterRoutes()
	return nil
}

// initServer 初始化服务器
func (app *Application) initServer() error {
	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.HTTPPort),
		Handler: app.router.GetEngine(),
	}
	app.shutdown.Register(app.server)
	return nil
}

func (app *Application) initRedis() error {
	app.redis = db.NewRedisClient(app.config.RedisAddr, app.config.RedisPassword, 0)
	if err := app.redis.Ping(app.ctx); err != nil {
		return err
	}
	app.shutdown.Register(app.redis)
	return nil
}
