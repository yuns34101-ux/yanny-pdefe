package service

import (
	"errors"
	"time"
	"yanny-service/internal/middleware"
	"yanny-service/internal/model"
	"yanny-service/internal/repository"

	"gorm.io/gorm"
)

// FindOrCreateUser 查找或创建用户
// inviterUserID 为邀请人用户 ID（分享裂变），仅在创建新用户时生效，已存在用户不覆盖归属关系
// sessionKey 为微信 code2Session 返回的会话密钥，每次登录都会轮换，用于后续手机号解密
func FindOrCreateUser(mpAccountID uint64, openid, unionid, nickname, avatarURL, ip string, inviterUserID uint64, sessionKey string) (*model.User, bool, error) {
	user, err := repository.FindUserByOpenid(mpAccountID, openid)
	if err == nil {
		// 用户已存在，更新登录信息 + 轮换 session_key
		now := time.Now()
		repository.UpdateUserLoginInfo(user.ID, ip, &now)
		if sessionKey != "" {
			repository.UpdateUserSessionKey(user.ID, sessionKey)
			user.SessionKey = sessionKey
		}
		return user, false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	// 创建新用户
	now := time.Now()
	user = &model.User{
		MpAccountID: mpAccountID,
		Openid:      openid,
		Unionid:     unionid,
		SessionKey:  sessionKey,
		Nickname:    nickname,
		AvatarURL:   avatarURL,
		Status:      1,
		LastLoginAt: &now,
		LastLoginIP: ip,
	}
	if inviterUserID > 0 && inviterUserID != user.ID {
		user.InviterUserID = &inviterUserID
	}
	if err := repository.CreateUser(user); err != nil {
		return nil, false, err
	}
	return user, true, nil
}

// FindMpAccountByAppID 根据 AppID 查找小程序账号
func FindMpAccountByAppID(appID string) (*model.MpAccount, error) {
	return repository.FindMpAccountByAppID(appID)
}

// GenerateMpToken 生成小程序端 JWT
func GenerateMpToken(userID, mpAccountID uint64) (string, error) {
	return middleware.GenerateMpToken(userID, mpAccountID)
}

// CreateComment 创建评论（一级或回复）
func CreateComment(mpAccountID, videoID, userID uint64, content string, parentID, replyToUserID *uint64) (*model.Comment, error) {
	comment := &model.Comment{
		MpAccountID:   mpAccountID,
		VideoID:       videoID,
		UserID:        userID,
		ParentID:      parentID,
		ReplyToUserID: replyToUserID,
		Content:       content,
		Status:        1,
	}
	// 一级评论：root_id 为 NULL，写入后回填自己的 ID
	// 二级回复：root_id 为父评论的 root_id（或 parent_id 自己如果是一级）
	if parentID != nil && *parentID > 0 {
		parent, err := repository.FindCommentByID(*parentID)
		if err == nil {
			if parent.RootID != nil {
				comment.RootID = parent.RootID
			} else {
				comment.RootID = parentID
			}
			// 更新父评论回复数
			rootID := comment.RootID
			if rootID == nil {
				rootID = parentID
			}
			_ = repository.IncrementCommentReplyCount(*rootID)
		}
	}

	if err := repository.CreateComment(comment); err != nil {
		return nil, err
	}

	// 一级评论：回填 root_id
	if parentID == nil || *parentID == 0 {
		repository.UpdateCommentRootID(comment.ID, comment.ID)
		comment.RootID = &comment.ID
	}

	// 更新视频评论数
	_ = repository.UpdateVideoCount(videoID, "comment_count", 1)

	return comment, nil
}

// ToggleLike 切换点赞状态，返回当前是否已点赞
func ToggleLike(mpAccountID, userID uint64, targetType string, targetID uint64) (bool, error) {
	existing, err := repository.FindLike(mpAccountID, userID, targetType, targetID)
	if err != nil {
		// 不存在，创建点赞
		like := &model.Like{
			MpAccountID: mpAccountID,
			UserID:      userID,
			TargetType:  targetType,
			TargetID:    targetID,
			Status:      1,
		}
		if err := repository.CreateLike(like); err != nil {
			return false, err
		}
		// 更新目标计数
		if targetType == "video" {
			_ = repository.UpdateVideoCount(targetID, "like_count", 1)
		} else if targetType == "comment" {
			_ = repository.IncrementCommentLikeCount(targetID)
		}
		return true, nil
	}

	// 存在，切换状态
	newStatus := int8(1)
	if existing.Status == 1 {
		newStatus = 0
	}
	repository.UpdateLikeStatus(existing.ID, newStatus)

	delta := 1
	if newStatus == 0 {
		delta = -1
	}
	if targetType == "video" {
		_ = repository.UpdateVideoCount(targetID, "like_count", delta)
	}

	return newStatus == 1, nil
}

// ToggleFavorite 切换收藏状态
func ToggleFavorite(mpAccountID, userID, videoID uint64) (bool, error) {
	existing, err := repository.FindFavorite(mpAccountID, userID, videoID)
	if err != nil {
		fav := &model.Favorite{
			MpAccountID: mpAccountID,
			UserID:      userID,
			VideoID:     videoID,
			Status:      1,
		}
		if err := repository.CreateFavorite(fav); err != nil {
			return false, err
		}
		_ = repository.UpdateVideoCount(videoID, "collect_count", 1)
		return true, nil
	}

	newStatus := int8(1)
	if existing.Status == 1 {
		newStatus = 0
	}
	repository.UpdateFavoriteStatus(existing.ID, newStatus)

	delta := 1
	if newStatus == 0 {
		delta = -1
	}
	_ = repository.UpdateVideoCount(videoID, "collect_count", delta)

	return newStatus == 1, nil
}

// RecordShare 记录分享
func RecordShare(mpAccountID, userID, videoID uint64, shareType string) error {
	share := &model.Share{
		MpAccountID: mpAccountID,
		UserID:      userID,
		VideoID:     videoID,
		ShareType:   shareType,
	}
	if err := repository.CreateShare(share); err != nil {
		return err
	}
	_ = repository.UpdateVideoCount(videoID, "share_count", 1)
	return nil
}

// GetInteractionStatus 获取用户对视频的互动状态
func GetInteractionStatus(mpAccountID, userID, videoID uint64) map[string]bool {
	liked := false
	if like, err := repository.FindLike(mpAccountID, userID, "video", videoID); err == nil && like.Status == 1 {
		liked = true
	}
	favored := false
	if fav, err := repository.FindFavorite(mpAccountID, userID, videoID); err == nil && fav.Status == 1 {
		favored = true
	}
	followed := false
	if video, err := repository.FindVideoByID(videoID); err == nil {
		if follow, err := repository.FindFollow(mpAccountID, userID, video.EntityID); err == nil && follow.Status == 1 {
			followed = true
		}
	}
	return map[string]bool{
		"liked":    liked,
		"favored":  favored,
		"followed": followed,
	}
}

// ToggleFollow 切换关注状态
// entityID 为 0 时兜底为 mpAccountID 当前绑定的默认主体
func ToggleFollow(mpAccountID, userID, entityID uint64) (bool, error) {
	entityID, err := ResolveEntityID(entityID, mpAccountID)
	if err != nil {
		return false, err
	}

	existing, err := repository.FindFollow(mpAccountID, userID, entityID)
	if err != nil {
		follow := &model.Follow{
			MpAccountID: mpAccountID,
			UserID:      userID,
			EntityID:    entityID,
			Status:      1,
		}
		if err := repository.CreateFollow(follow); err != nil {
			return false, err
		}
		return true, nil
	}

	newStatus := int8(1)
	if existing.Status == 1 {
		newStatus = 0
	}
	if err := repository.UpdateFollowStatus(existing.ID, newStatus); err != nil {
		return false, err
	}
	return newStatus == 1, nil
}
