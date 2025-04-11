package rbac

import (
	"context"
	"io"
	"testing"

	"shortLink/proto/userpb"
	"shortLink/userservice/logger"
	"shortLink/userservice/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MockRBACRepository 模拟RBACRepository
type MockRBACRepository struct {
	mock.Mock
}

// AssignPermissionToRole implements repository.RBACRepository.
func (m *MockRBACRepository) AssignPermissionToRole(ctx context.Context, roleID uint, permissionID uint) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

// AssignRoleToUser implements repository.RBACRepository.
func (m *MockRBACRepository) AssignRoleToUser(ctx context.Context, userID uint, roleID uint) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

// CreatePermission implements repository.RBACRepository.
func (m *MockRBACRepository) CreatePermission(ctx context.Context, permission *model.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

// CreateRole implements repository.RBACRepository.
func (m *MockRBACRepository) CreateRole(ctx context.Context, role *model.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

// GetPermissionByID implements repository.RBACRepository.
func (m *MockRBACRepository) GetPermissionByID(ctx context.Context, id uint) (*model.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

// GetRoleByID implements repository.RBACRepository.
func (m *MockRBACRepository) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

// HasPermission implements repository.RBACRepository.
func (m *MockRBACRepository) HasPermission(ctx context.Context, userID uint, resource string, action string) (bool, error) {
	args := m.Called(ctx, userID, resource, action)
	return args.Bool(0), args.Error(1)
}

func (m *MockRBACRepository) GetUserRoles(ctx context.Context, userID uint) ([]model.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Role), args.Error(1)
}

func (m *MockRBACRepository) GetRolePermissions(ctx context.Context, roleID uint) ([]model.Permission, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Permission), args.Error(1)
}

// MockRBACCacheRepository 模拟RBACCacheRepository
type MockRBACCacheRepository struct {
	mock.Mock
}

// DeleteRolePermissions implements repository.RBACCacheRepository.
func (m *MockRBACCacheRepository) DeleteRolePermissions(roleID uint) error {
	args := m.Called(roleID)
	return args.Error(0)
}

func (m *MockRBACCacheRepository) GetRolePermissions(roleID uint) ([]model.Permission, error) {
	args := m.Called(roleID)
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *MockRBACCacheRepository) SetRolePermissions(roleID uint, permissions []model.Permission) error {
	args := m.Called(roleID, permissions)
	return args.Error(0)
}

// 初始化测试环境的函数
func setupTestLogger() {
	// 创建一个测试用的logger，输出到内存而不是控制台或文件
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(io.Discard), zapcore.DebugLevel)
	logger.Log = zap.New(core)
}
func TestRBACService_GetUserRoles(t *testing.T) {
	// 初始化日志系统
	setupTestLogger()
	// 初始化mock对象
	rbacRepo := new(MockRBACRepository)
	cacheRepo := new(MockRBACCacheRepository)
	service := NewRBACService(rbacRepo, cacheRepo)

	// 测试用例
	t.Run("成功获取用户角色", func(t *testing.T) {
		expectedRoles := []model.Role{{
			ID:          1,
			Name:        "admin",
			Description: "管理员",
		}}

		rbacRepo.On("GetUserRoles", mock.Anything, uint(1)).Return(expectedRoles, nil)

		roles, err := service.GetUserRoles(context.Background(), &userpb.GetUserRolesRequest{UserId: 1})
		assert.NoError(t, err)
		assert.Equal(t, uint32(1), roles.Roles[0].Id)
		assert.Equal(t, "admin", roles.Roles[0].Name)
		assert.Equal(t, "管理员", roles.Roles[0].Description)
		rbacRepo.AssertExpectations(t)
	})

	t.Run("获取用户角色失败", func(t *testing.T) {
		rbacRepo.On("GetUserRoles", mock.Anything, uint(2)).Return([]model.Role{}, assert.AnError)

		roles, err := service.GetUserRoles(context.Background(), &userpb.GetUserRolesRequest{UserId: 2})
		assert.Error(t, err)
		assert.Empty(t, roles.Roles)
		rbacRepo.AssertExpectations(t)
	})
}

func TestRBACService_GetRolePermissions(t *testing.T) {
	// 初始化日志系统
	setupTestLogger()
	// 初始化mock对象
	rbacRepo := new(MockRBACRepository)
	cacheRepo := new(MockRBACCacheRepository)
	service := NewRBACService(rbacRepo, cacheRepo)

	// 测试用例
	t.Run("从数据库获取权限成功", func(t *testing.T) {
		expectedPerms := []model.Permission{{
			ID:          2,
			Name:        "delete",
			Description: "删除权限",
			Resource:    "user",
			Action:      "delete",
		}}

		rbacRepo.On("GetRolePermissions", mock.Anything, uint(2)).Return(expectedPerms, nil)

		perms, err := service.GetRolePermissions(context.Background(), &userpb.GetRolePermissionsRequest{RoleId: 2})
		assert.NoError(t, err)
		assert.Equal(t, uint32(2), perms.Permissions[0].Id)
		assert.Equal(t, "delete", perms.Permissions[0].Name)
		assert.Equal(t, "user", perms.Permissions[0].Resource)
		assert.Equal(t, "delete", perms.Permissions[0].Action)
		rbacRepo.AssertExpectations(t)
	})

	t.Run("获取角色权限失败", func(t *testing.T) {
		rbacRepo.On("GetRolePermissions", mock.Anything, uint(3)).Return([]model.Permission{}, assert.AnError)

		perms, err := service.GetRolePermissions(context.Background(), &userpb.GetRolePermissionsRequest{RoleId: 3})
		assert.Error(t, err)
		assert.Empty(t, perms.Permissions)
		rbacRepo.AssertExpectations(t)
	})
}

func TestRBACService_CheckPermission(t *testing.T) {
	// 初始化日志系统
	setupTestLogger()
	// 初始化mock对象
	rbacRepo := new(MockRBACRepository)
	cacheRepo := new(MockRBACCacheRepository)
	service := NewRBACService(rbacRepo, cacheRepo)

	// 测试用例
	t.Run("用户有权限", func(t *testing.T) {
		rbacRepo.On("HasPermission", mock.Anything, uint(1), "user", "create").Return(true, nil)

		resp, err := service.CheckPermission(context.Background(), &userpb.CheckPermissionRequest{
			UserId:   1,
			Resource: "user",
			Action:   "create",
		})
		assert.NoError(t, err)
		assert.True(t, resp.HasPermission)
		rbacRepo.AssertExpectations(t)
	})

	t.Run("用户无权限", func(t *testing.T) {
		rbacRepo.On("HasPermission", mock.Anything, uint(1), "user", "delete").Return(false, nil)

		resp, err := service.CheckPermission(context.Background(), &userpb.CheckPermissionRequest{
			UserId:   1,
			Resource: "user",
			Action:   "delete",
		})
		assert.NoError(t, err)
		assert.False(t, resp.HasPermission)
		rbacRepo.AssertExpectations(t)
	})

	t.Run("检查权限出错", func(t *testing.T) {
		rbacRepo.On("HasPermission", mock.Anything, uint(2), "user", "create").Return(false, assert.AnError)

		resp, err := service.CheckPermission(context.Background(), &userpb.CheckPermissionRequest{
			UserId:   2,
			Resource: "user",
			Action:   "create",
		})
		assert.Error(t, err)
		assert.False(t, resp.HasPermission)
		rbacRepo.AssertExpectations(t)
	})
}

func TestRBACService_CreateRole(t *testing.T) {
	// 初始化日志系统
	setupTestLogger()
	// 初始化mock对象
	rbacRepo := new(MockRBACRepository)
	cacheRepo := new(MockRBACCacheRepository)
	service := NewRBACService(rbacRepo, cacheRepo)

	// 测试用例
	t.Run("创建角色成功", func(t *testing.T) {
		rbacRepo.On("CreateRole", mock.Anything, mock.MatchedBy(func(role *model.Role) bool {
			return role.Name == "editor" && role.Description == "编辑者"
		})).Run(func(args mock.Arguments) {
			role := args.Get(1).(*model.Role)
			role.ID = 3 // 模拟数据库自动生成ID
		}).Return(nil)

		resp, err := service.CreateRole(context.Background(), &userpb.CreateRoleRequest{
			Name:        "editor",
			Description: "编辑者",
		})

		assert.NoError(t, err)
		assert.Equal(t, uint32(3), resp.Role.Id)
		assert.Equal(t, "editor", resp.Role.Name)
		assert.Equal(t, "编辑者", resp.Role.Description)
		assert.Contains(t, resp.Message, "成功")
		rbacRepo.AssertExpectations(t)
	})

	t.Run("创建角色失败", func(t *testing.T) {
		rbacRepo.On("CreateRole", mock.Anything, mock.MatchedBy(func(role *model.Role) bool {
			return role.Name == "duplicate" && role.Description == "重复角色"
		})).Return(assert.AnError)

		resp, err := service.CreateRole(context.Background(), &userpb.CreateRoleRequest{
			Name:        "duplicate",
			Description: "重复角色",
		})

		assert.Error(t, err)
		assert.Contains(t, resp.Message, "失败")
		rbacRepo.AssertExpectations(t)
	})
}

func TestRBACService_CreatePermission(t *testing.T) {
	// 初始化日志系统
	setupTestLogger()
	// 初始化mock对象
	rbacRepo := new(MockRBACRepository)
	cacheRepo := new(MockRBACCacheRepository)
	service := NewRBACService(rbacRepo, cacheRepo)

	// 测试用例
	t.Run("创建权限成功", func(t *testing.T) {
		rbacRepo.On("CreatePermission", mock.Anything, mock.MatchedBy(func(perm *model.Permission) bool {
			return perm.Name == "update" &&
				perm.Description == "更新权限" &&
				perm.Resource == "user" &&
				perm.Action == "update"
		})).Run(func(args mock.Arguments) {
			perm := args.Get(1).(*model.Permission)
			perm.ID = 3 // 模拟数据库自动生成ID
		}).Return(nil)

		resp, err := service.CreatePermission(context.Background(), &userpb.CreatePermissionRequest{
			Name:        "update",
			Description: "更新权限",
			Resource:    "user",
			Action:      "update",
		})

		assert.NoError(t, err)
		assert.Equal(t, uint32(3), resp.Permission.Id)
		assert.Equal(t, "update", resp.Permission.Name)
		assert.Equal(t, "更新权限", resp.Permission.Description)
		assert.Equal(t, "user", resp.Permission.Resource)
		assert.Equal(t, "update", resp.Permission.Action)
		assert.Contains(t, resp.Message, "成功")
		rbacRepo.AssertExpectations(t)
	})

	t.Run("创建权限失败", func(t *testing.T) {
		rbacRepo.On("CreatePermission", mock.Anything, mock.MatchedBy(func(perm *model.Permission) bool {
			return perm.Name == "duplicate" && perm.Resource == "user" && perm.Action == "read"
		})).Return(assert.AnError)

		resp, err := service.CreatePermission(context.Background(), &userpb.CreatePermissionRequest{
			Name:        "duplicate",
			Description: "重复权限",
			Resource:    "user",
			Action:      "read",
		})

		assert.Error(t, err)
		assert.Contains(t, resp.Message, "失败")
		rbacRepo.AssertExpectations(t)
	})
}

func TestRBACService_AssignRoleToUser(t *testing.T) {
	// 初始化日志系统
	setupTestLogger()
	// 初始化mock对象
	rbacRepo := new(MockRBACRepository)
	cacheRepo := new(MockRBACCacheRepository)
	service := NewRBACService(rbacRepo, cacheRepo)

	// 测试用例
	t.Run("分配角色给用户成功", func(t *testing.T) {
		rbacRepo.On("AssignRoleToUser", mock.Anything, uint(1), uint(2)).Return(nil)

		resp, err := service.AssignRoleToUser(context.Background(), &userpb.AssignRoleToUserRequest{
			UserId: 1,
			RoleId: 2,
		})

		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "成功")
		rbacRepo.AssertExpectations(t)
	})

	t.Run("分配角色给用户失败", func(t *testing.T) {
		rbacRepo.On("AssignRoleToUser", mock.Anything, uint(3), uint(4)).Return(assert.AnError)

		resp, err := service.AssignRoleToUser(context.Background(), &userpb.AssignRoleToUserRequest{
			UserId: 3,
			RoleId: 4,
		})

		assert.Error(t, err)
		assert.Contains(t, resp.Message, "失败")
		rbacRepo.AssertExpectations(t)
	})
}

func TestRBACService_AssignPermissionToRole(t *testing.T) {
	// 初始化日志系统
	setupTestLogger()
	// 初始化mock对象
	rbacRepo := new(MockRBACRepository)
	cacheRepo := new(MockRBACCacheRepository)
	service := NewRBACService(rbacRepo, cacheRepo)

	// 测试用例
	t.Run("分配权限给角色成功", func(t *testing.T) {
		rbacRepo.On("AssignPermissionToRole", mock.Anything, uint(1), uint(2)).Return(nil)
		cacheRepo.On("DeleteRolePermissions", uint(1)).Return(nil)

		resp, err := service.AssignPermissionToRole(context.Background(), &userpb.AssignPermissionToRoleRequest{
			RoleId:       1,
			PermissionId: 2,
		})

		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "成功")
		rbacRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
	})

	t.Run("分配权限给角色失败", func(t *testing.T) {
		rbacRepo.On("AssignPermissionToRole", mock.Anything, uint(3), uint(4)).Return(assert.AnError)

		resp, err := service.AssignPermissionToRole(context.Background(), &userpb.AssignPermissionToRoleRequest{
			RoleId:       3,
			PermissionId: 4,
		})

		assert.Error(t, err)
		assert.Contains(t, resp.Message, "失败")
		rbacRepo.AssertExpectations(t)
	})

	t.Run("分配权限成功但缓存删除失败", func(t *testing.T) {
		rbacRepo.On("AssignPermissionToRole", mock.Anything, uint(5), uint(6)).Return(nil)
		cacheRepo.On("DeleteRolePermissions", uint(5)).Return(assert.AnError)

		resp, err := service.AssignPermissionToRole(context.Background(), &userpb.AssignPermissionToRoleRequest{
			RoleId:       5,
			PermissionId: 6,
		})

		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "成功")
		rbacRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
	})
}
