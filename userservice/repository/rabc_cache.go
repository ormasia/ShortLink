package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"shortLink/userservice/model"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	RolePermissionPrefix = "role_permissions:"
	CacheExpiration      = time.Hour * 24
)

type RBACCacheRepository interface {
	SetRolePermissions(roleID uint, permissions []model.Permission) error
	GetRolePermissions(roleID uint) ([]model.Permission, error)
	DeleteRolePermissions(roleID uint) error
}

type RBACRedisCache struct {
	rdb        *redis.Client
	ctx        context.Context
	expiration time.Duration
}

// NewRedisCache 创建一个新的RedisCache实例
func NewRBACRedisCache(rdb *redis.Client, expiration time.Duration) *RBACRedisCache {
	return &RBACRedisCache{
		rdb:        rdb,
		ctx:        context.Background(),
		expiration: expiration,
	}
}

func (r *RBACRedisCache) SetRolePermissions(roleID uint, permissions []model.Permission) error {
	if r.rdb == nil {
		return fmt.Errorf("Redis未初始化")
	}

	key := fmt.Sprintf("%s%d", RolePermissionPrefix, roleID)
	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("序列化权限数据失败: %v", err)
	}

	err = r.rdb.Set(context.Background(), key, string(data), CacheExpiration).Err()
	if err != nil {
		return fmt.Errorf("设置角色权限缓存失败: %v", err)
	}

	return nil
}

// GetRolePermissions 从缓存中获取角色的权限信息
func (r *RBACRedisCache) GetRolePermissions(roleID uint) ([]model.Permission, error) {
	if r.rdb == nil {
		return nil, fmt.Errorf("Redis未初始化")
	}

	key := fmt.Sprintf("%s%d", RolePermissionPrefix, roleID)
	data, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var permissions []model.Permission
	err = json.Unmarshal([]byte(data), &permissions)
	if err != nil {
		return nil, fmt.Errorf("反序列化权限数据失败: %v", err)
	}

	return permissions, nil
}

// DeleteRolePermissions 从缓存中删除角色的权限信息
func (r *RBACRedisCache) DeleteRolePermissions(roleID uint) error {
	if r.rdb == nil {
		return fmt.Errorf("Redis未初始化")
	}

	key := fmt.Sprintf("%s%d", RolePermissionPrefix, roleID)
	err := r.rdb.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("删除角色权限缓存失败: %v", err)
	}

	return nil
}
