package service

import (
	"context"
	"errors"
	"fmt"

	"blog/internal/utils"

	"github.com/gogf/gf/v2/frame/g"
)

type sToken struct{}

var Token = sToken{}

// SetToken 将token存储到Redis
// key格式: token:{userID}:{tokenHash}
func (s *sToken) SetToken(ctx context.Context, userID int, token string) error {
	// 获取token过期时间（单位：小时）
	expiresHours := utils.GetTokenExpiresHours(ctx)
	
	// 生成Redis key
	key := s.getTokenKey(userID, token)
	
	g.Log().Infof(ctx, "正在存储token到Redis, key: %s, 过期时间: %d秒", key, expiresHours*3600)
	
	// 存储token到Redis，设置过期时间
	redis := g.Redis()
	result, err := redis.Do(ctx, "SETEX", key, expiresHours*3600, "1")
	if err != nil {
		g.Log().Errorf(ctx, "存储token到Redis失败: %v", err)
		return fmt.Errorf("存储token到Redis失败: %v", err)
	}
	
	g.Log().Infof(ctx, "token存储到Redis成功, 结果: %v", result)
	return nil
}

// CheckToken 检查token是否在Redis中存在且有效
func (s *sToken) CheckToken(ctx context.Context, userID int, token string) (bool, error) {
	// 生成Redis key
	key := s.getTokenKey(userID, token)
	
	// 从Redis获取token
	redis := g.Redis()
	value, err := redis.Do(ctx, "GET", key)
	if err != nil {
		// Redis错误，记录日志但不影响功能（降级到JWT验证）
		g.Log().Warning(ctx, "从Redis获取token失败:", err)
		return false, err
	}
	
	// 如果Redis中没有token，返回false
	if value == nil || value.IsEmpty() {
		return false, nil
	}
	
	return true, nil
}

// DeleteToken 从Redis删除token（退出登录时使用）
func (s *sToken) DeleteToken(ctx context.Context, userID int, token string) error {
	// 生成Redis key
	key := s.getTokenKey(userID, token)
	
	// 从Redis删除token
	redis := g.Redis()
	_, err := redis.Do(ctx, "DEL", key)
	if err != nil {
		return fmt.Errorf("从Redis删除token失败: %v", err)
	}
	
	return nil
}

// DeleteUserTokens 删除用户的所有token（强制下线时使用）
func (s *sToken) DeleteUserTokens(ctx context.Context, userID int) error {
	// 生成匹配所有用户token的pattern
	pattern := fmt.Sprintf("token:%d:*", userID)
	
	// 查找所有匹配的key
	redis := g.Redis()
	keys, err := redis.Do(ctx, "KEYS", pattern)
	if err != nil {
		return fmt.Errorf("查找用户token失败: %v", err)
	}
	
	// 如果没有找到key，直接返回
	if keys == nil || keys.IsEmpty() {
		return nil
	}
	
	// 批量删除所有匹配的key
	keysArray := keys.Array()
	if len(keysArray) > 0 {
		_, err = redis.Do(ctx, "DEL", keysArray...)
		if err != nil {
			return fmt.Errorf("删除用户token失败: %v", err)
		}
	}
	
	return nil
}

// RefreshToken 刷新token的过期时间
func (s *sToken) RefreshToken(ctx context.Context, userID int, token string) error {
	// 检查token是否存在
	exists, err := s.CheckToken(ctx, userID, token)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("token不存在")
	}
	
	// 获取token过期时间（单位：小时）
	expiresHours := utils.GetTokenExpiresHours(ctx)
	
	// 生成Redis key
	key := s.getTokenKey(userID, token)
	
	// 刷新token过期时间
	redis := g.Redis()
	_, err = redis.Do(ctx, "EXPIRE", key, expiresHours*3600)
	if err != nil {
		return fmt.Errorf("刷新token过期时间失败: %v", err)
	}
	
	return nil
}

// getTokenKey 生成Redis key
// 使用token的前16个字符作为hash，避免key过长
func (s *sToken) getTokenKey(userID int, token string) string {
	// 使用token的前32个字符作为标识（JWT token通常很长）
	tokenHash := token
	if len(token) > 32 {
		tokenHash = token[:32]
	}
	return fmt.Sprintf("token:%d:%s", userID, tokenHash)
}

// GetTokenTTL 获取token的剩余过期时间（单位：秒）
func (s *sToken) GetTokenTTL(ctx context.Context, userID int, token string) (int64, error) {
	// 生成Redis key
	key := s.getTokenKey(userID, token)
	
	// 获取token的剩余过期时间
	redis := g.Redis()
	ttl, err := redis.Do(ctx, "TTL", key)
	if err != nil {
		return 0, fmt.Errorf("获取token过期时间失败: %v", err)
	}
	
	return ttl.Int64(), nil
}
