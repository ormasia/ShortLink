package service

import (
	"context"
	"fmt"
	"shortLink/payservice/logger"
	"shortLink/payservice/model"
	"shortLink/proto/orderpb"
	"shortLink/proto/paymentpb"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

// PaymentService 实现支付相关的gRPC服务
type PaymentService struct {
	paymentpb.UnimplementedPaymentServiceServer
	DB     *gorm.DB
	Config *PaymentServiceConfig
}

// PaymentServiceConfig 支付服务配置
type PaymentServiceConfig struct {
	OrderServiceAddr string // 订单服务地址
}

// NewPaymentService 创建支付服务实例
func NewPaymentService(db *gorm.DB, config *PaymentServiceConfig) *PaymentService {
	return &PaymentService{
		DB:     db,
		Config: config,
	}
}

// CreatePayment 创建支付订单
func (s *PaymentService) CreatePayment(ctx context.Context, req *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	// 开启事务
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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
	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("创建支付订单失败",
			"order_id", req.OrderId,
			"error", err)
		return nil, fmt.Errorf("创建支付订单失败: %v", err)
	}

	// 根据支付方式生成支付链接和二维码
	paymentURL := fmt.Sprintf("/pay/%s", paymentID)
	qrCode := fmt.Sprintf("/qr/%s", paymentID)

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Log.Error("提交事务失败",
			"payment_id", paymentID,
			"error", err)
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	logger.Log.Info("创建支付订单成功",
		"payment_id", paymentID,
		"order_id", req.OrderId)

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
			"payment_id", req.PaymentId,
			"error", err)
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
			"payment_id", req.PaymentId,
			"error", err)
		return nil, fmt.Errorf("记录支付回调失败: %v", err)
	}

	// 查询支付订单信息
	var payment model.Payment
	if err := tx.Where("payment_id = ?", req.PaymentId).First(&payment).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("查询支付订单失败",
			"payment_id", req.PaymentId,
			"error", err)
		return nil, fmt.Errorf("查询支付订单失败: %v", err)
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
			"payment_id", req.PaymentId,
			"error", err)
		return nil, fmt.Errorf("更新支付订单状态失败: %v", err)
	}

	// 如果支付成功，通知订单服务更新订单状态
	if req.Status == "SUCCESS" {
		// 连接订单服务
		conn, err := grpc.Dial(s.Config.OrderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			tx.Rollback()
			logger.Log.Error("连接订单服务失败",
				"order_service_addr", s.Config.OrderServiceAddr,
				"error", err)
			return nil, fmt.Errorf("连接订单服务失败: %v", err)
		}
		defer conn.Close()

		// 创建订单服务客户端
		orderClient := orderpb.NewOrderServiceClient(conn)

		// 调用订单服务更新订单支付状态
		_, err = orderClient.UpdateOrderPaymentStatus(ctx, &orderpb.UpdateOrderPaymentStatusRequest{
			OrderId:   payment.OrderID,
			PaymentId: req.PaymentId,
			Status:    "PAID",
		})

		if err != nil {
			tx.Rollback()
			logger.Log.Error("通知订单服务更新订单状态失败",
				"order_id", payment.OrderID,
				"error", err)
			return nil, fmt.Errorf("通知订单服务更新订单状态失败: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Log.Error("提交事务失败",
			"payment_id", req.PaymentId,
			"error", err)
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	logger.Log.Info("处理支付回调成功",
		"payment_id", req.PaymentId,
		"status", req.Status)

	return &paymentpb.PaymentCallbackResponse{
		Success: true,
		Message: "处理成功",
	}, nil
}

// CancelPayment 取消支付订单
func (s *PaymentService) CancelPayment(ctx context.Context, req *paymentpb.CancelPaymentRequest) (*paymentpb.CancelPaymentResponse, error) {
	// 开启事务
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新支付订单状态
	result := tx.Model(&model.Payment{}).Where("payment_id = ? AND status = ?", req.PaymentId, "PENDING").Update("status", "CANCELLED")
	if result.Error != nil {
		tx.Rollback()
		logger.Log.Error("取消支付订单失败",
			"payment_id", req.PaymentId,
			"error", result.Error)
		return nil, fmt.Errorf("取消支付订单失败: %v", result.Error)
	}

	// 检查是否有记录被更新
	if result.RowsAffected == 0 {
		tx.Rollback()
		logger.Log.Warn("取消支付订单失败：订单不存在或状态不是PENDING",
			"payment_id", req.PaymentId)
		return nil, fmt.Errorf("订单不存在或状态不是PENDING")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Log.Error("提交事务失败",
			"payment_id", req.PaymentId,
			"error", err)
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	logger.Log.Info("取消支付订单成功",
		"payment_id", req.PaymentId)

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
			"payment_id", req.PaymentId,
			"error", err)
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
