package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// 初始化JWT密钥
func init() {
	// 优先从环境变量读取，如果没有则从配置文件读取，最后使用默认值
	secret := GetConfigString(nil, "JWT_SECRET", "jwt.secret", "your-secret-key-change-in-production")
	jwtSecret = []byte(secret)
}

// Claims JWT claims结构
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID int, username string) (string, error) {
	// 设置过期时间（默认7天）
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// ValidateToken 验证token是否有效
func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}




