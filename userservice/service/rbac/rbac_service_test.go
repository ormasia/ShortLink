package rbac

import (
	"context"
	"testing"

	"shortLink/userservice/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestRBACService_GetUserRoles(t *testing.T) {
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

		roles, err := service.getUserRoles(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedRoles, roles)
		rbacRepo.AssertExpectations(t)
	})

	t.Run("获取用户角色失败", func(t *testing.T) {
		rbacRepo.On("GetUserRoles", mock.Anything, uint(2)).Return([]model.Role{}, assert.AnError)

		roles, err := service.getUserRoles(context.Background(), 2)
		assert.Error(t, err)
		assert.Empty(t, roles)
		rbacRepo.AssertExpectations(t)
	})
}

func TestRBACService_GetRolePermissions(t *testing.T) {
	// 初始化mock对象
	rbacRepo := new(MockRBACRepository)
	cacheRepo := new(MockRBACCacheRepository)
	service := NewRBACService(rbacRepo, cacheRepo)

	// 测试用例
	t.Run("从缓存获取权限成功", func(t *testing.T) {
		expectedPerms := []model.Permission{{
			ID:          1,
			Name:        "create",
			Description: "创建权限",
			Resource:    "user",
			Action:      "create",
		}}

		cacheRepo.On("GetRolePermissions", uint(1)).Return(expectedPerms, nil)

		perms, err := service.getRolePermissions(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedPerms, perms)
		cacheRepo.AssertExpectations(t)
		rbacRepo.AssertNotCalled(t, "GetRolePermissions")
	})

	t.Run("从数据库获取权限成功", func(t *testing.T) {
		expectedPerms := []model.Permission{{
			ID:          2,
			Name:        "delete",
			Description: "删除权限",
			Resource:    "user",
			Action:      "delete",
		}}

		cacheRepo.On("GetRolePermissions", uint(2)).Return([]model.Permission{}, assert.AnError)
		rbacRepo.On("GetRolePermissions", mock.Anything, uint(2)).Return(expectedPerms, nil)
		cacheRepo.On("SetRolePermissions", uint(2), expectedPerms).Return(nil)

		perms, err := service.getRolePermissions(context.Background(), 2)
		assert.NoError(t, err)
		assert.Equal(t, expectedPerms, perms)
		cacheRepo.AssertExpectations(t)
		rbacRepo.AssertExpectations(t)
	})
}

func TestRBACService_HasPermission(t *testing.T) {
	// 初始化mock对象
	rbacRepo := new(MockRBACRepository)
	cacheRepo := new(MockRBACCacheRepository)
	service := NewRBACService(rbacRepo, cacheRepo)

	// 测试用例
	t.Run("用户有权限", func(t *testing.T) {
		roles := []model.Role{{
			ID:          1,
			Name:        "admin",
			Description: "管理员",
		}}

		perms := []model.Permission{{
			ID:          1,
			Name:        "create",
			Description: "创建权限",
			Resource:    "user",
			Action:      "create",
		}}

		rbacRepo.On("GetUserRoles", mock.Anything, uint(1)).Return(roles, nil)
		cacheRepo.On("GetRolePermissions", uint(1)).Return(perms, nil)

		hasPerm := service.hasPermission(context.Background(), 1, "user", "create")
		assert.True(t, hasPerm)
		rbacRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
	})

	t.Run("用户无权限", func(t *testing.T) {
		roles := []model.Role{{
			ID:          1,
			Name:        "admin",
			Description: "管理员",
		}}

		perms := []model.Permission{{
			ID:          1,
			Name:        "create",
			Description: "创建权限",
			Resource:    "user",
			Action:      "create",
		}}

		rbacRepo.On("GetUserRoles", mock.Anything, uint(1)).Return(roles, nil)
		cacheRepo.On("GetRolePermissions", uint(1)).Return(perms, nil)

		hasPerm := service.hasPermission(context.Background(), 1, "user", "delete")
		assert.False(t, hasPerm)
		rbacRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
	})
}
