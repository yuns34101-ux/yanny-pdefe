package middleware

import (
	"context"
	"fmt"
	"time"
	"yanny-service/internal/database"
	"yanny-service/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter 基于 Redis 滑动窗口的限流中间件
//
//	window: 时间窗口
//	limit:  窗口内最大请求数
//	keyFn:  限流 key 生成函数（nil 则使用客户端 IP）
func RateLimiter(window time.Duration, limit int, keyFn func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string
		if keyFn != nil {
			key = keyFn(c)
		} else {
			key = fmt.Sprintf("yanny:ratelimit:%s", c.ClientIP())
		}

		ctx := context.Background()
		now := time.Now().UnixMilli()
		windowStart := now - window.Milliseconds()

		pipe := database.Redis.Pipeline()

		// 移除窗口外的记录
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
		// 统计当前窗口内请求数
		countCmd := pipe.ZCard(ctx, key)
		// 添加当前请求
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: fmt.Sprintf("%d:%s", now, c.ClientIP()),
		})
		// 设置 key 过期
		pipe.Expire(ctx, key, window+time.Second)

		_, err := pipe.Exec(ctx)
		if err != nil {
			// Redis 故障时降级放行
			c.Next()
			return
		}

		count, _ := countCmd.Result()
		if count >= int64(limit) {
			dto.ErrorWithStatus(c, 429, dto.ErrCodeParamInvalid, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// API 级别的限流策略预设

// AdminAPILimiter 管理后台接口限流（每 IP 每分钟 120 次）
func AdminAPILimiter() gin.HandlerFunc {
	return RateLimiter(1*time.Minute, 120, func(c *gin.Context) string {
		return fmt.Sprintf("yanny:ratelimit:admin:%s", c.ClientIP())
	})
}

// MpAPILimiter 小程序接口限流（每 IP 每分钟 300 次）
func MpAPILimiter() gin.HandlerFunc {
	return RateLimiter(1*time.Minute, 300, func(c *gin.Context) string {
		return fmt.Sprintf("yanny:ratelimit:mp:%s", c.ClientIP())
	})
}

// LoginRateLimiter 登录接口限流（每 IP 每分钟 5 次，防暴力破解）
func LoginRateLimiter() gin.HandlerFunc {
	return RateLimiter(1*time.Minute, 5, func(c *gin.Context) string {
		return fmt.Sprintf("yanny:ratelimit:login:%s", c.ClientIP())
	})
}
