# GitHub 工作流配置文档（默认禁用状态）

本目录包含了 Domain MAX 项目的完整 GitHub Actions 工作流配置，提供了从开发到生产的全自动化 DevOps 流程。

⚠️ **重要提示**：为了避免意外触发和资源消耗，所有工作流默认处于禁用状态，需要手动启用。

## 🔒 工作流启用机制

### 默认禁用原因

1. **避免意外触发**：防止在项目初期或配置不完整时自动运行
2. **资源控制**：避免不必要的 GitHub Actions 使用量消耗
3. **安全考虑**：确保在环境配置完成后再启用自动化流程
4. **灵活控制**：允许用户根据需要选择性启用工作流

### 启用方法

#### 方法一：手动触发（推荐新用户）

1. 访问 GitHub 仓库的 Actions 页面
2. 选择要运行的工作流
3. 点击 "Run workflow" 按钮
4. **重要**：将 "启用工作流（默认禁用）" 选项设置为 `true`
5. 配置其他参数后运行

#### 方法二：修改工作流文件（启用自动触发）

如需恢复自动触发，请修改工作流文件：

```yaml
# 例如：恢复 CI/CD 自动触发
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]
  workflow_dispatch:
    # 保留手动触发配置
```

#### 方法三：环境变量控制

在工作流中设置环境变量：

```yaml
env:
  WORKFLOW_ENABLED: true # 设置为 true 启用
```

## 📋 工作流总览

| 工作流文件                 | 默认状态 | 启用方式 | 主要功能       | 原触发频率 |
| -------------------------- | -------- | -------- | -------------- | ---------- |
| `ci-cd.yml`                | 🔒 禁用  | 手动     | 持续集成和部署 | 代码变更时 |
| `release.yml`              | 🔒 禁用  | 手动     | 版本发布管理   | 版本发布时 |
| `security-scan.yml`        | 🔒 禁用  | 手动     | 安全扫描       | 每日/按需  |
| `performance-test.yml`     | 🔒 禁用  | 手动     | 性能测试       | 每周/按需  |
| `dependency-updates.yml`   | 🔒 禁用  | 手动     | 依赖更新       | 每周       |
| `backup.yml`               | 🔒 禁用  | 手动     | 数据备份       | 每日       |
| `monitoring.yml`           | 🔒 禁用  | 手动     | 系统监控       | 每 5 分钟  |
| `database-maintenance.yml` | 🔒 禁用  | 手动     | 数据库维护     | 每日/每周  |
| `auto-update.yml`          | 🔒 禁用  | 手动     | 自动更新       | 每日       |
| `Dockerbuild.yml`          | 🔒 禁用  | 手动     | Docker 构建    | Push 时    |

**图例**：

- 🔒 禁用：需要手动启用才能运行
- ✅ 启用：配置完成后可自动运行
- 📝 手动：仅支持手动触发模式

## 🔄 核心工作流详解

### 1. CI/CD Pipeline (`ci-cd.yml`)

**功能**：主要的持续集成和部署工作流  
**默认状态**：🔒 禁用

**启用方法**：

1. 手动触发时设置 `enable_workflow: true`
2. 或修改文件恢复 `push` 和 `pull_request` 触发器

**原触发条件**：

- Push 到 main/develop 分支
- Pull Request
- 手动触发

**主要任务**：

- 代码质量检查
- 单元测试和集成测试
- 安全扫描
- Docker 镜像构建
- 多环境部署

**环境支持**：

- Staging（自动部署）
- Production（需要 approval）

### 2. Release Management (`release.yml`)

**功能**：自动化版本发布流程  
**默认状态**：🔒 禁用

**启用方法**：

1. 手动触发时设置 `enable_workflow: true`
2. 或修改文件恢复 `push: tags` 触发器

**原触发条件**：

- 创建版本标签（v*.*.\*）
- 手动触发

**主要任务**：

- 版本验证
- 安全扫描
- 构建发布版本
- 生成 Release Notes
- 生产部署
- 通知发送

### 3. Security Scanning (`security-scan.yml`)

**功能**：全面的安全扫描

**扫描类型**：

- 静态代码分析（SAST）
- 依赖漏洞扫描
- 容器安全扫描
- 基础设施安全检查

**扫描工具**：

- gosec（Go 安全）
- CodeQL（代码分析）
- Trivy（容器/依赖）
- Semgrep（多语言 SAST）

### 4. Performance Testing (`performance-test.yml`)

**功能**：自动化性能测试

**测试类型**：

- 负载测试
- 压力测试
- 峰值测试
- 稳定性测试

**测试框架**：K6

### 5. Backup System (`backup.yml`)

**功能**：自动化数据备份

**备份内容**：

- 数据库备份（增量/完整）
- 配置文件备份
- 容器镜像备份

**存储方式**：

- GitHub Artifacts（短期）
- AWS S3（长期）

### 6. Monitoring (`monitoring.yml`)

**功能**：系统健康监控

**监控项目**：

- 应用健康检查
- 性能指标监控
- 安全状态检查
- 资源使用监控

**报警机制**：Slack 通知

### 7. Database Maintenance (`database-maintenance.yml`)

**功能**：数据库维护

**维护任务**：

- 统计信息更新
- 过期数据清理
- 索引优化
- VACUUM 操作
- 数据完整性检查

### 8. Auto Updates (`auto-update.yml`)

**功能**：自动更新管理

**更新类型**：

- Go 依赖更新
- Node.js 依赖更新
- Docker 镜像更新
- 安全漏洞修复

**更新方式**：自动创建 PR

## 🚀 快速开始

### 0. 工作流启用准备

**⚠️ 首次使用必读**：

1. **检查工作流状态**：所有工作流默认禁用，避免意外触发
2. **准备环境配置**：确保必要的 Secrets 已配置
3. **选择启用方式**：
   - 新用户：建议先手动测试各个工作流
   - 经验用户：可以修改文件启用自动触发
4. **逐步启用**：建议按以下顺序启用工作流：
   ```
   1. ci-cd.yml（基础 CI/CD）
   2. security-scan.yml（安全扫描）
   3. backup.yml（数据备份）
   4. monitoring.yml（系统监控）
   5. 其他工作流按需启用
   ```

### 1. 环境配置

在 GitHub 仓库设置中配置以下 Secrets：

```bash
# 数据库配置
PROD_DB_HOST=your-prod-db-host
PROD_DB_USER=your-prod-db-user
PROD_DB_PASSWORD=your-prod-db-password
PROD_DB_NAME=your-prod-db-name

STAGING_DB_HOST=your-staging-db-host
STAGING_DB_USER=your-staging-db-user
STAGING_DB_PASSWORD=your-staging-db-password
STAGING_DB_NAME=your-staging-db-name

# 应用URL
PROD_APP_URL=https://your-prod-domain.com
PROD_API_URL=https://api-prod-domain.com
STAGING_APP_URL=https://staging-domain.com
STAGING_API_URL=https://api-staging-domain.com

# 容器注册表
DOCKER_REGISTRY=ghcr.io
DOCKER_USERNAME=your-username
DOCKER_PASSWORD=your-token

# 备份配置
AWS_ACCESS_KEY_ID=your-aws-key
AWS_SECRET_ACCESS_KEY=your-aws-secret
AWS_REGION=us-east-1
BACKUP_BUCKET=your-backup-bucket

# 通知配置
SLACK_WEBHOOK_URL=your-slack-webhook

# GitHub Token（用于自动更新）
PAT_TOKEN=your-personal-access-token
```

### 2. 工作流启用

**重要**：所有工作流默认禁用，需要手动启用！

#### 方法一：手动触发（推荐新用户）

1. 访问 GitHub 仓库的 Actions 页面
2. 选择要运行的工作流
3. 点击 "Run workflow" 按钮
4. **关键步骤**：将 "启用工作流（默认禁用）" 设置为 `true`
5. 配置其他参数并运行

#### 方法二：恢复自动触发（经验用户）

编辑工作流文件，取消注释触发器：

```yaml
# 示例：启用 CI/CD 自动触发
on:
  push: # 取消注释
    branches: [main, develop]
  pull_request: # 取消注释
    branches: [main, develop]
  workflow_dispatch: # 保留手动触发
```

#### 方法三：批量启用脚本

创建脚本批量修改工作流文件：

```bash
#!/bin/bash
# 启用所有工作流自动触发
find .github/workflows -name "*.yml" -exec sed -i 's/# schedule:/schedule:/g' {} \;
find .github/workflows -name "*.yml" -exec sed -i 's/# push:/push:/g' {} \;
```

### 3. 手动触发

所有工作流都支持手动触发（唯一的启用方式）：

1. 访问 GitHub 仓库的 Actions 页面
2. 选择要运行的工作流
3. 点击 "Run workflow" 按钮
4. **必须设置**：`enable_workflow: true`
5. 配置其他参数并运行

**注意**：如果不设置 `enable_workflow: true`，工作流将跳过所有步骤！

## 🔧 配置定制

### 0. 工作流状态管理

#### 检查工作流状态

```bash
# 检查哪些工作流已启用自动触发
grep -l "schedule:" .github/workflows/*.yml
grep -l "push:" .github/workflows/*.yml

# 检查哪些工作流仍处于禁用状态
grep -l "# schedule:" .github/workflows/*.yml
grep -l "# push:" .github/workflows/*.yml
```

#### 启用特定工作流

```bash
# 仅启用 CI/CD 工作流
sed -i 's/# push:/push:/g' .github/workflows/ci-cd.yml
sed -i 's/# pull_request:/pull_request:/g' .github/workflows/ci-cd.yml

# 仅启用监控工作流
sed -i 's/# schedule:/schedule:/g' .github/workflows/monitoring.yml
```

#### 禁用特定工作流

```bash
# 禁用自动更新工作流
sed -i 's/schedule:/# schedule:/g' .github/workflows/auto-update.yml
```

### 1. 调整执行频率

修改`cron`表达式来调整定时任务的执行频率：

```yaml
schedule:
  # 每天凌晨2点执行
  - cron: "0 2 * * *"
  # 每周一上午8点执行
  - cron: "0 8 * * 1"
  # 每5分钟执行一次
  - cron: "*/5 * * * *"
```

### 2. 环境特定配置

根据不同环境调整配置：

```yaml
# 生产环境特殊处理
- name: Production specific steps
  if: github.ref == 'refs/heads/main'
  run: |
    echo "Production environment"

# Staging环境配置
- name: Staging specific steps
  if: github.ref == 'refs/heads/develop'
  run: |
    echo "Staging environment"
```

### 3. 通知配置

配置 Slack 通知：

```yaml
- name: Slack Notification
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    text: "Deployment completed"
  env:
    SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
```

## 📊 监控和报告

### 1. 工作流状态监控

- GitHub Actions 页面查看运行状态
- 通过 Slack 接收实时通知
- 查看生成的报告文件

### 2. 报告类型

- **CI/CD 报告**：构建、测试、部署状态
- **安全报告**：漏洞扫描结果
- **性能报告**：性能测试指标
- **监控报告**：系统健康状态
- **备份报告**：备份执行状态
- **更新报告**：依赖更新状态

### 3. 报告获取

报告通过以下方式获取：

- GitHub Actions Artifacts
- Slack 通知
- 工作流运行日志

## 🛠️ 故障排除

### 1. 常见问题

**问题**：工作流执行失败
**解决**：

1. 检查 Secrets 配置是否正确
2. 验证权限设置
3. 查看具体错误日志

**问题**：数据库连接失败
**解决**：

1. 确认数据库连接信息
2. 检查网络访问权限
3. 验证数据库服务状态

**问题**：Docker 构建失败
**解决**：

1. 检查 Dockerfile 语法
2. 验证基础镜像可用性
3. 确认构建上下文

### 2. 调试技巧

1. **启用调试日志**：

```yaml
- name: Debug step
  run: |
    set -x  # 启用调试模式
    echo "Debug information"
```

2. **添加检查点**：

```yaml
- name: Check environment
  run: |
    echo "Current user: $(whoami)"
    echo "Working directory: $(pwd)"
    echo "Environment variables:"
    env | sort
```

3. **条件执行**：

```yaml
- name: Debug on failure
  if: failure()
  run: |
    echo "Previous step failed"
    # 添加调试命令
```

## 🔐 安全最佳实践

### 1. Secrets 管理

- 使用 GitHub Secrets 存储敏感信息
- 定期轮换访问密钥
- 最小权限原则

### 2. 权限控制

- 限制工作流的执行权限
- 使用 environment protection rules
- 开启 branch protection

### 3. 安全扫描

- 启用所有安全扫描工作流
- 及时处理安全告警
- 定期更新扫描工具

## 📈 性能优化

### 1. 并行执行

- 合理使用`needs`关键字
- 独立任务并行运行
- 避免不必要的依赖

### 2. 缓存策略

- 使用 GitHub Actions 缓存
- 缓存依赖和构建产物
- 定期清理过期缓存

### 3. 资源优化

- 选择适当的 runner 类型
- 优化 Docker 镜像大小
- 减少不必要的步骤

## 🚀 扩展功能

### 1. 添加新的工作流

1. 创建新的`.yml`文件
2. 定义触发条件和任务
3. 测试和部署

### 2. 集成第三方服务

- 添加更多监控工具
- 集成代码质量平台
- 连接项目管理工具

### 3. 自定义 Actions

- 开发可复用的自定义 Actions
- 封装常用操作
- 提高工作流的可维护性

---

## 📝 总结

这套 GitHub Actions 工作流提供了：

✅ **完整的 CI/CD 流程** - 从代码提交到生产部署  
✅ **全面的安全保障** - 多层次安全扫描和漏洞修复  
✅ **自动化运维** - 监控、备份、维护全自动化  
✅ **智能更新管理** - 依赖和安全更新自动化  
✅ **详细的监控报告** - 实时状态监控和报告生成  
✅ **安全的默认配置** - 所有工作流默认禁用，避免意外触发

### 🔒 安全特性

- **默认禁用策略**：防止意外触发和资源浪费
- **手动控制机制**：用户可以精确控制工作流运行
- **灵活启用方式**：支持单次运行和永久启用
- **环境隔离保护**：多环境配置分离管理

### 🚀 使用建议

1. **新项目**：建议逐步启用工作流，先测试再自动化
2. **生产环境**：确保所有配置完整后再启用关键工作流
3. **资源优化**：根据项目需要选择性启用工作流
4. **监控管理**：定期检查工作流运行状态和资源使用

通过这些工作流，Domain MAX 项目可以实现真正的 DevOps 自动化，确保代码质量、系统安全和运维效率，同时保持对自动化流程的完全控制。
