/*
用户表,
*/
package model

import (
	"time"

	"gorm.io/gorm"
)

type URLMapping struct {
	ShortURL    string `gorm:"primaryKey"`
	OriginalURL string `gorm:"not null"`
	UserID      string
	CreateTime  time.Time `gorm:"autoCreateTime"`
}

func (URLMapping) TableName() string {
	return "url_mapping" // 显式指定表名
}

type User struct {
	ID        uint           `gorm:"primaryKey"`
	Username  string         `gorm:"uniqueIndex;size:50;not null"`
	Password  string         `gorm:"size:255;not null"`
	Nickname  string         `gorm:"size:50"`
	Email     string         `gorm:"size:100"`
	Role      string         `gorm:"size:20;default:user"`
	Status    int            `gorm:"default:0"` // 0 = 正常，1 = 禁用
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Roles     []Role         `gorm:"many2many:user_roles;" json:"roles"`
}

func (User) TableName() string {
	return "users"
}

// Role 角色表
type Role struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:50;not null;unique" json:"name"`
	Description string         `gorm:"size:200" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Users       []User         `gorm:"many2many:user_roles;" json:"users"` //这个是 GORM 的 反向关系映射，和你在 User 中的 Roles []Role 是 对称存在 的。
}

// Permission 权限表
type Permission struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:50;not null;unique" json:"name"`
	Description string         `gorm:"size:200" json:"description"`
	Resource    string         `gorm:"size:50;not null" json:"resource"` // 资源，如 users, links
	Action      string         `gorm:"size:50;not null" json:"action"`   // 操作，如 create, read, update, delete
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Roles       []Role         `gorm:"many2many:role_permissions;" json:"roles"`
}

// UserRole 用户-角色关联表
type UserRole struct {
	UserID    uint      `gorm:"primarykey" json:"user_id"`
	RoleID    uint      `gorm:"primarykey" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// RolePermission 角色-权限关联表
type RolePermission struct {
	RoleID       uint      `gorm:"primarykey" json:"role_id"`
	PermissionID uint      `gorm:"primarykey" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}
