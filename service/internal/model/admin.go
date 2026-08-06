package model

import "time"

// Admin 管理员
type Admin struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string     `gorm:"size:64;not null;uniqueIndex:uk_username" json:"username"`
	Password    string     `gorm:"size:200;not null" json:"-"`
	RealName    string     `gorm:"size:50;not null;default:''" json:"real_name"`
	AvatarURL   string     `gorm:"size:500;not null;default:''" json:"avatar_url"`
	Status      int8       `gorm:"not null;default:1" json:"status"`
	LastLoginAt *time.Time `gorm:"null" json:"last_login_at"`
	LastLoginIP string     `gorm:"size:50;not null;default:''" json:"last_login_ip"`
	CreatedAt   time.Time  `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null;autoUpdateTime" json:"updated_at"`

	// 关联
	Roles []Role `gorm:"many2many:admin_roles" json:"roles,omitempty"`
}

func (Admin) TableName() string { return "admins" }

// Role 角色
type Role struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Code        string    `gorm:"size:50;not null;uniqueIndex:uk_code" json:"code"`
	Description string    `gorm:"size:200;not null;default:''" json:"description"`
	IsSystem    int8      `gorm:"not null;default:0" json:"is_system"`
	Status      int8      `gorm:"not null;default:1" json:"status"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`

	// 关联
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
}

func (Role) TableName() string { return "roles" }

// Permission 权限
type Permission struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Code        string    `gorm:"size:100;not null;uniqueIndex:uk_code" json:"code"`
	Module      string    `gorm:"size:50;not null;index:idx_module" json:"module"`
	Action      string    `gorm:"size:50;not null" json:"action"`
	Description string    `gorm:"size:200;not null;default:''" json:"description"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}

func (Permission) TableName() string { return "permissions" }

// AdminEntity 管理员-主体绑定（数据级权限）
type AdminEntity struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID   uint64    `gorm:"not null;uniqueIndex:uk_admin_entity" json:"admin_id"`
	EntityID  uint64    `gorm:"not null;uniqueIndex:uk_admin_entity;index:idx_entity" json:"entity_id"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}

func (AdminEntity) TableName() string { return "admin_entities" }

// AdminRole 管理员-角色关联
type AdminRole struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID   uint64    `gorm:"not null;uniqueIndex:uk_admin_role" json:"admin_id"`
	RoleID    uint64    `gorm:"not null;uniqueIndex:uk_admin_role;index:idx_role" json:"role_id"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}

func (AdminRole) TableName() string { return "admin_roles" }

// RolePermission 角色-权限关联
type RolePermission struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID       uint64    `gorm:"not null;uniqueIndex:uk_role_perm" json:"role_id"`
	PermissionID uint64    `gorm:"not null;uniqueIndex:uk_role_perm;index:idx_permission" json:"permission_id"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}

func (RolePermission) TableName() string { return "role_permissions" }
