package middleware

import (
	"context"
	"fmt"
	"time"
	"yanny-service/internal/database"

	"github.com/gin-gonic/gin"
)

// ========== 视频访问防刷 ==========

const (
	// 同一用户/IP 对同一视频的观看间隔（秒），短于此间隔不计入播放量
	viewDebounceSeconds = 30
	// 每用户每小时最大视频播放次数
	maxViewsPerUserPerHour = 300
	// 每 IP 每小时最大视频播放次数
	maxViewsPerIPPerHour = 500
)

// VideoViewAntiAbuse 视频播放防刷中间件
// 1. 同用户/同IP 对同一视频短时间内重复请求去重
// 2. 单用户/IP 超频播放拦截
func VideoViewAntiAbuse() gin.HandlerFunc {
	return func(c *gin.Context) {
		videoID := c.Param("id")
		if videoID == "" {
			c.Next()
			return
		}

		ip := c.ClientIP()
		// 注意：JWT 中间件存储的是 uint64，不能用 GetString
		userID := GetUserID(c)
		ctx := context.Background()

		// 1. 去重：同一标识对同一视频 30 秒内只计一次
		identifier := fmt.Sprintf("%d", userID)
		if userID == 0 {
			identifier = ip
		}
		dedupKey := fmt.Sprintf("yanny:video:view:%s:%s", identifier, videoID)
		exists, _ := database.Redis.Exists(ctx, dedupKey).Result()
		if exists > 0 {
			// 不拦截播放，但标记为重复（不计入播放量）
			c.Set("skip_view_count", true)
			c.Next()
			return
		}
		database.Redis.Set(ctx, dedupKey, "1", viewDebounceSeconds*time.Second)

		// 2. 频率检查：用户维度（仅登录用户）
		if userID > 0 {
			if !checkViewRateLimit(ctx, identifier, maxViewsPerUserPerHour) {
				c.Set("skip_view_count", true)
				c.Next()
				return
			}
		}

		// 3. 频率检查：IP 维度
		if !checkViewRateLimit(ctx, ip, maxViewsPerIPPerHour) {
			c.Set("skip_view_count", true)
			c.Next()
			return
		}

		c.Next()
	}
}

// checkViewRateLimit 检查观看频率是否超限
func checkViewRateLimit(ctx context.Context, identifier string, limit int) bool {
	now := time.Now()
	windowKey := fmt.Sprintf("yanny:video:ratelimit:%s:%d", identifier, now.Hour())
	count, err := database.Redis.Incr(ctx, windowKey).Result()
	if err != nil {
		return true // Redis 故障降级放行
	}
	database.Redis.Expire(ctx, windowKey, 2*time.Hour)

	return count <= int64(limit)
}
