# OAuth 前端回调模式实现说明

## 实现概述

已实现前端回调模式的OAuth登录流程，GitHub认证完成后回调到前端页面，前端获取code后调用后端API，后端返回JWT token。

## 流程说明

```
1. 用户点击GitHub登录按钮
   ↓
2. 前端保存provider到sessionStorage
   ↓
3. 前端获取授权URL并跳转到GitHub
   ↓
4. GitHub授权后回调到前端：http://localhost:5173/oauth/callback?code=xxx
   ↓
5. 前端OAuthCallback页面接收code
   ↓
6. 前端调用后端API：GET /oauth/github/callback?code=xxx
   ↓
7. 后端处理OAuth回调，生成JWT token
   ↓
8. 后端返回JSON：{token, username, nickname, avatar}
   ↓
9. 前端保存token到localStorage
   ↓
10. 前端跳转到/admin/dashboard
```

## 后端修改

### 1. 添加JWT依赖
```bash
go get github.com/golang-jwt/jwt/v5
```

### 2. 创建JWT工具函数
- 文件：`backend/internal/utils/jwt.go`
- 功能：生成和验证JWT token
- 配置：从`config.yaml`读取`jwt.secret`

### 3. 修改配置文件
- 文件：`backend/manifest/config/config.yaml`
- 添加JWT配置
- 修改OAuth回调地址为前端地址：`http://localhost:5173/oauth/callback`

### 4. 修改OAuth服务
- 文件：`backend/internal/service/oauth.go`
- 使用前端回调地址生成授权URL
- 生成JWT token而不是简单字符串

### 5. 修改OAuth控制器
- 文件：`backend/internal/controller/oauth/oauth.go`
- Callback方法返回JSON而不是重定向
- 移除手动路由注册

## 前端修改

### 1. 修改登录页面
- 文件：`fronted/src/views/LoginNew.vue`
- 在点击OAuth登录按钮时保存provider到sessionStorage

### 2. 修改OAuth回调页面
- 文件：`fronted/src/views/OAuthCallback.vue`
- 从URL获取code参数
- 从sessionStorage获取provider
- 调用后端API获取token
- 保存token并跳转到admin页面

### 3. 清理AdminLayout
- 文件：`fronted/src/layout/AdminLayout.vue`
- 移除OAuth回调处理逻辑（现在由OAuthCallback页面处理）

## 配置说明

### GitHub OAuth App配置

在GitHub的OAuth App设置中，将Authorization callback URL设置为：
```
http://localhost:5173/oauth/callback
```

### 后端配置文件

```yaml
jwt:
  secret: "your-secret-key-change-in-production-min-32-chars"

oauth:
  frontendUrl: "http://localhost:5173"
  github:
    clientId: "your_client_id"
    clientSecret: "your_client_secret"
    redirectUri: "http://localhost:5173/oauth/callback"
```

## API接口

### 获取OAuth授权URL
```
GET /oauth/:provider/login
```

响应：
```json
{
  "redirectUrl": "https://github.com/login/oauth/authorize?..."
}
```

### OAuth回调处理
```
GET /oauth/:provider/callback?code=xxx&state=xxx
```

响应：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "github_123456",
  "nickname": "User Name",
  "avatar": "https://avatars.githubusercontent.com/..."
}
```

## JWT Token说明

- Token有效期：7天
- 包含信息：user_id, username
- 签名算法：HS256
- 签发者：blog-system

## 安全性

1. **State参数**：建议在生成授权URL时生成随机state，并在回调时验证
2. **JWT密钥**：生产环境必须修改`jwt.secret`为强密钥
3. **HTTPS**：生产环境必须使用HTTPS
4. **Token存储**：前端使用localStorage存储，生产环境可考虑使用httpOnly cookie

## 测试步骤

1. 确保后端服务运行在`http://localhost:8000`
2. 确保前端服务运行在`http://localhost:5173`
3. 在GitHub OAuth App中配置回调地址为`http://localhost:5173/oauth/callback`
4. 访问`http://localhost:5173/login-new`
5. 点击GitHub登录按钮
6. 完成GitHub授权
7. 应该自动跳转到admin页面

## 注意事项

1. **Provider识别**：当前通过sessionStorage传递provider，如果用户直接访问回调页面可能无法识别provider
2. **错误处理**：GitHub可能返回error参数，前端已处理
3. **Token验证**：后端提供了`utils.ParseToken`和`utils.ValidateToken`用于验证token

## 后续优化建议

1. 实现state参数验证，防止CSRF攻击
2. 添加token刷新机制
3. 实现token黑名单（登出时）
4. 添加请求拦截器，自动在请求头中添加token
5. 实现token过期自动刷新




