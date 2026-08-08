package model

import (
	"time"

	"gorm.io/gorm"
)

// Entity 运营主体
type Entity struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	LogoURL      string         `gorm:"size:500;not null;default:''" json:"logo_url"`
	Description  string         `gorm:"size:500;not null;default:''" json:"description"`
	ContactPhone string         `gorm:"size:20;not null;default:''" json:"contact_phone"`
	ContactEmail string         `gorm:"size:100;not null;default:''" json:"contact_email"`
	Address      string         `gorm:"size:300;not null;default:''" json:"address"`
	Latitude     *float64       `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude    *float64       `gorm:"type:decimal(10,7)" json:"longitude"`
	Extra        *string        `gorm:"type:json" json:"extra"`
	SortOrder    int            `gorm:"not null;default:0" json:"sort_order"`
	Status       int8           `gorm:"not null;default:1" json:"status"`
	CreatedAt    time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Entity) TableName() string { return "entities" }

// MpAccount 小程序账号
type MpAccount struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	AppID       string         `gorm:"size:64;not null;uniqueIndex:uk_app_id" json:"app_id"`
	AppSecret   string         `gorm:"size:128;not null" json:"-"` // 禁止 JSON 输出
	AppName     string         `gorm:"size:100;not null" json:"app_name"`
	AppIcon     string         `gorm:"size:500;not null;default:''" json:"app_icon"`
	Description string         `gorm:"size:300;not null;default:''" json:"description"`
	Status      int8           `gorm:"not null;default:1" json:"status"`
	CreatedAt   time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (MpAccount) TableName() string { return "mp_accounts" }

// EntityMpBinding 主体-小程序绑定
type EntityMpBinding struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	EntityID    uint64    `gorm:"not null;uniqueIndex:uk_entity_mp" json:"entity_id"`
	MpAccountID uint64    `gorm:"not null;uniqueIndex:uk_entity_mp;index:idx_mp_account" json:"mp_account_id"`
	IsDefault   int8      `gorm:"not null;default:0" json:"is_default"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (EntityMpBinding) TableName() string { return "entity_mp_bindings" }
