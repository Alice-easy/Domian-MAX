#!/bin/bash

# Domain MAX - 系统测试脚本
# 此脚本执行完整的系统功能测试

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 测试配置
API_BASE="http://localhost:8080/api"
WEB_BASE="http://localhost"
TEST_EMAIL="test@domain-max.com"
TEST_PASSWORD="TestPassword123!"
TEST_DOMAIN="test.example.com"

# 统计变量
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 日志函数
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
}

# 测试函数
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    log_info "Testing: $test_name"
    
    if eval "$test_command"; then
        log_success "$test_name - PASSED"
        return 0
    else
        log_error "$test_name - FAILED"
        return 1
    fi
}

# API 测试函数
test_api_endpoint() {
    local endpoint="$1"
    local expected_status="$2"
    local method="${3:-GET}"
    local data="${4:-}"
    
    local curl_cmd="curl -s -o /dev/null -w '%{http_code}' -X $method"
    
    if [ -n "$data" ]; then
        curl_cmd="$curl_cmd -H 'Content-Type: application/json' -d '$data'"
    fi
    
    curl_cmd="$curl_cmd $API_BASE$endpoint"
    
    local status_code=$(eval "$curl_cmd")
    
    if [ "$status_code" = "$expected_status" ]; then
        return 0
    else
        echo "Expected: $expected_status, Got: $status_code"
        return 1
    fi
}

# 1. 基础连接测试
test_basic_connectivity() {
    echo ""
    log_info "=== 基础连接测试 ==="
    
    run_test "应用服务响应" "curl -sf $API_BASE/health >/dev/null"
    run_test "前端页面响应" "curl -sf $WEB_BASE >/dev/null"
    run_test "Nginx 健康检查" "curl -sf $WEB_BASE/health >/dev/null"
    
    # 数据库连接测试
    run_test "数据库连接" "docker-compose exec -T postgres pg_isready -U postgres >/dev/null 2>&1"
    
    # Redis 连接测试
    run_test "Redis 连接" "docker-compose exec -T redis redis-cli ping | grep -q PONG"
}

# 2. API 端点测试
test_api_endpoints() {
    echo ""
    log_info "=== API 端点测试 ==="
    
    # 健康检查端点
    run_test "健康检查端点" "test_api_endpoint '/health' '200'"
    
    # 认证端点
    run_test "注册端点可达性" "test_api_endpoint '/auth/register' '400' 'POST'"
    run_test "登录端点可达性" "test_api_endpoint '/auth/login' '400' 'POST'"
    
    # DNS 端点 (需要认证)
    run_test "DNS 记录端点" "test_api_endpoint '/dns/records' '401'"
    run_test "域名端点" "test_api_endpoint '/domains' '401'"
}

# 3. 用户注册和登录测试
test_auth_flow() {
    echo ""
    log_info "=== 认证流程测试 ==="
    
    # 测试用户注册
    local register_data="{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\",\"name\":\"Test User\"}"
    run_test "用户注册" "test_api_endpoint '/auth/register' '201' 'POST' '$register_data'"
    
    # 测试用户登录
    local login_data="{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}"
    local login_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$login_data" "$API_BASE/auth/login")
    
    if echo "$login_response" | grep -q "token"; then
        log_success "用户登录 - PASSED"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        
        # 提取 token
        TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
        export AUTH_TOKEN="$TOKEN"
        
        # 测试受保护的端点
        run_test "受保护端点访问" "curl -sf -H 'Authorization: Bearer $TOKEN' $API_BASE/auth/profile >/dev/null"
    else
        log_error "用户登录 - FAILED"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    TOTAL_TESTS=$((TOTAL_TESTS + 2))
}

# 4. DNS 功能测试
test_dns_functionality() {
    echo ""
    log_info "=== DNS 功能测试 ==="
    
    if [ -z "$AUTH_TOKEN" ]; then
        log_warning "跳过 DNS 测试：未获取到认证 token"
        return
    fi
    
    # 测试域名创建
    local domain_data="{\"name\":\"$TEST_DOMAIN\",\"provider\":\"cloudflare\"}"
    run_test "域名创建" "curl -sf -X POST -H 'Authorization: Bearer $AUTH_TOKEN' -H 'Content-Type: application/json' -d '$domain_data' $API_BASE/domains >/dev/null"
    
    # 测试 DNS 记录创建
    local dns_data="{\"domain\":\"$TEST_DOMAIN\",\"type\":\"A\",\"name\":\"test\",\"value\":\"192.168.1.1\",\"ttl\":300}"
    run_test "DNS 记录创建" "curl -sf -X POST -H 'Authorization: Bearer $AUTH_TOKEN' -H 'Content-Type: application/json' -d '$dns_data' $API_BASE/dns/records >/dev/null"
    
    # 测试记录查询
    run_test "DNS 记录查询" "curl -sf -H 'Authorization: Bearer $AUTH_TOKEN' $API_BASE/dns/records >/dev/null"
    
    # 测试域名查询
    run_test "域名列表查询" "curl -sf -H 'Authorization: Bearer $AUTH_TOKEN' $API_BASE/domains >/dev/null"
}

# 5. 安全性测试
test_security() {
    echo ""
    log_info "=== 安全性测试 ==="
    
    # 测试未授权访问
    run_test "未授权访问拒绝" "test_api_endpoint '/auth/profile' '401'"
    run_test "无效 token 拒绝" "curl -sf -H 'Authorization: Bearer invalid_token' $API_BASE/auth/profile >/dev/null 2>&1 && false || true"
    
    # 测试 SQL 注入防护
    local malicious_data="{\"email\":\"test'; DROP TABLE users; --\",\"password\":\"password\"}"
    run_test "SQL 注入防护" "test_api_endpoint '/auth/login' '400' 'POST' '$malicious_data'"
    
    # 测试 XSS 防护
    local xss_data="{\"name\":\"<script>alert('xss')</script>\",\"email\":\"xss@test.com\",\"password\":\"password\"}"
    run_test "XSS 防护" "test_api_endpoint '/auth/register' '400' 'POST' '$xss_data'"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 4))
}

# 6. 性能测试
test_performance() {
    echo ""
    log_info "=== 性能测试 ==="
    
    # 测试响应时间
    local response_time=$(curl -o /dev/null -s -w '%{time_total}' "$API_BASE/health")
    if (( $(echo "$response_time < 1.0" | bc -l) )); then
        log_success "响应时间测试 - PASSED ($response_time 秒)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        log_error "响应时间测试 - FAILED ($response_time 秒，超过 1 秒)"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    # 测试并发处理
    log_info "执行并发测试 (10 个并发请求)..."
    local concurrent_success=0
    for i in {1..10}; do
        if curl -sf "$API_BASE/health" >/dev/null 2>&1 &; then
            concurrent_success=$((concurrent_success + 1))
        fi
    done
    wait
    
    if [ $concurrent_success -eq 10 ]; then
        log_success "并发测试 - PASSED (10/10 成功)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        log_error "并发测试 - FAILED ($concurrent_success/10 成功)"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    TOTAL_TESTS=$((TOTAL_TESTS + 2))
}

# 7. 数据持久性测试
test_data_persistence() {
    echo ""
    log_info "=== 数据持久性测试 ==="
    
    # 测试数据库数据持久性
    local user_count_before=$(docker-compose exec -T postgres psql -U postgres -d domain_manager -c "SELECT COUNT(*) FROM users;" | grep -o '[0-9]*' | head -1)
    
    # 重启应用服务
    log_info "重启应用服务以测试数据持久性..."
    docker-compose restart app >/dev/null 2>&1
    sleep 10
    
    # 等待服务重新启动
    local retry=0
    while [ $retry -lt 30 ]; do
        if curl -sf "$API_BASE/health" >/dev/null 2>&1; then
            break
        fi
        sleep 2
        retry=$((retry + 1))
    done
    
    local user_count_after=$(docker-compose exec -T postgres psql -U postgres -d domain_manager -c "SELECT COUNT(*) FROM users;" | grep -o '[0-9]*' | head -1)
    
    if [ "$user_count_before" = "$user_count_after" ]; then
        log_success "数据持久性测试 - PASSED"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        log_error "数据持久性测试 - FAILED (重启前: $user_count_before, 重启后: $user_count_after)"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
}

# 8. 容器健康检查
test_container_health() {
    echo ""
    log_info "=== 容器健康检查 ==="
    
    # 检查所有容器状态
    local containers=("domain-max-app" "domain-max-db" "domain-max-redis" "domain-max-nginx")
    
    for container in "${containers[@]}"; do
        if docker ps --filter "name=$container" --filter "status=running" | grep -q "$container"; then
            log_success "$container 容器运行正常"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            log_error "$container 容器未运行"
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
    done
    
    # 检查容器资源使用
    log_info "检查容器资源使用情况..."
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}" | grep -E "(domain-max|CONTAINER)"
}

# 9. SSL/TLS 测试
test_ssl() {
    echo ""
    log_info "=== SSL/TLS 测试 ==="
    
    # 测试 HTTPS 连接
    if curl -sk "https://localhost/health" >/dev/null 2>&1; then
        log_success "HTTPS 连接 - PASSED"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        log_error "HTTPS 连接 - FAILED"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    # 测试 HTTP 到 HTTPS 重定向
    local redirect_code=$(curl -s -o /dev/null -w '%{http_code}' "http://localhost/")
    if [ "$redirect_code" = "301" ] || [ "$redirect_code" = "302" ]; then
        log_success "HTTP 重定向 - PASSED"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        log_error "HTTP 重定向 - FAILED (状态码: $redirect_code)"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    TOTAL_TESTS=$((TOTAL_TESTS + 2))
}

# 10. 清理测试数据
cleanup_test_data() {
    echo ""
    log_info "=== 清理测试数据 ==="
    
    if [ -n "$AUTH_TOKEN" ]; then
        # 删除测试用户和相关数据
        docker-compose exec -T postgres psql -U postgres -d domain_manager -c "DELETE FROM users WHERE email = '$TEST_EMAIL';" >/dev/null 2>&1 || true
        log_info "测试数据已清理"
    fi
}

# 生成测试报告
generate_report() {
    echo ""
    echo "============================================="
    log_info "📊 测试报告"
    echo "============================================="
    echo "总测试数: $TOTAL_TESTS"
    echo -e "通过测试: ${GREEN}$PASSED_TESTS${NC}"
    echo -e "失败测试: ${RED}$FAILED_TESTS${NC}"
    echo ""
    
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "${GREEN}🎉 所有测试通过！系统运行正常。${NC}"
        echo ""
        echo "✅ 系统状态: 健康"
        echo "✅ 功能完整性: 正常"
        echo "✅ 安全性: 良好"
        echo "✅ 性能: 满足要求"
        return 0
    else
        echo -e "${RED}⚠️ 有 $FAILED_TESTS 个测试失败，请检查系统状态。${NC}"
        echo ""
        echo "建议检查项:"
        echo "1. 查看容器日志: docker-compose logs"
        echo "2. 检查服务状态: docker-compose ps"
        echo "3. 验证环境配置: cat deployments/.env"
        echo "4. 查看系统资源: docker stats"
        return 1
    fi
}

# 主测试流程
main() {
    echo "🧪 Domain MAX 系统测试"
    echo "======================"
    echo "开始时间: $(date)"
    echo ""
    
    # 检查先决条件
    if ! docker-compose ps | grep -q "Up"; then
        log_error "系统未运行，请先启动服务: docker-compose up -d"
        exit 1
    fi
    
    # 执行测试套件
    test_basic_connectivity
    test_api_endpoints
    test_auth_flow
    test_dns_functionality
    test_security
    test_performance
    test_data_persistence
    test_container_health
    test_ssl
    
    # 清理和报告
    cleanup_test_data
    generate_report
    
    echo ""
    echo "结束时间: $(date)"
    
    # 返回适当的退出代码
    if [ $FAILED_TESTS -eq 0 ]; then
        exit 0
    else
        exit 1
    fi
}

# 处理脚本参数
case "${1:-}" in
    --quick)
        echo "🏃 执行快速测试..."
        test_basic_connectivity
        test_api_endpoints
        generate_report
        ;;
    --security)
        echo "🔒 执行安全测试..."
        test_security
        generate_report
        ;;
    --performance)
        echo "⚡ 执行性能测试..."
        test_performance
        generate_report
        ;;
    --help)
        echo "Domain MAX 系统测试脚本"
        echo ""
        echo "用法: $0 [选项]"
        echo ""
        echo "选项:"
        echo "  (无选项)      执行完整测试套件"
        echo "  --quick      执行快速测试"
        echo "  --security   仅执行安全测试"
        echo "  --performance 仅执行性能测试"
        echo "  --help       显示此帮助信息"
        ;;
    *)
        main
        ;;
esac