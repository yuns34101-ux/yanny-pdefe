package repository

import (
	"yanny-service/internal/database"
	"yanny-service/internal/model"
)

// ========== 主体 ==========

func FindEntityByID(id uint64) (*model.Entity, error) {
	var entity model.Entity
	err := database.DB.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func ListEntities(keyword string, status *int8, page, pageSize int) ([]model.Entity, int64, error) {
	var entities []model.Entity
	var total int64
	db := database.DB.Model(&model.Entity{})
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("sort_order DESC, id DESC").Find(&entities).Error
	return entities, total, err
}

func CreateEntity(entity *model.Entity) error {
	return database.DB.Create(entity).Error
}

func UpdateEntity(id uint64, updates map[string]interface{}) error {
	return database.DB.Model(&model.Entity{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteEntity(id uint64) error {
	return database.DB.Delete(&model.Entity{}, id).Error
}

// ========== 小程序账号 ==========

func FindMpAccountByID(id uint64) (*model.MpAccount, error) {
	var mp model.MpAccount
	err := database.DB.First(&mp, id).Error
	if err != nil {
		return nil, err
	}
	return &mp, nil
}

func FindMpAccountByAppID(appID string) (*model.MpAccount, error) {
	var mp model.MpAccount
	err := database.DB.Where("app_id = ?", appID).First(&mp).Error
	if err != nil {
		return nil, err
	}
	return &mp, nil
}

func ListMpAccounts(page, pageSize int) ([]model.MpAccount, int64, error) {
	var mps []model.MpAccount
	var total int64
	db := database.DB.Model(&model.MpAccount{})
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id DESC").Find(&mps).Error
	return mps, total, err
}

func CreateMpAccount(mp *model.MpAccount) error {
	return database.DB.Create(mp).Error
}

func UpdateMpAccount(id uint64, updates map[string]interface{}) error {
	return database.DB.Model(&model.MpAccount{}).Where("id = ?", id).Updates(updates).Error
}

// ========== 主体-小程序绑定 ==========

func BindEntityMp(binding *model.EntityMpBinding) error {
	return database.DB.Create(binding).Error
}

func UnbindEntityMp(entityID, mpAccountID uint64) error {
	return database.DB.Where("entity_id = ? AND mp_account_id = ?", entityID, mpAccountID).
		Delete(&model.EntityMpBinding{}).Error
}

func FindBindingsByEntity(entityID uint64) ([]model.EntityMpBinding, error) {
	var bindings []model.EntityMpBinding
	err := database.DB.Where("entity_id = ?", entityID).Find(&bindings).Error
	return bindings, err
}

func FindBindingsByMp(mpAccountID uint64) ([]model.EntityMpBinding, error) {
	var bindings []model.EntityMpBinding
	err := database.DB.Where("mp_account_id = ?", mpAccountID).Find(&bindings).Error
	return bindings, err
}

// ========== CDN 配置 ==========

func ListCdnConfigs(mpAccountID uint64, page, pageSize int) ([]model.CdnConfig, int64, error) {
	var configs []model.CdnConfig
	var total int64
	db := database.DB.Model(&model.CdnConfig{}).Where("mp_account_id = ?", mpAccountID)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id DESC").Find(&configs).Error
	return configs, total, err
}

func CreateCdnConfig(cfg *model.CdnConfig) error {
	return database.DB.Create(cfg).Error
}

func UpdateCdnConfig(id uint64, updates map[string]interface{}) error {
	return database.DB.Model(&model.CdnConfig{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteCdnConfig(id uint64) error {
	return database.DB.Delete(&model.CdnConfig{}, id).Error
}
