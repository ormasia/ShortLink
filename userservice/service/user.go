package service

import (
	"context"
	"shortLink/proto/userpb"
	"shortLink/userservice/cache"
	"shortLink/userservice/logger"
	"shortLink/userservice/model"
	"shortLink/userservice/pkg/jwt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
}

func (s *UserService) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {

	// 校验重复
	var user model.User
	if err := model.GetDB().Where("username = ?", req.Username).First(&user).Error; err != gorm.ErrRecordNotFound {
		logger.Log.Warn("用户注册失败：用户名已存在",
			zap.String("username", req.Username),
			zap.String("email", req.Email),
		)
		return &userpb.RegisterResponse{Message: "用户名已存在"}, err
	}

	// 加密密码
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// 写入数据库
	user = model.User{
		Username: req.Username,
		Password: string(hash),
		Nickname: req.Nickname,
		Email:    req.Email,
	}
	if err := model.GetDB().Create(&user).Error; err != nil {
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

func (s *UserService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {

	var user model.User
	// 查询用户
	if err := model.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
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
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
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

	// 生成 JWT
	token, _ := jwt.GenerateToken(user.ID, user.Role, 24*time.Hour)

	// 写入redis
	cache.Set(token, strconv.FormatUint(uint64(user.ID), 10))

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
