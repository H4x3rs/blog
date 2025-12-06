# OAuth 第三方登录配置说明

本文档说明如何配置 GitHub、Google 和微信登录功能。

## 后端配置

### 1. 修改配置文件

编辑 `backend/manifest/config/config.yaml`，添加 OAuth 配置：

```yaml
oauth:
  github:
    clientId: "your_github_client_id"
    clientSecret: "your_github_client_secret"
    redirectUri: "http://localhost:8000/oauth/github/callback"
  google:
    clientId: "your_google_client_id"
    clientSecret: "your_google_client_secret"
    redirectUri: "http://localhost:8000/oauth/google/callback"
  wechat:
    appId: "your_wechat_app_id"
    appSecret: "your_wechat_app_secret"
    redirectUri: "http://localhost:8000/oauth/wechat/callback"
```

### 2. GitHub OAuth 配置

1. 访问 [GitHub Developer Settings](https://github.com/settings/developers)
2. 点击 "New OAuth App"
3. 填写应用信息：
   - Application name: 你的应用名称
   - Homepage URL: `http://localhost:8000`
   - Authorization callback URL: `http://localhost:8000/oauth/github/callback`
4. 创建后获取 Client ID 和 Client Secret
5. 将 Client ID 和 Client Secret 填入配置文件

### 3. Google OAuth 配置

1. 访问 [Google Cloud Console](https://console.cloud.google.com/)
2. 创建新项目或选择现有项目
3. 启用 Google+ API
4. 创建 OAuth 2.0 客户端 ID：
   - 应用类型：Web 应用
   - 已获授权的重定向 URI：`http://localhost:8000/oauth/google/callback`
5. 获取 Client ID 和 Client Secret
6. 将 Client ID 和 Client Secret 填入配置文件

### 4. 微信登录配置

1. 访问 [微信开放平台](https://open.weixin.qq.com/)
2. 注册开发者账号并创建网站应用
3. 填写应用信息：
   - 网站名称：你的网站名称
   - 网站域名：`localhost:8000`
   - 授权回调域：`localhost:8000`
4. 获取 AppID 和 AppSecret
5. 将 AppID 和 AppSecret 填入配置文件

## 前端配置

前端代码已经集成 OAuth 登录功能，无需额外配置。

### 使用方式

1. 访问登录页面：`http://localhost:5173/login-new`
2. 点击对应的社交登录按钮（GitHub、Google 或微信）
3. 系统会自动跳转到第三方登录页面
4. 授权后会自动回调并完成登录

## API 接口

### 获取 OAuth 登录 URL

```
GET /oauth/:provider/login
```

参数：
- `provider`: 登录提供商（github、google、wechat）

响应：
```json
{
  "redirectUrl": "https://github.com/login/oauth/authorize?..."
}
```

### OAuth 回调处理

```
GET /oauth/:provider/callback?code=xxx&state=xxx
```

参数：
- `provider`: 登录提供商
- `code`: 授权码
- `state`: 状态码

响应：
```json
{
  "token": "oauth-token-github-123",
  "username": "github_123456",
  "nickname": "User Name",
  "avatar": "https://avatars.githubusercontent.com/..."
}
```

## 注意事项

1. **生产环境配置**：
   - 将 `redirectUri` 中的 `localhost:8000` 替换为实际域名
   - 确保 HTTPS 配置正确

2. **安全性**：
   - 不要将 Client Secret 提交到代码仓库
   - 使用环境变量或密钥管理服务存储敏感信息

3. **用户创建**：
   - OAuth 登录会自动创建用户账号
   - 用户名格式：`{provider}_{openid}`
   - OAuth 用户密码为空

4. **Token 管理**：
   - 当前实现使用简化的 token 生成方式
   - 生产环境建议使用 JWT 或其他安全的 token 机制

## 测试

1. 确保后端服务运行在 `http://localhost:8000`
2. 确保前端服务运行在 `http://localhost:5173`
3. 访问登录页面测试 OAuth 登录功能

## 故障排查

1. **无法跳转到第三方登录页面**：
   - 检查配置文件中的 Client ID 是否正确
   - 检查 redirectUri 是否与第三方平台配置一致

2. **回调失败**：
   - 检查回调 URL 是否正确配置
   - 检查 Client Secret 是否正确
   - 查看后端日志获取详细错误信息

3. **用户创建失败**：
   - 检查数据库连接是否正常
   - 检查用户表结构是否正确






