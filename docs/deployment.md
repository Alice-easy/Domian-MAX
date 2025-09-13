# Domain MAX 后端 API 部署指南

本文档介绍 Domain MAX 后端 API 服务的 VPS 部署方式，适用于前后端分离架构。

> � **架构说明**：前端部署在 Cloudflare Pages，后端 API 部署在 VPS，数据库使用远程服务
>
> 📖 **完整分离部署**：查看 [前后端分离部署指南](separation-deployment.md)

## 📋 VPS 要求

**推荐配置**

- CPU: 1-2 核心
- 内存: 1-2GB RAM
- 存储: 10GB SSD
- 操作系统: Ubuntu 20.04+

**必需软件**

- Go 1.23+
- Nginx (反向代理)
- Git

## 🚀 API 服务部署

### 1. 服务器准备

```bash
# 连接到 VPS
ssh user@your-vps-ip

# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装 Go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装其他依赖
sudo apt install -y git nginx certbot python3-certbot-nginx
```

### 2. 部署代码

```bash
# 克隆项目
git clone https://github.com/your-username/domain-max.git
cd domain-max

# 配置环境变量
cp .env.example .env
nano .env
```

**环境变量配置**：

```bash
# API 服务配置
APP_MODE=production
PORT=8080
BASE_URL=https://api.yourdomain.com

# 远程数据库配置（PlanetScale/Supabase/自建）
DB_TYPE=postgres
DB_HOST=your-remote-db-host
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_strong_password
DB_NAME=domain_manager
DB_SSL_MODE=require

# CORS 配置（允许前端域名）
CORS_ALLOWED_ORIGINS=https://your-app.pages.dev,https://yourdomain.com

# 安全配置
JWT_SECRET=your_64_character_jwt_secret_here
ENCRYPTION_KEY=your_64_character_encryption_key_here
```

### 3. 构建和启动

```bash
# 构建 API 服务
go build -o domain-max-api ./cmd/api-server

# 测试运行
./domain-max-api
```

### 4. 配置系统服务

```bash
# 创建 systemd 服务
sudo nano /etc/systemd/system/domain-max-api.service
```

```ini
[Unit]
Description=Domain MAX API Server
After=network.target

[Service]
Type=simple
User=ubuntu
Group=ubuntu
WorkingDirectory=/home/ubuntu/domain-max
ExecStart=/home/ubuntu/domain-max/domain-max-api
Restart=always
RestartSec=5
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
```

```bash
# 启用服务
sudo systemctl daemon-reload
sudo systemctl enable domain-max-api
sudo systemctl start domain-max-api
sudo systemctl status domain-max-api
```

### 5. 配置 Nginx 反向代理

```bash
sudo nano /etc/nginx/sites-available/domain-max-api
```

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

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

        # CORS 配置
        add_header Access-Control-Allow-Origin "https://your-app.pages.dev" always;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
        add_header Access-Control-Allow-Headers "Origin, Content-Type, Authorization" always;
    }
}
```

```bash
# 启用配置
sudo ln -s /etc/nginx/sites-available/domain-max-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 6. 配置 SSL 证书

```bash
# 获取 Let's Encrypt 证书
sudo certbot --nginx -d api.yourdomain.com

# 设置自动续期
echo "0 12 * * * /usr/bin/certbot renew --quiet" | sudo crontab -
```

### 7. 验证部署

```bash
# API 健康检查
curl https://api.yourdomain.com/health

# 检查服务状态
sudo systemctl status domain-max-api
```

## 📊 维护和监控

### 服务管理

```bash
# 查看服务状态
sudo systemctl status domain-max-api

# 重启服务
sudo systemctl restart domain-max-api

# 查看日志
sudo journalctl -u domain-max-api -f
```

### 更新部署

```bash
# 拉取最新代码
git pull origin main

# 重新构建
go build -o domain-max-api ./cmd/api-server

# 重启服务
sudo systemctl restart domain-max-api

# 验证更新
curl https://api.yourdomain.com/health
```

### 备份

```bash
# 备份配置和可执行文件
tar -czf api_backup_$(date +%Y%m%d_%H%M%S).tar.gz .env domain-max-api
```

## 🔧 故障排除

### 常见问题

**1. API 服务无法启动**

```bash
# 检查配置文件
cat .env

# 查看详细错误
sudo journalctl -u domain-max-api -n 50

# 手动启动测试
./domain-max-api
```

**2. CORS 错误**

```bash
# 检查 CORS 配置
curl -H "Origin: https://your-app.pages.dev" \
     -H "Access-Control-Request-Method: POST" \
     -X OPTIONS \
     https://api.yourdomain.com/api/v1/auth/login
```

**3. SSL 证书问题**

```bash
# 检查证书状态
sudo certbot certificates

# 手动续期测试
sudo certbot renew --dry-run
```

## 📚 相关文档

- [前后端分离完整部署指南](separation-deployment.md)
- [项目 README](../README.md)
