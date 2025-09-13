#!/bin/bash

# Domain MAX 开发环境启动脚本
# 前后端分离开发模式（本地开发 + Cloudflare Pages预览）

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Domain MAX 开发环境启动脚本 ===${NC}"
echo -e "${CYAN}架构：前后端分离开发模式${NC}"
echo

# 检查依赖
check_dependencies() {
    echo -e "${YELLOW}🔍 检查开发依赖...${NC}"
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go 未安装，请先安装 Go 1.23+${NC}"
        exit 1
    fi
    
    # 检查Node.js
    if ! command -v node &> /dev/null; then
        echo -e "${RED}❌ Node.js 未安装，请先安装 Node.js 18+${NC}"
        exit 1
    fi
    
    # 检查npm
    if ! command -v npm &> /dev/null; then
        echo -e "${RED}❌ npm 未安装${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ 依赖检查完成${NC}"
    echo -e "   Go: $(go version | cut -d' ' -f3)"
    echo -e "   Node.js: $(node --version)"
    echo -e "   npm: $(npm --version)"
}

# 检查环境配置
check_env() {
    echo -e "${YELLOW}🔧 检查环境配置...${NC}"
    
    # 检查后端环境配置
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            echo -e "${YELLOW}⚠️  后端未找到 .env 文件，复制示例配置...${NC}"
            cp .env.example .env
            echo -e "${CYAN}📝 请编辑 .env 文件设置开发环境配置${NC}"
        else
            echo -e "${RED}❌ 未找到 .env.example 文件${NC}"
            exit 1
        fi
    fi
    
    # 检查前端环境配置
    if [ ! -f "web/.env.development" ]; then
        echo -e "${YELLOW}⚠️  前端未找到开发环境配置，创建默认配置...${NC}"
        cat > web/.env.development << EOF
# 开发环境配置
VITE_API_BASE_URL=http://localhost:8080
VITE_BACKEND_DOMAIN=localhost:8080
NODE_ENV=development
EOF
    fi
    
    echo -e "${GREEN}✅ 环境配置检查完成${NC}"
}

# 安装前端依赖
install_frontend_deps() {
    echo -e "${YELLOW}📦 安装前端依赖...${NC}"
    
    cd web
    if [ ! -d "node_modules" ]; then
        echo -e "${CYAN}正在安装前端依赖，这可能需要几分钟...${NC}"
        npm install
    else
        echo -e "${CYAN}检查依赖更新...${NC}"
        npm ci
    fi
    cd ..
    
    echo -e "${GREEN}✅ 前端依赖安装完成${NC}"
}

# 安装后端依赖
install_backend_deps() {
    echo -e "${YELLOW}📦 安装后端依赖...${NC}"
    
    echo -e "${CYAN}正在下载Go模块...${NC}"
    go mod tidy
    
    echo -e "${GREEN}✅ 后端依赖安装完成${NC}"
}

# 启动后端开发服务器
start_backend() {
    echo -e "${YELLOW}🚀 启动后端API服务器...${NC}"
    
    # 设置开发环境变量
    export APP_MODE=development
    export PORT=8080
    export CORS_ALLOWED_ORIGINS="http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:5173"
    
    echo -e "${CYAN}后端服务将运行在: http://localhost:8080${NC}"
    echo -e "${CYAN}API健康检查: http://localhost:8080/health${NC}"
    echo
    
    # 运行后端服务
    go run ./cmd/api-server
}

# 启动前端开发服务器
start_frontend() {
    echo -e "${YELLOW}🌐 启动前端开发服务器...${NC}"
    
    cd web
    echo -e "${CYAN}前端服务将运行在: http://localhost:5173${NC}"
    echo
    
    # 运行前端服务
    npm run dev
}

# 并行启动前后端
start_both() {
    echo -e "${YELLOW}🚀 并行启动前后端服务...${NC}"
    echo
    
    # 检查是否安装了并行工具
    if command -v tmux &> /dev/null; then
        start_with_tmux
    elif command -v screen &> /dev/null; then
        start_with_screen
    else
        start_manual
    fi
}

# 使用tmux启动
start_with_tmux() {
    echo -e "${CYAN}使用tmux启动服务...${NC}"
    
    # 创建新的tmux会话
    tmux new-session -d -s domain-max
    
    # 分割窗口
    tmux split-window -h
    
    # 在左侧窗口启动后端
    tmux send-keys -t domain-max:0.0 "export APP_MODE=development PORT=8080 CORS_ALLOWED_ORIGINS='http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:5173' && go run ./cmd/api-server" Enter
    
    # 在右侧窗口启动前端
    tmux send-keys -t domain-max:0.1 "cd web && npm run dev" Enter
    
    # 连接到会话
    echo -e "${GREEN}✅ 服务已启动！${NC}"
    echo -e "${CYAN}后端: http://localhost:8080${NC}"
    echo -e "${CYAN}前端: http://localhost:5173${NC}"
    echo
    echo -e "${YELLOW}连接到tmux会话查看日志: tmux attach -t domain-max${NC}"
    echo -e "${YELLOW}退出tmux: Ctrl+B 然后按 D${NC}"
    echo -e "${YELLOW}停止服务: tmux kill-session -t domain-max${NC}"
    
    tmux attach -t domain-max
}

# 使用screen启动
start_with_screen() {
    echo -e "${CYAN}使用screen启动服务...${NC}"
    
    # 启动后端
    screen -dmS domain-max-api bash -c "export APP_MODE=development PORT=8080 && go run ./cmd/api-server"
    
    # 启动前端
    screen -dmS domain-max-web bash -c "cd web && npm run dev"
    
    echo -e "${GREEN}✅ 服务已在后台启动！${NC}"
    echo -e "${CYAN}后端: http://localhost:8080${NC}"
    echo -e "${CYAN}前端: http://localhost:5173${NC}"
    echo
    echo -e "${YELLOW}查看后端日志: screen -r domain-max-api${NC}"
    echo -e "${YELLOW}查看前端日志: screen -r domain-max-web${NC}"
    echo -e "${YELLOW}停止服务:${NC}"
    echo -e "  screen -S domain-max-api -X quit"
    echo -e "  screen -S domain-max-web -X quit"
}

# 手动启动指导
start_manual() {
    echo -e "${YELLOW}⚠️  未检测到tmux或screen，需要手动启动${NC}"
    echo
    echo -e "${CYAN}请打开两个终端窗口：${NC}"
    echo
    echo -e "${WHITE}终端1 (后端API):${NC}"
    echo -e "  cd $(pwd)"
    echo -e "  export APP_MODE=development PORT=8080"
    echo -e "  export CORS_ALLOWED_ORIGINS='http://localhost:3000,http://localhost:5173'"
    echo -e "  go run ./cmd/api-server"
    echo
    echo -e "${WHITE}终端2 (前端Web):${NC}"
    echo -e "  cd $(pwd)/web"
    echo -e "  npm run dev"
    echo
    echo -e "${GREEN}启动后访问:${NC}"
    echo -e "${CYAN}  后端API: http://localhost:8080${NC}"
    echo -e "${CYAN}  前端Web: http://localhost:5173${NC}"
}

# 停止服务
stop_services() {
    echo -e "${YELLOW}🛑 停止开发服务...${NC}"
    
    # 停止tmux会话
    if tmux has-session -t domain-max 2>/dev/null; then
        tmux kill-session -t domain-max
        echo -e "${GREEN}✅ tmux会话已停止${NC}"
    fi
    
    # 停止screen会话
    screen -S domain-max-api -X quit 2>/dev/null || true
    screen -S domain-max-web -X quit 2>/dev/null || true
    
    # 停止可能的进程
    pkill -f "go run ./cmd/api-server" 2>/dev/null || true
    pkill -f "npm run dev" 2>/dev/null || true
    
    echo -e "${GREEN}✅ 开发服务已停止${NC}"
}

# 显示状态
show_status() {
    echo -e "${YELLOW}📊 开发服务状态:${NC}"
    echo
    
    # 检查后端端口
    if lsof -i :8080 &>/dev/null; then
        echo -e "${GREEN}✅ 后端API服务: http://localhost:8080${NC}"
    else
        echo -e "${RED}❌ 后端API服务未运行${NC}"
    fi
    
    # 检查前端端口
    if lsof -i :5173 &>/dev/null; then
        echo -e "${GREEN}✅ 前端Web服务: http://localhost:5173${NC}"
    else
        echo -e "${RED}❌ 前端Web服务未运行${NC}"
    fi
    
    # 检查tmux会话
    if tmux has-session -t domain-max 2>/dev/null; then
        echo -e "${CYAN}📺 tmux会话运行中: domain-max${NC}"
    fi
}

# 显示帮助
show_help() {
    echo -e "${WHITE}Domain MAX 开发环境管理工具${NC}"
    echo -e "${CYAN}适用于前后端分离开发模式${NC}"
    echo
    echo -e "${YELLOW}用法: $0 [命令]${NC}"
    echo
    echo -e "${WHITE}命令:${NC}"
    echo -e "  ${GREEN}start${NC}      - 启动完整开发环境（默认）"
    echo -e "  ${GREEN}backend${NC}    - 仅启动后端API服务"
    echo -e "  ${GREEN}frontend${NC}   - 仅启动前端Web服务"
    echo -e "  ${GREEN}install${NC}    - 安装所有依赖"
    echo -e "  ${GREEN}stop${NC}       - 停止所有开发服务"
    echo -e "  ${GREEN}status${NC}     - 显示服务状态"
    echo -e "  ${GREEN}help${NC}       - 显示帮助信息"
    echo
    echo -e "${WHITE}示例:${NC}"
    echo -e "  $0              # 启动完整开发环境"
    echo -e "  $0 start        # 启动完整开发环境"
    echo -e "  $0 backend      # 仅启动后端"
    echo -e "  $0 frontend     # 仅启动前端"
    echo -e "  $0 install      # 安装依赖"
    echo -e "  $0 stop         # 停止服务"
    echo
    echo -e "${WHITE}开发流程:${NC}"
    echo -e "  1. 首次运行: $0 install"
    echo -e "  2. 启动开发: $0 start"
    echo -e "  3. 前端访问: http://localhost:5173"
    echo -e "  4. API文档: http://localhost:8080"
}

# 主函数
main() {
    case "${1:-start}" in
        "start")
            check_dependencies
            check_env
            install_backend_deps
            install_frontend_deps
            start_both
            ;;
        "backend"|"api")
            check_dependencies
            check_env
            install_backend_deps
            start_backend
            ;;
        "frontend"|"web")
            check_dependencies
            check_env
            install_frontend_deps
            start_frontend
            ;;
        "install"|"setup")
            check_dependencies
            check_env
            install_backend_deps
            install_frontend_deps
            echo -e "${GREEN}🎉 所有依赖安装完成！${NC}"
            echo -e "${CYAN}运行 $0 start 启动开发环境${NC}"
            ;;
        "stop")
            stop_services
            ;;
        "status")
            show_status
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            echo -e "${RED}❌ 未知命令: $1${NC}"
            echo
            show_help
            exit 1
            ;;
    esac
}

main "$@"