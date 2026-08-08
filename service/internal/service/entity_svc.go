package service

import (
	"errors"
	"yanny-service/internal/model"
	"yanny-service/internal/repository"
	"yanny-service/internal/util/qiniu"

	"gorm.io/gorm"
)

// ========== 主体管理 ==========

// CreateEntity 创建主体
func CreateEntity(name, logoURL, desc, phone, email, addr string, lat, lng *float64, extra string, sortOrder int, status int8) (*model.Entity, error) {
	entity := &model.Entity{
		Name:         name,
		LogoURL:      logoURL,
		Description:  desc,
		ContactPhone: phone,
		ContactEmail: email,
		Address:      addr,
		Latitude:     lat,
		Longitude:    lng,
		SortOrder:    sortOrder,
		Status:       status,
	}
	if status == 0 {
		entity.Status = 1
	}
	if extra != "" {
		entity.Extra = &extra
	}
	if err := repository.CreateEntity(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// UpdateEntity 更新主体
func UpdateEntity(id uint64, name, logoURL, desc, phone, email, addr string, lat, lng *float64, extra string, sortOrder int, status int8) error {
	if _, err := repository.FindEntityByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("主体不存在")
		}
		return err
	}
	updates := map[string]interface{}{
		"name": name, "logo_url": logoURL, "description": desc,
		"contact_phone": phone, "contact_email": email, "address": addr,
		"latitude": lat, "longitude": lng,
		"sort_order": sortOrder, "status": status,
	}
	if extra != "" {
		updates["extra"] = extra
	} else {
		updates["extra"] = nil
	}
	return repository.UpdateEntity(id, updates)
}

// ========== 小程序账号 ==========

// CreateMpAccount 创建小程序账号
func CreateMpAccount(appID, appSecret, appName, appIcon, desc string, status int8) (*model.MpAccount, error) {
	mp := &model.MpAccount{
		AppID:       appID,
		AppSecret:   appSecret,
		AppName:     appName,
		AppIcon:     appIcon,
		Description: desc,
		Status:      status,
	}
	if status == 0 {
		mp.Status = 1
	}
	if err := repository.CreateMpAccount(mp); err != nil {
		return nil, err
	}
	return mp, nil
}

// ========== 绑定管理 ==========

// BindEntityMp 绑定主体和小程序
func BindEntityMp(entityID, mpAccountID uint64, isDefault int8) error {
	// 校验主体和小程序都存在
	if _, err := repository.FindEntityByID(entityID); err != nil {
		return errors.New("主体不存在")
	}
	if _, err := repository.FindMpAccountByID(mpAccountID); err != nil {
		return errors.New("小程序账号不存在")
	}

	binding := &model.EntityMpBinding{
		EntityID:    entityID,
		MpAccountID: mpAccountID,
		IsDefault:   isDefault,
	}
	return repository.BindEntityMp(binding)
}

// UnbindEntityMp 解绑
func UnbindEntityMp(entityID, mpAccountID uint64) error {
	return repository.UnbindEntityMp(entityID, mpAccountID)
}

// ========== 小程序端主体信息 ==========

// EntityWithStats 主体信息 + 实时统计
type EntityWithStats struct {
	*model.Entity
	VideoCount    int64 `json:"video_count"`
	FollowerCount int64 `json:"follower_count"`
	Followed      bool  `json:"followed"`
}

// ResolveEntityID 解析主体 ID：entityID 优先使用；为 0 时通过 mpAccountID 反查默认绑定的主体
// 小程序账号与主体目前是 1 对 1 关系，客户端提交 entity_id=0 时兜底为当前绑定主体，避免"主体不存在"异常
func ResolveEntityID(entityID, mpAccountID uint64) (uint64, error) {
	if entityID > 0 {
		return entityID, nil
	}
	if mpAccountID == 0 {
		return 0, errors.New("主体不存在")
	}
	bindings, err := repository.FindBindingsByMp(mpAccountID)
	if err != nil {
		return 0, err
	}
	for _, b := range bindings {
		if b.IsDefault == 1 {
			return b.EntityID, nil
		}
	}
	if len(bindings) > 0 {
		return bindings[0].EntityID, nil
	}
	return 0, errors.New("主体不存在")
}

// GetEntityForMp 小程序端获取主体详情及统计信息
// entityID 优先使用；为 0 时通过 mpAccountID 反查默认绑定的主体
// userID 为 0 表示游客，followed 恒为 false
func GetEntityForMp(entityID, mpAccountID, userID uint64) (*EntityWithStats, error) {
	entityID, err := ResolveEntityID(entityID, mpAccountID)
	if err != nil {
		return nil, err
	}

	entity, err := repository.FindEntityByID(entityID)
	if err != nil {
		return nil, err
	}
	if entity.LogoURL != "" {
		entity.LogoURL = qiniu.SignImageURL(entity.LogoURL)
	}
	videoCount, err := repository.CountVideosByEntity(entityID)
	if err != nil {
		return nil, err
	}
	followerCount, err := repository.CountFollowersByEntity(entityID)
	if err != nil {
		return nil, err
	}
	followed := false
	if userID > 0 {
		if follow, err := repository.FindFollow(mpAccountID, userID, entityID); err == nil && follow.Status == 1 {
			followed = true
		}
	}
	return &EntityWithStats{
		Entity:        entity,
		VideoCount:    videoCount,
		FollowerCount: followerCount,
		Followed:      followed,
	}, nil
}
