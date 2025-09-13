# Domain MAX - 域名管理系统

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Architecture](https://img.shields.io/badge/architecture-frontend%2Fbackend%20separated-brightgreen.svg)
![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)
![React Version](https://img.shields.io/badge/react-18+-blue.svg)
![Cloudflare](https://img.shields.io/badge/frontend-Cloudflare%20Pages-orange.svg)
![VPS](https://img.shields.io/badge/backend-VPS%20API-green.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

现代化的域名与 DNS 管理系统，采用**前后端分离架构**，前端部署在 Cloudflare Pages，后端 API 部署在 VPS，支持多种云数据库。

## 🏗️ 系统架构

```
┌─────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐
│   Cloudflare Pages  │    │     Your VPS        │    │   Remote Database   │
│   (前端 React SPA)   │    │   (后端 Go API)      │    │  (PostgreSQL/MySQL) │
├─────────────────────┤    ├─────────────────────┤    ├─────────────────────┤
│ • Global CDN        │    │ • RESTful API       │    │ • PlanetScale       │
│ • Static Assets     │───▶│ • JWT Auth          │───▶│ • Supabase          │
│ • React Router      │    │ • DNS Management    │    │ • AWS RDS           │
│ • Auto HTTPS        │    │ • CORS Enabled      │    │ • 自建 PostgreSQL    │
└─────────────────────┘    └─────────────────────┘    └─────────────────────┘
```

### 🌟 架构优势

- **🚀 全球加速**: 前端通过 Cloudflare CDN 全球加速访问
- **💰 成本优化**: 前端免费托管，后端 VPS 成本可控
- **🛡️ 高可用性**: 分离部署降低单点故障风险
- **⚡ 高性能**: 静态资源 CDN 缓存，API 服务独立优化
- **🔧 易维护**: 前后端独立开发、部署和扩展

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

### 🛠️ 开发环境

#### 系统要求

- **Go 1.23+** - 后端 API 开发
- **Node.js 18+** - 前端开发
- **远程数据库** - PlanetScale/Supabase/AWS RDS 等

#### 快速启动开发环境

```bash
# 1. 克隆项目
git clone https://github.com/your-username/domain-max.git
cd domain-max

# 2. 安装所有依赖
./dev.sh install

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env 设置数据库连接等

# 4. 启动开发环境（前后端同时启动）
./dev.sh start
```

开发环境启动后：

- **前端**: http://localhost:5173
- **后端 API**: http://localhost:8080
- **API 文档**: http://localhost:8080/health

#### 分别启动前后端

```bash
# 仅启动后端API
./dev.sh backend

# 仅启动前端
./dev.sh frontend

# 构建API服务
./build-api.sh build

# 开发模式运行API
./build-api.sh dev
```

### 🚀 生产部署

采用前后端分离架构，分别部署到不同平台：

#### 部署架构选择

| 组件         | 推荐平台             | 特点                     |
| ------------ | -------------------- | ------------------------ |
| **前端**     | Cloudflare Pages     | 全球 CDN、免费、自动构建 |
| **后端 API** | VPS (Ubuntu)         | 完全控制、成本可控       |
| **数据库**   | PlanetScale/Supabase | 托管服务、高可用         |

#### 🌐 前端部署 (Cloudflare Pages)

1. **连接 GitHub 仓库**

   - 登录 [Cloudflare Dashboard](https://dash.cloudflare.com)
   - Pages → Create project → Connect Git

2. **配置构建设置**

   ```yaml
   Build command: cd web && npm ci && npm run build
   Build output directory: web/dist
   Environment variables:
     NODE_ENV: production
     VITE_API_BASE_URL: https://api.yourdomain.com
     VITE_BACKEND_DOMAIN: api.yourdomain.com
   ```

3. **自定义域名（可选）**
   - Pages Settings → Custom domains
   - 添加你的域名

#### 🖥️ 后端部署 (VPS)

1. **VPS 环境准备**

   ```bash
   # 连接VPS
   ssh user@your-vps-ip

   # 安装依赖
   sudo apt update && sudo apt upgrade -y
   sudo apt install -y git nginx certbot python3-certbot-nginx

   # 安装Go
   wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

2. **部署 API 服务**

   ```bash
   # 克隆项目
   git clone https://github.com/your-username/domain-max.git
   cd domain-max

   # 配置环境变量
   cp .env.vps .env
   nano .env  # 编辑配置

   # 构建API
   ./build-api.sh cross

   # 配置系统服务
   sudo cp domain-max-api-linux /usr/local/bin/domain-max-api
   sudo chmod +x /usr/local/bin/domain-max-api
   ```

3. **配置 Nginx 和 SSL**

   ```bash
   # 配置Nginx反向代理
   sudo nano /etc/nginx/sites-available/domain-max-api

   # 启用配置
   sudo ln -s /etc/nginx/sites-available/domain-max-api /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl reload nginx

   # 配置SSL证书
   sudo certbot --nginx -d api.yourdomain.com
   ```

#### 🗄️ 数据库设置

**选择 1: PlanetScale (推荐)**

```bash
# 1. 注册 https://planetscale.com/
# 2. 创建数据库
# 3. 获取连接信息，配置环境变量：
DB_TYPE=mysql
DB_HOST=your-db.planetscale.com
DB_PORT=3306
DB_SSL_MODE=require
```

**选择 2: Supabase**

```bash
# 1. 注册 https://supabase.com/
# 2. 创建项目
# 3. 配置环境变量：
DB_TYPE=postgres
DB_HOST=db.your-project.supabase.co
DB_PORT=5432
DB_SSL_MODE=require
```

### 📖 详细部署指南

- **[前后端分离完整部署](docs/separation-deployment.md)** - 完整的分离架构部署指南
- **[后端 API 部署](docs/deployment.md)** - VPS 上的 API 服务部署

### 🔧 开发工具

```bash
# 开发环境管理
./dev.sh start     # 启动完整开发环境
./dev.sh stop      # 停止开发服务
./dev.sh status    # 查看服务状态

# API构建工具
./build-api.sh build    # 构建API服务
./build-api.sh cross    # 构建跨平台版本
./build-api.sh dev      # 开发模式运行
./build-api.sh clean    # 清理构建产物
```

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

## 📋 前后端分离架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Cloudflare Pages│    │   VPS Go API    │    │  Remote Database│
│   React SPA     ├────┤   RESTful API   ├────┤  PlanetScale    │
│   全球 CDN      │    │   CORS 支持     │    │  Supabase      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │              ┌─────────────────┐              │
         │              │   Nginx Proxy   │              │
         └──────────────┤   SSL 证书      ├──────────────┘
                        │   负载均衡      │
                        └─────────────────┘
```

### 服务架构

- **前端服务** - React SPA 部署在 Cloudflare Pages，全球 CDN 加速
- **后端服务** - Go API 服务器部署在 VPS，支持 CORS 跨域
- **数据库** - 远程数据库服务（PlanetScale MySQL、Supabase PostgreSQL 等）
- **代理服务** - Nginx 反向代理，提供 SSL 支持和负载均衡
- **域名解析** - Cloudflare DNS 管理，支持多级域名分发

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

- **Cloudflare Pages** - 全球 CDN 静态网站托管
- **VPS 服务器** - API 服务部署
- **Nginx** - 反向代理和负载均衡
- **Let's Encrypt** - 免费 SSL 证书
- **远程数据库** - 托管数据库服务

## 📖 详细文档

- 📚 [部署指南](docs/deployment.md) - 前后端分离部署步骤和配置
- ️ [系统架构](docs/architecture.md) - 详细的系统设计文档

## 🔧 开发指南

### 快速启动开发环境

使用我们提供的开发脚本一键启动：

```bash
# 一键启动前后端开发环境
./dev.sh

# 或者手动分别启动
# 终端1: 启动后端API服务器
cd domain-max && go run cmd/api-server/main.go

# 终端2: 启动前端开发服务器
cd web && npm run dev
```

### 构建部署

```bash
# 构建后端API
./build-api.sh

# 构建前端SPA
cd web && npm run build

# 部署到VPS
# 1. 上传后端二进制文件到VPS
# 2. 前端构建产物推送到Cloudflare Pages
# 3. 配置Nginx反向代理
```

## 🏥 运维管理

### VPS 服务器管理

```bash
# 使用systemd管理API服务
sudo systemctl start domain-max-api
sudo systemctl stop domain-max-api
sudo systemctl restart domain-max-api
sudo systemctl status domain-max-api

# 查看服务日志
sudo journalctl -u domain-max-api -f

# 手动启动（调试模式）
./domain-max-api
```

### 健康检查

```bash
# API健康检查
curl https://your-api-domain.com/api/health

# 前端访问检查
curl https://your-frontend-domain.com

# 检查服务状态
sudo systemctl status domain-max-api

# 检查端口占用
sudo netstat -tlnp | grep :8080
```

### 远程数据库管理

```bash
# PlanetScale连接示例
# 使用提供的连接字符串连接数据库

# Supabase连接示例
# 通过Web界面或SQL编辑器管理数据库

# 本地数据库迁移测试
go run cmd/migrate/main.go
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

````bash
# 检查端口占用
## 🛠️ 故障排除

### 常见问题

#### 前端无法连接后端API

```bash
# 检查CORS配置
curl -v -H "Origin: https://your-frontend.pages.dev" \
  https://your-api-domain.com/api/health

# 检查环境变量
grep -E "^(VITE_|REACT_)" web/.env

# 验证API地址配置
cat web/.env | grep VITE_API_URL
````

#### VPS 部署 API 服务启动失败

```bash
# 检查端口占用
sudo netstat -tlnp | grep :8080
sudo lsof -i :8080

# 检查服务日志
sudo journalctl -u domain-max-api -f

# 检查环境变量文件
cat .env.vps

# 检查文件权限
ls -la domain-max-api
chmod +x domain-max-api
```

#### 数据库连接失败

```bash
# 测试PlanetScale连接
mysql -h your-host -P 3306 -u your-user -p your-database

# 测试Supabase连接
psql "postgresql://user:pass@host:port/dbname?sslmode=require"

# 检查数据库配置
grep -E "^(DB_|DATABASE_)" .env.vps

# 验证网络连接
ping your-database-host
telnet your-database-host 3306
```

#### Cloudflare Pages 构建失败

```bash
# 本地测试构建
cd web && npm run build

# 检查构建配置
cat web/package.json | grep scripts -A 10

# 验证环境变量
# 在Cloudflare Pages设置中检查环境变量配置

# 检查依赖版本
npm ls
```

### 性能优化

```bash
# 前端性能分析
cd web && npm run build -- --analyze

# API性能监控
# 配置适当的监控和日志记录

# 数据库查询优化
# 使用数据库提供商的性能监控工具

# CDN缓存验证
curl -I https://your-frontend.pages.dev/
# 检查Cache-Control和CF-Cache-Status头
```

## 🌐 部署架构优势

### Cloudflare Pages 优势

- ✅ 全球 CDN 加速，访问速度快
- ✅ 自动 HTTPS 和 SSL 证书
- ✅ Git 集成，自动构建部署
- ✅ 免费额度充足
- ✅ 高可用性和容错能力

### VPS API 服务优势

- ✅ 完全控制服务器环境
- ✅ 成本可控，性能可预测
- ✅ 支持复杂业务逻辑
- ✅ 易于监控和调试
- ✅ 数据安全可控

### 远程数据库优势

- ✅ 专业数据库管理
- ✅ 自动备份和恢复
- ✅ 高可用性保障
- ✅ 按需扩展
- ✅ 安全防护完善

## 🤝 贡献指南

我们欢迎社区贡献！请遵循以下步骤：

1. **Fork** 项目仓库
2. **创建** 功能分支 (`git checkout -b feature/amazing-feature`)
3. **提交** 更改 (`git commit -m 'Add amazing feature'`)
4. **推送** 分支 (`git push origin feature/amazing-feature`)
5. **创建** Pull Request

### 开发规范

- 遵循 Go 代码规范
- 前端使用 TypeScript 和 React 最佳实践
- 编写单元测试和集成测试
- 更新相关文档
- 确保所有测试通过

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 获取帮助

- 📖 查看项目文档
- 🐛 报告问题和 Bug
- 💬 参与社区讨论
- 📧 技术支持咨询

## 🙏 致谢

感谢所有为本项目做出贡献的开发者和社区成员！

---

**Domain MAX** 致力于为用户提供最佳的域名管理体验，采用现代化的前后端分离架构，结合 Cloudflare 全球 CDN 和 VPS 部署，实现高性能、高可用、低成本的域名管理解决方案。

_Built with ❤️ by Domain MAX Team_

**Domain MAX** - 让域名管理更简单、更安全、更高效！

Made with ❤️ by Domain MAX Team
