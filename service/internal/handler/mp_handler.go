package handler

import (
	"yanny-service/internal/dto"
	"yanny-service/internal/model"
	"yanny-service/internal/service"
	"yanny-service/internal/util/wechat"

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
	var mpAccount *model.MpAccount
	if req.AppID != "" {
		if mp, err := service.FindMpAccountByAppID(req.AppID); err == nil {
			mpAccountID = mp.ID
			mpAccount = mp
		}
	}

	var openid, unionid, nickname, avatarURL, sessionKey string

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
		if mpAccount == nil {
			dto.Error(c, dto.ErrCodeMpAccountNotFound, "小程序账号不存在")
			return
		}
		wxResp, err := wechat.Code2Session(mpAccount.AppID, mpAccount.AppSecret, req.Code)
		if err != nil {
			dto.Error(c, dto.ErrCodeLoginFailed, "登录失败："+err.Error())
			return
		}
		openid = wxResp.Openid
		unionid = wxResp.Unionid
		sessionKey = wxResp.SessionKey
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

	// 查找或创建用户（inviter_user_id 仅新用户首次登录生效，实现分享裂变的单层邀请归属）
	user, isNew, err := service.FindOrCreateUser(mpAccountID, openid, unionid, nickname, avatarURL, c.ClientIP(), req.InviterUserID, sessionKey)
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
