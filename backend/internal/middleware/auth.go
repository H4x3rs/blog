package middleware

import (
	"blog/internal/service"
	"blog/internal/utils"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Auth 认证中间件，从请求头中提取JWT token并解析用户信息
// 同时验证token是否在Redis中存在（如果Redis验证失败，降级到JWT验证）
func Auth(r *ghttp.Request) {
	// 从请求头中获取token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		r.Middleware.Next()
		return
	}

	// 提取Bearer token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		r.Middleware.Next()
		return
	}

	token := parts[1]
	
	// 解析token
	claims, err := utils.ParseToken(token)
	if err != nil {
		r.Middleware.Next()
		return
	}

	// 验证token是否在Redis中存在
	ctx := r.GetCtx()
	exists, err := service.Token.CheckToken(ctx, claims.UserID, token)
	if err != nil {
		// Redis连接错误，记录日志，不设置用户信息
		g.Log().Warning(ctx, "Redis验证token失败:", err)
		r.Middleware.Next()
		return
	}
	if !exists {
		// Redis中没有token，说明用户未登录或已下线
		r.Middleware.Next()
		return
	}

	// Redis验证通过，将用户信息存储到请求上下文中
	r.SetCtxVar("user_id", claims.UserID)
	r.SetCtxVar("username", claims.Username)

	r.Middleware.Next()
}

// GetUserID 从请求上下文中获取用户ID
func GetUserID(r *ghttp.Request) int {
	return r.GetCtxVar("user_id", 0).Int()
}

// GetUsername 从请求上下文中获取用户名
func GetUsername(r *ghttp.Request) string {
	return r.GetCtxVar("username", "").String()
}








