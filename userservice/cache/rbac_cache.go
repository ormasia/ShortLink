package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"shortLink/userservice/model"
	"time"
)

const (
	RolePermissionPrefix = "role_permissions:"
	CacheExpiration      = time.Hour * 24
)

// SetRolePermissions 将角色的权限信息存入缓存
func SetRolePermissions(roleID uint, permissions []model.Permission) error {
	if rdb == nil {
		return fmt.Errorf("Redis未初始化")
	}

	key := fmt.Sprintf("%s%d", RolePermissionPrefix, roleID)
	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("序列化权限数据失败: %v", err)
	}

	err = rdb.Set(context.Background(), key, string(data), CacheExpiration).Err()
	if err != nil {
		return fmt.Errorf("设置角色权限缓存失败: %v", err)
	}

	return nil
}

// GetRolePermissions 从缓存中获取角色的权限信息
func GetRolePermissions(roleID uint) ([]model.Permission, error) {
	if rdb == nil {
		return nil, fmt.Errorf("Redis未初始化")
	}

	key := fmt.Sprintf("%s%d", RolePermissionPrefix, roleID)
	data, err := rdb.Get(context.Background(), key).Result()
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
func DeleteRolePermissions(roleID uint) error {
	if rdb == nil {
		return fmt.Errorf("Redis未初始化")
	}

	key := fmt.Sprintf("%s%d", RolePermissionPrefix, roleID)
	err := rdb.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("删除角色权限缓存失败: %v", err)
	}

	return nil
}
