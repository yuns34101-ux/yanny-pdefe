package dto

// CreateEntityRequest 创建主体
type CreateEntityRequest struct {
	Name         string   `json:"name" binding:"required,max=100"`
	LogoURL      string   `json:"logo_url"`
	Description  string   `json:"description"`
	ContactPhone string   `json:"contact_phone"`
	ContactEmail string   `json:"contact_email"`
	Address      string   `json:"address"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Extra        string   `json:"extra"`
	SortOrder    int      `json:"sort_order"`
	Status       int8     `json:"status"`
}

// UpdateEntityRequest 更新主体
type UpdateEntityRequest struct {
	Name         string   `json:"name" binding:"max=100"`
	LogoURL      string   `json:"logo_url"`
	Description  string   `json:"description"`
	ContactPhone string   `json:"contact_phone"`
	ContactEmail string   `json:"contact_email"`
	Address      string   `json:"address"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Extra        string   `json:"extra"`
	SortOrder    int      `json:"sort_order"`
	Status       int8     `json:"status"`
}

// EntityListQuery 主体列表查询
type EntityListQuery struct {
	Pagination
	Keyword string `json:"keyword" form:"keyword"`
	Status  *int8  `json:"status" form:"status"`
}

// CreateMpAccountRequest 创建小程序账号
type CreateMpAccountRequest struct {
	AppID       string `json:"app_id" binding:"required,max=64"`
	AppSecret   string `json:"app_secret" binding:"required,max=128"`
	AppName     string `json:"app_name" binding:"required,max=100"`
	AppIcon     string `json:"app_icon"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
}

// CreateCdnConfigRequest 创建 CDN 配置
type CreateCdnConfigRequest struct {
	MpAccountID uint64 `json:"mp_account_id" binding:"required"`
	Provider    string `json:"provider"`
	AccessKey   string `json:"access_key" binding:"required"`
	SecretKey   string `json:"secret_key" binding:"required"`
	Bucket      string `json:"bucket" binding:"required"`
	Domain      string `json:"domain" binding:"required"`
	Region      string `json:"region"`
	CallbackURL string `json:"callback_url"`
	Status      int8   `json:"status"`
}

// BindEntityMpRequest 绑定主体-小程序
type BindEntityMpRequest struct {
	EntityID    uint64 `json:"entity_id" binding:"required"`
	MpAccountID uint64 `json:"mp_account_id" binding:"required"`
	IsDefault   int8   `json:"is_default"`
}
