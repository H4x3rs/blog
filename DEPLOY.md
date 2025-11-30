# 博客系统部署文档

本文档说明如何部署博客系统，包括数据库初始化、环境配置和应用程序启动。

## 目录

- [环境要求](#环境要求)
- [快速开始](#快速开始)
- [环境变量配置](#环境变量配置)
- [数据库部署](#数据库部署)
- [应用配置](#应用配置)
- [启动应用](#启动应用)
- [验证部署](#验证部署)

## 环境要求

### 必需软件

- **MySQL** 5.7+ 或 8.0+
- **Go** 1.18+ (后端)
- **Node.js** 16+ 和 npm/yarn (前端)
- **Redis** 5.0+ (可选，用于缓存)

### 系统要求

- Linux/macOS/Windows
- 至少 2GB 可用内存
- 至少 10GB 可用磁盘空间

## 快速开始

### 方法一：交互式部署（推荐）

使用交互式部署脚本，自动完成所有配置和编译：

```bash
# Linux/macOS
chmod +x scripts/deploy.sh
./scripts/deploy.sh

# Windows (推荐使用 PowerShell)
powershell -ExecutionPolicy Bypass -File scripts\deploy.ps1

# Windows (CMD，如果 PowerShell 不可用)
scripts\deploy.bat
```

**注意：** 如果在 PowerShell 中执行 `deploy.bat` 出现中文乱码，请使用 `deploy.ps1` 脚本。

脚本会引导您完成：
1. 数据库配置
2. 管理员账户配置
3. 网站信息配置
4. 生产环境域名配置（可选）
5. OAuth配置（可选）
6. 自动生成 `.env` 文件
7. 部署数据库
8. 编译后端（可选）
9. 编译前端（可选）

### 方法二：手动配置

#### 1. 克隆项目

```bash
git clone <repository-url>
cd blog
```

#### 2. 配置环境变量

复制环境变量模板文件：

```bash
# Linux/macOS
cp .env.example .env

# Windows
copy .env.example .env
```

编辑 `.env` 文件，设置必要的环境变量（见下方说明）。

#### 3. 部署数据库

使用提供的部署脚本：

```bash
# Linux/macOS
chmod +x scripts/deploy.sh
./scripts/deploy.sh

# Windows
scripts\deploy.bat
```

#### 4. 配置应用

编辑 `backend/manifest/config/config.yaml`，设置数据库连接和其他配置。

#### 5. 编译应用

```bash
# 编译后端
cd backend
go build -o blog main.go

# 编译前端
cd ../fronted
npm install
npm run build
```

#### 6. 启动应用

```bash
# 启动后端
cd backend
./blog

# 启动前端（开发模式）
cd fronted
npm run dev
```

## 环境变量配置

### 数据库相关

创建 `.env` 文件（在项目根目录），包含以下环境变量：

```bash
# ============================================
# 数据库配置
# ============================================
DB_NAME=go_blog
DB_USER=root
DB_PASSWORD=your_password_here
DB_HOST=127.0.0.1
DB_PORT=3306

# ============================================
# 管理员账户配置
# ============================================
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_secure_password_here
ADMIN_NICKNAME=管理员
ADMIN_EMAIL=admin@example.com

# ============================================
# 网站配置
# ============================================
SITE_NAME=Blog System
SITE_TITLE=Blog System
SITE_SUBTITLE=分享编程心得，记录技术成长
SITE_DESC=分享编程心得，记录技术成长
SITE_KEYWORDS=Go, Vue, Cloud Native
SITE_COPYRIGHT=© 2024 Blog System

# ============================================
# 生产环境配置
# ============================================
# 前端域名（生产环境）
FRONTEND_DOMAIN=https://blog.example.com
# 后端API域名（生产环境）
BACKEND_DOMAIN=https://api.example.com
# 部署环境: development 或 production
DEPLOY_ENV=development

# ============================================
# JWT配置
# ============================================
JWT_SECRET=your-jwt-secret-key-min-32-chars-change-in-production

# ============================================
# OAuth配置（可选）
# ============================================
OAUTH_FRONTEND_URL=http://localhost:5173
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
WECHAT_APP_ID=your_wechat_app_id
WECHAT_APP_SECRET=your_wechat_app_secret

# ============================================
# Redis配置（可选）
# ============================================
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_DB=1
```

### 环境变量说明

| 变量名 | 说明 | 必需 | 默认值 |
|--------|------|------|--------|
| `DB_NAME` | 数据库名称 | 是 | `go_blog` |
| `DB_USER` | 数据库用户名 | 是 | `root` |
| `DB_PASSWORD` | 数据库密码 | 是 | - |
| `DB_HOST` | 数据库主机 | 是 | `127.0.0.1` |
| `DB_PORT` | 数据库端口 | 否 | `3306` |
| `ADMIN_USERNAME` | 管理员用户名 | 是 | `admin` |
| `ADMIN_PASSWORD` | 管理员密码 | 是 | - |
| `ADMIN_NICKNAME` | 管理员昵称 | 否 | `管理员` |
| `ADMIN_EMAIL` | 管理员邮箱 | 否 | `admin@example.com` |
| `SITE_NAME` | 网站名称 | 否 | `Blog System` |
| `FRONTEND_DOMAIN` | 前端域名（生产环境） | 否 | - |
| `BACKEND_DOMAIN` | 后端API域名（生产环境） | 否 | - |
| `DEPLOY_ENV` | 部署环境 | 否 | `development` |
| `JWT_SECRET` | JWT密钥（至少32字符） | 是 | - |

## 数据库部署

### 使用交互式部署脚本（推荐）

#### Linux/macOS

```bash
chmod +x scripts/deploy.sh
./scripts/deploy.sh
```

#### Windows

**推荐使用 PowerShell 脚本（解决中文乱码）：**
```powershell
powershell -ExecutionPolicy Bypass -File scripts\deploy.ps1
```

**或使用批处理脚本：**
```cmd
scripts\deploy.bat
```

脚本会自动：
1. **交互式配置**：引导您输入所有必要的配置信息
2. **生成 .env 文件**：自动创建环境变量配置文件
3. **替换 SQL 模板**：将环境变量替换到 SQL 模板中
4. **执行数据库初始化**：自动创建数据库和表结构
5. **编译后端**（可选）：编译 Go 后端程序
6. **编译前端**（可选）：编译 Vue 前端应用

### 手动部署数据库

如果您已经配置好 `.env` 文件，可以直接执行数据库部署部分：

```bash
# 生成SQL文件
envsubst < backend/manifest/sql/deploy.sql.template > backend/manifest/sql/deploy.sql

# 执行SQL
mysql -u $DB_USER -p$DB_PASSWORD < backend/manifest/sql/deploy.sql
```

### 方法二：手动部署

1. 设置环境变量（或导出到当前shell）：
```bash
export DB_NAME=go_blog
export DB_USER=root
export DB_PASSWORD=your_password
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD=your_secure_password
# ... 其他变量
```

2. 生成部署SQL：
```bash
# Linux/macOS
envsubst < backend/manifest/sql/deploy.sql.template > backend/manifest/sql/deploy.sql

# Windows (需要安装 Git Bash 或 WSL)
envsubst < backend/manifest/sql/deploy.sql.template > backend/manifest/sql/deploy.sql
```

3. 执行SQL：
```bash
mysql -u $DB_USER -p$DB_PASSWORD < backend/manifest/sql/deploy.sql
```

## 应用配置

### 后端配置

编辑 `backend/manifest/config/config.yaml`：

```yaml
server:
  address: ":8000"

database:
  default:
    link: "mysql:${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local"
    debug: true

jwt:
  secret: "${JWT_SECRET}"

oauth:
  frontendUrl: "${OAUTH_FRONTEND_URL}"
  github:
    clientId: "${GITHUB_CLIENT_ID}"
    clientSecret: "${GITHUB_CLIENT_SECRET}"
  # ... 其他OAuth配置
```

**注意：** GoFrame 支持环境变量替换，使用 `${VAR_NAME}` 格式。

### 前端配置

前端配置主要在 `fronted/src/utils/request.js` 中，API 基础URL可以通过环境变量配置：

创建 `fronted/.env`：

```bash
VITE_API_BASE_URL=http://localhost:8000/api
```

## 编译和启动应用

### 开发环境

#### 启动后端

```bash
cd backend
go mod download
go run main.go
```

后端默认运行在 `http://localhost:8000`

#### 启动前端

```bash
cd fronted
npm install
npm run dev
```

前端默认运行在 `http://localhost:5173`

### 生产环境部署

#### 编译后端

```bash
cd backend
go mod download
go build -o blog -ldflags "-s -w" main.go
```

编译后的可执行文件位于 `backend/blog`（Linux/macOS）或 `backend/blog.exe`（Windows）

#### 编译前端

```bash
cd fronted
npm install

# 如果配置了生产环境域名，创建 .env.production
cat > .env.production <<EOF
VITE_API_BASE_URL=https://api.example.com/api
VITE_FRONTEND_URL=https://blog.example.com
EOF

npm run build
```

编译后的静态文件位于 `fronted/dist` 目录

#### 部署步骤

1. **部署前端**：
   - 将 `fronted/dist` 目录内容上传到 Web 服务器（Nginx/Apache）
   - 配置域名指向该目录

2. **部署后端**：
   - 将编译后的 `blog` 或 `blog.exe` 文件上传到服务器
   - 配置 `backend/manifest/config/config.yaml` 中的数据库连接
   - 启动后端服务：
     ```bash
     ./blog
     ```

3. **配置反向代理**（Nginx 示例）：
   ```nginx
   server {
       listen 80;
       server_name blog.example.com;
       
       # 前端静态文件
       location / {
           root /path/to/fronted/dist;
           try_files $uri $uri/ /index.html;
       }
       
       # 后端API代理
       location /api {
           proxy_pass http://localhost:8000;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

## 验证部署

### 1. 检查数据库

```bash
mysql -u $DB_USER -p$DB_PASSWORD -e "USE $DB_NAME; SHOW TABLES;"
```

应该看到所有表都已创建。

### 2. 检查管理员账户

使用部署时设置的管理员账户登录：
- 用户名：`${ADMIN_USERNAME}`
- 密码：`${ADMIN_PASSWORD}`

### 3. 访问应用

- 前端：http://localhost:5173
- 后端API：http://localhost:8000/api
- API文档：http://localhost:8000/swagger

## 安全建议

1. **生产环境必须修改的配置：**
   - `ADMIN_PASSWORD` - 使用强密码
   - `JWT_SECRET` - 至少32字符的随机字符串
   - `DB_PASSWORD` - 使用强密码
   - OAuth Client Secret - 不要泄露

2. **不要将 `.env` 文件提交到版本控制：**
   ```bash
   echo ".env" >> .gitignore
   ```

3. **使用 HTTPS：**
   - 生产环境必须使用 HTTPS
   - 配置 SSL 证书

4. **定期备份数据库：**
   ```bash
   mysqldump -u $DB_USER -p$DB_PASSWORD $DB_NAME > backup.sql
   ```

## 故障排查

### 数据库连接失败

1. 检查 MySQL 服务是否运行
2. 检查环境变量中的数据库配置
3. 检查防火墙设置

### 权限错误

1. 确保数据库用户有创建数据库的权限
2. 检查 SQL 脚本执行权限

### 端口占用

1. 检查 8000 端口（后端）是否被占用
2. 检查 5173 端口（前端）是否被占用
3. 修改 `config.yaml` 中的端口配置

## 获取帮助

如遇到问题，请检查：
1. 日志文件
2. 数据库错误日志
3. 应用控制台输出

## 更新日志

- v1.0 - 初始版本，支持环境变量配置

