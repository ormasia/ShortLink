package scripts

import (
	"context"
	"shortLink/userservice/logger"
	"shortLink/userservice/model"
	"shortLink/userservice/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitRolesAndPermissions 初始化角色和权限
func InitRolesAndPermissions(db *gorm.DB) error {
	rbacRepo := repository.NewGormRBACRepository(db)

	// 创建管理员角色
	adminRole := &model.Role{
		Name:        "admin",
		Description: "系统管理员",
	}
	if err := db.FirstOrCreate(adminRole, model.Role{Name: "admin"}).Error; err != nil {
		logger.Log.Error("创建管理员角色失败", zap.Error(err))
		return err
	}

	// 创建普通用户角色
	userRole := &model.Role{
		Name:        "user",
		Description: "普通用户",
	}
	if err := db.FirstOrCreate(userRole, model.Role{Name: "user"}).Error; err != nil {
		logger.Log.Error("创建普通用户角色失败", zap.Error(err))
		return err
	}

	// 创建管理员权限
	adminPermissions := []model.Permission{
		{
			Name:        "admin:access",
			Description: "访问后台管理",
			Resource:    "admin",
			Action:      "access",
		},
		{
			Name:        "user:manage",
			Description: "管理用户",
			Resource:    "user",
			Action:      "manage",
		},
		{
			Name:        "link:manage",
			Description: "管理短链接",
			Resource:    "link",
			Action:      "manage",
		},
		{
			Name:        "role:manage",
			Description: "管理角色",
			Resource:    "role",
			Action:      "manage",
		},
		{
			Name:        "permission:manage",
			Description: "管理权限",
			Resource:    "permission",
			Action:      "manage",
		},
	}

	// 创建普通用户权限
	userPermissions := []model.Permission{
		{
			Name:        "link:create",
			Description: "创建短链接",
			Resource:    "link",
			Action:      "create",
		},
		{
			Name:        "link:read",
			Description: "查看短链接",
			Resource:    "link",
			Action:      "read",
		},
		{
			Name:        "link:update",
			Description: "更新短链接",
			Resource:    "link",
			Action:      "update",
		},
		{
			Name:        "link:delete",
			Description: "删除短链接",
			Resource:    "link",
			Action:      "delete",
		},
	}

	// 保存权限并关联到角色
	for _, perm := range adminPermissions {
		if err := db.FirstOrCreate(&perm, model.Permission{Name: perm.Name}).Error; err != nil {
			logger.Log.Error("创建管理员权限失败", zap.Error(err))
			return err
		}
		// 关联权限到管理员角色
		if err := rbacRepo.AssignPermissionToRole(context.Background(), adminRole.ID, perm.ID); err != nil {
			logger.Log.Error("关联管理员权限失败", zap.Error(err))
			return err
		}
	}

	for _, perm := range userPermissions {
		if err := db.FirstOrCreate(&perm, model.Permission{Name: perm.Name}).Error; err != nil {
			logger.Log.Error("创建用户权限失败", zap.Error(err))
			return err
		}
		// 关联权限到普通用户角色
		if err := rbacRepo.AssignPermissionToRole(context.Background(), userRole.ID, perm.ID); err != nil {
			logger.Log.Error("关联用户权限失败", zap.Error(err))
			return err
		}
	}

	logger.Log.Info("角色和权限初始化成功")
	return nil
}
