package model

import (
	"time"

	"gorm.io/gorm"
)

// Order 订单模型
type Order struct {
	gorm.Model
	OrderID      string    `gorm:"uniqueIndex;size:32" json:"order_id"`  // 订单ID
	UserID       string    `gorm:"size:32;index" json:"user_id"`         // 用户ID
	Amount       int64     `json:"amount"`                               // 订单金额（分）
	Status       string    `gorm:"size:16;index" json:"status"`          // 订单状态：CREATED, PAID, CANCELLED, COMPLETED
	PaymentID    string    `gorm:"size:32;index" json:"payment_id"`      // 支付ID
	PaymentTime  time.Time `json:"payment_time"`                         // 支付时间
	Currency     string    `gorm:"size:8;default:'CNY'" json:"currency"` // 货币类型
	Description  string    `gorm:"size:256" json:"description"`          // 订单描述
	ClientIP     string    `gorm:"size:64" json:"client_ip"`             // 客户端IP
	DeliveryInfo string    `gorm:"type:text" json:"delivery_info"`       // 配送信息（JSON）
	Extra        string    `gorm:"type:text" json:"extra"`               // 额外信息（JSON）
}

// OrderItem 订单项模型
type OrderItem struct {
	gorm.Model
	OrderID     string `gorm:"size:32;index" json:"order_id"`   // 订单ID
	ProductID   string `gorm:"size:32;index" json:"product_id"` // 商品ID
	ProductName string `gorm:"size:128" json:"product_name"`    // 商品名称
	Quantity    int    `json:"quantity"`                        // 数量
	UnitPrice   int64  `json:"unit_price"`                      // 单价（分）
	TotalPrice  int64  `json:"total_price"`                     // 总价（分）
	Attributes  string `gorm:"type:text" json:"attributes"`     // 商品属性（JSON）
}
