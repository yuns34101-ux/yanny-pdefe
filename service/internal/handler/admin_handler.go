package handler

import (
	"strconv"
	"strings"
	"time"
	"yanny-service/internal/config"
	"yanny-service/internal/dto"
	"yanny-service/internal/middleware"
	"yanny-service/internal/model"
	"yanny-service/internal/repository"
	"yanny-service/internal/service"
	"yanny-service/internal/util/qiniu"

	"github.com/gin-gonic/gin"
)

// getEntityScope 获取当前管理员的数据范围（entity ID 列表）
// 超管返回 nil（代表全部），普通管理员返回分配的主体 ID 列表
func getEntityScope(c *gin.Context) []uint64 {
	adminID := middleware.GetAdminID(c)
	if adminID == 0 || repository.IsSuperAdmin(adminID) {
		return nil // nil = 全部
	}
	ids, _ := repository.GetAdminEntityIDs(adminID)
	return ids
}

// signEntityLogoURLs 给实体列表的 logo_url 统一签名
func signEntityLogoURLs(entities []model.Entity) {
	for i := range entities {
		if entities[i].LogoURL != "" {
			entities[i].LogoURL = qiniu.SignImageURL(entities[i].LogoURL)
		}
	}
}

// signMpIconURLs 给小程序列表的 app_icon 统一签名
func signMpIconURLs(mps []model.MpAccount) {
	for i := range mps {
		if mps[i].AppIcon != "" {
			mps[i].AppIcon = qiniu.SignImageURL(mps[i].AppIcon)
		}
	}
}

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "请填写用户名和密码")
		return
	}

	token, admin, err := service.AdminLogin(req.Username, req.Password, c.ClientIP())
	if err != nil {
		dto.Error(c, dto.ErrCodeLoginFailed, err.Error())
		return
	}

	dto.Success(c, dto.AdminLoginResponse{
		Token:    token,
		AdminID:  admin.ID,
		Username: admin.Username,
		RealName: admin.RealName,
	})
}

// AdminInfo 获取当前管理员信息
func AdminInfo(c *gin.Context) {
	adminID := middleware.GetAdminID(c)
	info, err := service.GetAdminInfo(adminID)
	if err != nil {
		dto.Error(c, dto.ErrCodeAdminNotFound, "管理员不存在")
		return
	}
	dto.Success(c, info)
}

// AdminLogout 管理员登出（JWT 加入黑名单）
func AdminLogout(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	if tokenStr != "" {
		middleware.AddTokenBlacklist(tokenStr, time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireHours)*time.Hour))
	}
	dto.Success(c, nil)
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}

	adminID := middleware.GetAdminID(c)
	if err := service.ChangePassword(adminID, req.OldPassword, req.NewPassword); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, err.Error())
		return
	}
	dto.Success(c, nil)
}

// ========== 管理员 CRUD ==========

// ListAdmins 管理员列表（含角色和主体分配）
func ListAdmins(c *gin.Context) {
	var p dto.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p = dto.DefaultPagination()
	}

	admins, total, err := repository.ListAdmins(p.Page, p.PageSize)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	// 填充每个管理员的实体 ID 列表（用于前端展示）
	type AdminWithEntities struct {
		model.Admin
		EntityIDs []uint64 `json:"entity_ids"`
	}
	result := make([]AdminWithEntities, len(admins))
	for i := range admins {
		result[i].Admin = admins[i]
		result[i].EntityIDs, _ = repository.GetAdminEntityIDs(admins[i].ID)
	}
	dto.SuccessPage(c, result, p.Page, p.PageSize, total)
}

// ========== 主体 CRUD ==========

// CreateEntity 创建主体
func CreateEntity(c *gin.Context) {
	var req dto.CreateEntityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效："+err.Error())
		return
	}
	entity, err := service.CreateEntity(req.Name, qiniu.StripQuery(req.LogoURL), req.Description,
		req.ContactPhone, req.ContactEmail, req.Address, req.Extra,
		req.SortOrder, req.Status)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, entity)
}

// ListEntities 主体列表
func ListEntities(c *gin.Context) {
	var q dto.EntityListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		q.Page, q.PageSize = 1, 20
	}
	scope := getEntityScope(c)
	entities, total, err := repository.ListEntities(q.Keyword, q.Status, q.Page, q.PageSize, scope)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	signEntityLogoURLs(entities)
	dto.SuccessPage(c, entities, q.Page, q.PageSize, total)
}

// GetEntity 获取主体详情
func GetEntity(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	entity, err := repository.FindEntityByID(id)
	if err != nil {
		dto.Error(c, dto.ErrCodeEntityNotFound, "主体不存在")
		return
	}
	signEntityLogoURLs([]model.Entity{*entity})
	dto.Success(c, entity)
}

// UpdateEntity 更新主体
func UpdateEntity(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	var req dto.UpdateEntityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	if err := service.UpdateEntity(id, req.Name, qiniu.StripQuery(req.LogoURL), req.Description,
		req.ContactPhone, req.ContactEmail, req.Address, req.Extra,
		req.SortOrder, req.Status); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

// DeleteEntity 删除主体
func DeleteEntity(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	if err := repository.DeleteEntity(id); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

// ========== 小程序账号 CRUD ==========

// CreateMpAccount 创建小程序账号
func CreateMpAccount(c *gin.Context) {
	var req dto.CreateMpAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	mp, err := service.CreateMpAccount(req.AppID, req.AppSecret, req.AppName,
		qiniu.StripQuery(req.AppIcon), req.Description, req.Status)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, mp)
}

// ListMpAccounts 小程序账号列表
func ListMpAccounts(c *gin.Context) {
	var p dto.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p = dto.DefaultPagination()
	}
	mps, total, err := repository.ListMpAccounts(p.Page, p.PageSize)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	signMpIconURLs(mps)
	dto.SuccessPage(c, mps, p.Page, p.PageSize, total)
}

// ========== 绑定管理 ==========

// ListEntityBindings 查询主体的绑定列表
func ListEntityBindings(c *gin.Context) {
	entityID, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	bindings, err := repository.FindBindingsByEntity(entityID)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, bindings)
}

// BindEntityMp 绑定主体-小程序
func BindEntityMp(c *gin.Context) {
	var req dto.BindEntityMpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	if err := service.BindEntityMp(req.EntityID, req.MpAccountID, req.IsDefault); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, err.Error())
		return
	}
	dto.Success(c, nil)
}

// UnbindEntityMp 解绑
func UnbindEntityMp(c *gin.Context) {
	var req dto.BindEntityMpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	if err := repository.UnbindEntityMp(req.EntityID, req.MpAccountID); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

// ========== CDN 配置 ==========

// CreateCdnConfig 创建 CDN 配置
func CreateCdnConfig(c *gin.Context) {
	var req dto.CreateCdnConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	cfg := &model.CdnConfig{
		MpAccountID: req.MpAccountID,
		Provider:    req.Provider,
		AccessKey:   req.AccessKey,
		SecretKey:   req.SecretKey,
		Bucket:      req.Bucket,
		Domain:      req.Domain,
		Region:      req.Region,
		CallbackURL: req.CallbackURL,
		Status:      req.Status,
	}
	if cfg.Provider == "" {
		cfg.Provider = "qiniu"
	}
	if cfg.Status == 0 {
		cfg.Status = 1
	}
	if err := repository.CreateCdnConfig(cfg); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, cfg)
}

// ========== 小程序账号补充 ==========

// UpdateMpAccount 更新小程序账号
func UpdateMpAccount(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	if iconURL, ok := req["app_icon"].(string); ok {
		req["app_icon"] = qiniu.StripQuery(iconURL)
	}
	if err := repository.UpdateMpAccount(id, req); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

// ========== CDN 配置补充 ==========

// ListCdnConfigs CDN 配置列表
func ListCdnConfigs(c *gin.Context) {
	mpAccountID, _ := parseUintQuery(c, "mp_account_id")
	var p dto.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p = dto.DefaultPagination()
	}
	configs, total, err := repository.ListCdnConfigs(mpAccountID, p.Page, p.PageSize)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.SuccessPage(c, configs, p.Page, p.PageSize, total)
}

// UpdateCdnConfig 更新 CDN 配置
func UpdateCdnConfig(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	if err := repository.UpdateCdnConfig(id, req); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

// DeleteCdnConfig 删除 CDN 配置
func DeleteCdnConfig(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	if err := repository.DeleteCdnConfig(id); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

// ========== 用户管理 ==========

// ListUsers 用户列表
func ListUsers(c *gin.Context) {
	var q struct {
		dto.Pagination
		MpAccountID uint64 `form:"mp_account_id"`
		Phone       string `form:"phone"`
	}
	if err := c.ShouldBindQuery(&q); err != nil {
		q.Page, q.PageSize = 1, 20
	}
	users, total, err := repository.ListUsers(q.MpAccountID, q.Phone, q.Page, q.PageSize)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.SuccessPage(c, users, q.Page, q.PageSize, total)
}

// UpdateUserStatus 切换用户状态（启用/禁用）
func UpdateUserStatus(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	if err := repository.UpdateUserStatus(id, req.Status); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

// ========== 管理员 CRUD 补充 ==========

// CreateAdmin 创建管理员
func CreateAdmin(c *gin.Context) {
	var req struct {
		Username  string   `json:"username" binding:"required"`
		Password  string   `json:"password" binding:"required"`
		RealName  string   `json:"real_name"`
		RoleIDs   []uint64 `json:"role_ids"`
		EntityIDs []uint64 `json:"entity_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	admin, err := service.CreateAdmin(req.Username, req.Password, req.RealName, req.RoleIDs)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	if len(req.EntityIDs) > 0 {
		repository.AssignAdminEntities(admin.ID, req.EntityIDs)
	}
	dto.Success(c, admin)
}

// UpdateAdmin 更新管理员
func UpdateAdmin(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	var req struct {
		RealName  string   `json:"real_name"`
		Status    *int8    `json:"status"`
		RoleIDs   []uint64 `json:"role_ids"`
		EntityIDs []uint64 `json:"entity_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	if err := service.UpdateAdmin(id, req.RealName, req.Status, req.RoleIDs); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	if len(req.EntityIDs) > 0 {
		repository.AssignAdminEntities(id, req.EntityIDs)
	}
	dto.Success(c, nil)
}

// DeleteAdmin 删除管理员
func DeleteAdmin(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	if err := repository.DeleteAdmin(id); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

// ========== 角色管理 ==========

// ListRoles 角色列表
func ListRoles(c *gin.Context) {
	roles, err := repository.ListRoles()
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, roles)
}

// GetRolePermissions 获取角色权限
func GetRolePermissions(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	perms, err := repository.FindRolePermissions(id)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, perms)
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	var req struct {
		Name          string   `json:"name" binding:"required"`
		Code          string   `json:"code" binding:"required"`
		Description   string   `json:"description"`
		PermissionIDs []uint64 `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	role, err := service.CreateRole(req.Name, req.Code, req.Description, req.PermissionIDs)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, role)
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		dto.Error(c, dto.ErrCodeParamInvalid, "ID 参数格式无效")
		return
	}
	var req struct {
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		PermissionIDs []uint64 `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效")
		return
	}
	if err := service.UpdateRole(id, req.Name, req.Description, req.PermissionIDs); err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, nil)
}

// parseUintParam 从 URL 路径参数解析 uint64，解析失败返回 0 + false
func parseUintParam(c *gin.Context, name string) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

// parseUintQuery 从 Query 参数解析 uint64
func parseUintQuery(c *gin.Context, name string) (uint64, bool) {
	v := c.Query(name)
	if v == "" {
		return 0, false
	}
	id, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}