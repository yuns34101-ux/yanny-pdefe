package handler

import (
	"yanny-service/internal/dto"
	"yanny-service/internal/middleware"
	"yanny-service/internal/repository"

	"github.com/gin-gonic/gin"
)

// MpReportView 上报视频播放事件（防刷由 VideoViewAntiAbuse 中间件保证）
func MpReportView(c *gin.Context) {
	var req struct {
		VideoID       uint64 `json:"video_id" binding:"required"`
		WatchDuration uint   `json:"watch_duration"`
		IsComplete    int8   `json:"is_complete"`
		Source        string `json:"source"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}

	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	if skipped, _ := c.Get("skip_view_count"); skipped == true {
		dto.Success(c, nil)
		return
	}

	log := repository.BuildViewLog(mpAccountID, userID, req.VideoID, req.WatchDuration, req.IsComplete, req.Source, c.ClientIP())
	if err := repository.InsertViewLog(log); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}

	_ = repository.IncrementVideoViewCount(req.VideoID)

	dto.Success(c, nil)
}

// MpReportAction 上报行为事件
func MpReportAction(c *gin.Context) {
	var req struct {
		EventType  string `json:"event_type" binding:"required"`
		TargetType string `json:"target_type"`
		TargetID   uint64 `json:"target_id"`
		PagePath   string `json:"page_path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}

	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	log := repository.BuildActionLog(mpAccountID, userID, req.EventType, req.TargetType, req.TargetID, req.PagePath, c.ClientIP())
	_ = repository.InsertActionLog(log)

	dto.Success(c, nil)
}
