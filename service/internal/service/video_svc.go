package service

import (
	"errors"
	"yanny-service/internal/model"
	"yanny-service/internal/repository"
	"yanny-service/internal/util/qiniu"

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

// ensureProtocol 补全缺少的 https:// 前缀
func ensureProtocol(rawURL string) string {
	if len(rawURL) > 0 && !(len(rawURL) >= 7 && rawURL[:7] == "http://") &&
		!(len(rawURL) >= 8 && rawURL[:8] == "https://") {
		return "https://" + rawURL
	}
	return rawURL
}

// signVideoURLs 给视频的封面和视频 URL 加七牛签名（仅对七牛域名生效）
func signVideoURLs(v *model.Video) {
	if v.CoverURL != "" {
		v.CoverURL = qiniu.GenerateCoverURL(ensureProtocol(v.CoverURL))
	}
	if v.VideoURL != "" {
		v.VideoURL = qiniu.GenerateVideoURL(ensureProtocol(v.VideoURL))
	}
}

// signVideoListURLs 给视频列表批量加签名
func signVideoListURLs(videos []model.Video) {
	for i := range videos {
		signVideoURLs(&videos[i])
	}
}

// GetVideosForMp 小程序端获取视频列表
func GetVideosForMp(mpAccountID, entityID, categoryID uint64, page, pageSize int) ([]model.Video, int64, error) {
	videos, total, err := repository.ListVideosForMp(mpAccountID, entityID, categoryID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	signVideoListURLs(videos)
	return videos, total, nil
}

// GetVideosForAdmin 管理后台获取视频列表
func GetVideosForAdmin(mpAccountID, entityID, categoryID uint64, keyword string, status *int8, page, pageSize int, entityIDs []uint64) ([]model.Video, int64, error) {
	videos, total, err := repository.ListVideosForAdmin(mpAccountID, entityID, categoryID, keyword, status, page, pageSize, entityIDs)
	if err != nil {
		return nil, 0, err
	}
	signVideoListURLs(videos)
	return videos, total, nil
}

// GetVideoDetail 获取视频详情
func GetVideoDetail(videoID uint64) (*model.Video, error) {
	v, err := repository.FindVideoByID(videoID)
	if err != nil {
		return nil, err
	}
	signVideoURLs(v)
	return v, nil
}
