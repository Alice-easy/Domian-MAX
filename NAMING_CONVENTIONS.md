# 项目命名规范 (Naming Conventions)

本文档定义了 Domain-MAX 项目的文件和目录命名规范，确保代码库的一致性和可维护性。

## 📁 目录命名规范

### Go 项目目录

- **全小写字母**，使用有意义的英文单词
- **多个单词使用下划线分隔** （如果必要）
- **语义清晰**，避免缩写

```
✅ 正确示例：
internal/
├── api/
├── config/
├── constants/
├── database/
├── middleware/
├── models/
├── providers/
├── services/
└── utils/

❌ 错误示例：
Internal/
├── API/
├── cfg/
├── db-conn/
└── svc/
```

### 前端项目目录

- **小驼峰命名**（React 组件目录可使用 PascalCase）
- **语义清晰**

```
✅ 正确示例：
frontend/src/
├── components/
├── pages/
├── stores/
└── utils/

❌ 错误示例：
frontend/src/
├── Components/
├── Pages/
└── Stores/
```

## 📄 文件命名规范

### Go 源文件

- **全小写字母**
- **多个单词使用下划线分隔**
- **文件名应反映其主要功能**

```
✅ 正确示例：
main.go
handlers.go
middleware.go
refresh_token.go
token_manager.go
generate_config.go

❌ 错误示例：
Main.go
handlers-api.go
generateConfig.go
refreshToken.go
```

### 前端文件

- **React 组件**: PascalCase (大驼峰)
- **工具函数**: camelCase (小驼峰)
- **配置文件**: 遵循工具约定

```
✅ 正确示例：
App.tsx
Dashboard.tsx
AdminLayout.tsx
authStore.ts
api.ts
vite.config.ts
package.json

❌ 错误示例：
app.tsx
admin_layout.tsx
auth-store.ts
```

### 配置和文档文件

- **配置文件**: 遵循对应工具的约定
- **文档文件**: 全大写 Markdown 文件，小写其他文档

```
✅ 正确示例：
docker-compose.yml
package.json
tsconfig.json
README.md
DEPLOYMENT.md
OPERATIONS.md
LICENSE
Makefile

❌ 错误示例：
Docker-Compose.yml
Package.json
readme.md
license.txt
```

## 🚫 避免的命名模式

### 通用规则

- **不要使用**拼写错误（如 `Domian` 应为 `Domain`）
- **不要混合**命名风格（如 `getUserInfo_fromDB`）
- **不要使用**没有意义的缩写（如 `usr`, `cfg`, `svc`）
- **不要在 Go 中**使用连字符分隔文件名

### 特殊字符使用

- **Go 文件**: 只允许字母、数字、下划线
- **前端文件**: 允许点号（如 `vite.config.ts`）
- **配置文件**: 遵循工具约定（如 `docker-compose.yml`）

## 📦 包和模块命名

### Go 包名

- **全小写**
- **简短但有意义**
- **不使用下划线或连字符**

```
✅ 正确示例：
package main
package config
package middleware
package providers

❌ 错误示例：
package Main
package Config
package middle_ware
package dns-providers
```

### 前端模块

- **遵循 npm 约定**
- **使用 kebab-case** (连字符分隔)

## 🔧 工具和脚本

### 脚本文件

- **Go 脚本**: 使用下划线分隔 (`generate_config.go`)
- **Shell 脚本**: 使用连字符分隔 (`build-docker.sh`)

### 可执行文件

- **使用连字符分隔**
- **语义清晰的名称**

```
✅ 正确示例：
domain-manager
backup-tool

❌ 错误示例：
dm
bkup
DomainManager
```

## 📝 Git 相关

### 分支命名

- **使用连字符分隔**
- **包含类型前缀**

```
✅ 正确示例：
feature/user-authentication
bugfix/login-error
hotfix/security-patch

❌ 错误示例：
userAuth
fix_login
SecurityPatch
```

### 标签命名

- **语义化版本号**

```
✅ 正确示例：
v1.0.0
v1.2.3-beta
v2.0.0-rc.1

❌ 错误示例：
version1
release-1.0
stable
```

## 🎯 最佳实践

1. **保持一致性**: 在整个项目中使用相同的命名风格
2. **语义清晰**: 文件名应该清楚表达其用途
3. **避免缩写**: 除非是广泛认知的缩写（如 `HTTP`, `API`）
4. **考虑国际化**: 使用英文命名，避免特殊字符
5. **遵循社区约定**: 遵循 Go、React 等技术栈的命名约定

## 🔄 定期检查

建议定期审查项目中的命名，确保：

- 新增文件遵循命名规范
- 没有命名不一致的地方
- 删除不再使用的文件

---

**最后更新**: 2025 年 9 月 11 日  
**维护者**: Domain-MAX 开发团队
