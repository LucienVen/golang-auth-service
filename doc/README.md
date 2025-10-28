# 项目文档索引

## 📚 文档结构

本项目的文档采用模块化组织结构，按功能和用途进行分类。

### 🎯 项目核心文档 (`project/`)

- **[项目开发计划](project/plan.md)** - 详细的项目实施计划和MVP功能规划
- **[架构分析报告](project/architecture-analysis.md)** - 项目技术架构和依赖关系分析

### 🛠️ 开发规范 (`development/`)

- **[变更日志](development/changelog.md)** - 项目开发记录和checkpoint机制
- **[开发指南](development/claude.md)** - Claude Code配置和开发规范

### 📦 功能模块文档 (`modules/`)

#### 认证模块 (`modules/auth/`)
- 用户认证和授权相关文档

#### 用户模块 (`modules/user/`)
- 用户管理和功能相关文档

#### 安全模块 (`modules/security/`)
- 安全策略和实现相关文档

#### 工具模块 (`modules/utils/`)
- 通用工具和库相关文档

### 🌐 API文档 (`api/`)
- REST API接口文档和规范

### 🚀 部署文档 (`deployment/`)
- 部署配置和运维相关文档

## 📖 文档使用指南

### 开发流程
1. 首先阅读 [项目开发计划](project/plan.md) 了解整体规划
2. 查看 [架构分析报告](project/architecture-analysis.md) 理解技术架构
3. 参考 [开发指南](development/claude.md) 遵循编码规范
4. 在 [变更日志](development/changelog.md) 中跟踪开发进度

### 模块开发
- 各功能模块的详细文档位于对应的 `modules/` 子目录中
- 每个模块包含设计思路、实现细节和使用说明
- API接口文档统一在 `api/` 目录中管理

### 文档维护
- 所有文档变更需在 [变更日志](development/changelog.md) 中记录
- 新增功能模块需创建对应的文档目录
- 保持文档与代码同步更新

---

**最后更新**: 2025-10-28
**维护者**: Claude Code Development Team