package model

import (
	"time"

	"gorm.io/gorm"
)

// Product 产品模型
type Product struct {
	gorm.Model
	ProductID   string `gorm:"uniqueIndex;size:32" json:"product_id"` // 产品ID
	Name        string `gorm:"size:128" json:"name"`                  // 产品名称
	Description string `gorm:"size:512" json:"description"`           // 产品描述
	Price       int64  `json:"price"`                                 // 产品价格（分）
	Status      string `gorm:"size:16;index" json:"status"`           // 产品状态：ACTIVE, INACTIVE
	Category    string `gorm:"size:64;index" json:"category"`         // 产品类别
	Attributes  string `gorm:"type:text" json:"attributes"`           // 产品属性（JSON）
}

// Inventory 库存模型
type Inventory struct {
	gorm.Model
	ProductID   string    `gorm:"uniqueIndex;size:32" json:"product_id"` // 产品ID
	TotalStock  int       `json:"total_stock"`                           // 总库存
	Available   int       `json:"available"`                             // 可用库存
	Reserved    int       `json:"reserved"`                              // 已预留库存
	LastUpdated time.Time `json:"last_updated"`                          // 最后更新时间
}

// InventoryLock 库存锁定记录
type InventoryLock struct {
	gorm.Model
	LockID      string    `gorm:"uniqueIndex;size:32" json:"lock_id"` // 锁定ID
	OrderID     string    `gorm:"size:32;index" json:"order_id"`      // 订单ID
	ProductID   string    `gorm:"size:32;index" json:"product_id"`    // 产品ID
	Quantity    int       `json:"quantity"`                           // 锁定数量
	Status      string    `gorm:"size:16" json:"status"`              // 锁定状态：LOCKED, RELEASED, CONFIRMED
	LockTime    time.Time `json:"lock_time"`                          // 锁定时间
	ReleaseTime time.Time `json:"release_time"`                       // 释放时间
	ExpireTime  time.Time `json:"expire_time"`                        // 过期时间
}
