# Domain MAX - 生产环境配置指南

## 🎯 生产环境准备清单

### 1. 安全配置 ✅

#### SSL/TLS 证书

- [ ] 获取有效的 SSL 证书 (Let's Encrypt 或商业 CA)
- [ ] 配置证书自动续期
- [ ] 启用 HSTS 和其他安全标头
- [ ] 禁用 HTTP，强制 HTTPS

#### 密钥和密码安全

- [ ] 生成强随机密码 (最少 16 字符)
- [ ] 使用密钥管理服务 (如 HashiCorp Vault)
- [ ] 定期轮换密钥
- [ ] 启用数据库 SSL 连接

#### 网络安全

- [ ] 配置防火墙规则
- [ ] 限制数据库端口访问
- [ ] 配置 VPN 或专用网络
- [ ] 启用 DDoS 防护

### 2. 性能优化 ⚡

#### 应用层优化

```yaml
# 生产环境 docker-compose.override.yml
version: "3.8"
services:
  app:
    deploy:
      replicas: 3 # 多实例部署
      resources:
        limits:
          cpus: "4"
          memory: 2G
        reservations:
          cpus: "1"
          memory: 512M
    environment:
      - GOMAXPROCS=4
      - LOG_LEVEL=warn # 减少日志输出
```

#### 数据库优化

```sql
-- PostgreSQL 生产配置优化
ALTER SYSTEM SET max_connections = 200;
ALTER SYSTEM SET shared_buffers = '512MB';
ALTER SYSTEM SET effective_cache_size = '2GB';
ALTER SYSTEM SET work_mem = '8MB';
ALTER SYSTEM SET maintenance_work_mem = '128MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '32MB';
ALTER SYSTEM SET max_wal_size = '2GB';
SELECT pg_reload_conf();
```

#### Redis 缓存优化

```conf
# redis-production.conf
maxmemory 1gb
maxmemory-policy allkeys-lru
save 900 1
save 300 10
save 60 10000
tcp-keepalive 300
timeout 0
```

### 3. 监控和日志 📊

#### 应用监控

```yaml
# monitoring/docker-compose.yml
version: "3.8"
services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"
      - "--storage.tsdb.path=/prometheus"
      - "--web.console.libraries=/etc/prometheus/console_libraries"
      - "--web.console.templates=/etc/prometheus/consoles"

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=secure_password
    volumes:
      - grafana_data:/var/lib/grafana

volumes:
  grafana_data:
```

#### 日志聚合

```yaml
# ELK Stack for log aggregation
elasticsearch:
  image: docker.elastic.co/elasticsearch/elasticsearch:7.15.0
  environment:
    - discovery.type=single-node
    - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
  volumes:
    - elasticsearch_data:/usr/share/elasticsearch/data

logstash:
  image: docker.elastic.co/logstash/logstash:7.15.0
  volumes:
    - ./logstash.conf:/usr/share/logstash/pipeline/logstash.conf
  depends_on:
    - elasticsearch

kibana:
  image: docker.elastic.co/kibana/kibana:7.15.0
  ports:
    - "5601:5601"
  environment:
    - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
  depends_on:
    - elasticsearch
```

### 4. 高可用性配置 🏗️

#### 负载均衡器配置

```nginx
# nginx-ha.conf
upstream backend_servers {
    least_conn;
    server app1:8080 weight=3 max_fails=3 fail_timeout=30s;
    server app2:8080 weight=3 max_fails=3 fail_timeout=30s;
    server app3:8080 weight=2 max_fails=3 fail_timeout=30s;
    keepalive 32;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    location / {
        proxy_pass http://backend_servers;
        proxy_next_upstream error timeout invalid_header http_500 http_502 http_503;
        proxy_next_upstream_tries 3;
        proxy_next_upstream_timeout 10s;
    }
}
```

#### 数据库主从复制

```yaml
# PostgreSQL Master-Slave setup
services:
  postgres-master:
    image: postgres:15-alpine
    environment:
      - POSTGRES_REPLICATION_USER=replicator
      - POSTGRES_REPLICATION_PASSWORD=repl_password
    command: |
      postgres 
      -c wal_level=replica 
      -c max_wal_senders=3 
      -c max_replication_slots=3

  postgres-slave:
    image: postgres:15-alpine
    environment:
      - PGUSER=postgres
      - POSTGRES_MASTER_SERVICE=postgres-master
    command: |
      bash -c '
      if [ ! -s "/var/lib/postgresql/data/PG_VERSION" ]; then
        pg_basebackup -h postgres-master -D /var/lib/postgresql/data -U replicator -v -P -W
        echo "standby_mode = on" >> /var/lib/postgresql/data/recovery.conf
        echo "primary_conninfo = \"host=postgres-master port=5432 user=replicator\"" >> /var/lib/postgresql/data/recovery.conf
      fi
      postgres
      '
```

### 5. 备份和恢复策略 💾

#### 自动化备份脚本

```bash
#!/bin/bash
# production-backup.sh

set -e

BACKUP_DIR="/backups"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# 创建备份目录
mkdir -p "$BACKUP_DIR/$DATE"

# 数据库备份
echo "🗄️ Backing up database..."
docker-compose exec -T postgres pg_dump -U postgres -Fc domain_manager > \
  "$BACKUP_DIR/$DATE/database.dump"

# Redis 备份
echo "💾 Backing up Redis..."
docker-compose exec redis redis-cli BGSAVE
sleep 5
docker cp domain-max-redis:/data/dump.rdb "$BACKUP_DIR/$DATE/redis.rdb"

# 配置文件备份
echo "⚙️ Backing up configurations..."
tar -czf "$BACKUP_DIR/$DATE/configs.tar.gz" \
  deployments/.env \
  deployments/nginx/ \
  deployments/ssl/

# 压缩备份
echo "📦 Compressing backup..."
cd "$BACKUP_DIR"
tar -czf "$DATE.tar.gz" "$DATE/"
rm -rf "$DATE/"

# 清理旧备份
echo "🧹 Cleaning old backups..."
find "$BACKUP_DIR" -name "*.tar.gz" -mtime +$RETENTION_DAYS -delete

# 上传到云存储 (可选)
if [ -n "$AWS_S3_BUCKET" ]; then
  echo "☁️ Uploading to S3..."
  aws s3 cp "$DATE.tar.gz" "s3://$AWS_S3_BUCKET/backups/"
fi

echo "✅ Backup completed: $DATE.tar.gz"
```

#### Cron 任务配置

```cron
# /etc/crontab
# 每日凌晨 2 点执行备份
0 2 * * * root /opt/domain-max/scripts/production-backup.sh

# 每小时检查服务健康状态
0 * * * * root /opt/domain-max/scripts/health-check.sh

# 每天轮转日志
0 0 * * * root /opt/domain-max/scripts/rotate-logs.sh
```

### 6. 健康检查和告警 🚨

#### 健康检查脚本

```bash
#!/bin/bash
# health-check.sh

check_service() {
    local service=$1
    local endpoint=$2
    local max_retries=3
    local retry=0

    while [ $retry -lt $max_retries ]; do
        if curl -sf "$endpoint" >/dev/null 2>&1; then
            echo "✅ $service is healthy"
            return 0
        fi
        retry=$((retry + 1))
        sleep 10
    done

    echo "❌ $service is unhealthy"
    return 1
}

# 检查各服务
UNHEALTHY_SERVICES=""

if ! check_service "Application" "http://localhost:8080/api/health"; then
    UNHEALTHY_SERVICES="$UNHEALTHY_SERVICES Application"
fi

if ! check_service "Nginx" "http://localhost/health"; then
    UNHEALTHY_SERVICES="$UNHEALTHY_SERVICES Nginx"
fi

if ! docker-compose exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
    echo "❌ Database is unhealthy"
    UNHEALTHY_SERVICES="$UNHEALTHY_SERVICES Database"
fi

# 发送告警
if [ -n "$UNHEALTHY_SERVICES" ]; then
    MESSAGE="🚨 Domain MAX Alert: Services unhealthy: $UNHEALTHY_SERVICES"

    # 发送邮件告警
    if [ -n "$ALERT_EMAIL" ]; then
        echo "$MESSAGE" | mail -s "Domain MAX Alert" "$ALERT_EMAIL"
    fi

    # 发送 Slack 通知
    if [ -n "$SLACK_WEBHOOK" ]; then
        curl -X POST -H 'Content-type: application/json' \
          --data "{\"text\":\"$MESSAGE\"}" \
          "$SLACK_WEBHOOK"
    fi

    exit 1
fi

echo "✅ All services are healthy"
```

### 7. 安全加固 🔒

#### 系统级安全

```bash
# 防火墙配置
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable

# 禁用不必要的服务
systemctl disable apache2
systemctl disable sendmail
systemctl disable telnet

# 配置 fail2ban
apt-get install fail2ban
```

#### 应用安全配置

```yaml
# docker-compose.security.yml
version: "3.8"
services:
  app:
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
    user: "1001:1001"
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE

  postgres:
    security_opt:
      - no-new-privileges:true
    user: "999:999"
    cap_drop:
      - ALL
```

### 8. 性能调优 📈

#### 应用级调优

```go
// main.go 生产环境配置
func main() {
    // 设置 Go runtime 参数
    runtime.GOMAXPROCS(runtime.NumCPU())

    // 配置 GC
    debug.SetGCPercent(100)

    // 配置连接池
    db.SetMaxOpenConns(50)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(time.Hour)

    // 启用 pprof (仅在需要时)
    if os.Getenv("ENABLE_PPROF") == "true" {
        go func() {
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
    }
}
```

#### Nginx 性能调优

```nginx
# nginx-performance.conf
worker_processes auto;
worker_rlimit_nofile 65535;

events {
    worker_connections 65535;
    use epoll;
    multi_accept on;
}

http {
    # 启用缓存
    proxy_cache_path /var/cache/nginx/api levels=1:2 keys_zone=api_cache:10m
                     max_size=1g inactive=60m use_temp_path=off;

    # 压缩优化
    gzip_comp_level 6;
    gzip_min_length 1000;

    # 连接保持
    keepalive_timeout 65;
    keepalive_requests 1000;

    # 缓存配置
    location /api/ {
        proxy_cache api_cache;
        proxy_cache_valid 200 302 5m;
        proxy_cache_valid 404 1m;
        proxy_cache_use_stale error timeout updating http_500 http_502 http_503 http_504;
        proxy_cache_lock on;
        add_header X-Cache-Status $upstream_cache_status;
    }
}
```

### 9. 容器编排 (Kubernetes) ☸️

#### Kubernetes 部署配置

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: domain-max-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: domain-max
  template:
    metadata:
      labels:
        app: domain-max
    spec:
      containers:
        - name: domain-max
          image: domain-max:latest
          ports:
            - containerPort: 8080
          env:
            - name: DB_HOST
              value: "postgres-service"
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "1Gi"
              cpu: "1000m"
          livenessProbe:
            httpGet:
              path: /api/health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /api/health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
```

### 10. 灾难恢复计划 🔄

#### 恢复流程

```bash
#!/bin/bash
# disaster-recovery.sh

# 1. 停止所有服务
docker-compose down

# 2. 恢复数据库
echo "🔄 Restoring database..."
docker-compose up -d postgres
sleep 30
docker-compose exec -T postgres pg_restore -U postgres -d domain_manager < backup/database.dump

# 3. 恢复 Redis 数据
echo "🔄 Restoring Redis..."
docker cp backup/redis.rdb domain-max-redis:/data/dump.rdb
docker-compose restart redis

# 4. 恢复配置文件
echo "🔄 Restoring configurations..."
tar -xzf backup/configs.tar.gz -C /

# 5. 启动所有服务
echo "🚀 Starting all services..."
docker-compose up -d

# 6. 验证恢复
echo "✅ Verifying recovery..."
sleep 60
./scripts/health-check.sh
```

## 📋 生产环境检查清单

### 部署前检查

- [ ] SSL 证书已配置并有效
- [ ] 所有密码和密钥已更新为生产值
- [ ] 防火墙规则已配置
- [ ] 监控系统已设置
- [ ] 备份策略已实施
- [ ] 告警机制已配置
- [ ] 负载测试已完成
- [ ] 安全扫描已通过

### 部署后验证

- [ ] 所有服务健康检查通过
- [ ] API 端点响应正常
- [ ] 数据库连接正常
- [ ] SSL 证书正确配置
- [ ] 性能指标正常
- [ ] 日志记录正常
- [ ] 备份功能正常
- [ ] 监控告警正常

---

**Domain MAX 生产环境配置指南** | 版本 1.0.0
