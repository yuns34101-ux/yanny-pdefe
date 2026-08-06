package dto

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLoginResponse 登录响应
type AdminLoginResponse struct {
	Token    string `json:"token"`
	AdminID  uint64 `json:"admin_id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
}

// AdminInfoResponse 管理员信息
type AdminInfoResponse struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	RealName  string `json:"real_name"`
	AvatarURL string `json:"avatar_url"`
	Status    int8   `json:"status"`
	Roles     []RoleSimple `json:"roles"`
	Perms     []string     `json:"perms"`
}

// RoleSimple 角色简版
type RoleSimple struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// ChangePasswordRequest 修改密码
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}
