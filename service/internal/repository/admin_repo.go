package repository

import (
	"yanny-service/internal/database"
	"yanny-service/internal/model"

	"gorm.io/gorm"
)

// FindAdminByUsername 按用户名查询管理员
func FindAdminByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := database.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindAdminByID 按 ID 查询管理员
func FindAdminByID(id uint64) (*model.Admin, error) {
	var admin model.Admin
	err := database.DB.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindAdminWithRoles 查询管理员及其角色
func FindAdminWithRoles(adminID uint64) (*model.Admin, error) {
	var admin model.Admin
	err := database.DB.Preload("Roles", "status = ?", 1).First(&admin, adminID).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindAdminPermissions 查询管理员的所有权限编码
func FindAdminPermissions(adminID uint64) ([]string, error) {
	var perms []string
	err := database.DB.Raw(`
		SELECT DISTINCT p.code
		FROM admin_roles ar
		JOIN role_permissions rp ON ar.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		JOIN roles r ON ar.role_id = r.id
		WHERE ar.admin_id = ? AND r.status = 1
	`, adminID).Scan(&perms).Error
	return perms, err
}

// UpdateAdminLoginInfo 更新管理员登录信息
func UpdateAdminLoginInfo(adminID uint64, ip string) error {
	return database.DB.Model(&model.Admin{}).Where("id = ?", adminID).
		Updates(map[string]interface{}{
			"last_login_at": database.DB.Raw("NOW()"),
			"last_login_ip": ip,
		}).Error
}

// UpdateAdminPassword 修改管理员密码
func UpdateAdminPassword(adminID uint64, newPassword string) error {
	return database.DB.Model(&model.Admin{}).Where("id = ?", adminID).
		Update("password", newPassword).Error
}

// ListAdmins 管理员列表
func ListAdmins(page, pageSize int) ([]model.Admin, int64, error) {
	var admins []model.Admin
	var total int64
	db := database.DB.Model(&model.Admin{}).Preload("Roles")
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&admins).Error
	return admins, total, err
}

// CreateAdmin 创建管理员
func CreateAdmin(admin *model.Admin) error {
	return database.DB.Create(admin).Error
}

// UpdateAdmin 更新管理员字段
func UpdateAdmin(id uint64, updates map[string]interface{}) error {
	return database.DB.Model(&model.Admin{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteAdmin 删除管理员
func DeleteAdmin(id uint64) error {
	// 同时清理角色和主体绑定
	database.DB.Where("admin_id = ?", id).Delete(&model.AdminRole{})
	database.DB.Where("admin_id = ?", id).Delete(&model.AdminEntity{})
	return database.DB.Delete(&model.Admin{}, id).Error
}

// ========== 数据级权限 ==========

// GetAdminEntityIDs 获取管理员可管理的主体 ID 列表
func GetAdminEntityIDs(adminID uint64) ([]uint64, error) {
	var ids []uint64
	err := database.DB.Model(&model.AdminEntity{}).
		Where("admin_id = ?", adminID).Pluck("entity_id", &ids).Error
	return ids, err
}

// AssignAdminEntities 分配管理员可管理的主体
func AssignAdminEntities(adminID uint64, entityIDs []uint64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("admin_id = ?", adminID).Delete(&model.AdminEntity{}).Error; err != nil {
			return err
		}
		for _, eid := range entityIDs {
			if err := tx.Create(&model.AdminEntity{AdminID: adminID, EntityID: eid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// IsSuperAdmin 判断是否为超级管理员
func IsSuperAdmin(adminID uint64) bool {
	var count int64
	database.DB.Raw(`SELECT COUNT(1) FROM admin_roles ar
		JOIN roles r ON ar.role_id = r.id
		WHERE ar.admin_id = ? AND r.code = 'super_admin' AND r.status = 1`, adminID).Scan(&count)
	return count > 0
}

// AssignRoles 给管理员分配角色
func AssignRoles(adminID uint64, roleIDs []uint64) error {
	if err := database.DB.Where("admin_id = ?", adminID).Delete(&model.AdminRole{}).Error; err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return nil
	}
	records := make([]model.AdminRole, len(roleIDs))
	for i, rid := range roleIDs {
		records[i] = model.AdminRole{AdminID: adminID, RoleID: rid}
	}
	return database.DB.Create(&records).Error
}

// ========== 角色管理 ==========

// ListRoles 角色列表（含权限）
func ListRoles() ([]model.Role, error) {
	var roles []model.Role
	err := database.DB.Preload("Permissions").Order("id ASC").Find(&roles).Error
	return roles, err
}

// FindRolePermissions 获取角色权限
func FindRolePermissions(roleID uint64) ([]model.Permission, error) {
	var perms []model.Permission
	err := database.DB.Raw(`
		SELECT p.* FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
	`, roleID).Scan(&perms).Error
	return perms, err
}

// CreateRole 创建角色
func CreateRole(role *model.Role, permissionIDs []uint64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}
		for _, pid := range permissionIDs {
			if err := tx.Create(&model.RolePermission{RoleID: role.ID, PermissionID: pid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateRole 更新角色字段
func UpdateRole(id uint64, updates map[string]interface{}) error {
	return database.DB.Model(&model.Role{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateRolePermissions 更新角色权限
func UpdateRolePermissions(roleID uint64, permissionIDs []uint64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}
		for _, pid := range permissionIDs {
			if err := tx.Create(&model.RolePermission{RoleID: roleID, PermissionID: pid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
