package service

import (
	"errors"
	"time"
	"yanny-service/internal/config"
	"yanny-service/internal/model"
	"yanny-service/internal/repository"
	"yanny-service/internal/util/qiniu"

	"gorm.io/gorm"
)

// 视频处理规则
// 七牛存储桶未开通实时转码（fop）服务，开启会导致请求被拒绝（this fop is blocked, please use pfop service）
// 暂时注释，直接签名原片 URL 播放
// const videoProcessRule = "avthumb/mp4/s/1280x720/vb/2000k/ab/128k/acodec/aac"

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

// signVideoURLs 给视频的封面和视频 URL 加七牛签名 + 统一处理规则（admin/mp 共用）
// 封面 → SignImageURL（与其他图片统一），视频 → SignURL + 转码规则
func signVideoURLs(v *model.Video) {
	if v.CoverURL != "" {
		v.CoverURL = qiniu.SignImageURL(v.CoverURL)
	}
	if v.VideoURL != "" {
		ttl := time.Duration(config.AppConfig.Qiniu.GetMediaTTL()) * time.Hour
		v.VideoURL = qiniu.SignURL(qiniu.EnsureProtocol(v.VideoURL), "", ttl)
	}
}

// signVideoListURLs 给视频列表批量加签名
func signVideoListURLs(videos []model.Video) {
	for i := range videos {
		signVideoURLs(&videos[i])
	}
}

// signVideoListCoverOnly 列表场景只签封面，video_url 留空
// 列表页（首页宫格、播放页滑动导航）不会直接播放，video_url 由播放页按需通过 GetVideoDetail 单独拉取
func signVideoListCoverOnly(videos []model.Video) {
	for i := range videos {
		if videos[i].CoverURL != "" {
			videos[i].CoverURL = qiniu.SignImageURL(videos[i].CoverURL)
		}
		videos[i].VideoURL = ""
	}
}

// GetVideosForMp 小程序端获取视频列表（不含 video_url，需播放时用 GetVideoDetail 单独获取）
func GetVideosForMp(mpAccountID, entityID, categoryID uint64, page, pageSize int) ([]model.Video, int64, error) {
	videos, total, err := repository.ListVideosForMp(mpAccountID, entityID, categoryID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	signVideoListCoverOnly(videos)
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

// GetFavoriteVideosForMp 我的收藏视频列表
func GetFavoriteVideosForMp(mpAccountID, userID uint64, page, pageSize int) ([]model.Video, int64, error) {
	videos, total, err := repository.ListFavoriteVideosByUser(mpAccountID, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	signVideoListURLs(videos)
	return videos, total, nil
}

// GetHistoryVideosForMp 观看历史视频列表（按视频去重，取最近一次播放）
func GetHistoryVideosForMp(mpAccountID, userID uint64, page, pageSize int) ([]model.Video, int64, error) {
	videos, total, err := repository.ListHistoryVideosByUser(mpAccountID, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	signVideoListURLs(videos)
	return videos, total, nil
}
