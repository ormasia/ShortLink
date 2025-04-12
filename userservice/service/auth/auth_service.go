package auth

import (
	"shortLink/userservice/pkg/jwt"
	"shortLink/userservice/repository"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthService 定义认证服务接口
type AuthService interface {
	VerifyPassword(hashedPassword, password string) error
	GenerateToken(userID uint, role string) (string, error)
	GetUserIDFromToken(token string) (uint, error)
	InvalidateToken(token string) error
}

// DefaultAuthService 实现默认的认证服务
type DefaultAuthService struct {
	tokenCacheRepository repository.TokenCache
}

// NewDefaultAuthService 创建一个新的DefaultAuthService实例
func NewDefaultAuthService(tokenCacheRepository repository.TokenCache) *DefaultAuthService {
	return &DefaultAuthService{
		tokenCacheRepository: tokenCacheRepository,
	}
}

// VerifyPassword 验证密码
func (s *DefaultAuthService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateToken 生成JWT令牌
func (s *DefaultAuthService) GenerateToken(userID uint, role string) (string, error) {
	// 生成令牌
	token, err := jwt.GenerateToken(userID, role, 24*time.Hour)
	if err != nil {
		return "", err
	}

	// 保存令牌到存储库
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	err = s.tokenCacheRepository.Set(token, userIDStr)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserIDFromToken 从令牌中获取用户ID
func (s *DefaultAuthService) GetUserIDFromToken(token string) (uint, error) {
	userIDStr, err := s.tokenCacheRepository.Get(token)
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(userID), nil
}

// InvalidateToken 使令牌失效
func (s *DefaultAuthService) InvalidateToken(token string) error {
	return s.tokenCacheRepository.Delete(token)
}
