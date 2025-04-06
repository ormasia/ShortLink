package service

import (
	"context"
	"log"
	"shortLink/pkg/jwt"
	"shortLink/proto/userpb"
	"shortLink/userservice/cache"
	"shortLink/userservice/model"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
}

func (s *UserService) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	log.Printf("📨 注册请求: %v", req.Username)
	// "TODO": 写入数据库、校验重复
	var user model.User
	if err := model.GetDB().Where("username = ?", req.Username).First(&user).Error; err != gorm.ErrRecordNotFound {
		return &userpb.RegisterResponse{Message: "用户名已存在"}, err
	}

	// 加密密码
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user = model.User{
		Username: req.Username,
		Password: string(hash),
		Nickname: req.Nickname,
		Email:    req.Email,
	}
	// 写入数据库
	if err := model.GetDB().Create(&user).Error; err != nil {
		return &userpb.RegisterResponse{Message: "注册失败"}, err
	}

	return &userpb.RegisterResponse{Message: "注册成功"}, nil
}

func (s *UserService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	log.Printf("🔐 登录请求: %v", req.Username)
	// TODO: 验证密码、生成 JWT
	var user model.User
	// 查询用户
	if err := model.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		return &userpb.LoginResponse{
			Token: "",
			User:  nil,
		}, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &userpb.LoginResponse{
			Token: "",
			User:  nil,
		}, err
	}

	// 生成 JWT
	token, _ := jwt.GenerateToken(user.ID, user.Role, 24*time.Hour)
	// 写入redis
	cache.Set(token, strconv.FormatUint(uint64(user.ID), 10))
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
