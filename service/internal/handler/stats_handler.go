package handler

import (
	"fmt"
	"yanny-service/internal/dto"
	"yanny-service/internal/repository"

	"github.com/gin-gonic/gin"
)

// GetDashboardStats 数据看板汇总
func GetDashboardStats(c *gin.Context) {
	scope := getEntityScope(c)

	// 平台汇总
	platform, _ := repository.GetPlatformDailyStats(scope, 7)
	// 热门视频 Top10
	topVideos, _ := repository.GetTopVideos(scope, 10)
	// 地域分布
	regions, _ := repository.GetRegionStats(scope, 7)

	dto.Success(c, gin.H{
		"platform":  platform,
		"top_videos": topVideos,
		"regions":   regions,
	})
}

// GetTrendData 播放量趋势
func GetTrendData(c *gin.Context) {
	scope := getEntityScope(c)
	days := 7
	if d := c.Query("days"); d != "" {
		fmt.Sscanf(d, "%d", &days)
	}
	trend, _ := repository.GetPlatformDailyStats(scope, days)
	dto.Success(c, trend)
}

// GetTopVideos 热门视频排行
func GetTopVideos(c *gin.Context) {
	scope := getEntityScope(c)
	limit := 10
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	videos, _ := repository.GetTopVideos(scope, limit)
	dto.Success(c, videos)
}

// GetRegionStats 地域分布
func GetRegionStats(c *gin.Context) {
	scope := getEntityScope(c)
	days := 7
	if d := c.Query("days"); d != "" {
		fmt.Sscanf(d, "%d", &days)
	}
	regions, _ := repository.GetRegionStats(scope, days)
	dto.Success(c, regions)
}
