# Domain MAX 前后端分离部署指南

## 🏗️ 架构概览

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

> � **优势**：前端全球 CDN 加速，后端 VPS 完全控制，数据库托管服务高可用

## 🗄️ 第一步：配置远程数据库

### PlanetScale (推荐)

```bash
# 1. 注册 https://planetscale.com/
# 2. 创建数据库，获取连接信息
# 3. 配置环境变量：
DB_TYPE=mysql
DB_HOST=your-db.planetscale.com
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSL_MODE=require
```

### Supabase

```bash
# 1. 注册 https://supabase.com/
# 2. 创建项目，获取数据库连接信息
# 3. 配置环境变量：
DB_TYPE=postgres
DB_HOST=db.your-project.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=postgres
DB_SSL_MODE=require
```

## 🚀 第二步：部署后端到 VPS

### 快速部署

```bash
# 连接 VPS
ssh user@your-vps-ip

# 安装依赖
sudo apt update && sudo apt upgrade -y
sudo apt install -y git nginx certbot python3-certbot-nginx

# 安装 Go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 部署代码
git clone https://github.com/your-username/domain-max.git
cd domain-max
```

### 环境配置

```bash
cp .env.example .env
nano .env
```

```bash
# API 服务配置
APP_MODE=production
PORT=8080
BASE_URL=https://api.yourdomain.com

# 远程数据库配置
DB_TYPE=postgres
DB_HOST=your-remote-db-host
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_strong_password
DB_NAME=domain_manager
DB_SSL_MODE=require

# CORS 配置（允许前端域名访问）
CORS_ALLOWED_ORIGINS=https://your-app.pages.dev,https://yourdomain.com

# 安全配置
JWT_SECRET=your_64_character_jwt_secret_here
ENCRYPTION_KEY=your_64_character_encryption_key_here
```

### 服务部署

```bash
# 构建 API 服务
go build -o domain-max-api ./cmd/api-server

# 配置系统服务
sudo tee /etc/systemd/system/domain-max-api.service > /dev/null <<EOF
[Unit]
Description=Domain MAX API Server
After=network.target

[Service]
Type=simple
User=$USER
Group=$USER
WorkingDirectory=$PWD
ExecStart=$PWD/domain-max-api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable domain-max-api
sudo systemctl start domain-max-api
```

### Nginx 配置

```bash
sudo tee /etc/nginx/sites-available/domain-max-api > /dev/null <<EOF
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

sudo ln -s /etc/nginx/sites-available/domain-max-api /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```

### SSL 证书

```bash
sudo certbot --nginx -d api.yourdomain.com
echo "0 12 * * * /usr/bin/certbot renew --quiet" | sudo crontab -
```

## 🌐 第三步：部署前端到 Cloudflare Pages

### 连接 GitHub

1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com)
2. 进入 **Pages** → **Create a project**
3. 连接 GitHub 仓库
4. 配置构建设置：

```yaml
Build command: cd web && npm ci && npm run build
Build output directory: web/dist
Root directory: (leave empty)
Environment variables:
  NODE_ENV: production
  VITE_API_BASE_URL: https://api.yourdomain.com
  VITE_BACKEND_DOMAIN: api.yourdomain.com
```

### 自定义域名（可选）

1. Pages 项目设置 → **Custom domains**
2. 添加域名 `yourdomain.com`
3. 更新后端 CORS 配置包含新域名

## ✅ 第四步：验证部署

### API 测试

```bash
# 健康检查
curl https://api.yourdomain.com/health

# 登录接口测试
curl -X POST https://api.yourdomain.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
```

### 前端测试

1. 访问 `https://your-app.pages.dev`
2. 测试登录功能
3. 检查浏览器开发者工具确认无 CORS 错误

### CORS 验证

```bash
curl -H "Origin: https://your-app.pages.dev" \
     -H "Access-Control-Request-Method: POST" \
     -X OPTIONS \
     https://api.yourdomain.com/api/v1/auth/login
```

## 🔧 维护和监控

### 后端维护

```bash
# 服务状态
sudo systemctl status domain-max-api

# 查看日志
sudo journalctl -u domain-max-api -f

# 重启服务
sudo systemctl restart domain-max-api

# 更新部署
git pull origin main
go build -o domain-max-api ./cmd/api-server
sudo systemctl restart domain-max-api
```

### 前端更新

```bash
# Cloudflare Pages 自动构建
# 每次 push 到 main 分支自动更新
git push origin main
```

### 数据库备份

- **PlanetScale/Supabase**：自动备份
- **自建数据库**：

```bash
pg_dump -h your-db-host -U your-user -d domain_manager > backup_$(date +%Y%m%d).sql
```

## 🎯 性能优化建议

### Cloudflare 优化

- ✅ 启用 Brotli 压缩
- ✅ 配置缓存规则
- ✅ 启用 Auto Minify
- ✅ 配置 Page Rules

### API 优化

```nginx
# Nginx 压缩配置
gzip on;
gzip_types text/plain text/css application/json application/javascript;

# 静态文件缓存
location /api/v1/static/ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

## 🔒 安全检查清单

- ✅ API 域名使用 HTTPS
- ✅ CORS 仅允许信任域名
- ✅ 前端不暴露敏感信息
- ✅ 数据库使用 SSL 连接
- ✅ 定期更新系统和依赖

---

🎉 **部署完成！** 享受现代化的前后端分离架构带来的性能和维护优势！
