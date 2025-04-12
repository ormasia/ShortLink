package user

import (
	"context"
	"errors"
	"io"
	"shortLink/proto/userpb"
	"shortLink/userservice/logger"
	"shortLink/userservice/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 模拟UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// 模拟AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) VerifyPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func (m *MockAuthService) GenerateToken(userID uint, role string) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GetUserIDFromToken(token string) (uint, error) {
	args := m.Called(token)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockAuthService) InvalidateToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

// 初始化测试环境的函数
func setupTestLogger() {
	// 创建一个测试用的logger，输出到内存而不是控制台或文件
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(io.Discard), zapcore.DebugLevel)
	logger.Log = zap.New(core)
}

// 测试Register方法
func TestUserService_Register(t *testing.T) {
	// 初始化日志系统
	setupTestLogger()

	// 创建模拟对象
	mockUserRepo := new(MockUserRepository)
	mockAuthService := new(MockAuthService)

	// 创建UserService实例
	userService := NewUserService(mockUserRepo, mockAuthService)

	// 测试用例1: 用户名已存在
	t.Run("Username already exists", func(t *testing.T) {
		// 设置模拟行为
		existingUser := &model.User{ID: 1, Username: "existinguser"}
		mockUserRepo.On("FindByUsername", "existinguser").Return(existingUser, nil).Once()

		// 执行测试
		req := &userpb.RegisterRequest{
			Username: "existinguser",
			Password: "password123",
			Nickname: "Existing User",
			Email:    "existing@example.com",
		}
		resp, err := userService.Register(context.Background(), req)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, "用户名已存在", resp.Message)

		// 验证模拟对象的调用
		mockUserRepo.AssertExpectations(t)
	})

	// 测试用例2: 数据库查询错误
	t.Run("Database query error", func(t *testing.T) {
		// 设置模拟行为
		dbError := errors.New("database error")
		mockUserRepo.On("FindByUsername", "newuser").Return(nil, dbError).Once()

		// 执行测试
		req := &userpb.RegisterRequest{
			Username: "newuser",
			Password: "password123",
			Nickname: "New User",
			Email:    "new@example.com",
		}
		resp, err := userService.Register(context.Background(), req)

		// 验证结果
		assert.Equal(t, dbError, err)
		assert.Equal(t, "注册失败", resp.Message)

		// 验证模拟对象的调用
		mockUserRepo.AssertExpectations(t)
	})

	// 测试用例3: 注册成功
	t.Run("Registration success", func(t *testing.T) {
		// 设置模拟行为
		mockUserRepo.On("FindByUsername", "newuser").Return(nil, gorm.ErrRecordNotFound).Once()
		mockUserRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil).Once()

		// 执行测试
		req := &userpb.RegisterRequest{
			Username: "newuser",
			Password: "password123",
			Nickname: "New User",
			Email:    "new@example.com",
		}
		resp, err := userService.Register(context.Background(), req)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, "注册成功", resp.Message)

		// 验证模拟对象的调用
		mockUserRepo.AssertExpectations(t)
	})
}

// 测试Login方法
func TestUserService_Login(t *testing.T) {
	// 初始化日志系统
	setupTestLogger()

	// 创建模拟对象
	mockUserRepo := new(MockUserRepository)
	mockAuthService := new(MockAuthService)

	// 创建UserService实例
	userService := NewUserService(mockUserRepo, mockAuthService)

	// 测试用例1: 用户不存在
	t.Run("User not found", func(t *testing.T) {
		// 设置模拟行为
		userNotFoundErr := gorm.ErrRecordNotFound
		mockUserRepo.On("FindByUsername", "nonexistentuser").Return(nil, userNotFoundErr).Once()

		// 执行测试
		req := &userpb.LoginRequest{
			Username: "nonexistentuser",
			Password: "password123",
		}
		resp, err := userService.Login(context.Background(), req)

		// 验证结果
		assert.Equal(t, userNotFoundErr, err)
		assert.Empty(t, resp.Token)
		assert.Nil(t, resp.User)

		// 验证模拟对象的调用
		mockUserRepo.AssertExpectations(t)
	})

	// 测试用例2: 密码错误
	t.Run("Incorrect password", func(t *testing.T) {
		// 设置模拟行为
		user := &model.User{
			ID:       1,
			Username: "testuser",
			Password: "hashedpassword",
			Role:     "user",
		}
		passwordErr := bcrypt.ErrMismatchedHashAndPassword
		mockUserRepo.On("FindByUsername", "testuser").Return(user, nil).Once()
		mockAuthService.On("VerifyPassword", "hashedpassword", "wrongpassword").Return(passwordErr).Once()

		// 执行测试
		req := &userpb.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}
		resp, err := userService.Login(context.Background(), req)

		// 验证结果
		assert.Equal(t, passwordErr, err)
		assert.Empty(t, resp.Token)
		assert.Nil(t, resp.User)

		// 验证模拟对象的调用
		mockUserRepo.AssertExpectations(t)
		mockAuthService.AssertExpectations(t)
	})

	// 测试用例3: 登录成功
	t.Run("Login success", func(t *testing.T) {
		// 设置模拟行为
		user := &model.User{
			ID:       1,
			Username: "testuser",
			Password: "hashedpassword",
			Nickname: "Test User",
			Email:    "test@example.com",
			Role:     "user",
		}
		mockUserRepo.On("FindByUsername", "testuser").Return(user, nil).Once()
		mockAuthService.On("VerifyPassword", "hashedpassword", "correctpassword").Return(nil).Once()
		mockAuthService.On("GenerateToken", uint(1), "user").Return("valid-token", nil).Once()

		// 执行测试
		req := &userpb.LoginRequest{
			Username: "testuser",
			Password: "correctpassword",
		}
		resp, err := userService.Login(context.Background(), req)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, "valid-token", resp.Token)
		assert.NotNil(t, resp.User)
		assert.Equal(t, uint32(1), resp.User.UserId)
		assert.Equal(t, "testuser", resp.User.Username)
		assert.Equal(t, "Test User", resp.User.Nickname)
		assert.Equal(t, "test@example.com", resp.User.Email)
		assert.Equal(t, "user", resp.User.Role)

		// 验证模拟对象的调用
		mockUserRepo.AssertExpectations(t)
		mockAuthService.AssertExpectations(t)
	})
}

// 测试Logout方法
func TestUserService_Logout(t *testing.T) {
	// 初始化日志系统
	setupTestLogger()

	// 创建模拟对象
	mockUserRepo := new(MockUserRepository)
	mockAuthService := new(MockAuthService)

	// 创建UserService实例
	userService := NewUserService(mockUserRepo, mockAuthService)

	// 测试用例1: 无效的令牌
	t.Run("Invalid token", func(t *testing.T) {
		// 设置模拟行为
		tokenErr := errors.New("invalid token")
		mockAuthService.On("GetUserIDFromToken", "invalid-token").Return(uint(0), tokenErr).Once()

		// 执行测试
		req := &userpb.LogoutRequest{
			Token: "invalid-token",
		}
		resp, err := userService.Logout(context.Background(), req)

		// 验证结果
		assert.Equal(t, tokenErr, err)
		assert.Equal(t, "登出失败：无效的令牌", resp.Message)

		// 验证模拟对象的调用
		mockAuthService.AssertExpectations(t)
	})

	// 测试用例2: 删除令牌错误
	t.Run("Token deletion error", func(t *testing.T) {
		// 设置模拟行为
		deleteErr := errors.New("token deletion error")
		mockAuthService.On("GetUserIDFromToken", "valid-token").Return(uint(1), nil).Once()
		mockAuthService.On("InvalidateToken", "valid-token").Return(deleteErr).Once()

		// 执行测试
		req := &userpb.LogoutRequest{
			Token: "valid-token",
		}
		resp, err := userService.Logout(context.Background(), req)

		// 验证结果
		assert.Equal(t, deleteErr, err)
		assert.Equal(t, "登出失败", resp.Message)

		// 验证模拟对象的调用
		mockAuthService.AssertExpectations(t)
	})

	// 测试用例3: 登出成功
	t.Run("Logout success", func(t *testing.T) {
		// 设置模拟行为
		mockAuthService.On("GetUserIDFromToken", "valid-token").Return(uint(1), nil).Once()
		mockAuthService.On("InvalidateToken", "valid-token").Return(nil).Once()

		// 执行测试
		req := &userpb.LogoutRequest{
			Token: "valid-token",
		}
		resp, err := userService.Logout(context.Background(), req)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, "登出成功", resp.Message)

		// 验证模拟对象的调用
		mockAuthService.AssertExpectations(t)
	})
}
