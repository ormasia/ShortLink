package service

import (
	"context"
	"fmt"
	"shortLink/proto/userpb"
	"shortLink/userservice/cache"
	"shortLink/userservice/logger"
	"shortLink/userservice/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RBACService 实现RBAC相关的gRPC服务
type RBACService struct {
	userpb.UnimplementedRBACServiceServer
	DB *gorm.DB
}

// NewRBACService 创建RBAC服务实例
func NewRBACService(db *gorm.DB) *RBACService {
	return &RBACService{DB: db}
}

// GetUserRoles 获取用户的所有角色
func GetUserRoles(db *gorm.DB, userID uint) ([]model.Role, error) {
	var user model.User
	if err := db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Roles, nil
}

// GetRolePermissions 获取角色的所有权限
func GetRolePermissions(db *gorm.DB, roleID uint) ([]model.Permission, error) {
	// 优先从缓存中获取权限信息
	permissions, err := cache.GetRolePermissions(roleID)
	if err == nil {
		// 将缓存中的权限类型转换为服务所需的类型
		var servicePermissions []model.Permission
		for _, p := range permissions {
			servicePermissions = append(servicePermissions, p)
		}
		return servicePermissions, nil
	}

	// 缓存未命中，从数据库获取
	var permissions1 []model.Permission
	// 通过中间表查询角色对应的权限
	if err := db.Table("permissions").Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").Where("role_permissions.role_id = ?", roleID).Find(&permissions1).Error; err != nil {
		return nil, err
	}

	// 更新缓存
	// 将 model.Permission 类型转换为 cache 包需要的类型
	var cachePermissions []model.Permission
	for _, p := range permissions1 {
		cachePermissions = append(cachePermissions, p)
	}
	// 将 model.Permission 转换为 cache 包需要的类型
	var cacheModelPermissions []model.Permission
	for _, p := range cachePermissions {
		cacheModelPermissions = append(cacheModelPermissions, model.Permission(p))
	}
	if err := cache.SetRolePermissions(roleID, cacheModelPermissions); err != nil {
		// 缓存更新失败只记录日志，不影响正常流程
		logger.Log.Warn("更新角色权限缓存失败", zap.Error(err))
	}

	return permissions1, nil
}

// HasPermission 检查用户是否有指定的权限
func HasPermission(db *gorm.DB, userID uint, requiredResource, requiredAction string) bool {
	roles, err := GetUserRoles(db, userID)
	if err != nil {
		logger.Log.Warn("获取用户角色失败", zap.Uint("user_id", userID), zap.Error(err))
		return false
	}

	for _, role := range roles {
		permissions, err := GetRolePermissions(db, role.ID)
		if err != nil {
			logger.Log.Warn("获取角色权限失败", zap.Uint("role_id", role.ID), zap.Error(err))
			continue
		}

		for _, perm := range permissions {
			if perm.Resource == requiredResource && perm.Action == requiredAction {
				return true
			}
		}
	}

	return false
}

// CheckPermission 实现gRPC接口，检查用户是否有指定的权限
func (s *RBACService) CheckPermission(ctx context.Context, req *userpb.CheckPermissionRequest) (*userpb.CheckPermissionResponse, error) {
	userID := uint(req.UserId)
	resource := req.Resource
	action := req.Action

	hasPermission := HasPermission(s.DB, userID, resource, action)

	var message string
	if hasPermission {
		message = "用户有权限执行该操作"
		logger.Log.Info("权限检查通过",
			zap.Uint("user_id", userID),
			zap.String("resource", resource),
			zap.String("action", action))
	} else {
		message = fmt.Sprintf("用户没有权限执行该操作: %s %s", resource, action)
		logger.Log.Warn("权限检查失败",
			zap.Uint("user_id", userID),
			zap.String("resource", resource),
			zap.String("action", action))
	}

	return &userpb.CheckPermissionResponse{
		HasPermission: hasPermission,
		Message:       message,
	}, nil
}

// GetUserRoles 实现gRPC接口，获取用户的所有角色
func (s *RBACService) GetUserRoles(ctx context.Context, req *userpb.GetUserRolesRequest) (*userpb.GetUserRolesResponse, error) {
	userID := uint(req.UserId)

	roles, err := GetUserRoles(s.DB, userID)
	if err != nil {
		logger.Log.Error("获取用户角色失败", zap.Uint("user_id", userID), zap.Error(err))
		return &userpb.GetUserRolesResponse{}, err
	}

	var roleInfos []*userpb.RoleInfo
	for _, role := range roles {
		roleInfos = append(roleInfos, &userpb.RoleInfo{
			Id:          uint32(role.ID),
			Name:        role.Name,
			Description: role.Description,
		})
	}

	return &userpb.GetUserRolesResponse{
		Roles: roleInfos,
	}, nil
}

// GetRolePermissions 实现gRPC接口，获取角色的所有权限
func (s *RBACService) GetRolePermissions(ctx context.Context, req *userpb.GetRolePermissionsRequest) (*userpb.GetRolePermissionsResponse, error) {
	roleID := uint(req.RoleId)

	permissions, err := GetRolePermissions(s.DB, roleID)
	if err != nil {
		logger.Log.Error("获取角色权限失败", zap.Uint("role_id", roleID), zap.Error(err))
		return &userpb.GetRolePermissionsResponse{}, err
	}

	var permissionInfos []*userpb.PermissionInfo
	for _, perm := range permissions {
		permissionInfos = append(permissionInfos, &userpb.PermissionInfo{
			Id:          uint32(perm.ID),
			Name:        perm.Name,
			Description: perm.Description,
			Resource:    perm.Resource,
			Action:      perm.Action,
		})
	}

	return &userpb.GetRolePermissionsResponse{
		Permissions: permissionInfos,
	}, nil
}

// CreateRole 实现gRPC接口，创建新角色
func (s *RBACService) CreateRole(ctx context.Context, req *userpb.CreateRoleRequest) (*userpb.CreateRoleResponse, error) {
	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.DB.Create(&role).Error; err != nil {
		logger.Log.Error("创建角色失败", zap.String("name", req.Name), zap.Error(err))
		return &userpb.CreateRoleResponse{
			Message: "创建角色失败",
		}, err
	}

	logger.Log.Info("创建角色成功", zap.String("name", role.Name), zap.Uint("id", role.ID))
	return &userpb.CreateRoleResponse{
		Role: &userpb.RoleInfo{
			Id:          uint32(role.ID),
			Name:        role.Name,
			Description: role.Description,
		},
		Message: "创建角色成功",
	}, nil
}

// CreatePermission 实现gRPC接口，创建新权限
func (s *RBACService) CreatePermission(ctx context.Context, req *userpb.CreatePermissionRequest) (*userpb.CreatePermissionResponse, error) {
	permission := model.Permission{
		Name:        req.Name,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
	}

	if err := s.DB.Create(&permission).Error; err != nil {
		logger.Log.Error("创建权限失败", zap.String("name", req.Name), zap.Error(err))
		return &userpb.CreatePermissionResponse{
			Message: "创建权限失败",
		}, err
	}

	logger.Log.Info("创建权限成功",
		zap.String("name", permission.Name),
		zap.Uint("id", permission.ID),
		zap.String("resource", permission.Resource),
		zap.String("action", permission.Action))
	return &userpb.CreatePermissionResponse{
		Permission: &userpb.PermissionInfo{
			Id:          uint32(permission.ID),
			Name:        permission.Name,
			Description: permission.Description,
			Resource:    permission.Resource,
			Action:      permission.Action,
		},
		Message: "创建权限成功",
	}, nil
}

// AssignRoleToUser 实现gRPC接口，为用户分配角色
func (s *RBACService) AssignRoleToUser(ctx context.Context, req *userpb.AssignRoleToUserRequest) (*userpb.AssignRoleToUserResponse, error) {
	userRole := model.UserRole{
		UserID: uint(req.UserId),
		RoleID: uint(req.RoleId),
	}

	if err := s.DB.Create(&userRole).Error; err != nil {
		logger.Log.Error("分配角色给用户失败",
			zap.Uint("user_id", userRole.UserID),
			zap.Uint("role_id", userRole.RoleID),
			zap.Error(err))
		return &userpb.AssignRoleToUserResponse{
			Message: "分配角色给用户失败",
		}, err
	}

	logger.Log.Info("分配角色给用户成功",
		zap.Uint("user_id", userRole.UserID),
		zap.Uint("role_id", userRole.RoleID))
	return &userpb.AssignRoleToUserResponse{
		Message: "分配角色给用户成功",
	}, nil
}

// AssignPermissionToRole 实现gRPC接口，为角色分配权限
func (s *RBACService) AssignPermissionToRole(ctx context.Context, req *userpb.AssignPermissionToRoleRequest) (*userpb.AssignPermissionToRoleResponse, error) {
	rolePermission := model.RolePermission{
		RoleID:       uint(req.RoleId),
		PermissionID: uint(req.PermissionId),
	}

	if err := s.DB.Create(&rolePermission).Error; err != nil {
		logger.Log.Error("分配权限给角色失败",
			zap.Uint("role_id", rolePermission.RoleID),
			zap.Uint("permission_id", rolePermission.PermissionID),
			zap.Error(err))
		return &userpb.AssignPermissionToRoleResponse{
			Message: "分配权限给角色失败",
		}, err
	}

	// 删除角色权限缓存，确保下次获取时能获取到最新的权限
	if err := cache.DeleteRolePermissions(rolePermission.RoleID); err != nil {
		logger.Log.Warn("删除角色权限缓存失败",
			zap.Uint("role_id", rolePermission.RoleID),
			zap.Error(err))
	}

	logger.Log.Info("分配权限给角色成功",
		zap.Uint("role_id", rolePermission.RoleID),
		zap.Uint("permission_id", rolePermission.PermissionID))
	return &userpb.AssignPermissionToRoleResponse{
		Message: "分配权限给角色成功",
	}, nil
}
