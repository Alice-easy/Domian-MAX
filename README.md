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

- **Docker** 20.10+
- **Docker Compose** 2.0+
- **内存** 2GB+
- **磁盘空间** 5GB+

### 一键部署

```bash
# 1. 克隆项目
git clone https://github.com/your-repo/domain-max.git
cd domain-max

# 2. 一键部署
./scripts/deploy-complete.sh
```

部署脚本会自动：

- ✅ 复制环境配置文件
- ✅ 生成 SSL 证书
- ✅ 构建应用镜像
- ✅ 启动所有服务
- ✅ 执行健康检查

### 访问应用

部署完成后，您可以通过以下地址访问：

- **前端界面**: https://localhost
- **API 接口**: https://localhost/api
- **管理后台**: https://localhost/admin

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
                                │
                     ┌─────────────────┐
                     │  Nginx 代理     │
                     │  SSL 终止       │
                     │  负载均衡       │
                     └─────────────────┘
```

### 服务组件

| 服务         | 端口    | 描述                             |
| ------------ | ------- | -------------------------------- |
| **nginx**    | 80, 443 | 反向代理、SSL 终止、静态文件服务 |
| **app**      | 8080    | Go 后端应用服务                  |
| **postgres** | 5432    | PostgreSQL 数据库                |
| **redis**    | 6379    | Redis 缓存和会话存储             |

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

- **Docker** - 容器化部署
- **Nginx** - 反向代理和负载均衡
- **Let's Encrypt** - 免费 SSL 证书

## 📖 详细文档

- 📚 [文档中心](docs/) - 完整的项目文档导航
- 🚀 [部署指南](docs/deployment.md) - 完整的部署步骤和配置
- 🏭 [生产环境指南](docs/production-guide.md) - 生产环境优化和安全配置
- 🏗️ [系统架构](docs/architecture.md) - 详细的系统设计文档

## 🔧 开发指南

### 本地开发环境

```bash
# 1. 启动开发环境
docker-compose -f docker-compose.dev.yml up -d

# 2. 前端开发
cd web
npm install
npm run dev

# 3. 后端开发
cd cmd/server
go run main.go
```

### 测试

```bash
# 运行完整系统测试
./scripts/system-test.sh

# 快速测试
./scripts/system-test.sh --quick

# 安全测试
./scripts/system-test.sh --security

# 性能测试
./scripts/system-test.sh --performance
```

## 🏥 运维管理

### 服务管理

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose stop

# 重启服务
docker-compose restart

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 健康检查

```bash
# 应用健康检查
curl https://localhost/api/health

# 服务健康检查
./scripts/deploy-complete.sh --check-health

# 容器状态检查
docker stats
```

### 备份恢复

```bash
# 创建备份
./scripts/backup.sh

# 恢复数据
./scripts/restore.sh backup_file.tar.gz
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

#### 服务无法启动

```bash
# 检查容器状态
docker-compose ps

# 查看错误日志
docker-compose logs app

# 重新构建
docker-compose build --no-cache
```

#### 数据库连接失败

```bash
# 检查数据库状态
docker-compose exec postgres pg_isready -U postgres

# 查看数据库日志
docker-compose logs postgres

# 检查网络连通性
docker-compose exec app ping postgres
```

#### SSL 证书问题

```bash
# 重新生成证书
./scripts/generate-ssl.sh

# 检查证书有效性
openssl x509 -in deployments/ssl/nginx-selfsigned.crt -text -noout
```

### 性能问题诊断

```bash
# 查看资源使用
docker stats

# 分析慢查询
docker-compose exec postgres psql -U postgres -c "SELECT * FROM pg_stat_activity WHERE state = 'active';"

# 检查网络延迟
curl -w "@curl-format.txt" -o /dev/null -s https://localhost/api/health
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
