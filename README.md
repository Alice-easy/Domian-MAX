# Domain MAX - 二级域名分发管理系统

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)
![React Version](https://img.shields.io/badge/react-18+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

一个现代化的二级域名分发管理系统，支持多 DNS 提供商，提供完整的域名管理、DNS 记录管理和用户权限控制功能。

## ✨ 核心功能

### 🌐 多 DNS 提供商支持

- **Cloudflare** - 全球 CDN 和 DNS 服务
- **阿里云 DNS** - 企业级 DNS 解析服务
- **腾讯云 DNS** - 高可用 DNS 解析
- **华为云 DNS** - 智能 DNS 解析
- **DNSPod** - 专业 DNS 服务
- **AWS Route53** - 亚马逊云 DNS 服务

### 👤 用户管理系统

- 用户注册、登录、密码重置
- 基于 JWT 的身份认证
- 角色权限控制（普通用户/管理员）
- 用户资料管理

### 🏗️ 域名管理

- 域名添加、删除、修改
- 多 DNS 提供商配置
- 域名状态监控
- 批量操作支持

### 📋 DNS 记录管理

- 支持所有常用记录类型（A、AAAA、CNAME、MX、TXT 等）
- 批量导入/导出 DNS 记录
- 记录模板管理
- 操作历史追踪

### 🔒 安全特性

- 数据加密存储
- SQL 注入防护
- XSS 攻击防护
- CSRF 保护
- 速率限制
- 操作日志审计

## 🚀 快速开始

### 系统要求

- **Go 1.23+** - 后端开发环境
- **Node.js 18+** - 前端开发环境
- **PostgreSQL 14+** - 数据库服务
- **Redis 7+** - 缓存服务（可选）
- **内存** 2GB+
- **磁盘空间** 2GB+

### 环境准备

#### 1. 安装依赖软件

**Ubuntu/Debian:**

```bash
# 安装 Go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 安装 Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安装 PostgreSQL
sudo apt-get install postgresql postgresql-contrib

# 安装 Redis（可选）
sudo apt-get install redis-server
```

**macOS:**

```bash
# 使用 Homebrew 安装
brew install go node postgresql redis
```

**Windows:**

```bash
# 使用 Scoop 安装（推荐）
scoop install go nodejs postgresql redis

# 或下载官方安装包
# Go: https://golang.org/dl/
# Node.js: https://nodejs.org/
# PostgreSQL: https://www.postgresql.org/download/
```

#### 2. 配置数据库

```bash
# 启动 PostgreSQL 服务
sudo systemctl start postgresql  # Linux
brew services start postgresql   # macOS
# Windows: 通过服务管理器启动

# 创建数据库用户和数据库
sudo -u postgres psql
CREATE USER domain_user WITH PASSWORD 'your_password';
CREATE DATABASE domain_manager OWNER domain_user;
GRANT ALL PRIVILEGES ON DATABASE domain_manager TO domain_user;
\q
```

#### 3. 配置环境变量

```bash
# 复制环境配置文件
cp .env.example .env

# 编辑环境配置
vi .env
```

必须配置的环境变量：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=domain_user
DB_PASSWORD=your_password
DB_NAME=domain_manager

# 应用配置
PORT=8080
JWT_SECRET=your-super-secret-jwt-key
ENCRYPTION_KEY=your-32-byte-encryption-key

# Redis 配置（可选）
REDIS_HOST=localhost
REDIS_PORT=6379
```

### 本地开发部署

```bash
# 1. 克隆项目
git clone https://github.com/your-repo/domain-max.git
cd domain-max

# 2. 安装依赖
make install

# 3. 构建项目
make build

# 4. 初始化数据库
make db-migrate

# 5. 启动应用
make dev
```

### 生产环境部署

```bash
# 1. 构建生产版本
make build-all

# 2. 配置生产环境变量
cp .env .env.production
# 编辑生产配置...

# 3. 启动应用
./domain-max

# 或者在后台运行
nohup ./domain-max > app.log 2>&1 &
```

### 访问应用

应用启动后，您可以通过以下地址访问：

- **前端界面**: http://localhost:8080
- **API 接口**: http://localhost:8080/api
- **健康检查**: http://localhost:8080/api/health

## 📋 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React 前端    │    │   Go 后端 API   │    │  PostgreSQL DB  │
│   TypeScript    ├────┤   RESTful API   ├────┤   数据存储      │
│   响应式设计    │    │   JWT 认证      │    │   ACID 事务     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │              ┌─────────────────┐              │
         │              │   Redis 缓存    │              │
         └──────────────┤   会话存储      ├──────────────┘
                        │   频率限制      │
                        └─────────────────┘
```

### 服务架构

- **前端服务** - React 应用已构建并内嵌到 Go 二进制文件中
- **后端服务** - Go 单体应用，内置静态文件服务
- **数据库** - PostgreSQL 独立部署
- **缓存** - Redis 独立部署（可选）

## 🛠️ 技术栈

### 后端技术

- **Go 1.23+** - 高性能服务端语言
- **Gin** - 轻量级 Web 框架
- **GORM** - ORM 数据库操作
- **JWT** - 身份认证
- **Bcrypt** - 密码加密
- **PostgreSQL** - 可靠的关系型数据库
- **Redis** - 内存缓存数据库

### 前端技术

- **React 18** - 现代化前端框架
- **TypeScript** - 类型安全的 JavaScript
- **Vite** - 快速构建工具
- **Tailwind CSS** - 实用优先的 CSS 框架
- **React Router** - 客户端路由
- **Axios** - HTTP 客户端

### 基础设施

- **内嵌静态服务** - Go 应用内置前端静态文件服务
- **PostgreSQL** - 关系型数据库
- **Redis** - 缓存和会话存储（可选）

## 📖 详细文档

- 📚 [文档中心](docs/) - 完整的项目文档导航
- 🚀 [部署指南](docs/deployment.md) - 完整的部署步骤和配置
- 🏭 [生产环境指南](docs/production-guide.md) - 生产环境优化和安全配置
- 🏗️ [系统架构](docs/architecture.md) - 详细的系统设计文档

## 🔧 开发指南

### 本地开发环境

```bash
# 1. 安装依赖
make install

# 2. 启动开发环境（分离模式）
# 终端1: 启动后端开发服务器
make dev

# 终端2: 启动前端开发服务器（热重载）
make dev-web
```

### 构建和测试

```bash
# 构建项目
make build

# 运行测试
make test

# 代码检查
make lint

# 生成测试覆盖率报告
make test-coverage
```

## 🏥 运维管理

### 服务管理

```bash
# 启动应用
./domain-max

# 后台运行
nohup ./domain-max > app.log 2>&1 &

# 停止应用（查找进程ID）
ps aux | grep domain-max
kill <PID>

# 或使用脚本管理
# 创建服务脚本 /etc/systemd/system/domain-max.service
sudo systemctl start domain-max
sudo systemctl stop domain-max
sudo systemctl restart domain-max
```

### 健康检查

```bash
# 应用健康检查
curl http://localhost:8080/api/health

# 检查进程状态
ps aux | grep domain-max

# 检查端口占用
netstat -tlnp | grep :8080

# 检查日志
tail -f app.log
```

### 数据库管理

```bash
# 连接数据库
psql -h localhost -U domain_user -d domain_manager

# 备份数据库
pg_dump -h localhost -U domain_user domain_manager > backup.sql

# 恢复数据库
psql -h localhost -U domain_user domain_manager < backup.sql

# 数据库迁移
make db-migrate
```

## 🔒 安全特性

### 数据安全

- **加密存储** - 敏感数据 AES 加密
- **密码安全** - Bcrypt 哈希算法
- **SQL 注入防护** - 参数化查询
- **XSS 防护** - 输入验证和输出编码

### 网络安全

- **HTTPS 强制** - 所有通信加密
- **HSTS 启用** - 防止协议降级
- **安全标头** - CSP、X-Frame-Options 等
- **速率限制** - 防止暴力攻击

### 认证授权

- **JWT 认证** - 无状态身份验证
- **角色权限** - 基于角色的访问控制
- **会话管理** - 安全的会话控制
- **操作审计** - 完整的操作日志

## 📊 性能特性

### 应用性能

- **缓存策略** - Redis 多层缓存
- **连接池** - 数据库连接优化
- **异步处理** - 后台任务队列
- **压缩传输** - Gzip 内容压缩

### 数据库优化

- **索引优化** - 关键字段索引
- **查询优化** - SQL 性能调优
- **连接池** - 连接复用机制
- **读写分离** - 主从数据库架构

### 前端优化

- **代码分割** - 按需加载
- **资源压缩** - 静态资源优化
- **CDN 加速** - 全球内容分发
- **缓存策略** - 浏览器缓存优化

## 🚨 故障排除

### 常见问题

#### 应用无法启动

```bash
# 检查端口占用
netstat -tlnp | grep :8080
lsof -i :8080

# 检查配置文件
cat .env

# 检查日志
tail -f app.log

# 检查权限
ls -la domain-max
chmod +x domain-max
```

#### 数据库连接失败

```bash
# 检查数据库服务状态
sudo systemctl status postgresql  # Linux
brew services list | grep postgresql  # macOS

# 测试数据库连接
psql -h localhost -U domain_user -d domain_manager

# 检查数据库配置
grep -E "^(DB_|POSTGRES_)" .env

# 重启数据库服务
sudo systemctl restart postgresql
```

#### 前端页面无法访问

```bash
# 检查静态文件是否存在
ls -la web/dist/

# 重新构建前端
cd web && npm run build

# 检查服务器路由配置
curl -v http://localhost:8080/
```

### 性能问题诊断

```bash
# 检查系统资源
top
htop
free -h
df -h

# 检查应用性能
# 安装 pprof
go tool pprof http://localhost:8080/debug/pprof/profile

# 数据库性能分析
psql -U domain_user -d domain_manager -c "
SELECT query, calls, total_time, mean_time
FROM pg_stat_statements
ORDER BY total_time DESC
LIMIT 10;"

# 检查网络延迟
curl -w "%{time_total}" -o /dev/null -s http://localhost:8080/api/health
```

## 🤝 贡献指南

我们欢迎社区贡献！请遵循以下步骤：

1. **Fork** 项目仓库
2. **创建** 功能分支 (`git checkout -b feature/amazing-feature`)
3. **提交** 更改 (`git commit -m 'Add amazing feature'`)
4. **推送** 分支 (`git push origin feature/amazing-feature`)
5. **创建** Pull Request

### 开发规范

- 遵循 Go 代码规范
- 编写单元测试
- 更新相关文档
- 确保所有测试通过

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 获取帮助

- 📖 查看 [文档](docs/)
- 🐛 报告 [Issues](https://github.com/your-repo/domain-max/issues)
- 💬 加入 [讨论](https://github.com/your-repo/domain-max/discussions)
- 📧 邮件联系: support@domain-max.com

## 🙏 致谢

感谢所有为本项目做出贡献的开发者和社区成员！

---

**Domain MAX** - 让域名管理更简单、更安全、更高效！

Made with ❤️ by Domain MAX Team
