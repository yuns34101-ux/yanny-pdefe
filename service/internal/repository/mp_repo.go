package repository

import (
	"yanny-service/internal/database"
	"yanny-service/internal/model"
)

// ========== 评论 ==========

func CreateComment(c *model.Comment) error {
	return database.DB.Create(c).Error
}

func FindCommentByID(id uint64) (*model.Comment, error) {
	var c model.Comment
	err := database.DB.First(&c, id).Error
	return &c, err
}

func UpdateCommentRootID(commentID, rootID uint64) error {
	return database.DB.Model(&model.Comment{}).Where("id = ?", commentID).
		Update("root_id", rootID).Error
}

func ListRootComments(videoID uint64, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64
	db := database.DB.Model(&model.Comment{}).
		Where("video_id = ? AND parent_id IS NULL AND status = 1", videoID)
	db.Count(&total)
	err := db.Preload("User").Preload("ReplyToUser").
		Order("is_top DESC, created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&comments).Error
	return comments, total, err
}

func ListReplies(commentID uint64, limit int) ([]model.Comment, int64, error) {
	// commentID 作为 root_id 来查所有二级回复
	var replies []model.Comment
	var total int64
	db := database.DB.Model(&model.Comment{}).
		Where("root_id = ? AND parent_id IS NOT NULL AND status = 1", commentID)
	db.Count(&total)
	err := db.Preload("User").Preload("ReplyToUser").
		Order("created_at ASC").Limit(limit).Find(&replies).Error
	return replies, total, err
}

func IncrementCommentReplyCount(commentID uint64) error {
	return database.DB.Model(&model.Comment{}).Where("id = ?", commentID).
		UpdateColumn("reply_count", database.DB.Raw("reply_count + 1")).Error
}

func IncrementCommentLikeCount(commentID uint64) error {
	return database.DB.Model(&model.Comment{}).Where("id = ?", commentID).
		UpdateColumn("like_count", database.DB.Raw("like_count + 1")).Error
}

// ========== 点赞 ==========

func FindLike(mpAccountID, userID uint64, targetType string, targetID uint64) (*model.Like, error) {
	var like model.Like
	err := database.DB.Where("mp_account_id = ? AND user_id = ? AND target_type = ? AND target_id = ?",
		mpAccountID, userID, targetType, targetID).First(&like).Error
	return &like, err
}

func CreateLike(like *model.Like) error {
	return database.DB.Create(like).Error
}

func UpdateLikeStatus(id uint64, status int8) error {
	return database.DB.Model(&model.Like{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status}).Error
}

// ========== 收藏 ==========

func FindFavorite(mpAccountID, userID, videoID uint64) (*model.Favorite, error) {
	var fav model.Favorite
	err := database.DB.Where("mp_account_id = ? AND user_id = ? AND video_id = ?",
		mpAccountID, userID, videoID).First(&fav).Error
	return &fav, err
}

func CreateFavorite(fav *model.Favorite) error {
	return database.DB.Create(fav).Error
}

func UpdateFavoriteStatus(id uint64, status int8) error {
	return database.DB.Model(&model.Favorite{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status}).Error
}

func ListFavoritesByUser(mpAccountID, userID uint64, page, pageSize int) ([]model.Favorite, int64, error) {
	var favs []model.Favorite
	var total int64
	db := database.DB.Model(&model.Favorite{}).
		Where("mp_account_id = ? AND user_id = ? AND status = 1", mpAccountID, userID)
	db.Count(&total)
	err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&favs).Error
	return favs, total, err
}

// ========== 分享 ==========

func CreateShare(share *model.Share) error {
	return database.DB.Create(share).Error
}

// ========== 关注 ==========

func FindFollow(mpAccountID, userID, entityID uint64) (*model.Follow, error) {
	var follow model.Follow
	err := database.DB.Where("mp_account_id = ? AND user_id = ? AND entity_id = ?",
		mpAccountID, userID, entityID).First(&follow).Error
	return &follow, err
}

func CreateFollow(follow *model.Follow) error {
	return database.DB.Create(follow).Error
}

func UpdateFollowStatus(id uint64, status int8) error {
	return database.DB.Model(&model.Follow{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status}).Error
}

func CountFollowersByEntity(entityID uint64) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Follow{}).
		Where("entity_id = ? AND status = 1", entityID).Count(&count).Error
	return count, err
}

// ========== 埋点 ==========

func BuildViewLog(mpAccountID uint64, userID uint64, videoID uint64, watchDuration uint, isComplete int8, source, ip string) *model.ViewLog {
	log := &model.ViewLog{
		MpAccountID:   mpAccountID,
		VideoID:       videoID,
		WatchDuration: watchDuration,
		IsComplete:    isComplete,
		Source:        source,
		IP:            ip,
		Province:      "未知",
		City:          "未知",
	}
	if userID > 0 {
		log.UserID = &userID
	}
	return log
}

func InsertViewLog(log *model.ViewLog) error {
	return database.DB.Create(log).Error
}

func BuildActionLog(mpAccountID uint64, userID uint64, eventType, targetType string, targetID uint64, pagePath, ip string) *model.ActionLog {
	log := &model.ActionLog{
		MpAccountID: mpAccountID,
		EventType:   eventType,
		TargetType:  targetType,
		TargetID:    targetID,
		PagePath:    pagePath,
		IP:          ip,
		Province:    "未知",
		City:        "未知",
	}
	if userID > 0 {
		log.UserID = &userID
	}
	return log
}

func InsertActionLog(log *model.ActionLog) error {
	return database.DB.Create(log).Error
}
