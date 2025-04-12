package user

import (
	"context"
	"fmt"
	"shortLink/proto/userpb"
	"shortLink/userservice/logger"
	"shortLink/userservice/model"
	"shortLink/userservice/repository"
	"shortLink/userservice/service/auth"
	"strconv"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 实现用户服务
type UserService struct {
	userpb.UnimplementedUserServiceServer
	userRepository repository.UserRepository
	authService    auth.AuthService
}

// NewUserService 创建一个新的UserService实例
func NewUserService(userRepository repository.UserRepository, authService auth.AuthService) *UserService {
	return &UserService{
		userRepository: userRepository,
		authService:    authService,
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	// 检查用户名是否已存在
	_, err := s.userRepository.FindByUsername(req.Username)
	if err == nil {
		// 用户已存在
		logger.Log.Warn("用户注册失败：用户名已存在",
			zap.String("username", req.Username),
			zap.String("email", req.Email),
		)
		return &userpb.RegisterResponse{Message: "用户名已存在"}, nil
	} else if err != gorm.ErrRecordNotFound {
		// 数据库错误
		logger.Log.Error("用户注册失败：数据库查询错误",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		return &userpb.RegisterResponse{Message: "注册失败"}, err
	}

	// 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("用户注册失败：密码加密错误",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		return &userpb.RegisterResponse{Message: "注册失败"}, err
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: string(hash),
		Nickname: req.Nickname,
		Email:    req.Email,
		Role:     "user", // 默认角色
	}

	// 保存用户
	err = s.userRepository.Create(user)
	if err != nil {
		logger.Log.Error("用户注册失败：数据库错误",
			zap.String("username", req.Username),
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return &userpb.RegisterResponse{Message: "注册失败"}, err
	}

	logger.Log.Info("用户注册成功",
		zap.String("uid", strconv.FormatUint(uint64(user.ID), 10)),
		zap.String("username", user.Username),
		zap.String("nickname", user.Nickname),
		zap.String("email", user.Email),
	)
	return &userpb.RegisterResponse{Message: "注册成功"}, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	// 查询用户
	user, err := s.userRepository.FindByUsername(req.Username)
	if err != nil {
		logger.Log.Warn("用户登录失败：用户不存在",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		return &userpb.LoginResponse{
			Token: "",
			User:  nil,
		}, err
	}

	// 验证密码
	err = s.authService.VerifyPassword(user.Password, req.Password)
	if err != nil {
		logger.Log.Warn("用户登录失败：密码错误",
			zap.String("uid", strconv.FormatUint(uint64(user.ID), 10)),
			zap.String("username", user.Username),
			zap.Error(err),
		)
		return &userpb.LoginResponse{
			Token: "",
			User:  nil,
		}, err
	}
	fmt.Println("user.ID", user.ID)
	// 生成令牌
	token, err := s.authService.GenerateToken(user.ID, user.Role)
	if err != nil {
		logger.Log.Error("用户登录失败：生成令牌错误",
			zap.String("uid", strconv.FormatUint(uint64(user.ID), 10)),
			zap.String("username", user.Username),
			zap.Error(err),
		)
		return &userpb.LoginResponse{
			Token: "",
			User:  nil,
		}, err
	}

	logger.Log.Info("用户登录成功",
		zap.String("uid", strconv.FormatUint(uint64(user.ID), 10)),
		zap.String("username", user.Username),
		zap.String("nickname", user.Nickname),
		zap.String("email", user.Email),
		zap.String("role", user.Role),
	)

	return &userpb.LoginResponse{
		Token: token,
		User: &userpb.UserInfo{
			UserId:   uint32(user.ID),
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}

// Logout 用户登出
func (s *UserService) Logout(ctx context.Context, req *userpb.LogoutRequest) (*userpb.LogoutResponse, error) {
	// 获取用户ID
	userID, err := s.authService.GetUserIDFromToken(req.Token)
	if err != nil {
		logger.Log.Warn("用户登出失败：无效的令牌",
			zap.String("token", req.Token),
			zap.Error(err),
		)
		return &userpb.LogoutResponse{Message: "登出失败：无效的令牌"}, err
	}

	// 删除令牌
	err = s.authService.InvalidateToken(req.Token)
	if err != nil {
		logger.Log.Error("用户登出失败：删除令牌错误",
			zap.String("uid", strconv.FormatUint(uint64(userID), 10)),
			zap.String("token", req.Token),
			zap.Error(err),
		)
		return &userpb.LogoutResponse{Message: "登出失败"}, err
	}

	logger.Log.Info("用户登出成功",
		zap.String("uid", strconv.FormatUint(uint64(userID), 10)),
		zap.String("token", req.Token),
	)

	return &userpb.LogoutResponse{
		Message: "登出成功",
	}, nil
}

// 注销用户
// TODO：1，删除用户，2，删除用户缓存token，3，异步删除短链 4，删除失败的操作 5，结果补偿
// README：注销用户
