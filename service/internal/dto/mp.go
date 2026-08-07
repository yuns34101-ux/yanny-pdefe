package dto

// ========== 小程序端请求/响应 ==========

// MpLoginRequest 小程序登录请求（多平台兼容）
type MpLoginRequest struct {
	Code          string `json:"code"`
	Platform      string `json:"platform"`      // wechat / douyin / h5 / mock
	AppID         string `json:"app_id"`        // 小程序 AppID（用于查找 mp_account_id）
	EncryptedData string `json:"encrypted_data"` // 微信手机号加密数据
	Iv            string `json:"iv"`             // 微信手机号加密向量
	Phone         string `json:"phone"`          // H5 手机号登录
	Nickname      string `json:"nickname"`       // H5 昵称
	AvatarURL     string `json:"avatar_url"`     // H5 头像
}

// MpLoginResponse 小程序登录响应
type MpLoginResponse struct {
	Token      string `json:"token"`
	UserID     uint64 `json:"user_id"`
	IsNewUser  bool   `json:"is_new_user"`
	Nickname   string `json:"nickname"`
	AvatarURL  string `json:"avatar_url"`
}

// MpUserInfoResponse 小程序用户信息
type MpUserInfoResponse struct {
	UserID    uint64 `json:"user_id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Phone     string `json:"phone"`
	Gender    int8   `json:"gender"`
}

// UpdateUserInfoRequest 更新用户信息
type UpdateUserInfoRequest struct {
	Nickname  string `json:"nickname" binding:"max=100"`
	AvatarURL string `json:"avatar_url" binding:"max=500"`
	Gender    int8   `json:"gender"`
}

// CreateCommentRequest 创建评论
type CreateCommentRequest struct {
	VideoID       uint64  `json:"video_id" binding:"required"`
	Content       string  `json:"content" binding:"required,max=1000"`
	ParentID      *uint64 `json:"parent_id"`
	ReplyToUserID *uint64 `json:"reply_to_user_id"`
}

// VideoListQuery 视频列表查询（小程序端）
type VideoListQuery struct {
	Pagination
	CategoryID uint64 `json:"category_id" form:"category_id"`
	EntityID   uint64 `json:"entity_id" form:"entity_id"`
}

// ========== 管理后台视频管理 ==========

// CreateVideoRequest 创建视频（管理后台）
type CreateVideoRequest struct {
	MpAccountID  uint64 `json:"mp_account_id" binding:"required"`
	EntityID     uint64 `json:"entity_id" binding:"required"`
	CategoryID   uint64 `json:"category_id" binding:"required"`
	Title        string `json:"title" binding:"max=200"`
	Description  string `json:"description"`
	CoverURL     string `json:"cover_url" binding:"required"`
	VideoURL     string `json:"video_url" binding:"required"`
	Duration     uint   `json:"duration"`
	Width        uint   `json:"width"`
	Height       uint   `json:"height"`
	FileSize     uint64 `json:"file_size"`
	Tags         string `json:"tags"`
	Status       int8   `json:"status"`
	IsRecommended int8  `json:"is_recommended"`
}

// UpdateVideoRequest 更新视频
type UpdateVideoRequest struct {
	Title         string `json:"title" binding:"max=200"`
	Description   string `json:"description"`
	CoverURL      string `json:"cover_url"`
	VideoURL      string `json:"video_url"`
	CategoryID    uint64 `json:"category_id"`
	Tags          string `json:"tags"`
	Status        int8   `json:"status"`
	IsRecommended int8   `json:"is_recommended"`
}
