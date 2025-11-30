package middleware

import (
	"blog/internal/utils"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Auth 认证中间件，从请求头中提取JWT token并解析用户信息
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

	// 将用户信息存储到请求上下文中
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




