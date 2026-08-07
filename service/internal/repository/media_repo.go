package repository

import (
	"yanny-service/internal/database"
	"yanny-service/internal/model"
)

// ========== 媒体资源去重 ==========

func FindMediaAssetByClientHash(mpAccountID uint64, clientHash string) (*model.MediaAsset, error) {
	var asset model.MediaAsset
	err := database.DB.Where("mp_account_id = ? AND client_hash = ?", mpAccountID, clientHash).
		First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func FindMediaAssetByContentHash(mpAccountID uint64, contentHash string) (*model.MediaAsset, error) {
	var asset model.MediaAsset
	err := database.DB.Where("mp_account_id = ? AND content_hash = ?", mpAccountID, contentHash).
		First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func CreateMediaAsset(asset *model.MediaAsset) error {
	return database.DB.Create(asset).Error
}
