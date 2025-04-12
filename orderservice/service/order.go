package service

import (
	"context"
	"fmt"
	"shortLink/orderservice/logger"
	"shortLink/orderservice/model"
	"shortLink/proto/orderpb"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OrderService 实现订单相关的gRPC服务
type OrderService struct {
	orderpb.UnimplementedOrderServiceServer
	DB *gorm.DB
}

// NewOrderService 创建订单服务实例
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{DB: db}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	// 开启事务
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 生成订单ID
	orderID := uuid.New().String()

	// 创建订单记录
	order := &model.Order{
		OrderID:     orderID,
		UserID:      req.UserId,
		Amount:      req.TotalAmount,
		Status:      "CREATED",
		Currency:    req.Currency,
		Description: req.Description,
		ClientIP:    req.ClientIp,
	}

	// 保存订单到数据库
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("创建订单失败",
			zap.String("user_id", req.UserId),
			zap.Error(err))
		return nil, fmt.Errorf("创建订单失败: %v", err)
	}

	// 创建订单项
	for _, item := range req.Items {
		orderItem := &model.OrderItem{
			OrderID:     orderID,
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
			Attributes:  item.Attributes,
		}

		if err := tx.Create(orderItem).Error; err != nil {
			tx.Rollback()
			logger.Log.Error("创建订单项失败",
				zap.String("order_id", orderID),
				zap.String("product_id", item.ProductId),
				zap.Error(err))
			return nil, fmt.Errorf("创建订单项失败: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Log.Error("提交事务失败",
			zap.String("order_id", orderID),
			zap.Error(err))
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	logger.Log.Info("创建订单成功",
		zap.String("order_id", orderID),
		zap.String("user_id", req.UserId))

	return &orderpb.CreateOrderResponse{
		OrderId: orderID,
		Status:  "CREATED",
	}, nil
}

// GetOrderDetail 获取订单详情
func (s *OrderService) GetOrderDetail(ctx context.Context, req *orderpb.GetOrderDetailRequest) (*orderpb.GetOrderDetailResponse, error) {
	// 查询订单信息
	var order model.Order
	if err := s.DB.Where("order_id = ?", req.OrderId).First(&order).Error; err != nil {
		logger.Log.Error("查询订单失败",
			zap.String("order_id", req.OrderId),
			zap.Error(err))
		return nil, fmt.Errorf("查询订单失败: %v", err)
	}

	// 查询订单项
	var orderItems []model.OrderItem
	if err := s.DB.Where("order_id = ?", req.OrderId).Find(&orderItems).Error; err != nil {
		logger.Log.Error("查询订单项失败",
			zap.String("order_id", req.OrderId),
			zap.Error(err))
		return nil, fmt.Errorf("查询订单项失败: %v", err)
	}

	// 构建响应
	response := &orderpb.GetOrderDetailResponse{
		Order: &orderpb.OrderDetail{
			OrderId:      order.OrderID,
			UserId:       order.UserID,
			TotalAmount:  order.Amount,
			Status:       order.Status,
			PaymentId:    order.PaymentID,
			CreatedTime:  order.CreatedAt.Format(time.RFC3339),
			PaymentTime:  order.PaymentTime.Format(time.RFC3339),
			Currency:     order.Currency,
			Description:  order.Description,
			ClientIp:     order.ClientIP,
			DeliveryInfo: order.DeliveryInfo,
		},
		Items: make([]*orderpb.OrderItemDetail, 0, len(orderItems)),
	}

	// 添加订单项
	for _, item := range orderItems {
		response.Items = append(response.Items, &orderpb.OrderItemDetail{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    int32(item.Quantity),
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
			Attributes:  item.Attributes,
		})
	}

	return response, nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	// 开启事务
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新订单状态
	result := tx.Model(&model.Order{}).Where("order_id = ? AND status = ?", req.OrderId, "CREATED").Update("status", "CANCELLED")
	if result.Error != nil {
		tx.Rollback()
		logger.Log.Error("取消订单失败",
			zap.String("order_id", req.OrderId),
			zap.Error(result.Error))
		return nil, fmt.Errorf("取消订单失败: %v", result.Error)
	}

	// 检查是否有记录被更新
	if result.RowsAffected == 0 {
		tx.Rollback()
		logger.Log.Warn("取消订单失败：订单不存在或状态不是CREATED",
			zap.String("order_id", req.OrderId))
		return nil, fmt.Errorf("订单不存在或状态不是CREATED")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Log.Error("提交事务失败",
			zap.String("order_id", req.OrderId),
			zap.Error(err))
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	logger.Log.Info("取消订单成功",
		zap.String("order_id", req.OrderId))

	return &orderpb.CancelOrderResponse{
		Success: true,
		Message: "取消成功",
	}, nil
}

// UpdateOrderPaymentStatus 更新订单支付状态
func (s *OrderService) UpdateOrderPaymentStatus(ctx context.Context, req *orderpb.UpdateOrderPaymentStatusRequest) (*orderpb.UpdateOrderPaymentStatusResponse, error) {
	// 开启事务
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新订单状态
	updates := map[string]interface{}{
		"status":       req.Status,
		"payment_id":   req.PaymentId,
		"payment_time": time.Now(),
	}

	result := tx.Model(&model.Order{}).Where("order_id = ?", req.OrderId).Updates(updates)
	if result.Error != nil {
		tx.Rollback()
		logger.Log.Error("更新订单支付状态失败",
			zap.String("order_id", req.OrderId),
			zap.Error(result.Error))
		return nil, fmt.Errorf("更新订单支付状态失败: %v", result.Error)
	}

	// 检查是否有记录被更新
	if result.RowsAffected == 0 {
		tx.Rollback()
		logger.Log.Warn("更新订单支付状态失败：订单不存在",
			zap.String("order_id", req.OrderId))
		return nil, fmt.Errorf("订单不存在")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Log.Error("提交事务失败",
			zap.String("order_id", req.OrderId),
			zap.Error(err))
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	logger.Log.Info("更新订单支付状态成功",
		zap.String("order_id", req.OrderId),
		zap.String("status", req.Status))

	return &orderpb.UpdateOrderPaymentStatusResponse{
		Success: true,
		Message: "更新成功",
	}, nil
}
