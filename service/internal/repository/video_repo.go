package repository

import (
	"yanny-service/internal/database"
	"yanny-service/internal/model"
)

// ========== 分类 ==========

func ListCategories(entityID, mpAccountID uint64) ([]model.VideoCategory, error) {
	var categories []model.VideoCategory
	err := database.DB.Where("entity_id = ? AND mp_account_id = ? AND status = 1", entityID, mpAccountID).
		Order("sort_order DESC").Find(&categories).Error
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

// ListVideosForAdmin 管理后台视频列表
func ListVideosForAdmin(mpAccountID, entityID, categoryID uint64, keyword string, status *int8, page, pageSize int) ([]model.Video, int64, error) {
	var videos []model.Video
	var total int64
	db := database.DB.Model(&model.Video{})
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

// IncrementVideoViewCount 增加视频播放量
func IncrementVideoViewCount(videoID uint64) error {
	return database.DB.Model(&model.Video{}).Where("id = ?", videoID).
		UpdateColumn("view_count", database.DB.Raw("view_count + 1")).Error
}

// UpdateVideoCount 更新视频计数器
func UpdateVideoCount(videoID uint64, field string, delta int) error {
	return database.DB.Model(&model.Video{}).Where("id = ?", videoID).
		UpdateColumn(field, database.DB.Raw(field+" + ?", delta)).Error
}
