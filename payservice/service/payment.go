package service

import (
	"context"
	"fmt"
	"shortLink/payservice/logger"
	"shortLink/payservice/model"
	"shortLink/proto/paymentpb"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PaymentService 实现支付相关的gRPC服务
type PaymentService struct {
	paymentpb.UnimplementedPaymentServiceServer
	DB *gorm.DB
}

// NewPaymentService 创建支付服务实例
func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{DB: db}
}

// CreatePayment 创建支付订单
func (s *PaymentService) CreatePayment(ctx context.Context, req *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	// 生成支付ID
	paymentID := uuid.New().String()

	// 创建支付订单记录
	payment := &model.Payment{
		PaymentID:     paymentID,
		OrderID:       req.OrderId,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		Status:        "PENDING",
		Description:   req.Description,
		NotifyURL:     req.NotifyUrl,
		ReturnURL:     req.ReturnUrl,
	}

	// 保存到数据库
	if err := s.DB.Create(payment).Error; err != nil {
		logger.Log.Error("创建支付订单失败",
			zap.String("order_id", req.OrderId),
			zap.Error(err))
		return nil, fmt.Errorf("创建支付订单失败: %v", err)
	}

	// TODO: 根据支付方式生成支付链接和二维码
	paymentURL := fmt.Sprintf("/pay/%s", paymentID)
	qrCode := fmt.Sprintf("/qr/%s", paymentID)

	logger.Log.Info("创建支付订单成功",
		zap.String("payment_id", paymentID),
		zap.String("order_id", req.OrderId))

	return &paymentpb.CreatePaymentResponse{
		PaymentId:  paymentID,
		PaymentUrl: paymentURL,
		QrCode:     qrCode,
		Status:     "PENDING",
	}, nil
}

// QueryPaymentStatus 查询支付状态
func (s *PaymentService) QueryPaymentStatus(ctx context.Context, req *paymentpb.QueryPaymentStatusRequest) (*paymentpb.QueryPaymentStatusResponse, error) {
	var payment model.Payment
	if err := s.DB.Where("payment_id = ?", req.PaymentId).First(&payment).Error; err != nil {
		logger.Log.Error("查询支付订单失败",
			zap.String("payment_id", req.PaymentId),
			zap.Error(err))
		return nil, fmt.Errorf("查询支付订单失败: %v", err)
	}

	return &paymentpb.QueryPaymentStatusResponse{
		PaymentId:     payment.PaymentID,
		Status:        payment.Status,
		PaidAmount:    payment.PaidAmount,
		PaidTime:      payment.PaidTime.Format(time.RFC3339),
		TransactionId: payment.TransactionID,
	}, nil
}

// HandlePaymentCallback 处理支付回调
func (s *PaymentService) HandlePaymentCallback(ctx context.Context, req *paymentpb.PaymentCallbackRequest) (*paymentpb.PaymentCallbackResponse, error) {
	// 开启事务
	tx := s.DB.Begin()

	// 记录回调信息
	callback := &model.PaymentCallback{
		PaymentID:     req.PaymentId,
		TransactionID: req.TransactionId,
		Status:        req.Status,
		Amount:        req.PaidAmount,
		PaidTime:      time.Now(),
		Sign:          req.Sign,
		Processed:     false,
	}

	if err := tx.Create(callback).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("记录支付回调失败",
			zap.String("payment_id", req.PaymentId),
			zap.Error(err))
		return nil, fmt.Errorf("记录支付回调失败: %v", err)
	}

	// 更新支付订单状态
	if err := tx.Model(&model.Payment{}).Where("payment_id = ?", req.PaymentId).Updates(map[string]interface{}{
		"status":         req.Status,
		"paid_amount":    req.PaidAmount,
		"paid_time":      time.Now(),
		"transaction_id": req.TransactionId,
	}).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("更新支付订单状态失败",
			zap.String("payment_id", req.PaymentId),
			zap.Error(err))
		return nil, fmt.Errorf("更新支付订单状态失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Log.Error("提交事务失败",
			zap.String("payment_id", req.PaymentId),
			zap.Error(err))
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	logger.Log.Info("处理支付回调成功",
		zap.String("payment_id", req.PaymentId),
		zap.String("status", req.Status))

	return &paymentpb.PaymentCallbackResponse{
		Success: true,
		Message: "处理成功",
	}, nil
}

// CancelPayment 取消支付订单
func (s *PaymentService) CancelPayment(ctx context.Context, req *paymentpb.CancelPaymentRequest) (*paymentpb.CancelPaymentResponse, error) {
	if err := s.DB.Model(&model.Payment{}).Where("payment_id = ? AND status = ?", req.PaymentId, "PENDING").Update("status", "CANCELLED").Error; err != nil {
		logger.Log.Error("取消支付订单失败",
			zap.String("payment_id", req.PaymentId),
			zap.Error(err))
		return nil, fmt.Errorf("取消支付订单失败: %v", err)
	}

	logger.Log.Info("取消支付订单成功",
		zap.String("payment_id", req.PaymentId))

	return &paymentpb.CancelPaymentResponse{
		Success: true,
		Message: "取消成功",
	}, nil
}

// GetPaymentDetail 获取支付订单详情
func (s *PaymentService) GetPaymentDetail(ctx context.Context, req *paymentpb.GetPaymentDetailRequest) (*paymentpb.GetPaymentDetailResponse, error) {
	var payment model.Payment
	if err := s.DB.Where("payment_id = ?", req.PaymentId).First(&payment).Error; err != nil {
		logger.Log.Error("获取支付订单详情失败",
			zap.String("payment_id", req.PaymentId),
			zap.Error(err))
		return nil, fmt.Errorf("获取支付订单详情失败: %v", err)
	}

	return &paymentpb.GetPaymentDetailResponse{
		Payment: &paymentpb.PaymentDetail{
			PaymentId:     payment.PaymentID,
			OrderId:       payment.OrderID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			PaymentMethod: payment.PaymentMethod,
			Status:        payment.Status,
			CreatedTime:   payment.CreatedAt.Format(time.RFC3339),
			PaidTime:      payment.PaidTime.Format(time.RFC3339),
			TransactionId: payment.TransactionID,
			Description:   payment.Description,
		},
	}, nil
}
