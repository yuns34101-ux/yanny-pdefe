package handler

import (
	"crypto/md5"
	"fmt"
	"yanny-service/internal/config"
	"yanny-service/internal/dto"
	"yanny-service/internal/middleware"
	"yanny-service/internal/repository"

	"github.com/gin-gonic/gin"
)

// MpReportView 上报视频播放事件（签名校验防伪造）
func MpReportView(c *gin.Context) {
	var req struct {
		VideoID       uint64 `json:"video_id" binding:"required"`
		WatchDuration uint   `json:"watch_duration"`
		IsComplete    int8   `json:"is_complete"`
		Source        string `json:"source"`
		Timestamp     int64  `json:"timestamp" binding:"required"`
		Nonce         string `json:"nonce" binding:"required"`
		Sign          string `json:"sign" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}

	// 签名校验：sign = MD5(video_id + timestamp + nonce + secret)
	if !validateSign(req.VideoID, req.Timestamp, req.Nonce, req.Sign) {
		dto.Error(c, dto.ErrCodeForbidden, "签名校验失败")
		return
	}

	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	if skipped, _ := c.Get("skip_view_count"); skipped == true {
		dto.Success(c, nil) // 去重，静默忽略
		return
	}

	log := repository.BuildViewLog(mpAccountID, userID, req.VideoID, req.WatchDuration, req.IsComplete, req.Source, c.ClientIP())
	if err := repository.InsertViewLog(log); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}

	// 同步更新视频表的播放次数（view_logs 只存明细，videos.view_count 供列表排序/展示）
	_ = repository.IncrementVideoViewCount(req.VideoID)

	dto.Success(c, nil)
}

// MpReportAction 上报行为事件（签名校验）
func MpReportAction(c *gin.Context) {
	var req struct {
		EventType  string `json:"event_type" binding:"required"`
		TargetType string `json:"target_type"`
		TargetID   uint64 `json:"target_id"`
		PagePath   string `json:"page_path"`
		Timestamp  int64  `json:"timestamp" binding:"required"`
		Nonce      string `json:"nonce" binding:"required"`
		Sign       string `json:"sign" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}

	if !validateSign(req.TargetID, req.Timestamp, req.Nonce, req.Sign) {
		dto.Error(c, dto.ErrCodeForbidden, "签名校验失败")
		return
	}

	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	log := repository.BuildActionLog(mpAccountID, userID, req.EventType, req.TargetType, req.TargetID, req.PagePath, c.ClientIP())
	_ = repository.InsertActionLog(log) // 异步写入，降级不阻塞

	dto.Success(c, nil)
}

// validateSign 校验客户端上报签名
// 算法：MD5(fmt.Sprintf("%d:%d:%s:%s", id, timestamp, nonce, secret))
func validateSign(id uint64, timestamp int64, nonce, sign string) bool {
	secret := config.AppConfig.JWT.Secret
	raw := fmt.Sprintf("%d:%d:%s:%s", id, timestamp, nonce, secret)
	expected := fmt.Sprintf("%x", md5.Sum([]byte(raw)))
	return expected == sign
}
