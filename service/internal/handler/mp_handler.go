package handler

import (
	"yanny-service/internal/dto"
	"yanny-service/internal/service"

	"github.com/gin-gonic/gin"
)

// MpLogin 小程序登录（多平台兼容：wechat / douyin / h5 / mock）
func MpLogin(c *gin.Context) {
	var req dto.MpLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 兼容无参数请求（静默登录），使用默认 mock 平台
		req.Platform = "mock"
	}

	if req.Platform == "" {
		req.Platform = "mock"
	}

	// 根据 AppID 查找 mp_account_id
	mpAccountID := uint64(1) // 默认
	if req.AppID != "" {
		if mp, err := service.FindMpAccountByAppID(req.AppID); err == nil {
			mpAccountID = mp.ID
		}
	}

	var openid, unionid, nickname, avatarURL string

	switch req.Platform {
	case "mock":
		// E2E 测试模式：用 code 作为 openid
		openid = "mock_" + req.Code
		if req.Code == "" {
			openid = "mock_test_user"
		}
		nickname = req.Nickname
		if nickname == "" {
			nickname = "测试用户"
		}
		avatarURL = req.AvatarURL

	case "wechat":
		// TODO: 调用微信 code2Session
		// wxResp := wechat.Code2Session(req.Code)
		// openid = wxResp.Openid
		// unionid = wxResp.Unionid
		openid = "wx_" + req.Code
		nickname = "微信用户"

	case "douyin":
		// TODO: 调用抖音 code2Session
		openid = "dy_" + req.Code
		nickname = "抖音用户"

	case "h5":
		openid = "h5_" + req.Phone
		nickname = req.Nickname
		avatarURL = req.AvatarURL

	default:
		openid = "unknown_" + req.Code
		nickname = "用户"
	}

	// 查找或创建用户
	user, isNew, err := service.FindOrCreateUser(mpAccountID, openid, unionid, nickname, avatarURL, c.ClientIP())
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, "登录失败："+err.Error())
		return
	}

	// 签发 JWT
	token, err := service.GenerateMpToken(user.ID, mpAccountID)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, "生成令牌失败")
		return
	}

	dto.Success(c, dto.MpLoginResponse{
		Token:     token,
		UserID:    user.ID,
		IsNewUser: isNew,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
	})
}
