package middleware

import (
	"context"
	"fmt"
	"time"
	"yanny-service/internal/database"
	"yanny-service/internal/dto"

	"github.com/gin-gonic/gin"
)

const (
	// IP 黑名单 Redis Key
	ipBlacklistKey = "yanny:blacklist:ip"
	// JWT 黑名单 Redis Key 前缀（登出后失效 token）
	tokenBlacklistPrefix = "yanny:blacklist:token:"
)

// ========== IP 黑名单 ==========

// AddIPBlacklist 添加 IP 到黑名单
// 使用独立 String key 存储，支持 TTL 自动过期
func AddIPBlacklist(ip string, duration time.Duration, reason string) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", ipBlacklistKey, ip)
	database.Redis.Set(ctx, key, reason, duration)
}

// RemoveIPBlacklist 从黑名单移除 IP
func RemoveIPBlacklist(ip string) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", ipBlacklistKey, ip)
	database.Redis.Del(ctx, key)
}

// IsIPBlacklisted 检查 IP 是否在黑名单中
func IsIPBlacklisted(ip string) bool {
	if database.Redis == nil {
		return false
	}
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", ipBlacklistKey, ip)
	exists, _ := database.Redis.Exists(ctx, key).Result()
	return exists > 0
}

// IPBlacklistMiddleware IP 黑名单拦截中间件
func IPBlacklistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if IsIPBlacklisted(ip) {
			dto.ErrorWithStatus(c, 403, dto.ErrCodeForbidden, "您的 IP 已被限制访问")
			c.Abort()
			return
		}
		c.Next()
	}
}

// ========== Token 黑名单（登出） ==========

// AddTokenBlacklist 将 JWT 加入黑名单（登出时调用）
func AddTokenBlacklist(tokenStr string, expireAt time.Time) {
	ctx := context.Background()
	key := tokenBlacklistPrefix + hashToken(tokenStr)
	duration := time.Until(expireAt)
	if duration <= 0 {
		return
	}
	database.Redis.Set(ctx, key, "1", duration)
}

// IsTokenBlacklisted 检查 JWT 是否在黑名单中
func IsTokenBlacklisted(tokenStr string) bool {
	if database.Redis == nil {
		return false // Redis 未初始化时降级放行（测试环境等）
	}
	ctx := context.Background()
	key := tokenBlacklistPrefix + hashToken(tokenStr)
	exists, _ := database.Redis.Exists(ctx, key).Result()
	return exists > 0
}

// hashToken 对 token 做简单哈希避免 key 过长
func hashToken(tokenStr string) string {
	if len(tokenStr) <= 32 {
		return tokenStr
	}
	return tokenStr[len(tokenStr)-32:]
}

// ========== 封禁触发 ==========

// BanIPOnViolation 当某 IP 在窗口内触发过多异常行为时自动封禁
// window: 检测窗口, threshold: 触发阈值, banDuration: 封禁时长
func BanIPOnViolation(ip, violationType string, window time.Duration, threshold int, banDuration time.Duration) {
	ctx := context.Background()
	key := fmt.Sprintf("yanny:violation:%s:%s", violationType, ip)

	count, _ := database.Redis.Incr(ctx, key).Result()
	if count == 1 {
		database.Redis.Expire(ctx, key, window)
	}

	if count >= int64(threshold) {
		AddIPBlacklist(ip, banDuration, fmt.Sprintf("自动封禁：%s 行为异常（%d 次/%v）", violationType, threshold, window))
		database.Redis.Del(ctx, key) // 清理计数器
	}
}
