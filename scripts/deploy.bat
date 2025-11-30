@echo off
REM =================================================================================
REM Blog System Deployment Script (Windows)
REM Function: Interactive Config -> Generate .env -> Deploy Database -> Build Frontend/Backend
REM =================================================================================

REM Set code page to UTF-8 to fix Chinese character encoding issues
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

REM Set console output encoding for PowerShell compatibility
if defined PSVersionTable (
    [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
)

REM 脚本目录
set "SCRIPT_DIR=%~dp0"
set "PROJECT_ROOT=%SCRIPT_DIR%.."

REM 配置文件路径
set "ENV_FILE=%PROJECT_ROOT%\.env"
set "SQL_TEMPLATE=%PROJECT_ROOT%\backend\manifest\sql\deploy.sql.template"
set "SQL_OUTPUT=%PROJECT_ROOT%\backend\manifest\sql\deploy.sql"
set "BACKEND_DIR=%PROJECT_ROOT%\backend"
set "FRONTEND_DIR=%PROJECT_ROOT%\fronted"

REM 设置控制台输出为UTF-8
echo.
echo ==========================================
echo 博客系统完整部署脚本
echo ==========================================
echo.

REM Check if .env file exists
if exist "%ENV_FILE%" (
    echo [WARN] .env file already exists
    set /p REBUILD="Reconfigure? (y/N): "
    if /i not "!REBUILD!"=="y" (
        echo [INFO] Using existing configuration...
        goto :load_env
    )
)

REM Interactive configuration
echo.
echo ==========================================
echo Interactive Configuration
echo ==========================================
echo.

REM Database configuration
echo --- Database Configuration ---
set /p DB_NAME="Database Name [go_blog]: "
if "!DB_NAME!"=="" set "DB_NAME=go_blog"

set /p DB_USER="Database User [root]: "
if "!DB_USER!"=="" set "DB_USER=root"

set /p DB_PASSWORD="Database Password: "
if "!DB_PASSWORD!"=="" (
    echo [ERROR] Database password cannot be empty
    exit /b 1
)

set /p DB_HOST="Database Host [127.0.0.1]: "
if "!DB_HOST!"=="" set "DB_HOST=127.0.0.1"

set /p DB_PORT="Database Port [3306]: "
if "!DB_PORT!"=="" set "DB_PORT=3306"
echo.

REM Admin account configuration
echo --- Admin Account Configuration ---
set /p ADMIN_USERNAME="Admin Username [admin]: "
if "!ADMIN_USERNAME!"=="" set "ADMIN_USERNAME=admin"

set /p ADMIN_PASSWORD="Admin Password: "
if "!ADMIN_PASSWORD!"=="" (
    echo [ERROR] Admin password cannot be empty
    exit /b 1
)

set /p ADMIN_NICKNAME="Admin Nickname [Admin]: "
if "!ADMIN_NICKNAME!"=="" set "ADMIN_NICKNAME=Admin"

set /p ADMIN_EMAIL="Admin Email [admin@example.com]: "
if "!ADMIN_EMAIL!"=="" set "ADMIN_EMAIL=admin@example.com"
echo.

REM Website configuration
echo --- Website Configuration ---
set /p SITE_NAME="Site Name [Blog System]: "
if "!SITE_NAME!"=="" set "SITE_NAME=Blog System"

set /p SITE_TITLE="Site Title [Blog System]: "
if "!SITE_TITLE!"=="" set "SITE_TITLE=Blog System"

set /p SITE_SUBTITLE="Site Subtitle [Share programming insights]: "
if "!SITE_SUBTITLE!"=="" set "SITE_SUBTITLE=Share programming insights"

set /p SITE_DESC="Site Description: "
if "!SITE_DESC!"=="" set "SITE_DESC=Share programming insights and technical growth"

set /p SITE_KEYWORDS="Keywords [Go, Vue, Cloud Native]: "
if "!SITE_KEYWORDS!"=="" set "SITE_KEYWORDS=Go, Vue, Cloud Native"

set /p SITE_COPYRIGHT="Copyright [© 2024 Blog System]: "
if "!SITE_COPYRIGHT!"=="" set "SITE_COPYRIGHT=© 2024 Blog System"
echo.

REM Production environment configuration
echo --- Production Environment Configuration ---
set /p FRONTEND_DOMAIN="Frontend Domain (e.g., https://blog.example.com, leave empty to skip): "
set /p BACKEND_DOMAIN="Backend API Domain (e.g., https://api.example.com, leave empty to skip): "
set /p DEPLOY_ENV="Deploy Environment [development]: "
if "!DEPLOY_ENV!"=="" set "DEPLOY_ENV=development"
echo.

REM JWT configuration
echo --- JWT Configuration ---
set /p JWT_SECRET="JWT Secret (leave empty to auto-generate): "
if "!JWT_SECRET!"=="" (
    REM Generate random secret (simple way)
    set "JWT_SECRET=!RANDOM!!RANDOM!!RANDOM!!RANDOM!!RANDOM!!RANDOM!!RANDOM!!RANDOM!"
)
echo.

REM OAuth configuration (optional)
echo --- OAuth Configuration (Optional, press Enter to skip) ---
set /p GITHUB_CLIENT_ID="GitHub Client ID: "
set /p GITHUB_CLIENT_SECRET="GitHub Client Secret: "
set /p GOOGLE_CLIENT_ID="Google Client ID: "
set /p GOOGLE_CLIENT_SECRET="Google Client Secret: "
set /p WECHAT_APP_ID="WeChat AppID: "
set /p WECHAT_APP_SECRET="WeChat AppSecret: "
echo.

REM Redis configuration
echo --- Redis Configuration (Optional) ---
set /p REDIS_HOST="Redis Host [127.0.0.1]: "
if "!REDIS_HOST!"=="" set "REDIS_HOST=127.0.0.1"

set /p REDIS_PORT="Redis Port [6379]: "
if "!REDIS_PORT!"=="" set "REDIS_PORT=6379"

set /p REDIS_DB="Redis Database [1]: "
if "!REDIS_DB!"=="" set "REDIS_DB=1"
echo.

REM 设置OAuth前端URL
if not "!FRONTEND_DOMAIN!"=="" (
    set "OAUTH_FRONTEND_URL=!FRONTEND_DOMAIN!"
) else (
    set "OAUTH_FRONTEND_URL=http://localhost:5173"
)

REM 生成.env文件
echo [INFO] 生成 .env 文件...
(
    echo # =================================================================================
    echo # 博客系统环境变量配置
    echo # 自动生成时间: %DATE% %TIME%
    echo # =================================================================================
    echo.
    echo # ============================================
    echo # 数据库配置
    echo # ============================================
    echo DB_NAME=!DB_NAME!
    echo DB_USER=!DB_USER!
    echo DB_PASSWORD=!DB_PASSWORD!
    echo DB_HOST=!DB_HOST!
    echo DB_PORT=!DB_PORT!
    echo.
    echo # ============================================
    echo # 管理员账户配置
    echo # ============================================
    echo ADMIN_USERNAME=!ADMIN_USERNAME!
    echo ADMIN_PASSWORD=!ADMIN_PASSWORD!
    echo ADMIN_NICKNAME=!ADMIN_NICKNAME!
    echo ADMIN_EMAIL=!ADMIN_EMAIL!
    echo.
    echo # ============================================
    echo # 网站配置
    echo # ============================================
    echo SITE_NAME=!SITE_NAME!
    echo SITE_TITLE=!SITE_TITLE!
    echo SITE_SUBTITLE=!SITE_SUBTITLE!
    echo SITE_DESC=!SITE_DESC!
    echo SITE_KEYWORDS=!SITE_KEYWORDS!
    echo SITE_COPYRIGHT=!SITE_COPYRIGHT!
    echo.
    echo # ============================================
    echo # 生产环境配置
    echo # ============================================
    echo FRONTEND_DOMAIN=!FRONTEND_DOMAIN!
    echo BACKEND_DOMAIN=!BACKEND_DOMAIN!
    echo DEPLOY_ENV=!DEPLOY_ENV!
    echo.
    echo # ============================================
    echo # JWT配置
    echo # ============================================
    echo JWT_SECRET=!JWT_SECRET!
    echo.
    echo # ============================================
    echo # OAuth配置（可选）
    echo # ============================================
    echo OAUTH_FRONTEND_URL=!OAUTH_FRONTEND_URL!
    echo GITHUB_CLIENT_ID=!GITHUB_CLIENT_ID!
    echo GITHUB_CLIENT_SECRET=!GITHUB_CLIENT_SECRET!
    echo GOOGLE_CLIENT_ID=!GOOGLE_CLIENT_ID!
    echo GOOGLE_CLIENT_SECRET=!GOOGLE_CLIENT_SECRET!
    echo WECHAT_APP_ID=!WECHAT_APP_ID!
    echo WECHAT_APP_SECRET=!WECHAT_APP_SECRET!
    echo.
    echo # ============================================
    echo # Redis配置（可选）
    echo # ============================================
    echo REDIS_HOST=!REDIS_HOST!
    echo REDIS_PORT=!REDIS_PORT!
    echo REDIS_DB=!REDIS_DB!
) > "%ENV_FILE%"

echo [INFO] .env 文件已生成: %ENV_FILE%

:load_env
REM 加载环境变量
echo [INFO] 加载环境变量...
for /f "usebackq tokens=1,* delims==" %%a in ("%ENV_FILE%") do (
    set "%%a=%%b"
)

REM 设置默认值
if not defined DB_HOST set "DB_HOST=127.0.0.1"
if not defined DB_PORT set "DB_PORT=3306"
if not defined ADMIN_NICKNAME set "ADMIN_NICKNAME=管理员"
if not defined ADMIN_EMAIL set "ADMIN_EMAIL=admin@example.com"
if not defined SITE_NAME set "SITE_NAME=Blog System"
if not defined SITE_TITLE set "SITE_TITLE=Blog System"
if not defined SITE_SUBTITLE set "SITE_SUBTITLE=分享编程心得，记录技术成长"
if not defined SITE_DESC set "SITE_DESC=分享编程心得，记录技术成长。分享 Go, Vue, 云原生等前沿技术心得。"
if not defined SITE_KEYWORDS set "SITE_KEYWORDS=Go, Vue, Cloud Native"
if not defined SITE_COPYRIGHT set "SITE_COPYRIGHT=© 2024 Blog System"
if not defined DEPLOY_ENV set "DEPLOY_ENV=development"
if not defined OAUTH_FRONTEND_URL set "OAUTH_FRONTEND_URL=http://localhost:5173"
if not defined REDIS_HOST set "REDIS_HOST=127.0.0.1"
if not defined REDIS_PORT set "REDIS_PORT=6379"
if not defined REDIS_DB set "REDIS_DB=1"

REM 检查必需的工具
echo [INFO] 检查必需的工具...
where mysql >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] MySQL 客户端未安装，请先安装 MySQL
    exit /b 1
)

where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [WARN] Go 未安装，将跳过后端编译
    set "SKIP_BACKEND_BUILD=1"
)

where npm >nul 2>&1
if %errorlevel% neq 0 (
    echo [WARN] npm 未安装，将跳过前端编译
    set "SKIP_FRONTEND_BUILD=1"
)

REM 生成SQL文件
echo [INFO] 生成部署SQL文件...
if not exist "%SQL_TEMPLATE%" (
    echo [ERROR] SQL模板文件不存在: %SQL_TEMPLATE%
    exit /b 1
)

REM 使用 PowerShell 进行环境变量替换
powershell -Command "$content = Get-Content '%SQL_TEMPLATE%' -Raw -Encoding UTF8; $content = $content -replace '\$\{DB_NAME\}', '%DB_NAME%'; $content = $content -replace '\$\{ADMIN_USERNAME\}', '%ADMIN_USERNAME%'; $content = $content -replace '\$\{ADMIN_PASSWORD\}', '%ADMIN_PASSWORD%'; $content = $content -replace '\$\{ADMIN_NICKNAME\}', '%ADMIN_NICKNAME%'; $content = $content -replace '\$\{ADMIN_EMAIL\}', '%ADMIN_EMAIL%'; $content = $content -replace '\$\{SITE_NAME\}', '%SITE_NAME%'; $content = $content -replace '\$\{SITE_TITLE\}', '%SITE_TITLE%'; $content = $content -replace '\$\{SITE_SUBTITLE\}', '%SITE_SUBTITLE%'; $content = $content -replace '\$\{SITE_DESC\}', '%SITE_DESC%'; $content = $content -replace '\$\{SITE_KEYWORDS\}', '%SITE_KEYWORDS%'; $content = $content -replace '\$\{SITE_COPYRIGHT\}', '%SITE_COPYRIGHT%'; Set-Content -Path '%SQL_OUTPUT%' -Value $content -Encoding UTF8"

if %errorlevel% neq 0 (
    echo [ERROR] SQL文件生成失败
    exit /b 1
)

echo [INFO] SQL文件已生成: %SQL_OUTPUT%

REM 测试数据库连接
echo [INFO] 测试数据库连接...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% -e "SELECT 1" >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] 数据库连接失败，请检查配置
    exit /b 1
)
echo [INFO] 数据库连接成功

REM 执行SQL脚本
echo [INFO] 执行SQL脚本...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% < "%SQL_OUTPUT%"
if %errorlevel% neq 0 (
    echo [ERROR] 数据库部署失败
    exit /b 1
)

echo [INFO] 数据库部署成功！

REM 编译后端
if not defined SKIP_BACKEND_BUILD (
    echo.
    set /p BUILD_BACKEND="是否编译后端？(Y/n): "
    if /i not "!BUILD_BACKEND!"=="n" (
        echo [INFO] 开始编译后端...
        cd /d "%BACKEND_DIR%"
        
        if not exist "go.mod" (
            echo [ERROR] go.mod 文件不存在，请检查后端目录
            exit /b 1
        )
        
        echo [INFO] 下载Go依赖...
        go mod download
        
        echo [INFO] 编译后端程序...
        go build -o blog.exe -ldflags "-s -w" main.go
        if %errorlevel% neq 0 (
            echo [ERROR] 后端编译失败
            exit /b 1
        )
        
        echo [INFO] 后端编译成功: %BACKEND_DIR%\blog.exe
        cd /d "%PROJECT_ROOT%"
    )
)

REM 编译前端
if not defined SKIP_FRONTEND_BUILD (
    echo.
    set /p BUILD_FRONTEND="是否编译前端？(Y/n): "
    if /i not "!BUILD_FRONTEND!"=="n" (
        echo [INFO] 开始编译前端...
        cd /d "%FRONTEND_DIR%"
        
        if not exist "package.json" (
            echo [ERROR] package.json 文件不存在，请检查前端目录
            exit /b 1
        )
        
        echo [INFO] 安装前端依赖...
        call npm install
        if %errorlevel% neq 0 (
            echo [ERROR] 前端依赖安装失败
            exit /b 1
        )
        
        REM 创建生产环境配置文件
        if not "!FRONTEND_DOMAIN!"=="" if not "!BACKEND_DOMAIN!"=="" (
            echo [INFO] 创建生产环境配置...
            (
                echo VITE_API_BASE_URL=!BACKEND_DOMAIN!/api
                echo VITE_FRONTEND_URL=!FRONTEND_DOMAIN!
            ) > .env.production
        )
        
        echo [INFO] 编译前端应用...
        call npm run build
        if %errorlevel% neq 0 (
            echo [ERROR] 前端编译失败
            exit /b 1
        )
        
        if exist "dist" (
            echo [INFO] 前端编译成功: %FRONTEND_DIR%\dist
        ) else (
            echo [ERROR] 前端编译失败，dist目录不存在
            exit /b 1
        )
        
        cd /d "%PROJECT_ROOT%"
    )
)

REM 显示部署信息
echo.
echo ==========================================
echo 部署完成！
echo ==========================================
echo.
echo [INFO] 数据库信息:
echo   数据库名: %DB_NAME%
echo   主机: %DB_HOST%:%DB_PORT%
echo   用户: %DB_USER%
echo.
echo [INFO] 管理员账户:
echo   用户名: %ADMIN_USERNAME%
echo   密码: %ADMIN_PASSWORD%
echo.
if not "!FRONTEND_DOMAIN!"=="" (
    echo [INFO] 生产环境配置:
    echo   前端域名: %FRONTEND_DOMAIN%
    echo   后端域名: %BACKEND_DOMAIN%
    echo   前端构建目录: %FRONTEND_DIR%\dist
    echo   后端可执行文件: %BACKEND_DIR%\blog.exe
    echo.
)
echo [WARN] 请妥善保管管理员密码和JWT密钥！
echo.
echo [INFO] 下一步：
echo   1. 配置Web服务器（Nginx/IIS）指向前端dist目录
echo   2. 配置反向代理将API请求转发到后端
echo   3. 启动后端服务: cd backend ^&^& blog.exe
echo.

endlocal
