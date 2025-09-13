#!/bin/bash

# Domain MAX API 构建脚本
# 适用于前后端分离架构（Cloudflare + VPS）

set -e

echo "=== Domain MAX API 构建脚本 ==="
echo "架构：前端 Cloudflare Pages + 后端 VPS API"
echo

# 检查Go环境
check_go() {
    echo "🔍 检查Go环境..."
    if ! command -v go &> /dev/null; then
        echo "❌ Go 未安装，请先安装 Go 1.23+"
        exit 1
    fi
    
    GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
    echo "✅ Go环境检查完成: $(go version)"
}

# 验证环境配置
check_env() {
    echo "🔧 检查环境配置..."
    
    if [ ! -f ".env" ]; then
        if [ ! -f ".env.example" ]; then
            echo "❌ 未找到 .env.example 文件"
            exit 1
        fi
        echo "⚠️  未找到 .env 文件，请复制 .env.example 并配置"
        echo "   cp .env.example .env"
        echo "   然后编辑 .env 文件设置您的配置"
        exit 1
    fi
    
    echo "✅ 环境配置检查完成"
}

# 构建API服务
build_api() {
    echo "🏗️  构建API服务..."
    
    # 下载Go依赖
    echo "📦 下载Go依赖..."
    go mod tidy
    
    # 构建API服务
    echo "🔨 构建API服务..."
    CGO_ENABLED=0 go build -ldflags="-w -s" -o domain-max-api ./cmd/api-server
    
    # 检查构建结果
    if [ ! -f "domain-max-api" ]; then
        echo "❌ 构建失败，未生成可执行文件"
        exit 1
    fi
    
    echo "✅ API构建完成"
}

# 构建不同平台版本
build_cross() {
    echo "🌐 构建跨平台版本..."
    
    # Linux版本（常用于VPS部署）
    echo "📦 构建Linux版本..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o domain-max-api-linux ./cmd/api-server
    
    # Windows版本
    echo "📦 构建Windows版本..."
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o domain-max-api.exe ./cmd/api-server
    
    echo "✅ 跨平台构建完成"
}

# 清理构建产物
clean() {
    echo "🧹 清理构建产物..."
    rm -f domain-max-api
    rm -f domain-max-api.exe
    rm -f domain-max-api-linux
    echo "✅ 清理完成"
}

# 测试API服务
test_api() {
    echo "🧪 测试API服务..."
    
    if [ ! -f "domain-max-api" ]; then
        echo "❌ 未找到可执行文件，请先构建"
        exit 1
    fi
    
    # 运行测试
    go test ./... -v
    
    echo "✅ 测试完成"
}

# 显示帮助
help() {
    echo "Domain MAX API 构建工具"
    echo "适用于前后端分离架构（Cloudflare Pages + VPS API）"
    echo
    echo "用法: $0 [选项]"
    echo
    echo "选项:"
    echo "  build       - 构建API服务 (默认)"
    echo "  cross       - 构建跨平台版本"
    echo "  clean       - 清理构建产物"
    echo "  test        - 运行测试"
    echo "  dev         - 开发模式运行"
    echo "  help        - 显示帮助信息"
    echo
    echo "示例:"
    echo "  $0            # 构建API服务"
    echo "  $0 build      # 构建API服务"
    echo "  $0 cross      # 构建Linux和Windows版本"
    echo "  $0 dev        # 开发模式运行"
    echo "  $0 clean      # 清理构建产物"
    echo
    echo "部署流程:"
    echo "  1. 本地构建: $0 cross"
    echo "  2. 上传到VPS: scp domain-max-api-linux user@vps:/path/"
    echo "  3. VPS运行: ./domain-max-api-linux"
}

# 开发模式运行
dev_run() {
    echo "🚀 开发模式启动API服务..."
    echo "前端请在另一个终端运行: cd web && npm run dev"
    echo "前端地址: http://localhost:5173"
    echo "API地址: http://localhost:8080"
    echo
    
    check_env
    
    # 设置开发环境
    export APP_MODE=development
    export PORT=8080
    
    # 直接运行Go程序
    go run ./cmd/api-server
}

# 主函数
main() {
    case "${1:-build}" in
        "build")
            check_go
            check_env
            build_api
            echo
            echo "🎉 构建完成！"
            echo "📁 可执行文件：./domain-max-api"
            echo "🚀 运行命令：./domain-max-api"
            echo "📖 部署指南：docs/separation-deployment.md"
            ;;
        "cross")
            check_go
            check_env
            go mod tidy
            build_cross
            echo
            echo "🎉 跨平台构建完成！"
            echo "📁 Linux版本：./domain-max-api-linux"
            echo "📁 Windows版本：./domain-max-api.exe"
            ;;
        "clean")
            clean
            ;;
        "test")
            check_go
            test_api
            ;;
        "dev")
            check_go
            dev_run
            ;;
        "help"|"-h"|"--help")
            help
            ;;
        *)
            echo "❌ 未知选项: $1"
            echo
            help
            exit 1
            ;;
    esac
}

main "$@"