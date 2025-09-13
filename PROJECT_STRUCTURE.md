# Domain MAX 项目结构

本文档描述了整理后的 Domain MAX 项目目录结构和文件组织。

## 📁 项目根目录

```
Domain-MAX/
├── .github/               # GitHub工作流和配置
│   └── workflows/         # CI/CD自动化工作流
├── cmd/                   # 应用程序入口
│   └── server/           # 主服务器应用
├── configs/              # 配置文件
├── deployments/          # 部署相关文件
├── docs/                 # 项目文档
├── pkg/                  # 可重用的包和模块
├── scripts/              # 构建和部署脚本
├── web/                  # 前端React应用
├── go.mod                # Go模块定义
├── go.sum                # Go模块依赖
├── LICENSE               # 许可证
├── Makefile              # 构建自动化
└── README.md             # 项目主文档
```

## 📂 详细目录结构

### .github/workflows/

GitHub Actions 工作流配置，提供完整的 DevOps 自动化：

```
.github/workflows/
├── auto-update.yml           # 自动更新工作流
├── backup.yml               # 备份工作流
├── ci-cd.yml               # 持续集成和部署
├── database-maintenance.yml # 数据库维护
├── dependency-updates.yml   # 依赖更新
├── monitoring.yml          # 系统监控
├── performance-test.yml    # 性能测试
├── release.yml            # 版本发布
├── security-scan.yml      # 安全扫描
└── README.md              # 工作流文档
```

### cmd/server/

应用程序主入口：

```
cmd/server/
└── main.go               # 主服务器启动文件
```

### configs/

系统配置文件：

```
configs/
└── init.sql              # 数据库初始化脚本
```

### deployments/

部署配置和容器化文件：

```
deployments/
├── .env.example          # 环境变量模板
├── docker-compose.yml    # Docker容器编排
├── Dockerfile           # 容器构建配置
└── ssl/                 # SSL证书目录
```

### docs/

项目文档目录：

```
docs/
├── README.md            # 文档中心导航
├── architecture.md     # 系统架构文档
├── deployment.md       # 部署指南（合并后的完整版）
└── production-guide.md # 生产环境指南
```

### pkg/

核心业务逻辑包：

```
pkg/
├── api/                 # API相关代码
│   ├── auth.go
│   └── simple_dns.go
├── auth/               # 认证模块
│   └── models/
│       └── user.go
├── config/             # 配置管理
│   └── config.go
├── database/           # 数据库操作
│   ├── connection.go
│   └── migration.go
├── dns/                # DNS管理模块
│   ├── models/
│   │   └── dns.go
│   └── providers/      # DNS提供商
│       ├── aliyun.go
│       ├── cloudflare.go
│       ├── dnspod.go
│       ├── factory.go
│       ├── interface.go
│       └── others.go
├── email/              # 邮件服务
│   └── models/
│       └── smtp.go
├── middleware/         # HTTP中间件
│   ├── auth.go
│   ├── cors.go
│   └── rate-limit.go
└── utils/              # 工具函数
    └── validation.go
```

### scripts/

构建和部署脚本：

```
scripts/
├── build.sh            # 完整构建脚本（合并后）
├── cleanup.sh          # 清理脚本
├── deploy-complete.sh  # 完整部署脚本（主要部署工具）
├── generate-ssl.sh     # SSL证书生成
└── system-test.sh      # 系统测试脚本
```

### web/

前端 React 应用：

```
web/
├── public/             # 静态资源
├── src/               # 源代码
│   ├── components/    # React组件
│   ├── pages/         # 页面组件
│   ├── stores/        # 状态管理
│   ├── types/         # TypeScript类型
│   └── utils/         # 前端工具
├── package.json       # Node.js依赖
├── tsconfig.json      # TypeScript配置
└── vite.config.ts     # Vite构建配置
```

## 🗂️ 文件整理说明

### 已删除的冗余文件

1. **重复的部署文档**：

   - ❌ `docs/deployment-guide.md` → 内容合并到 `docs/deployment.md`

2. **重复的构建脚本**：

   - ❌ `scripts/build-go.sh` → 功能合并到 `scripts/build.sh`
   - ❌ `scripts/test-build.sh` → 功能合并到 `scripts/build.sh`
   - ❌ `scripts/deploy.sh` → 使用更完整的 `scripts/deploy-complete.sh`

3. **重复的配置文件**：

   - ❌ `configs/env.example` → 使用 `deployments/.env.example`

4. **空目录**：
   - ❌ `pkg/admin/` → 空目录已删除

### 优化后的特点

1. **清晰的分层结构**：

   - 🎯 `cmd/` - 应用入口
   - 📦 `pkg/` - 业务逻辑
   - 🌐 `web/` - 前端应用
   - 📖 `docs/` - 文档中心
   - 🔧 `scripts/` - 工具脚本

2. **统一的文档体系**：

   - 📚 `docs/README.md` - 文档导航中心
   - 🚀 完整的部署指南
   - 🏗️ 详细的架构文档

3. **高效的脚本工具**：

   - 🔨 `build.sh` - 统一构建工具
   - 🚀 `deploy-complete.sh` - 一键部署
   - 🧪 `system-test.sh` - 系统测试

4. **完整的 DevOps 流程**：
   - ⚙️ 9 个 GitHub Actions 工作流
   - 🔒 全面的安全扫描
   - 📊 自动化监控和报告

## 📋 使用指南

### 快速开始

```bash
# 1. 克隆项目
git clone <repository-url>
cd Domain-MAX

# 2. 一键部署
./scripts/deploy-complete.sh
```

### 开发环境

```bash
# 构建前端
./scripts/build.sh web

# 构建后端
./scripts/build.sh server

# 完整构建
./scripts/build.sh all
```

### 文档导航

- **开始使用** → [README.md](../README.md)
- **文档中心** → [docs/README.md](../docs/README.md)
- **部署指南** → [docs/deployment.md](../docs/deployment.md)
- **系统架构** → [docs/architecture.md](../docs/architecture.md)

---

**项目结构文档** | Domain MAX v1.0 | 最后更新: 2024 年 12 月
