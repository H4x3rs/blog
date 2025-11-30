#!/bin/bash

# =================================================================================
# 博客系统完整部署脚本
# 功能：交互式配置 -> 生成.env -> 部署数据库 -> 编译前后端
# =================================================================================

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# 配置文件路径
ENV_FILE="$PROJECT_ROOT/.env"
ENV_EXAMPLE="$PROJECT_ROOT/.env.example"
SQL_TEMPLATE="$PROJECT_ROOT/backend/manifest/sql/deploy.sql.template"
SQL_OUTPUT="$PROJECT_ROOT/backend/manifest/sql/deploy.sql"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/fronted"

# 打印带颜色的消息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_question() {
    echo -e "${BLUE}[?]${NC} $1"
}

# 读取用户输入（带默认值）
read_input() {
    local prompt="$1"
    local default="$2"
    local var_name="$3"
    local is_password="$4"
    
    if [ -n "$default" ]; then
        prompt="$prompt [$default]: "
    else
        prompt="$prompt: "
    fi
    
    print_question "$prompt"
    
    if [ "$is_password" = "true" ]; then
        read -s input
        echo ""
    else
        read input
    fi
    
    if [ -z "$input" ] && [ -n "$default" ]; then
        input="$default"
    fi
    
    eval "$var_name='$input'"
}

# 检查必需的命令
check_requirements() {
    print_info "检查必需的工具..."
    
    local missing_tools=()
    
    if ! command -v mysql &> /dev/null; then
        missing_tools+=("mysql")
    fi
    
    if ! command -v envsubst &> /dev/null; then
        missing_tools+=("envsubst (gettext)")
    fi
    
    if ! command -v go &> /dev/null; then
        missing_tools+=("go")
    fi
    
    if ! command -v npm &> /dev/null && ! command -v yarn &> /dev/null; then
        missing_tools+=("npm 或 yarn")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        print_error "缺少必需的工具:"
        for tool in "${missing_tools[@]}"; do
            print_error "  - $tool"
        done
        print_info "安装说明："
        print_info "  MySQL: https://dev.mysql.com/downloads/"
        print_info "  gettext: sudo apt-get install gettext-base (Ubuntu/Debian)"
        print_info "  Go: https://golang.org/dl/"
        print_info "  Node.js: https://nodejs.org/"
        exit 1
    fi
    
    print_info "工具检查完成"
}

# 交互式配置
interactive_config() {
    echo ""
    print_info "=========================================="
    print_info "开始交互式配置"
    print_info "=========================================="
    echo ""
    
    # 数据库配置
    print_info "--- 数据库配置 ---"
    read_input "数据库名称" "go_blog" DB_NAME
    read_input "数据库用户名" "root" DB_USER
    read_input "数据库密码" "" DB_PASSWORD true
    read_input "数据库主机" "127.0.0.1" DB_HOST
    read_input "数据库端口" "3306" DB_PORT
    echo ""
    
    # 管理员账户配置
    print_info "--- 管理员账户配置 ---"
    read_input "管理员用户名" "admin" ADMIN_USERNAME
    read_input "管理员密码" "" ADMIN_PASSWORD true
    read_input "管理员昵称" "管理员" ADMIN_NICKNAME
    read_input "管理员邮箱" "admin@example.com" ADMIN_EMAIL
    echo ""
    
    # 网站配置
    print_info "--- 网站配置 ---"
    read_input "网站名称" "Blog System" SITE_NAME
    read_input "网站标题" "Blog System" SITE_TITLE
    read_input "网站副标题" "分享编程心得，记录技术成长" SITE_SUBTITLE
    read_input "网站描述" "分享编程心得，记录技术成长。分享 Go, Vue, 云原生等前沿技术心得。" SITE_DESC
    read_input "关键词" "Go, Vue, Cloud Native" SITE_KEYWORDS
    read_input "版权信息" "© 2024 Blog System" SITE_COPYRIGHT
    echo ""
    
    # 生产环境配置
    print_info "--- 生产环境配置 ---"
    read_input "前端域名（如：https://blog.example.com）" "" FRONTEND_DOMAIN
    read_input "后端API域名（如：https://api.example.com）" "" BACKEND_DOMAIN
    read_input "部署环境 (development/production)" "development" DEPLOY_ENV
    echo ""
    
    # JWT配置
    print_info "--- JWT配置 ---"
    if [ -z "$JWT_SECRET" ]; then
        JWT_SECRET=$(openssl rand -hex 32 2>/dev/null || head -c 32 /dev/urandom | base64 | tr -d '\n' | cut -c1-32)
    fi
    read_input "JWT密钥（留空自动生成）" "$JWT_SECRET" JWT_SECRET
    echo ""
    
    # OAuth配置（可选）
    print_info "--- OAuth配置（可选，直接回车跳过）---"
    read_input "GitHub Client ID" "" GITHUB_CLIENT_ID
    read_input "GitHub Client Secret" "" GITHUB_CLIENT_SECRET
    read_input "Google Client ID" "" GOOGLE_CLIENT_ID
    read_input "Google Client Secret" "" GOOGLE_CLIENT_SECRET
    read_input "微信 AppID" "" WECHAT_APP_ID
    read_input "微信 AppSecret" "" WECHAT_APP_SECRET
    echo ""
    
    # Redis配置（可选）
    print_info "--- Redis配置（可选）---"
    read_input "Redis主机" "127.0.0.1" REDIS_HOST
    read_input "Redis端口" "6379" REDIS_PORT
    read_input "Redis数据库" "1" REDIS_DB
    echo ""
    
    # 设置OAuth前端URL
    if [ -n "$FRONTEND_DOMAIN" ]; then
        OAUTH_FRONTEND_URL="$FRONTEND_DOMAIN"
    else
        OAUTH_FRONTEND_URL="http://localhost:5173"
    fi
}

# 生成.env文件
generate_env_file() {
    print_info "生成 .env 文件..."
    
    cat > "$ENV_FILE" <<EOF
# =================================================================================
# 博客系统环境变量配置
# 自动生成时间: $(date '+%Y-%m-%d %H:%M:%S')
# =================================================================================

# ============================================
# 数据库配置
# ============================================
DB_NAME=$DB_NAME
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT

# ============================================
# 管理员账户配置
# ============================================
ADMIN_USERNAME=$ADMIN_USERNAME
ADMIN_PASSWORD=$ADMIN_PASSWORD
ADMIN_NICKNAME=$ADMIN_NICKNAME
ADMIN_EMAIL=$ADMIN_EMAIL

# ============================================
# 网站配置
# ============================================
SITE_NAME=$SITE_NAME
SITE_TITLE=$SITE_TITLE
SITE_SUBTITLE=$SITE_SUBTITLE
SITE_DESC=$SITE_DESC
SITE_KEYWORDS=$SITE_KEYWORDS
SITE_COPYRIGHT=$SITE_COPYRIGHT

# ============================================
# 生产环境配置
# ============================================
FRONTEND_DOMAIN=$FRONTEND_DOMAIN
BACKEND_DOMAIN=$BACKEND_DOMAIN
DEPLOY_ENV=$DEPLOY_ENV

# ============================================
# JWT配置
# ============================================
JWT_SECRET=$JWT_SECRET

# ============================================
# OAuth配置（可选）
# ============================================
OAUTH_FRONTEND_URL=$OAUTH_FRONTEND_URL
GITHUB_CLIENT_ID=$GITHUB_CLIENT_ID
GITHUB_CLIENT_SECRET=$GITHUB_CLIENT_SECRET
GOOGLE_CLIENT_ID=$GOOGLE_CLIENT_ID
GOOGLE_CLIENT_SECRET=$GOOGLE_CLIENT_SECRET
WECHAT_APP_ID=$WECHAT_APP_ID
WECHAT_APP_SECRET=$WECHAT_APP_SECRET

# ============================================
# Redis配置（可选）
# ============================================
REDIS_HOST=$REDIS_HOST
REDIS_PORT=$REDIS_PORT
REDIS_DB=$REDIS_DB
EOF
    
    print_info ".env 文件已生成: $ENV_FILE"
}

# 加载环境变量
load_env() {
    if [ ! -f "$ENV_FILE" ]; then
        print_error ".env 文件不存在: $ENV_FILE"
        exit 1
    fi
    
    set -a
    source "$ENV_FILE"
    set +a
    
    # 设置默认值
    DB_HOST=${DB_HOST:-127.0.0.1}
    DB_PORT=${DB_PORT:-3306}
    ADMIN_NICKNAME=${ADMIN_NICKNAME:-管理员}
    ADMIN_EMAIL=${ADMIN_EMAIL:-admin@example.com}
    SITE_NAME=${SITE_NAME:-Blog System}
    SITE_TITLE=${SITE_TITLE:-Blog System}
    SITE_SUBTITLE=${SITE_SUBTITLE:-分享编程心得，记录技术成长}
    SITE_DESC=${SITE_DESC:-分享编程心得，记录技术成长。分享 Go, Vue, 云原生等前沿技术心得。}
    SITE_KEYWORDS=${SITE_KEYWORDS:-Go, Vue, Cloud Native}
    SITE_COPYRIGHT=${SITE_COPYRIGHT:-© 2024 Blog System}
    DEPLOY_ENV=${DEPLOY_ENV:-development}
    OAUTH_FRONTEND_URL=${OAUTH_FRONTEND_URL:-http://localhost:5173}
    REDIS_HOST=${REDIS_HOST:-127.0.0.1}
    REDIS_PORT=${REDIS_PORT:-6379}
    REDIS_DB=${REDIS_DB:-1}
}

# 生成SQL文件
generate_sql() {
    print_info "生成部署SQL文件..."
    
    if [ ! -f "$SQL_TEMPLATE" ]; then
        print_error "SQL模板文件不存在: $SQL_TEMPLATE"
        exit 1
    fi
    
    envsubst < "$SQL_TEMPLATE" > "$SQL_OUTPUT"
    print_info "SQL文件已生成: $SQL_OUTPUT"
}

# 执行数据库部署
deploy_database() {
    print_info "开始部署数据库..."
    
    # 测试数据库连接
    print_info "测试数据库连接..."
    if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" &> /dev/null; then
        print_error "数据库连接失败，请检查配置"
        exit 1
    fi
    print_info "数据库连接成功"
    
    # 执行SQL脚本
    print_info "执行SQL脚本..."
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < "$SQL_OUTPUT"; then
        print_info "数据库部署成功！"
    else
        print_error "数据库部署失败"
        exit 1
    fi
    
    # 验证部署
    print_info "验证部署..."
    table_count=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -D"$DB_NAME" -se "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='$DB_NAME'" 2>/dev/null)
    
    if [ "$table_count" -gt 0 ]; then
        print_info "数据库表创建成功，共 $table_count 个表"
    else
        print_warn "未检测到数据库表，请检查SQL执行结果"
    fi
}

# 编译后端
build_backend() {
    print_info "开始编译后端..."
    
    cd "$BACKEND_DIR"
    
    # 检查go.mod是否存在
    if [ ! -f "go.mod" ]; then
        print_error "go.mod 文件不存在，请检查后端目录"
        exit 1
    fi
    
    # 下载依赖
    print_info "下载Go依赖..."
    go mod download
    
    # 编译
    print_info "编译后端程序..."
    if go build -o blog -ldflags "-s -w" main.go; then
        print_info "后端编译成功: $BACKEND_DIR/blog"
    else
        print_error "后端编译失败"
        exit 1
    fi
    
    cd "$PROJECT_ROOT"
}

# 编译前端
build_frontend() {
    print_info "开始编译前端..."
    
    cd "$FRONTEND_DIR"
    
    # 检查package.json是否存在
    if [ ! -f "package.json" ]; then
        print_error "package.json 文件不存在，请检查前端目录"
        exit 1
    fi
    
    # 安装依赖
    print_info "安装前端依赖..."
    if command -v yarn &> /dev/null; then
        yarn install
    else
        npm install
    fi
    
    # 创建生产环境配置文件
    if [ -n "$FRONTEND_DOMAIN" ] && [ -n "$BACKEND_DOMAIN" ]; then
        print_info "创建生产环境配置..."
        cat > .env.production <<EOF
VITE_API_BASE_URL=$BACKEND_DOMAIN/api
VITE_FRONTEND_URL=$FRONTEND_DOMAIN
EOF
    fi
    
    # 编译
    print_info "编译前端应用..."
    if command -v yarn &> /dev/null; then
        yarn build
    else
        npm run build
    fi
    
    if [ -d "dist" ]; then
        print_info "前端编译成功: $FRONTEND_DIR/dist"
    else
        print_error "前端编译失败，dist目录不存在"
        exit 1
    fi
    
    cd "$PROJECT_ROOT"
}

# 显示部署信息
show_deployment_info() {
    echo ""
    print_info "=========================================="
    print_info "部署完成！"
    print_info "=========================================="
    echo ""
    print_info "数据库信息:"
    echo "  数据库名: $DB_NAME"
    echo "  主机: $DB_HOST:$DB_PORT"
    echo "  用户: $DB_USER"
    echo ""
    print_info "管理员账户:"
    echo "  用户名: $ADMIN_USERNAME"
    echo "  密码: $ADMIN_PASSWORD"
    echo ""
    if [ -n "$FRONTEND_DOMAIN" ]; then
        print_info "生产环境配置:"
        echo "  前端域名: $FRONTEND_DOMAIN"
        echo "  后端域名: $BACKEND_DOMAIN"
        echo "  前端构建目录: $FRONTEND_DIR/dist"
        echo "  后端可执行文件: $BACKEND_DIR/blog"
    fi
    echo ""
    print_warn "请妥善保管管理员密码和JWT密钥！"
    echo ""
    print_info "下一步："
    echo "  1. 配置Web服务器（Nginx/Apache）指向前端dist目录"
    echo "  2. 配置反向代理将API请求转发到后端"
    echo "  3. 启动后端服务: cd backend && ./blog"
    echo ""
}

# 主函数
main() {
    echo ""
    print_info "=========================================="
    print_info "博客系统完整部署脚本"
    print_info "=========================================="
    echo ""
    
    # 检查是否已有.env文件
    if [ -f "$ENV_FILE" ]; then
        print_warn ".env 文件已存在"
        read_input "是否重新配置？(y/N)" "N" REBUILD
        if [ "$REBUILD" != "y" ] && [ "$REBUILD" != "Y" ]; then
            print_info "使用现有配置..."
            load_env
        else
            interactive_config
            generate_env_file
        fi
    else
        interactive_config
        generate_env_file
    fi
    
    # 加载环境变量
    load_env
    
    # 部署步骤
    check_requirements
    generate_sql
    deploy_database
    
    # 询问是否编译
    echo ""
    read_input "是否编译后端？(Y/n)" "Y" BUILD_BACKEND
    if [ "$BUILD_BACKEND" != "n" ] && [ "$BUILD_BACKEND" != "N" ]; then
        build_backend
    fi
    
    read_input "是否编译前端？(Y/n)" "Y" BUILD_FRONTEND
    if [ "$BUILD_FRONTEND" != "n" ] && [ "$BUILD_FRONTEND" != "N" ]; then
        build_frontend
    fi
    
    show_deployment_info
}

# 执行主函数
main
