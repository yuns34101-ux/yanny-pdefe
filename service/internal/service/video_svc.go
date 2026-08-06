package service

import (
	"errors"
	"yanny-service/internal/model"
	"yanny-service/internal/repository"

	"gorm.io/gorm"
)

// CreateCategory 创建分类
func CreateCategory(entityID, mpAccountID uint64, name, iconURL string, sortOrder int) (*model.VideoCategory, error) {
	cat := &model.VideoCategory{
		EntityID:    entityID,
		MpAccountID: mpAccountID,
		Name:        name,
		IconURL:     iconURL,
		SortOrder:   sortOrder,
		Status:      1,
	}
	if err := repository.CreateCategory(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

// CreateVideo 创建视频
func CreateVideo(req struct {
	MpAccountID   uint64
	EntityID      uint64
	CategoryID    uint64
	Title         string
	Description   string
	CoverURL      string
	VideoURL      string
	Duration      uint
	Width         uint
	Height        uint
	FileSize      uint64
	Tags          string
	Status        int8
	IsRecommended int8
}) (*model.Video, error) {
	video := &model.Video{
		MpAccountID:   req.MpAccountID,
		EntityID:      req.EntityID,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Description:   req.Description,
		CoverURL:      req.CoverURL,
		VideoURL:      req.VideoURL,
		Duration:      req.Duration,
		Width:         req.Width,
		Height:        req.Height,
		FileSize:      req.FileSize,
		Tags:          req.Tags,
		Status:        req.Status,
		IsRecommended: req.IsRecommended,
	}
	if err := repository.CreateVideo(video); err != nil {
		return nil, err
	}
	return video, nil
}

// UpdateVideo 更新视频
func UpdateVideo(id uint64, updates map[string]interface{}) error {
	if _, err := repository.FindVideoByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("视频不存在")
		}
		return err
	}
	return repository.UpdateVideo(id, updates)
}

// GetVideosForMp 小程序端获取视频列表
func GetVideosForMp(mpAccountID, entityID, categoryID uint64, page, pageSize int) ([]model.Video, int64, error) {
	return repository.ListVideosForMp(mpAccountID, entityID, categoryID, page, pageSize)
}

// GetVideosForAdmin 管理后台获取视频列表
func GetVideosForAdmin(mpAccountID, entityID, categoryID uint64, keyword string, status *int8, page, pageSize int) ([]model.Video, int64, error) {
	return repository.ListVideosForAdmin(mpAccountID, entityID, categoryID, keyword, status, page, pageSize)
}

// GetVideoDetail 获取视频详情
func GetVideoDetail(videoID uint64) (*model.Video, error) {
	return repository.FindVideoByID(videoID)
}
