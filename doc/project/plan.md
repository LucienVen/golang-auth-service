# Go认证服务项目开发计划 (MVP阶段)

## 📋 项目概述
- **目标**: 开发通用认证服务，支持多种登录方式
- **架构**: Clean Architecture + Gin框架
- **数据库**: MySQL/PostgreSQL + Redis
- **阶段**: MVP阶段 - 核心功能已实现，待集成修复

## 🎯 项目现状分析 (2025-10-29 修订)

### 当前项目阶段
- **状态**: 集成阶段 - 核心功能完成，存在集成问题
- **完成度**: 约82% (修正前评估为85%)
- **优势**: 架构设计优秀，代码结构规范，核心功能基本完整
- **主要问题**: Controller层未集成到主应用流程

### 已实现功能 (82%)
- ✅ Clean Architecture分层架构
- ✅ 配置管理 (Viper + 环境变量)
- ✅ 日志系统 (Zap + 文件轮转)
- ✅ 数据库基础设施 (GORM + MySQL/PostgreSQL)
- ✅ JWT工具包 (生成/解析/刷新)
- ✅ Redis客户端配置
- ✅ 密码加密模块 (bcrypt, cost=12)
- ✅ 用户Repository层 (完整CRUD + 测试)
- ✅ 认证Service层 (基础实现)
- ✅ 会话Service层 (内存存储)
- ✅ AuthController (完整实现)
- ✅ JWT验证中间件
- ✅ 路由配置 (AuthRoutes完整)
- ✅ 健康检查接口

### 核心集成问题 (18%缺失)
- ❌ **AuthController未集成到Container** (阻塞性)
- ❌ **认证路由未连接到主应用** (阻塞性)
- ❌ **主应用无法处理认证请求** (阻塞性)
- ❌ 用户服务层存在备份文件需恢复
- ❌ Redis会话管理未完全实现
- ❌ 测试覆盖存在编译错误

## 🚀 详细实施顺序

### 实施顺序原则

基于项目的Clean Architecture架构和现有实现情况，我们采用**自底向上**的实施策略：

1. **依赖关系清晰**: 从基础设施层到接口层，无循环依赖
2. **渐进可见**: 每个阶段都有可测试的功能产出
3. **风险可控**: 优先实现核心功能，后续添加增强特性
4. **便于调试**: 每层完成后都可以独立测试

### 四阶段详细实施计划

#### 第一阶段：基础设施扩展 (1-2天) ⭐⭐⭐⭐⭐
**目标**: 为认证功能提供基础工具和安全组件

**任务1.1: 密码加密模块** (优先级: 最高)
- **文件**: `internal/utils/password.go`
- **测试**: `internal/utils/password_test.go`
- **依赖**: 无，可立即开始
- **关键功能**:
  ```go
  func HashPassword(password string) (string, error)
  func VerifyPassword(password, hash string) error
  func CheckPasswordStrength(password string) error
  ```
- **安全要求**: bcrypt算法，cost=12

**任务1.2: User实体扩展** (优先级: 最高)
- **文件**: `internal/entity/user.go`
- **依赖**: 无
- **关键变更**:
  - 添加 `PasswordHash` 字段替换明文密码
  - 扩展邮箱/手机号/用户名验证方法
  - 完善用户状态转移逻辑
  - 添加BeforeCreate GORM钩子

**任务1.3: 输入验证工具** (优先级: 高)
- **文件**: `internal/utils/validator.go`
- **测试**: `internal/utils/validator_test.go`
- **依赖**: 无
- **关键功能**:
  ```go
  func ValidateEmail(email string) error
  func ValidatePhone(phone string) error
  func ValidateUsername(username string) error
  func ValidatePassword(password string) error
  ```

**阶段产出**: ✅ 可独立测试的工具组件，为后续开发提供安全基础

#### 第二阶段：数据访问层实现 (2-3天) ⭐⭐⭐⭐
**目标**: 实现用户数据的持久化和查询功能

**任务2.1: 用户Repository接口定义** (优先级: 高)
- **文件**: `internal/repository/user_repository.go`
- **依赖**: User实体
- **接口设计**:
  ```go
  type UserRepository interface {
      Create(user *entity.User) error
      GetByID(id uint) (*entity.User, error)
      GetByEmail(email string) (*entity.User, error)
      GetByPhone(phone string) (*entity.User, error)
      GetByUsername(username string) (*entity.User, error)
      Update(user *entity.User) error
      Delete(id uint) error
      GetByAccount(account string) (*entity.User, error) // 统一查询接口
  }
  ```

**任务2.2: 用户Repository实现** (优先级: 高)
- **文件**: `internal/repository/user_repository_impl.go`
- **依赖**: Repository接口 + 数据库连接
- **关键实现**:
  - GORM操作封装
  - 错误处理和转换
  - 查询优化和索引利用
  - 软删除支持

**任务2.3: 数据库事务支持** (优先级: 中)
- **文件**: `internal/repository/transaction.go`
- **依赖**: Repository实现
- **功能**: 事务管理和回滚机制

**阶段产出**: ✅ 完整的用户数据访问层，支持所有CRUD操作

#### 第三阶段：业务逻辑层实现 (3-4天) ⭐⭐⭐⭐⭐
**目标**: 实现核心认证和用户管理业务逻辑

**任务3.1: 用户Service** (优先级: 最高)
- **文件**: `internal/service/user_service.go`
- **依赖**: UserRepository + 密码工具 + 验证工具
- **核心方法**:
  ```go
  func (s *UserService) Register(req *RegisterRequest) (*User, error)
  func (s *UserService) UpdateProfile(userID uint, req *UpdateProfileRequest) error
  func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error
  func (s *UserService) GetUserProfile(userID uint) (*User, error)
  ```

**任务3.2: 认证Service** (优先级: 最高)
- **文件**: `internal/service/auth_service.go`
- **依赖**: UserRepository + 密码工具 + JWT工具
- **核心方法**:
  ```go
  func (s *AuthService) Login(account, password string) (*LoginResponse, error)
  func (s *AuthService) Logout(token string) error
  func (s *AuthService) RefreshToken(refreshToken string) (*TokenResponse, error)
  func (s *AuthService) ValidateToken(token string) (*UserClaims, error)
  ```

**任务3.3: 会话Service** (优先级: 高)
- **文件**: `internal/service/session_service.go`
- **依赖**: Redis客户端 + JWT工具
- **核心功能**:
  - Token存储和管理
  - 黑名单机制
  - 会话状态跟踪
  - 多设备登录管理

**阶段产出**: ✅ 完整的认证业务逻辑，支持注册、登录、会话管理

#### 第四阶段：接口层实现 (2-3天) ⭐⭐⭐⭐
**目标**: 提供HTTP接口和安全验证

**任务4.1: JWT验证中间件** (优先级: 最高)
- **文件**: `internal/middleware/auth.go`
- **依赖**: 认证Service + 会话Service
- **关键功能**:
  ```go
  func JWTAuthMiddleware() gin.HandlerFunc
  func RequireAuth() gin.HandlerFunc
  func OptionalAuth() gin.HandlerFunc
  ```

**任务4.2: 认证Controller** (优先级: 最高)
- **文件**: `internal/controller/auth_controller.go`
- **依赖**: 用户Service + 认证Service + 会话Service
- **接口实现**:
  - POST /api/auth/register
  - POST /api/auth/login
  - POST /api/auth/logout
  - GET /api/auth/profile
  - POST /api/auth/refresh

**任务4.3: 路由配置集成** (优先级: 高)
- **文件**: `api/router.go`
- **依赖**: 认证Controller + JWT中间件
- **更新内容**:
  - 添加认证路由组
  - 注册中间件
  - 更新控制器容器

**阶段产出**: ✅ 完整可用的认证API接口

### 最小MVP功能集合

为了快速验证核心功能，建议按以下顺序实现最小功能集：

#### 核心MVP (1周)
1. **密码加密模块** → 2. **User实体扩展** → 3. **用户Repository**
4. **用户Service (仅注册功能)** → 5. **认证Service (仅登录功能)**
6. **JWT中间件** → 7. **认证Controller (基础接口)**

**MVP验证功能**:
- 用户注册 (/api/auth/register)
- 用户登录 (/api/auth/login)
- 获取用户信息 (/api/auth/profile)
- Token验证中间件

#### 完整MVP (2周)
在核心MVP基础上添加：
- 用户登出 (/api/auth/logout)
- Token刷新 (/api/auth/refresh)
- 输入验证
- 基础错误处理

### 技术实施要点

#### 错误处理策略
```go
// 统一错误类型
type BusinessError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// 常用错误码
const (
    ErrCodeUserNotFound     = 1001
    ErrCodeUserExists       = 1002
    ErrCodeInvalidPassword  = 1003
    ErrCodeInvalidToken     = 2001
    ErrCodeTokenExpired     = 2002
)
```

#### 安全考虑
- **密码安全**: bcrypt cost=12，加盐哈希
- **Token安全**: JWT + 黑名单机制
- **输入验证**: 严格的参数验证和清理
- **日志安全**: 避免记录敏感信息

#### 性能优化
- **数据库查询**: 合理使用索引，避免N+1查询
- **缓存策略**: 用户信息缓存，Token缓存
- **连接池**: 数据库和Redis连接池优化

### 验证和测试策略

#### 每阶段验证
1. **第一阶段**: 单元测试验证工具函数正确性
2. **第二阶段**: 集成测试验证数据库操作
3. **第三阶段**: 业务逻辑测试各种场景
4. **第四阶段**: API测试验证完整流程

#### 测试覆盖目标
- **工具层**: 95%+ 覆盖率
- **Repository层**: 90%+ 覆盖率
- **Service层**: 85%+ 覆盖率
- **Controller层**: API集成测试

这个实施顺序确保了：
- ✅ **依赖清晰**: 自底向上，无循环依赖
- ✅ **渐进可见**: 每个阶段都有可测试产出
- ✅ **风险可控**: 核心功能优先，增强功能后续
- ✅ **质量保证**: 每层都有对应的测试验证

## 🚀 开发路线图

### 阶段1: 基础认证核心 (高优先级) - 预计3-4天

#### 1.1 密码加密模块
- **目标**: 实现安全的密码哈希和验证
- **技术**: bcrypt加密算法
- **文件**: `internal/utils/password.go`
- **测试**: `internal/utils/password_test.go`

**实现要点**:
```go
// 密码哈希生成
func HashPassword(password string) (string, error)
// 密码验证
func VerifyPassword(password, hash string) error
// 密码强度检查
func CheckPasswordStrength(password string) error
```

#### 1.2 User实体扩展
- **目标**: 扩展用户模型支持多方式认证
- **文件**: `internal/entity/user.go`
- **变更内容**:
  - 添加密码哈希字段
  - 支持邮箱、手机号、用户名字段
  - 完善用户状态转换逻辑
  - 添加验证方法

#### 1.3 用户Repository层
- **目标**: 实现用户数据访问层
- **文件**: `internal/repository/user_repository.go`
- **功能**: CRUD操作、查询优化、事务支持

**核心方法**:
```go
func (r *UserRepository) Create(user *entity.User) error
func (r *UserRepository) GetByID(id uint) (*entity.User, error)
func (r *UserRepository) GetByEmail(email string) (*entity.User, error)
func (r *UserRepository) GetByPhone(phone string) (*entity.User, error)
func (r *UserRepository) GetByUsername(username string) (*entity.User, error)
func (r *UserRepository) Update(user *entity.User) error
func (r *UserRepository) Delete(id uint) error
```

#### 1.4 用户Service层
- **目标**: 实现用户业务逻辑
- **文件**: `internal/service/user_service.go`
- **功能**: 注册逻辑、登录验证、用户管理

#### 1.5 认证Controller
- **目标**: 实现认证相关HTTP接口
- **文件**: `internal/controller/auth_controller.go`
- **接口设计**:
  - `POST /api/auth/register` - 用户注册
  - `POST /api/auth/login` - 用户登录
  - `POST /api/auth/logout` - 用户登出
  - `GET /api/auth/profile` - 获取用户信息
  - `POST /api/auth/refresh` - 刷新token

#### 1.6 JWT验证中间件
- **目标**: 实现token验证和用户身份识别
- **文件**: `internal/middleware/auth.go`
- **功能**: token解析、用户查询、权限验证

#### 1.7 Redis会话管理
- **目标**: 实现token存储和会话管理
- **文件**: `internal/service/session_service.go`
- **功能**: token存储、黑名单管理、踢下线

### 阶段2: 多方式认证支持 (中优先级) - 预计2-3天

#### 2.1 输入验证模块
- **目标**: 实现统一的数据验证
- **文件**: `internal/utils/validator.go`
- **验证规则**:
  - 邮箱格式验证
  - 手机号格式验证
  - 用户名格式验证
  - 密码强度验证

#### 2.2 多方式登录支持
- **目标**: 支持邮箱、手机号、用户名登录
- **文件**: `internal/service/auth_service.go`
- **功能**: 智能识别登录方式、统一验证流程

#### 2.3 统一错误处理
- **目标**: 规范化错误响应
- **文件**: `internal/errors/errors.go`
- **功能**: 错误码定义、错误响应格式、错误日志

#### 2.4 API文档集成
- **目标**: 生成接口文档
- **技术**: Swagger/OpenAPI 3.0
- **文件**: `docs/swagger.go`

### 阶段3: 安全和完善 (中优先级) - 预计2天

#### 3.1 登录频率限制
- **目标**: 防止暴力破解攻击
- **文件**: `internal/middleware/rate_limiter.go`
- **技术**: Redis计数器 + 滑动窗口

#### 3.2 用户状态管理完善
- **目标**: 完善账户生命周期管理
- **功能**: 账户激活、禁用、注销逻辑

#### 3.3 核心功能测试
- **目标**: 保证代码质量
- **覆盖范围**: Service层、Controller层、工具类

#### 3.4 生产环境配置优化
- **目标**: 适配生产环境需求
- **内容**: 安全配置、性能调优、监控配置

### 阶段4: 部署和监控 (低优先级) - 预计1-2天

#### 4.1 Docker容器化
- **目标**: 支持容器化部署
- **文件**: `Dockerfile`, `docker-compose.yml`
- **功能**: 多阶段构建、健康检查

#### 4.2 监控和健康检查增强
- **目标**: 完善服务监控
- **功能**: 数据库连接监控、Redis连接监控、系统资源监控

#### 4.3 日志优化
- **目标**: 结构化日志和审计
- **功能**: 请求日志、审计日志、错误追踪

## 📊 接口设计规范

### 认证接口

#### 用户注册
```
POST /api/auth/register
Content-Type: application/json

Request:
{
    "username": "string",     // 可选
    "email": "string",        // 可选，与phone至少一个
    "phone": "string",        // 可选，与email至少一个
    "password": "string"
}

Response (201):
{
    "code": 201,
    "message": "注册成功",
    "data": {
        "user": {
            "id": 1,
            "username": "test",
            "email": "test@example.com",
            "phone": "13800138000",
            "status": "active",
            "created_at": "2024-01-01T00:00:00Z"
        },
        "token": "jwt_token_string",
        "expires_at": "2024-01-02T00:00:00Z"
    }
}
```

#### 用户登录
```
POST /api/auth/login
Content-Type: application/json

Request:
{
    "account": "string",      // 邮箱/手机号/用户名
    "password": "string"
}

Response (200):
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "user": {
            "id": 1,
            "username": "test",
            "email": "test@example.com",
            "phone": "13800138000",
            "status": "active"
        },
        "token": "jwt_token_string",
        "refresh_token": "refresh_token_string",
        "expires_at": "2024-01-02T00:00:00Z"
    }
}
```

#### 用户登出
```
POST /api/auth/logout
Authorization: Bearer {token}

Response (200):
{
    "code": 200,
    "message": "登出成功"
}
```

#### 获取用户信息
```
GET /api/auth/profile
Authorization: Bearer {token}

Response (200):
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": 1,
        "username": "test",
        "email": "test@example.com",
        "phone": "13800138000",
        "status": "active",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
    }
}
```

#### 刷新Token
```
POST /api/auth/refresh
Authorization: Bearer {refresh_token}

Response (200):
{
    "code": 200,
    "message": "刷新成功",
    "data": {
        "token": "new_jwt_token_string",
        "refresh_token": "new_refresh_token_string",
        "expires_at": "2024-01-02T00:00:00Z"
    }
}
```

### 错误响应格式
```
{
    "code": 400,
    "message": "参数错误",
    "error": "email格式不正确",
    "timestamp": "2024-01-01T00:00:00Z",
    "path": "/api/auth/register"
}
```

## 🔧 技术实现细节

### 认证流程设计

#### 注册流程
1. 接收注册请求
2. 验证输入参数格式
3. 检查邮箱/手机号/用户名是否已存在
4. 验证密码强度
5. 生成密码哈希
6. 创建用户记录
7. 生成JWT token和refresh token
8. 存储token到Redis
9. 返回用户信息和token

#### 登录流程
1. 接收登录请求
2. 验证输入参数
3. 识别登录方式 (邮箱/手机号/用户名)
4. 查询用户信息
5. 检查用户状态
6. 验证密码
7. 检查登录频率限制
8. 生成JWT token和refresh token
9. 存储token到Redis
10. 记录登录日志
11. 返回用户信息和token

#### 验证流程
1. 从请求头提取token
2. 解析JWT token
3. 检查token是否在Redis黑名单中
4. 查询用户信息
5. 检查用户状态
6. 将用户信息存入上下文
7. 继续处理请求

#### 登出流程
1. 从请求头提取token
2. 将token加入Redis黑名单
3. 清理用户会话信息
4. 记录登出日志
5. 返回成功响应

### 数据库设计

#### 用户表 (users)
```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NULL,
    email VARCHAR(100) UNIQUE NULL,
    phone VARCHAR(20) UNIQUE NULL,
    password_hash VARCHAR(255) NOT NULL,
    status ENUM('inactive', 'active', 'disabled', 'deleted') DEFAULT 'inactive',
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    INDEX idx_email (email),
    INDEX idx_phone (phone),
    INDEX idx_username (username),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);
```

### Redis键设计

#### Token存储
```
auth:token:{user_id} = {token_info}     // 用户当前token
auth:blacklist:{token_jti} = 1          // token黑名单
auth:refresh:{refresh_jti} = {user_id}  // refresh token映射
```

#### 频率限制
```
rate_limit:login:{ip} = {count,timestamp}  // IP登录频率限制
rate_limit:login:{account} = {count,timestamp}  // 账户登录频率限制
```

#### 会话管理
```
session:user:{user_id} = {session_info}     // 用户会话信息
session:device:{user_id}:{device_id} = {token_info}  // 设备token管理
```

## 📋 质量保证

### 测试策略
- **单元测试**: 覆盖Service层和工具类，目标覆盖率 >80%
- **集成测试**: 测试API接口完整流程
- **性能测试**: 测试并发登录和token验证性能

### 安全检查清单
- ✅ 密码使用bcrypt加密 (cost=12)
- ✅ JWT token设置合理过期时间
- ✅ 实现token黑名单机制
- ✅ 登录频率限制防止暴力破解
- ✅ 输入参数验证和清理
- ✅ 敏感信息不在日志中暴露
- ✅ HTTPS强制使用 (生产环境)

### 性能指标
- **注册响应时间**: < 200ms
- **登录响应时间**: < 150ms
- **token验证时间**: < 10ms
- **并发登录支持**: 1000 QPS
- **数据库连接池**: 最大50连接

## 🎯 里程碑和交付物

### 里程碑1: 核心认证功能 (第1周)
- **交付物**:
  - 完整的用户注册/登录/登出功能
  - JWT中间件和token管理
  - Redis会话管理
  - 基础API接口

### 里程碑2: 多方式认证 (第2周)
- **交付物**:
  - 支持邮箱/手机号/用户名登录
  - 完善的输入验证
  - 统一错误处理
  - API文档

### 里程碑3: 安全增强 (第3周)
- **交付物**:
  - 登录频率限制
  - 用户状态管理完善
  - 核心功能测试覆盖
  - 生产环境配置

### 里程碑4: 部署就绪 (第4周)
- **交付物**:
  - Docker容器化配置
  - 监控和健康检查
  - 完整的部署文档
  - 性能测试报告

## 📈 风险评估和缓解

### 技术风险
- **风险**: JWT token安全性
- **缓解**: 使用RS256算法、短过期时间、黑名单机制

- **风险**: 并发性能瓶颈
- **缓解**: Redis缓存、数据库连接池、异步处理

- **风险**: 数据一致性
- **缓解**: 事务处理、乐观锁、定期数据校验

### 业务风险
- **风险**: 用户数据泄露
- **缓解**: 密码加密、敏感数据脱敏、审计日志

- **风险**: 服务可用性
- **缓解**: 健康检查、优雅关闭、监控告警

## 📚 参考文档和资源

### 技术文档
- [Gin框架官方文档](https://gin-gonic.com/docs/)
- [GORM文档](https://gorm.io/docs/)
- [JWT最佳实践](https://auth0.com/docs/tokens/json-web-tokens)
- [Redis最佳实践](https://redis.io/documentation)

### 安全标准
- [OWASP认证备忘单](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [密码存储指南](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)

---

**最后更新**: 2025-10-29 (项目状态评估和计划修订)
**修订原因**: 发现Controller层集成问题，调整实施优先级
**负责人**: 开发团队
**审核状态**: 已修订 - 优先调整