package service

import (
	"errors"
	"yanny-service/internal/model"
	"yanny-service/internal/repository"

	"gorm.io/gorm"
)

// ========== 主体管理 ==========

// CreateEntity 创建主体
func CreateEntity(name, logoURL, desc, phone, email, addr, extra string, sortOrder int, status int8) (*model.Entity, error) {
	entity := &model.Entity{
		Name:         name,
		LogoURL:      logoURL,
		Description:  desc,
		ContactPhone: phone,
		ContactEmail: email,
		Address:      addr,
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
func UpdateEntity(id uint64, name, logoURL, desc, phone, email, addr, extra string, sortOrder int, status int8) error {
	if _, err := repository.FindEntityByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("主体不存在")
		}
		return err
	}
	updates := map[string]interface{}{
		"name": name, "logo_url": logoURL, "description": desc,
		"contact_phone": phone, "contact_email": email, "address": addr,
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
