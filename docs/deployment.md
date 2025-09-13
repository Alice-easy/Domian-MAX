# Domain MAX 部署指南

本文档详细介绍了 Domain MAX 的源码部署方式，包括本地开发和生产环境部署。

## 📋 部署前准备

### 系统要求

**最低配置**

- CPU: 1 核心
- 内存: 1GB RAM
- 存储: 10GB 可用空间
- 操作系统: Linux/macOS/Windows

**推荐配置**

- CPU: 2 核心以上
- 内存: 2GB RAM 以上
- 存储: 20GB 可用空间
- 操作系统: Ubuntu 20.04+ / CentOS 8+ / macOS 12+

### 依赖软件

**必需软件**

- Go 1.23+
- Node.js 18+
- PostgreSQL 12+ 或 MySQL 8.0+

**可选软件**

- Nginx (反向代理)
- PM2 (进程管理)
- Redis (缓存)

## 🚀 快速部署

### 1. 下载项目

```bash
git clone <repository-url>
cd domain-max
```

### 2. 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env` 文件，设置必要的配置：

```bash
# 数据库密码 (必需)
DB_PASSWORD=your_secure_password_here

# JWT密钥 (必需，生产环境建议64位以上)
JWT_SECRET=your_jwt_secret_key_here_at_least_64_characters_long

# 加密密钥 (必需，32字节十六进制)
ENCRYPTION_KEY=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef

# 可选：邮件配置
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTP_FROM=noreply@yourdomain.com
```

### 3. 安装依赖

```bash
# 安装Go依赖
go mod tidy

# 安装前端依赖
cd web && npm install && cd ..
```

### 4. 初始化数据库

```bash
# 使用提供的SQL脚本初始化数据库
# PostgreSQL
psql -U domain_user -d domain_manager -f init.sql

# MySQL
mysql -u domain_user -p domain_manager < init.sql
```

### 5. 构建应用

```bash
# 使用构建脚本
./scripts/build.sh

# 或者手动构建
cd web && npm run build && cd ..
go build -o domain-max ./cmd/server
```

### 6. 运行应用

```bash
./domain-max
```

### 7. 验证部署

```bash
# 健康检查
curl http://localhost:8080/api/health

# 检查应用日志
tail -f domain-max.log
```

### 8. 访问应用

- 应用地址: http://localhost:8080
- 默认管理员: admin@example.com / admin123

**⚠️ 重要：首次登录后请立即修改默认密码！**

## 🛠️ 本地开发部署

### 1. 环境准备

```bash
# 安装Go (如果未安装)
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 安装Node.js (如果未安装)
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 验证安装
go version
node --version
npm --version
```

### 2. 数据库准备

**PostgreSQL**

```bash
# 安装PostgreSQL
sudo apt-get install postgresql postgresql-contrib

# 创建数据库和用户
sudo -u postgres psql
CREATE DATABASE domain_manager;
CREATE USER domain_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE domain_manager TO domain_user;
\q
```

**MySQL**

```bash
# 安装MySQL
sudo apt-get install mysql-server

# 创建数据库和用户
sudo mysql
CREATE DATABASE domain_manager;
CREATE USER 'domain_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON domain_manager.* TO 'domain_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 3. 项目配置

```bash
# 克隆项目
git clone <repository-url>
cd domain-max

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件

# 安装依赖
go mod tidy
cd web && npm install && cd ..
```

### 4. 初始化数据库

```bash
# 使用提供的SQL脚本初始化数据库
# PostgreSQL
psql -U domain_user -d domain_manager -f init.sql

# MySQL
mysql -u domain_user -p domain_manager < init.sql
```

### 4. 构建和运行

```bash
# 使用构建脚本
./scripts/build.sh

# 或者手动构建
cd web && npm run build && cd ..
go build -o domain-max ./cmd/server

# 运行应用
./domain-max
```

## 🏭 生产环境部署

### 1. 服务器准备

```bash
# 更新系统
sudo apt-get update && sudo apt-get upgrade -y

# 安装必要软件
sudo apt-get install -y curl wget git nginx certbot python3-certbot-nginx

# 安装Go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 安装Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安装PM2 (进程管理器)
sudo npm install -g pm2
```

### 2. 安全配置

```bash
# 配置防火墙
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable

# 创建应用用户
sudo useradd -m -s /bin/bash domain-max
```

### 3. 部署应用

```bash
# 切换到应用用户
sudo su - domain-max

# 克隆项目
git clone <repository-url>
cd domain-max

# 配置生产环境变量
cp .env.example .env
```

编辑生产环境配置：

```bash
# 生产环境配置
ENVIRONMENT=production
BASE_URL=https://yourdomain.com
HTTP_PORT=8080

# 强密码配置
DB_PASSWORD=<strong-random-password>
JWT_SECRET=<64-character-random-string>
ENCRYPTION_KEY=<32-byte-hex-string>

# 邮件配置
SMTP_HOST=smtp.yourdomain.com
SMTP_PORT=587
SMTP_USER=noreply@yourdomain.com
SMTP_PASSWORD=<smtp-password>
SMTP_FROM=noreply@yourdomain.com
```

```bash
# 构建应用
./scripts/build.sh

# 配置PM2启动应用
pm2 start ./domain-max --name "domain-max"
pm2 save
pm2 startup
```

### 4. 配置反向代理

创建 Nginx 配置文件：

```bash
sudo nano /etc/nginx/sites-available/domain-max
```

```nginx
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;

    # 重定向到HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    # SSL配置
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    # SSL安全配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # 安全头
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # 代理到应用
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # 超时配置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 静态文件缓存
    location /static/ {
        proxy_pass http://localhost:8080;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/domain-max /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 5. 配置 SSL 证书

```bash
# 获取Let's Encrypt证书
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# 设置自动续期
sudo crontab -e
# 添加以下行
0 12 * * * /usr/bin/certbot renew --quiet
```

### 6. 配置进程管理

创建 systemd 服务文件：

```bash
sudo nano /etc/systemd/system/domain-max.service
```

```ini
[Unit]
Description=Domain MAX Application
After=network.target

[Service]
Type=simple
User=domain-max
Group=domain-max
WorkingDirectory=/home/domain-max/domain-max
ExecStart=/home/domain-max/domain-max/domain-max
Restart=always
RestartSec=5
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
```

启用服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable domain-max
sudo systemctl start domain-max
sudo systemctl status domain-max
```

## 📊 监控和维护

### 健康检查

```bash
# 检查应用状态
curl -f http://localhost:8080/api/health || echo "Service is down"

# 检查systemd服务状态
sudo systemctl status domain-max

# 检查PM2进程状态
pm2 status

# 查看应用日志
sudo journalctl -u domain-max -f
# 或者查看PM2日志
pm2 logs domain-max
```

### 备份策略

**数据库备份**

```bash
# PostgreSQL备份
pg_dump -U domain_user -d domain_manager > backup_$(date +%Y%m%d_%H%M%S).sql

# MySQL备份
mysqldump -u domain_user -p domain_manager > backup_$(date +%Y%m%d_%H%M%S).sql
```

**配置备份**

```bash
# 备份配置文件和可执行文件
tar -czf config_backup_$(date +%Y%m%d_%H%M%S).tar.gz .env domain-max init.sql
```

### 更新部署

```bash
# 拉取最新代码
git pull origin main

# 停止服务
sudo systemctl stop domain-max
# 或者使用PM2
pm2 stop domain-max

# 重新构建
./scripts/build.sh

# 启动服务
sudo systemctl start domain-max
# 或者使用PM2
pm2 restart domain-max

# 验证更新
curl http://localhost:8080/api/health
```

### 性能优化

**数据库优化**

```sql
-- PostgreSQL优化
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
SELECT pg_reload_conf();
```

**应用优化**

```bash
# 编译时优化
go build -ldflags="-s -w" -o domain-max ./cmd/server

# 设置Go运行时参数
export GOMAXPROCS=2
export GOGC=100

# 前端优化
cd web && npm run build -- --mode production
```

## 🔧 故障排除

### 常见问题

**1. 数据库连接失败**

```bash
# 检查数据库状态
sudo systemctl status postgresql
# 或者MySQL
sudo systemctl status mysql

# 测试数据库连接
psql -U domain_user -d domain_manager -c "SELECT 1;"
# 或者MySQL
mysql -u domain_user -p domain_manager -e "SELECT 1;"

# 检查网络连接
telnet localhost 5432  # PostgreSQL
telnet localhost 3306  # MySQL
```

**2. 应用启动失败**

```bash
# 检查应用日志
sudo journalctl -u domain-max -n 50

# 检查配置文件
cat .env

# 手动启动测试
./domain-max
```

**3. 前端资源加载失败**

```bash
# 检查构建输出
ls -la web/dist/

# 重新构建前端
cd web && npm run build && cd ..
go build -o domain-max ./cmd/server
```

**4. SSL 证书问题**

```bash
# 检查证书状态
sudo certbot certificates

# 手动续期
sudo certbot renew --dry-run
```

### 日志分析

```bash
# 应用日志
sudo journalctl -u domain-max -f
# 或者PM2日志
pm2 logs domain-max --lines 100

# 数据库日志
# PostgreSQL
sudo tail -f /var/log/postgresql/postgresql-*.log
# MySQL
sudo tail -f /var/log/mysql/error.log

# Nginx日志
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# 系统日志
sudo journalctl -f
```

### 性能调试

```bash
# 检查系统资源使用
top
htop
free -h
df -h

# 检查应用性能
ps aux | grep domain-max

# Go性能分析
# 在应用中启用pprof，然后访问
curl http://localhost:8080/debug/pprof/
go tool pprof http://localhost:8080/debug/pprof/profile

# 检查数据库性能
# PostgreSQL
psql -U domain_user -d domain_manager -c "SELECT * FROM pg_stat_activity;"
# MySQL
mysql -u domain_user -p domain_manager -e "SHOW PROCESSLIST;"

# 网络连接检查
netstat -tuln | grep :8080
ss -tuln | grep :8080
```

## 📚 参考资料

- [Go 官方文档](https://golang.org/doc/)
- [Node.js 文档](https://nodejs.org/docs/)
- [PostgreSQL 文档](https://www.postgresql.org/docs/)
- [MySQL 文档](https://dev.mysql.com/doc/)
- [Nginx 配置指南](https://nginx.org/en/docs/)
- [Let's Encrypt 文档](https://letsencrypt.org/docs/)
- [PM2 进程管理](https://pm2.keymetrics.io/docs/)
- [systemd 服务管理](https://systemd.io/)

---

如有部署问题，请查看[架构文档](architecture.md)或提交[Issue](../../issues)。
