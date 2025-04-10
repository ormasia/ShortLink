package model

import (
	"time"

	"gorm.io/gorm"
)

// Payment 支付订单模型
type Payment struct {
	gorm.Model
	PaymentID     string    `gorm:"uniqueIndex;size:32" json:"payment_id"` // 支付ID
	OrderID       string    `gorm:"size:32;index" json:"order_id"`         // 订单ID
	Amount        int64     `json:"amount"`                                // 支付金额（分）
	Currency      string    `gorm:"size:8" json:"currency"`                // 货币类型
	PaymentMethod string    `gorm:"size:16" json:"payment_method"`         // 支付方式
	Status        string    `gorm:"size:16;index" json:"status"`           // 支付状态
	Description   string    `gorm:"size:256" json:"description"`           // 订单描述
	NotifyURL     string    `gorm:"size:256" json:"notify_url"`            // 通知地址
	ReturnURL     string    `gorm:"size:256" json:"return_url"`            // 返回地址
	PaidAmount    int64     `json:"paid_amount"`                           // 实付金额
	PaidTime      time.Time `json:"paid_time"`                             // 支付时间
	TransactionID string    `gorm:"size:64" json:"transaction_id"`         // 支付平台交易号
	QRCode        string    `gorm:"size:256" json:"qr_code"`               // 二维码链接
	PaymentURL    string    `gorm:"size:256" json:"payment_url"`           // 支付链接
	Extra         string    `gorm:"type:text" json:"extra"`                // 额外信息（JSON）
}

// PaymentCallback 支付回调记录
type PaymentCallback struct {
	gorm.Model
	PaymentID     string    `gorm:"size:32;index" json:"payment_id"` // 支付ID
	TransactionID string    `gorm:"size:64" json:"transaction_id"`   // 支付平台交易号
	Status        string    `gorm:"size:16" json:"status"`           // 回调状态
	Amount        int64     `json:"amount"`                          // 支付金额
	PaidTime      time.Time `json:"paid_time"`                       // 支付时间
	RawData       string    `gorm:"type:text" json:"raw_data"`       // 原始回调数据
	Sign          string    `gorm:"size:256" json:"sign"`            // 签名
	Processed     bool      `json:"processed"`                       // 是否已处理
	ProcessTime   time.Time `json:"process_time"`                    // 处理时间
	Error         string    `gorm:"size:256" json:"error"`           // 处理错误信息
}

// PaymentRefund 退款记录
type PaymentRefund struct {
	gorm.Model
	RefundID      string    `gorm:"uniqueIndex;size:32" json:"refund_id"` // 退款ID
	PaymentID     string    `gorm:"size:32;index" json:"payment_id"`      // 支付ID
	Amount        int64     `json:"amount"`                               // 退款金额
	Reason        string    `gorm:"size:256" json:"reason"`               // 退款原因
	Status        string    `gorm:"size:16" json:"status"`                // 退款状态
	TransactionID string    `gorm:"size:64" json:"transaction_id"`        // 退款交易号
	RefundTime    time.Time `json:"refund_time"`                          // 退款时间
	Extra         string    `gorm:"type:text" json:"extra"`               // 额外信息
}
