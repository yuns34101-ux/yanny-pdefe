package handler

import (
	"errors"
	"log"
	"yanny-service/internal/dto"
	"yanny-service/internal/middleware"
	"yanny-service/internal/repository"
	"yanny-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ========== 主体信息 ==========

// MpGetEntity 小程序端获取主体详情及统计信息（游客可看）
func MpGetEntity(c *gin.Context) {
	entityID, _ := parseUintQuery(c, "entity_id")
	mpAccountID := middleware.GetMpAccountID(c)
	if mpAccountID == 0 {
		mpAccountID, _ = parseUintQuery(c, "mp_account_id")
	}
	userID := middleware.GetUserID(c)
	entity, err := service.GetEntityForMp(entityID, mpAccountID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "主体不存在" {
			dto.Error(c, dto.ErrCodeEntityNotFound, "主体不存在")
			return
		}
		log.Printf("MpGetEntity 内部错误: entity_id=%d mp_account_id=%d err=%v", entityID, mpAccountID, err)
		dto.Error(c, dto.ErrCodeInternal, "服务器内部错误")
		return
	}
	dto.Success(c, entity)
}

// ========== 评论 ==========

// MpCreateComment 发表评论（需登录）
func MpCreateComment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	comment, err := service.CreateComment(mpAccountID, req.VideoID, userID, req.Content, req.ParentID, req.ReplyToUserID)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, comment)
}

// MpListComments 评论列表（游客可看）
func MpListComments(c *gin.Context) {
	videoID, ok := parseUintQuery(c, "video_id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	var p dto.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p = dto.DefaultPagination()
	}
	comments, total, err := repository.ListRootComments(videoID, p.Page, p.PageSize)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	// 加载前 3 条回复
	for i := range comments {
		replies, _, _ := repository.ListReplies(comments[i].ID, 3)
		comments[i].Replies = replies
	}
	dto.SuccessPage(c, comments, p.Page, p.PageSize, total)
}

// MpListReplies 加载更多回复
func MpListReplies(c *gin.Context) {
	commentID, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	var p dto.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p = dto.DefaultPagination()
	}
	replies, _, err := repository.ListReplies(commentID, p.PageSize)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, replies)
}

// ========== 点赞 ==========

// MpToggleLike 切换点赞（需登录）
func MpToggleLike(c *gin.Context) {
	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint64 `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	liked, err := service.ToggleLike(mpAccountID, userID, req.TargetType, req.TargetID)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, gin.H{"liked": liked})
}

// ========== 收藏 ==========

// MpToggleFavorite 切换收藏（需登录）
func MpToggleFavorite(c *gin.Context) {
	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	var req struct {
		VideoID uint64 `json:"video_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	favored, err := service.ToggleFavorite(mpAccountID, userID, req.VideoID)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, gin.H{"favored": favored})
}

// MpMyFavorites 我的收藏列表（返回视频完整信息，供首页三列宫格复用）
func MpMyFavorites(c *gin.Context) {
	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	var p dto.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p = dto.DefaultPagination()
	}
	favs, total, err := service.GetFavoriteVideosForMp(mpAccountID, userID, p.Page, p.PageSize)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.SuccessPage(c, favs, p.Page, p.PageSize, total)
}

// MpMyHistory 观看历史列表（按视频去重，取最近一次播放，返回视频完整信息供首页三列宫格复用）
func MpMyHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	var p dto.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p = dto.DefaultPagination()
	}
	history, total, err := service.GetHistoryVideosForMp(mpAccountID, userID, p.Page, p.PageSize)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.SuccessPage(c, history, p.Page, p.PageSize, total)
}

// ========== 分享 ==========

// MpRecordShare 记录分享（需登录）
func MpRecordShare(c *gin.Context) {
	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	var req struct {
		VideoID   uint64 `json:"video_id" binding:"required"`
		ShareType string `json:"share_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	if err := service.RecordShare(mpAccountID, userID, req.VideoID, req.ShareType); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

// ========== 关注 ==========

// MpToggleFollow 切换关注（需登录）
func MpToggleFollow(c *gin.Context) {
	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	var req struct {
		EntityID uint64 `json:"entity_id"` // 0 时兜底为当前小程序账号绑定的默认主体
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	followed, err := service.ToggleFollow(mpAccountID, userID, req.EntityID)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, gin.H{"followed": followed})
}

// ========== 用户互动状态查询 ==========

// MpInteractionStatus 查询用户对视频的互动状态（点赞/收藏/是否已分享）
func MpInteractionStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	mpAccountID := middleware.GetMpAccountID(c)
	videoID, ok := parseUintQuery(c, "video_id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	status := service.GetInteractionStatus(mpAccountID, userID, videoID)
	dto.Success(c, status)
}
