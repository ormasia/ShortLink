package repository

import (
	"shortLink/userservice/model"

	"gorm.io/gorm"
)

// UserRepository 定义用户数据访问接口
type UserRepository interface {
	FindByUsername(username string) (*model.User, error)
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
}

// GormUserRepository 实现基于Gorm的用户数据访问
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository 创建一个新的GormUserRepository实例
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// FindByUsername 根据用户名查找用户
func (r *GormUserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建新用户
func (r *GormUserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID 根据ID查找用户
func (r *GormUserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
