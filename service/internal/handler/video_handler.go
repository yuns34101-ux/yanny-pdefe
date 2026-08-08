package handler

import (
	"yanny-service/internal/dto"
	"yanny-service/internal/repository"
	"yanny-service/internal/service"
	"yanny-service/internal/util/qiniu"

	"github.com/gin-gonic/gin"
)

// ========== 分类管理 ==========

func CreateCategory(c *gin.Context) {
	var req struct {
		EntityID    uint64 `json:"entity_id" binding:"required"`
		MpAccountID uint64 `json:"mp_account_id" binding:"required"`
		Name        string `json:"name" binding:"required,max=50"`
		IconURL     string `json:"icon_url"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	cat, err := service.CreateCategory(req.EntityID, req.MpAccountID,
		req.Name, req.IconURL, req.SortOrder)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, cat)
}

func ListCategories(c *gin.Context) {
	entityID, _ := parseUintQuery(c, "entity_id")
	mpAccountID, _ := parseUintQuery(c, "mp_account_id")
	categories, err := repository.ListCategories(entityID, mpAccountID)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, categories)
}

// MpListCategories 小程序端分类列表（entityID=0 时解析默认主体）
func MpListCategories(c *gin.Context) {
	entityID, _ := parseUintQuery(c, "entity_id")
	mpAccountID := c.GetUint64("mp_account_id")
	if mpAccountID == 0 {
		mpAccountID, _ = parseUintQuery(c, "mp_account_id")
	}
	if entityID == 0 {
		if resolved, err := service.ResolveEntityID(0, mpAccountID); err == nil {
			entityID = resolved
		}
	}
	categories, err := repository.ListCategories(entityID, mpAccountID)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, categories)
}

// ========== 视频管理（管理后台） ==========

func CreateVideo(c *gin.Context) {
	var req dto.CreateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效："+err.Error())
		return
	}
	video, err := service.CreateVideo(struct {
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
	}{
		MpAccountID:   req.MpAccountID,
		EntityID:      req.EntityID,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Description:   req.Description,
		CoverURL:      qiniu.StripQuery(req.CoverURL),
		VideoURL:      qiniu.StripQuery(req.VideoURL),
		Duration:      req.Duration,
		Width:         req.Width,
		Height:        req.Height,
		FileSize:      req.FileSize,
		Tags:          req.Tags,
		Status:        req.Status,
		IsRecommended: req.IsRecommended,
	})
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, video)
}

func ListVideos(c *gin.Context) {
	var q struct {
		dto.Pagination
		MpAccountID uint64 `form:"mp_account_id"`
		EntityID    uint64 `form:"entity_id"`
		CategoryID  uint64 `form:"category_id"`
		Keyword     string `form:"keyword"`
		Status      *int8  `form:"status"`
	}
	if err := c.ShouldBindQuery(&q); err != nil {
		q.Page, q.PageSize = 1, 20
	}
	videos, total, err := service.GetVideosForAdmin(q.MpAccountID, q.EntityID,
		q.CategoryID, q.Keyword, q.Status, q.Page, q.PageSize, getEntityScope(c))
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.SuccessPage(c, videos, q.Page, q.PageSize, total)
}

func UpdateVideo(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	var req dto.UpdateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.CoverURL != "" {
		updates["cover_url"] = qiniu.StripQuery(req.CoverURL)
	}
	if req.VideoURL != "" {
		updates["video_url"] = qiniu.StripQuery(req.VideoURL)
	}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}
	updates["is_recommended"] = req.IsRecommended

	if err := service.UpdateVideo(id, updates); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

func DeleteVideo(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	if err := repository.DeleteVideo(id); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

// ========== 小程序端视频接口 ==========

// MpGetVideos 小程序端视频列表（entityID=0 时解析当前小程序默认主体）
func MpGetVideos(c *gin.Context) {
	var q dto.VideoListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		q.Page, q.PageSize = 1, 20
	}
	mpAccountID := c.GetUint64("mp_account_id")
	if mpAccountID == 0 {
		mpAccountID, _ = parseUintQuery(c, "mp_account_id")
	}
	if mpAccountID == 0 {
		mpAccountID = 1
	}
	entityID := q.EntityID
	if entityID == 0 {
		entityID, _ = parseUintQuery(c, "entity_id")
	}
	// entityID 为 0 时解析默认主体，避免列出其他主体的视频
	if entityID == 0 {
		if resolved, err := service.ResolveEntityID(0, mpAccountID); err == nil {
			entityID = resolved
		}
	}
	videos, total, err := service.GetVideosForMp(mpAccountID, entityID, q.CategoryID, q.Page, q.PageSize)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.SuccessPage(c, videos, q.Page, q.PageSize, total)
}

// MpGetVideoDetail 小程序端视频详情
func MpGetVideoDetail(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeVideoNotFound, "视频不存在")
		return
	}
	video, err := service.GetVideoDetail(id)
	if err != nil {
		dto.Error(c, dto.ErrCodeVideoNotFound, "视频不存在")
		return
	}
	if video.Status != 1 {
		dto.Error(c, dto.ErrCodeVideoOffline, "视频已下架")
		return
	}
	dto.Success(c, video)
}
