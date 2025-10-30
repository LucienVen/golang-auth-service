# 中间件依赖管理规范

## 概述

本文档定义了 Go 认证服务项目中中间件依赖管理的规范和最佳实践，基于 Clean Architecture 原则和依赖注入模式。

## 核心原则

### 1. 单一职责原则
- **Container**: 只负责管理 Controller 实例
- **Service**: 负责业务逻辑，可被多个组件共享
- **Middleware**: 负责横切关注点（认证、日志、CORS等）
- **Router**: 负责路由配置和依赖协调

### 2. 依赖倒置原则
- 高层模块不依赖低层模块，都依赖抽象
- 具体实现依赖接口，而非接口依赖实现
- 依赖通过构造函数注入，而非内部获取

### 3. 依赖方向规则
```
Service → Controller
Service → Middleware
Router → Service → Controller
Router → Service → Middleware
```

## 架构模式

### 1. Container 模式
```go
// ✅ 正确：Container 只管理 Controller
type Container struct {
    Health *HealthController
    Auth   *AuthController
}

// ❌ 错误：Container 混合管理不同类型的组件
type Container struct {
    Auth        *AuthController
    AuthService service.AuthService // 违反单一职责
}
```

### 2. 依赖注入模式
```go
// ✅ 正确：在 Router 层统一创建和注入依赖
func NewRouter(appCtx *appcontext.AppContext) *Router {
    // 创建服务实例
    authService := service.NewAuthService(appCtx)

    // 创建中间件实例，注入服务依赖
    authMiddleware := middleware.NewAuthMiddleware(authService)

    // 创建控制器实例，注入服务依赖
    authController := controller.NewAuthController(authService)

    return &Router{
        engine:        gin.New(),
        controllers:   &controller.Container{Auth: authController},
        middlewares:   &middleware.Container{Auth: authMiddleware},
        authService:   authService,
    }
}
```

### 3. Service 共享模式
```go
// ✅ 正确：Service 可同时服务于多个组件
type AuthService interface {
    ValidateToken(ctx context.Context, token string) (*response.UserResponse, error)
    Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error)
    // ... 其他业务方法
}

// Controller 使用 Service
type AuthController struct {
    authService service.AuthService
}

// Middleware 使用相同的 Service
type AuthMiddleware struct {
    authService service.AuthService
}
```

## 禁止模式

### 1. 禁止循环依赖
```go
// ❌ 错误：Controller 和 Middleware 相互依赖
type AuthController struct {
    authMiddleware *middleware.AuthMiddleware // 错误依赖
}

type AuthMiddleware struct {
    authController *controller.AuthController // 错误依赖
}
```

### 2. 禁止跨层直接访问
```go
// ❌ 错误：Middleware 直接访问 Controller
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 不应该直接调用 Controller 方法
        m.authController.SomeMethod() // 错误
    }
}

// ❌ 错误：通过 Container 获取 Service
func (r *Router) setupAuthRoutes() {
    // 不应该通过 Container 获取 Service
    authMiddleware := middleware.NewAuthMiddleware(r.controllers.AuthService) // 错误
}
```

### 3. 禁止职责混乱
```go
// ❌ 错误：Container 持有不同类型的依赖
type Container struct {
    Controllers map[string]interface{}         // 类型不安全
    Services    map[string]interface{}         // 职责混乱
    Middlewares map[string]interface{}         // 职责混乱
}
```

## 最佳实践

### 1. 依赖创建和注入
```go
// ✅ 最佳实践：在应用启动时统一创建依赖
func InitializeApp(appCtx *appcontext.AppContext) (*Router, error) {
    // 1. 创建基础服务
    authService := service.NewAuthService(appCtx)

    // 2. 创建中间件
    authMiddleware := middleware.NewAuthMiddleware(authService)
    loggerMiddleware := middleware.NewLogger()

    // 3. 创建控制器
    authController := controller.NewAuthController(authService)
    healthController := controller.NewHealthController(appCtx.DB)

    // 4. 创建路由
    router := NewRouter(gin.New(), &controller.Container{
        Health: healthController,
        Auth:   authController,
    }, authMiddleware, loggerMiddleware)

    return router, nil
}
```

### 2. 接口隔离原则
```go
// ✅ 正确：Middleware 只依赖需要的接口
type TokenValidator interface {
    ValidateToken(ctx context.Context, token string) (*response.UserResponse, error)
}

type AuthMiddleware struct {
    tokenValidator TokenValidator // 最小接口依赖
}

// ✅ 正确：Controller 依赖完整的 Service 接口
type AuthController struct {
    authService service.AuthService // 完整服务接口
}
```

### 3. 错误处理和验证
```go
// ✅ 正确：在依赖创建时进行验证
func NewAuthMiddleware(authService service.AuthService) (*AuthMiddleware, error) {
    if authService == nil {
        return nil, errors.New("authService cannot be nil")
    }

    return &AuthMiddleware{
        authService: authService,
    }, nil
}
```

## 代码组织规范

### 1. 目录结构
```
internal/
├── controller/
│   ├── container.go      # Controller 容器
│   ├── auth_controller.go
│   └── health_controller.go
├── middleware/
│   ├── auth_middleware.go
│   ├── cors_middleware.go
│   └── logger_middleware.go
├── service/
│   ├── interfaces.go     # Service 接口定义
│   ├── auth_service.go
│   └── user_service.go
└── appcontext/
    └── appcontext.go     # 全局上下文
```

### 2. 命名规范
- **Container**: `XxxContainer` (如 `ControllerContainer`)
- **Service**: `XxxService` (如 `AuthService`)
- **Middleware**: `XxxMiddleware` (如 `AuthMiddleware`)
- **Interface**: 简洁明了 (如 `TokenValidator`)

### 3. 导入规范
```go
// ✅ 正确：按类型分组导入
import (
    // 标准库
    "context"
    "net/http"

    // 第三方库
    "github.com/gin-gonic/gin"

    // 项目内部包
    "github.com/LucienVen/golang-auth-service/internal/service"
    "github.com/LucienVen/golang-auth-service/internal/response"
)
```

## 测试规范

### 1. 依赖注入测试
```go
// ✅ 正确：使用 Mock 接口进行测试
func TestAuthMiddleware_RequireAuth(t *testing.T) {
    mockService := &MockAuthService{}
    middleware := NewAuthMiddleware(mockService)

    // 测试逻辑...
}

type MockAuthService struct {
    service.AuthService
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*response.UserResponse, error) {
    // Mock 实现
}
```

### 2. 集成测试
```go
// ✅ 正确：集成测试使用真实依赖
func TestAuthIntegration(t *testing.T) {
    appCtx := &appcontext.AppContext{
        DB:    setupTestDB(t),
        Redis: setupTestRedis(t),
    }

    router := NewRouter(appCtx)
    // 集成测试逻辑...
}
```

## 性能考虑

### 1. 单例模式
- Service 实例应该是单例，避免重复创建
- Middleware 实例应该是单例，在多个路由间共享
- Controller 实例可以是单例，因为它们是无状态的

### 2. 生命周期管理
```go
// ✅ 正确：明确生命周期
type Router struct {
    engine      *gin.Engine
    controllers *controller.Container  // 应用生命周期
    middlewares *middleware.Container  // 应用生命周期
    authService service.AuthService    // 应用生命周期
}
```

## 扩展指南

### 1. 添加新的中间件
1. 定义中间件结构体和构造函数
2. 在 Router 中创建中间件实例
3. 在相应路由中应用中间件
4. 确保依赖关系清晰

### 2. 添加新的服务
1. 定义服务接口
2. 实现服务逻辑
3. 在 Router 中创建服务实例
4. 注入到需要的服务组件中

### 3. 重构现有代码
1. 识别依赖关系问题
2. 重构为符合规范的依赖注入模式
3. 更新相关测试
4. 验证功能正确性

## 总结

遵循本规范可以确保：
- **清晰的依赖关系**: 避免循环依赖和职责混乱
- **良好的可测试性**: 依赖注入使单元测试更容易
- **代码可维护性**: 职责分离使代码更容易理解和修改
- **架构一致性**: 统一的模式便于团队协作

本规范适用于所有 Go 项目的中间件开发，特别是需要复杂依赖管理的 Web 服务。