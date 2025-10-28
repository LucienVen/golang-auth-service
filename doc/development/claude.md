# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 项目概述

这是一个基于 Go 的认证服务（golang-auth-service），提供用户认证和管理功能。服务使用 Gin 框架构建，支持 MySQL 和 PostgreSQL 数据库，并使用 Redis 进行缓存。

## 命令

### 构建和运行
```bash
# 直接运行服务
go run cmd/main.go

# 构建二进制文件
go build -o bin/auth-service cmd/main.go

# 运行测试
go test ./...

# 运行测试并生成覆盖率报告
go test -cover ./...

# 运行特定测试文件
go test ./internal/db/gorm_test.go

# 运行特定目录下的测试
go test ./internal/utils
```

### 开发命令
```bash
# 格式化代码
go fmt ./...

# 检查潜在问题
go vet ./...

# 管理依赖
go mod tidy

# 下载所有依赖
go mod download

# 验证依赖
go mod verify

# 整理并打包依赖
go mod tidy && go mod vendor
```

### 测试
- 测试文件遵循 `*_test.go` 模式
- 主要测试文件包括：
  - `internal/db/gorm_test.go` - 数据库测试
  - `internal/utils/utils_test.go` - 工具函数测试
  - `pkg/jwt/jwt_test.go` - JWT 令牌测试
  - `pkg/log/log_test.go` - 日志测试
  - `test/main_test.go` - 集成测试

## 架构

### 目录结构
- `cmd/` - 应用程序入口点
- `api/` - HTTP 处理器和路由层
- `config/` - 配置管理
- `internal/` - 私有应用程序代码
  - `app/` - 核心应用程序逻辑和生命周期管理
  - `appcontext/` - 依赖注入的应用程序上下文
  - `controller/` - HTTP 控制器
  - `db/` - 数据库层（MySQL、PostgreSQL、Redis）
  - `entity/` - 数据库实体和模型
  - `errors/` - 自定义错误类型
  - `middleware/` - HTTP 中间件（日志等）
  - `request/` - 请求 DTO
  - `response/` - 响应 DTO
  - `utils/` - 工具函数
  - `sql/mysql/data_gen/` - 模拟数据生成工具
- `pkg/` - 公共包
  - `jwt/` - JWT 令牌处理
  - `log/` - 日志工具
- `test/` - 集成测试和端到端测试

### 核心组件

#### 应用程序生命周期 (`internal/app/app.go`)
- `Application` 结构体管理整个服务生命周期
- 初始化序列：配置 → 日志 → 数据库 → Redis → 健康检查 → 路由 → 服务器
- 使用关闭管理器进行优雅关闭

#### 数据库层 (`internal/db/`)
- 通过 GORM 支持 MySQL 和 PostgreSQL
- 包含健康检查功能
- Redis 客户端用于缓存和会话管理

#### 认证 (`pkg/jwt/`)
- JWT 令牌生成和验证
- 用户认证的安全令牌处理

#### 配置 (`config/config.go`)
- 基于环境的配置加载
- 支持开发和生产模式
- 数据库和服务器配置

### 关键模式
- **依赖注入**：使用 `appcontext.AppContext` 管理依赖
- **优雅关闭**：通过 `ShutdownManager` 实现适当的清理
- **健康检查**：内置数据库和服务健康监控
- **中间件栈**：日志记录和其他横切关注点
- **模块化设计**：各层之间明确的关注点分离

### 数据库支持
- 主要：MySQL（默认）
- 备选：PostgreSQL
- 缓存：Redis 用于会话管理和性能

### 测试策略
- 核心工具（JWT、日志）的单元测试
- 使用测试数据库进行数据库测试
- `test/` 目录中的集成测试
- 测试场景的模拟数据生成工具

## 环境设置

服务需要：
- Go 1.19+
- MySQL 或 PostgreSQL 数据库
- Redis 服务器
- 配置的环境变量（参见 config 包）

## 注意事项
- 主入口点在 `cmd/main.go`
- 服务在通过 `config.HTTPPort` 配置的端口上启动
- 在某些地方使用中文注释和错误消息
- 实现了具有结构化输出的全面日志记录
