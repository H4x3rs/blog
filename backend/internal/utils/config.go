package utils

import (
	"os"

	"github.com/gogf/gf/v2/frame/g"
)

// GetConfigString 从环境变量或配置文件获取字符串配置
// 优先从环境变量读取，如果没有则从配置文件读取，最后使用默认值
func GetConfigString(ctx g.Ctx, envKey, configKey, defaultValue string) string {
	// 优先从环境变量读取
	if value := os.Getenv(envKey); value != "" {
		return value
	}
	// 从配置文件读取
	return g.Cfg().MustGet(ctx, configKey, defaultValue).String()
}

// GetConfigInt 从环境变量或配置文件获取整数配置
func GetConfigInt(ctx g.Ctx, envKey, configKey string, defaultValue int) int {
	// 优先从环境变量读取
	if value := os.Getenv(envKey); value != "" {
		if intValue := g.NewVar(value).Int(); intValue != 0 {
			return intValue
		}
	}
	// 从配置文件读取
	return g.Cfg().MustGet(ctx, configKey, defaultValue).Int()
}

// GetConfigBool 从环境变量或配置文件获取布尔配置
func GetConfigBool(ctx g.Ctx, envKey, configKey string, defaultValue bool) bool {
	// 优先从环境变量读取
	if value := os.Getenv(envKey); value != "" {
		return g.NewVar(value).Bool()
	}
	// 从配置文件读取
	return g.Cfg().MustGet(ctx, configKey, defaultValue).Bool()
}

