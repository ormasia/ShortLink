package repository

import (
	"context"
	"shortLink/userservice/model"

	"gorm.io/gorm"
)

// RBACRepository 定义RBAC数据访问接口
type RBACRepository interface {
	CreateRole(ctx context.Context, role *model.Role) error
	GetRoleByID(ctx context.Context, id uint) (*model.Role, error)
	GetUserRoles(ctx context.Context, userID uint) ([]model.Role, error)
	CreatePermission(ctx context.Context, permission *model.Permission) error
	GetPermissionByID(ctx context.Context, id uint) (*model.Permission, error)
	GetRolePermissions(ctx context.Context, roleID uint) ([]model.Permission, error)
	AssignRoleToUser(ctx context.Context, userID, roleID uint) error
	AssignPermissionToRole(ctx context.Context, roleID, permissionID uint) error
	HasPermission(ctx context.Context, userID uint, resource, action string) (bool, error)
}

// GormRBACRepository 实现基于Gorm的RBAC数据访问
type GormRBACRepository struct {
	db *gorm.DB
}

// NewGormRBACRepository 创建一个新的GormRBACRepository实例
func NewGormRBACRepository(db *gorm.DB) *GormRBACRepository {
	return &GormRBACRepository{db: db}
}

// CreateRole 创建新角色
func (r *GormRBACRepository) CreateRole(ctx context.Context, role *model.Role) error {
	return r.db.Create(role).Error
}

// GetRoleByID 根据ID获取角色
func (r *GormRBACRepository) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetUserRoles 获取用户的所有角色
func (r *GormRBACRepository) GetUserRoles(ctx context.Context, userID uint) ([]model.Role, error) {
	var user model.User
	if err := r.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Roles, nil
}

// CreatePermission 创建新权限
func (r *GormRBACRepository) CreatePermission(ctx context.Context, permission *model.Permission) error {
	return r.db.Create(permission).Error
}

// GetPermissionByID 根据ID获取权限
func (r *GormRBACRepository) GetPermissionByID(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	if err := r.db.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetRolePermissions 获取角色的所有权限
func (r *GormRBACRepository) GetRolePermissions(ctx context.Context, roleID uint) ([]model.Permission, error) {
	var permissions []model.Permission
	if err := r.db.Table("permissions").Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").Where("role_permissions.role_id = ?", roleID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// AssignRoleToUser 为用户分配角色
func (r *GormRBACRepository) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	userRole := model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(&userRole).Error
}

// AssignPermissionToRole 为角色分配权限
func (r *GormRBACRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID uint) error {
	rolePermission := model.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.db.Create(&rolePermission).Error
}

// HasPermission 检查用户是否有指定的权限
func (r *GormRBACRepository) HasPermission(ctx context.Context, userID uint, requiredResource, requiredAction string) (bool, error) {
	roles, err := r.GetUserRoles(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		permissions, err := r.GetRolePermissions(ctx, role.ID)
		if err != nil {
			continue
		}

		for _, perm := range permissions {
			if perm.Resource == requiredResource && perm.Action == requiredAction {
				return true, nil
			}
		}
	}

	return false, nil
}
