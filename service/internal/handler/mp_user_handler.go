package handler

import (
	"yanny-service/internal/dto"
	"yanny-service/internal/middleware"
	"yanny-service/internal/repository"
	"yanny-service/internal/util/wechat"

	"github.com/gin-gonic/gin"
)

// MpUpdatePhone 绑定微信手机号（需登录，用登录时留存的 session_key 解密）
func MpUpdatePhone(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req struct {
		EncryptedData string `json:"encrypted_data" binding:"required"`
		Iv            string `json:"iv" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}

	user, err := repository.FindUserByID(userID)
	if err != nil {
		dto.Error(c, dto.ErrCodeUserNotFound, "用户不存在")
		return
	}
	if user.SessionKey == "" {
		dto.Error(c, dto.ErrCodeParamInvalid, "登录状态已过期，请重新进入小程序后再试")
		return
	}

	phone, err := wechat.DecryptPhoneNumber(user.SessionKey, req.EncryptedData, req.Iv)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, "手机号解析失败："+err.Error())
		return
	}

	if err := repository.UpdateUserPhone(userID, phone); err != nil {
		dto.Error(c, dto.ErrCodeInternal, "手机号保存失败")
		return
	}

	dto.Success(c, gin.H{"phone": phone})
}

// MpUpdateUserInfo 更新用户昵称/头像（需登录，chooseAvatar + nickname 采集后调用）
func MpUpdateUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	if err := repository.UpdateUserInfo(userID, req.Nickname, req.AvatarURL); err != nil {
		dto.Error(c, dto.ErrCodeInternal, "保存失败")
		return
	}
	dto.Success(c, gin.H{"nickname": req.Nickname, "avatar_url": req.AvatarURL})
}

// MpGetUserInfo 获取当前登录用户信息（token 存在但内存中 userInfo 丢失时恢复用）
func MpGetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := repository.FindUserByID(userID)
	if err != nil {
		dto.Error(c, dto.ErrCodeUserNotFound, "用户不存在")
		return
	}
	dto.Success(c, dto.MpUserInfoResponse{
		UserID:    user.ID,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Phone:     user.Phone,
		Gender:    user.Gender,
	})
}

// MpGetUploadToken 获取七牛上传 Token（小程序端，仅图片，用于头像直传）
func MpGetUploadToken(c *gin.Context) {
	result, err := buildQiniuUploadToken()
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, result)
}
