package repository

import (
	"time"
	"yanny-service/internal/database"
)

// PlatformDailyStat 平台日统计
type PlatformDailyStat struct {
	StatDate   string `json:"stat_date"`
	TotalViews int64  `json:"total_views"`
	ActiveUser int64  `json:"active_users"`
	NewUsers   int64  `json:"new_users"`
}

// TopVideo 热门视频
type TopVideo struct {
	VideoID   uint64 `json:"video_id"`
	ViewCount int64  `json:"view_count"`
}

// RegionStat 地域统计
type RegionStat struct {
	Province   string `json:"province"`
	ViewCount  int64  `json:"view_count"`
	ViewUsers  int64  `json:"view_users"`
}

// GetPlatformDailyStats 平台日统计（直接查 view_logs + users 实时聚合）
func GetPlatformDailyStats(entityIDs []uint64, days int) ([]PlatformDailyStat, error) {
	var stats []PlatformDailyStat
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	db := database.DB.Table("view_logs").
		Select("DATE(created_at) as stat_date, COUNT(1) as total_views, COUNT(DISTINCT user_id) as active_user").
		Where("created_at >= ?", since)
	if len(entityIDs) > 0 {
		db = db.Where("entity_id IN ?", entityIDs)
	}
	db = db.Group("DATE(created_at)").Order("stat_date ASC")

	err := db.Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	// 补充新用户数
	for i := range stats {
		var newUsers int64
		database.DB.Table("users").
			Where("DATE(created_at) = ?", stats[i].StatDate).Count(&newUsers)
		stats[i].NewUsers = newUsers
	}
	return stats, nil
}

// GetTopVideos 热门视频排行
func GetTopVideos(entityIDs []uint64, limit int) ([]TopVideo, error) {
	var videos []TopVideo
	since := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	db := database.DB.Table("view_logs").
		Select("video_id, COUNT(1) as view_count").
		Where("created_at >= ?", since)
	if len(entityIDs) > 0 {
		db = db.Where("entity_id IN ?", entityIDs)
	}
	err := db.Group("video_id").Order("view_count DESC").Limit(limit).Scan(&videos).Error
	return videos, err
}

// GetRegionStats 地域分布统计
func GetRegionStats(entityIDs []uint64, days int) ([]RegionStat, error) {
	var regions []RegionStat
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	db := database.DB.Table("view_logs").
		Select("province, COUNT(1) as view_count, COUNT(DISTINCT user_id) as view_users").
		Where("created_at >= ? AND province != '' AND province != '未知'", since)
	if len(entityIDs) > 0 {
		db = db.Where("entity_id IN ?", entityIDs)
	}
	err := db.Group("province").Order("view_count DESC").Scan(&regions).Error
	return regions, err
}

// GetTodaySummary 今日汇总（指标卡）
func GetTodaySummary(entityIDs []uint64) map[string]int64 {
	today := time.Now().Format("2006-01-02")

	var views, activeUsers, newUsers, totalUsers int64

	// 今日播放
	db := database.DB.Table("view_logs").Where("DATE(created_at) = ?", today)
	if len(entityIDs) > 0 {
		db = db.Where("entity_id IN ?", entityIDs)
	}
	db.Count(&views)

	// 今日活跃
	db2 := database.DB.Table("view_logs").
		Select("COUNT(DISTINCT user_id)").
		Where("DATE(created_at) = ?", today)
	if len(entityIDs) > 0 {
		db2 = db2.Where("entity_id IN ?", entityIDs)
	}
	db2.Scan(&activeUsers)

	// 今日新增
	database.DB.Table("users").Where("DATE(created_at) = ?", today).Count(&newUsers)

	// 总用户
	database.DB.Table("users").Count(&totalUsers)

	return map[string]int64{
		"views":       views,
		"users":       activeUsers,
		"new_users":   newUsers,
		"total_users": totalUsers,
	}
}
