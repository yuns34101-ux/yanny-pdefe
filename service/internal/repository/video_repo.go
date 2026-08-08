package repository

import (
	"log"

	"yanny-service/internal/database"
	"yanny-service/internal/model"
)

// ========== 分类 ==========

func ListCategories(entityID, mpAccountID uint64) ([]model.VideoCategory, error) {
	var categories []model.VideoCategory
	db := database.DB.Model(&model.VideoCategory{}).Where("status = 1")
	if entityID > 0 {
		db = db.Where("entity_id = ?", entityID)
	}
	if mpAccountID > 0 {
		db = db.Where("mp_account_id = ?", mpAccountID)
	}
	err := db.Order("sort_order DESC").Find(&categories).Error
	return categories, err
}

func CreateCategory(cat *model.VideoCategory) error {
	return database.DB.Create(cat).Error
}

func UpdateCategory(id uint64, updates map[string]interface{}) error {
	return database.DB.Model(&model.VideoCategory{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteCategory(id uint64) error {
	return database.DB.Delete(&model.VideoCategory{}, id).Error
}

// ========== 视频 ==========

func FindVideoByID(id uint64) (*model.Video, error) {
	var video model.Video
	err := database.DB.First(&video, id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// ListVideosForMp 小程序端视频列表（已发布、按分类筛选、按发布时间倒序）
func ListVideosForMp(mpAccountID, entityID, categoryID uint64, page, pageSize int) ([]model.Video, int64, error) {
	var videos []model.Video
	var total int64
	db := database.DB.Model(&model.Video{}).
		Where("mp_account_id = ? AND status = 1", mpAccountID)
	if entityID > 0 {
		db = db.Where("entity_id = ?", entityID)
	}
	if categoryID > 0 {
		db = db.Where("category_id = ?", categoryID)
	}
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("is_recommended DESC, published_at DESC").Find(&videos).Error
	return videos, total, err
}

// ListVideosForAdmin 管理后台视频列表（支持 entityIDs 数据级过滤）
func ListVideosForAdmin(mpAccountID, entityID, categoryID uint64, keyword string, status *int8, page, pageSize int, entityIDs []uint64) ([]model.Video, int64, error) {
	var videos []model.Video
	var total int64
	db := database.DB.Model(&model.Video{})
	if len(entityIDs) > 0 {
		db = db.Where("entity_id IN ?", entityIDs)
	}
	if mpAccountID > 0 {
		db = db.Where("mp_account_id = ?", mpAccountID)
	}
	if entityID > 0 {
		db = db.Where("entity_id = ?", entityID)
	}
	if categoryID > 0 {
		db = db.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		db = db.Where("title LIKE ?", "%"+keyword+"%")
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id DESC").Find(&videos).Error
	return videos, total, err
}

func CreateVideo(video *model.Video) error {
	return database.DB.Create(video).Error
}

func UpdateVideo(id uint64, updates map[string]interface{}) error {
	return database.DB.Model(&model.Video{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteVideo(id uint64) error {
	return database.DB.Delete(&model.Video{}, id).Error
}

// IncrementVideoViewCount 增加视频播放量（原生 SQL，绕过 GORM 表达式层）
func IncrementVideoViewCount(videoID uint64) error {
	err := database.DB.Exec("UPDATE videos SET view_count = view_count + 1 WHERE id = ?", videoID).Error
	if err != nil {
		log.Printf("IncrementVideoViewCount 失败: id=%d err=%v", videoID, err)
	}
	return err
}

// UpdateVideoCount 更新视频计数器（原生 SQL Exec，绕过 GORM Raw()/clause.Expr 兼容性问题）
// field 参数来自调用方代码常量，非用户输入，无 SQL 注入风险
func UpdateVideoCount(videoID uint64, field string, delta int) error {
	err := database.DB.Exec(
		"UPDATE videos SET "+field+" = "+field+" + ? WHERE id = ?",
		delta, videoID,
	).Error
	if err != nil {
		log.Printf("UpdateVideoCount 失败: id=%d field=%s delta=%d err=%v", videoID, field, delta, err)
	}
	return err
}

// CountVideosByEntity 统计主体下已发布视频数
func CountVideosByEntity(entityID uint64) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Video{}).
		Where("entity_id = ? AND status = 1", entityID).Count(&count).Error
	return count, err
}

// ListFavoriteVideosByUser 我的收藏视频列表（关联 favorites 表，按收藏时间倒序）
func ListFavoriteVideosByUser(mpAccountID, userID uint64, page, pageSize int) ([]model.Video, int64, error) {
	var videos []model.Video
	var total int64
	db := database.DB.Model(&model.Video{}).
		Joins("JOIN favorites ON favorites.video_id = videos.id").
		Where("favorites.mp_account_id = ? AND favorites.user_id = ? AND favorites.status = 1", mpAccountID, userID)
	db.Count(&total)
	err := db.Order("favorites.created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&videos).Error
	return videos, total, err
}

// ListHistoryVideosByUser 观看历史视频列表（关联 view_logs 表，按视频去重取最近一次播放时间倒序）
func ListHistoryVideosByUser(mpAccountID, userID uint64, page, pageSize int) ([]model.Video, int64, error) {
	var videos []model.Video
	var total int64
	latest := database.DB.Table("view_logs").
		Select("video_id, MAX(created_at) as last_watched_at").
		Where("mp_account_id = ? AND user_id = ?", mpAccountID, userID).
		Group("video_id")
	db := database.DB.Model(&model.Video{}).
		Joins("JOIN (?) AS vl ON vl.video_id = videos.id", latest)
	db.Count(&total)
	err := db.Order("vl.last_watched_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&videos).Error
	return videos, total, err
}
