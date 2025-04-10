/*
用户表,
*/
package model

import (
	"time"

	"gorm.io/gorm"
)

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
}

func (User) TableName() string {
	return "users"
}
