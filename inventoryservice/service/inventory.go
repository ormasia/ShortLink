package service

import (
	"context"
	"fmt"
	"shortLink/inventoryservice/logger"
	"shortLink/inventoryservice/model"
	"shortLink/proto/inventorypb"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InventoryService 实现库存相关的gRPC服务
type InventoryService struct {
	inventorypb.UnimplementedInventoryServiceServer
	DB *gorm.DB
}

// NewInventoryService 创建库存服务实例
func NewInventoryService(db *gorm.DB) *InventoryService {
	return &InventoryService{DB: db}
}

// GetProductInventory 获取产品库存信息
func (s *InventoryService) GetProductInventory(ctx context.Context, req *inventorypb.GetProductInventoryRequest) (*inventorypb.GetProductInventoryResponse, error) {
	var inventory model.Inventory
	if err := s.DB.Where("product_id = ?", req.ProductId).First(&inventory).Error; err != nil {
		logger.Log.Error("查询产品库存失败",
			zap.String("product_id", req.ProductId),
			zap.Error(err))
		return nil, fmt.Errorf("查询产品库存失败: %v", err)
	}

	return &inventorypb.GetProductInventoryResponse{
		ProductId:  inventory.ProductID,
		TotalStock: int32(inventory.TotalStock),
		Available:  int32(inventory.Available),
		Reserved:   int32(inventory.Reserved),
	}, nil
}

// LockInventory 锁定库存
func (s *InventoryService) LockInventory(ctx context.Context, req *inventorypb.LockInventoryRequest) (*inventorypb.LockInventoryResponse, error) {
	// 开启事务
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查库存是否足够
	var inventory model.Inventory
	if err := tx.Where("product_id = ?", req.ProductId).First(&inventory).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("查询产品库存失败",
			zap.String("product_id", req.ProductId),
			zap.Error(err))
		return nil, fmt.Errorf("查询产品库存失败: %v", err)
	}

	// 检查可用库存是否足够
	if inventory.Available < int(req.Quantity) {
		tx.Rollback()
		logger.Log.Warn("库存不足",
			zap.String("product_id", req.ProductId),
			zap.Int("available", inventory.Available),
			zap.Int32("requested", req.Quantity))
		return nil, fmt.Errorf("库存不足")
	}

	// 更新库存
	if err := tx.Model(&model.Inventory{}).Where("product_id = ?", req.ProductId).Updates(map[string]interface{}{
		"available":    inventory.Available - int(req.Quantity),
		"reserved":     inventory.Reserved + int(req.Quantity),
		"last_updated": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("更新库存失败",
			zap.String("product_id", req.ProductId),
			zap.Error(err))
		return nil, fmt.Errorf("更新库存失败: %v", err)
	}

	// 创建库存锁定记录
	lockID := uuid.New().String()
	expireTime := time.Now().Add(time.Duration(req.ExpireSeconds) * time.Second)

	inventoryLock := &model.InventoryLock{
		LockID:     lockID,
		OrderID:    req.OrderId,
		ProductID:  req.ProductId,
		Quantity:   int(req.Quantity),
		Status:     "LOCKED",
		LockTime:   time.Now(),
		ExpireTime: expireTime,
	}

	if err := tx.Create(inventoryLock).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("创建库存锁定记录失败",
			zap.String("product_id", req.ProductId),
			zap.String("order_id", req.OrderId),
			zap.Error(err))
		return nil, fmt.Errorf("创建库存锁定记录失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Log.Error("提交事务失败",
			zap.String("product_id", req.ProductId),
			zap.Error(err))
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	logger.Log.Info("锁定库存成功",
		zap.String("lock_id", lockID),
		zap.String("product_id", req.ProductId),
		zap.String("order_id", req.OrderId),
		zap.Int32("quantity", req.Quantity))

	return &inventorypb.LockInventoryResponse{
		LockId:     lockID,
		Success:    true,
		ExpireTime: expireTime.Format(time.RFC3339),
	}, nil
}

// ReleaseInventory 释放库存
func (s *InventoryService) ReleaseInventory(ctx context.Context, req *inventorypb.ReleaseInventoryRequest) (*inventorypb.ReleaseInventoryResponse, error) {
	// 开启事务
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询锁定记录
	var lock model.InventoryLock
	if err := tx.Where("lock_id = ? AND status = ?", req.LockId, "LOCKED").First(&lock).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("查询库存锁定记录失败",
			zap.String("lock_id", req.LockId),
			zap.Error(err))
		return nil, fmt.Errorf("查询库存锁定记录失败: %v", err)
	}

	// 更新锁定记录状态
	if err := tx.Model(&model.InventoryLock{}).Where("lock_id = ?", req.LockId).Updates(map[string]interface{}{
		"status":       "RELEASED",
		"release_time": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("更新锁定记录状态失败",
			zap.String("lock_id", req.LockId),
			zap.Error(err))
		return nil, fmt.Errorf("更新锁定记录状态失败: %v", err)
	}

	// 恢复库存
	if err := tx.Model(&model.Inventory{}).Where("product_id = ?", lock.ProductID).Updates(map[string]interface{}{
		"available":    gorm.Expr("available + ?", lock.Quantity),
		"reserved":     gorm.Expr("reserved - ?", lock.Quantity),
		"last_updated": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("恢复库存失败",
			zap.String("product_id", lock.ProductID),
			zap.Error(err))
		return nil, fmt.Errorf("恢复库存失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Log.Error("提交事务失败",
			zap.String("lock_id", req.LockId),
			zap.Error(err))
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	logger.Log.Info("释放库存成功",
		zap.String("lock_id", req.LockId),
		zap.String("product_id", lock.ProductID),
		zap.Int("quantity", lock.Quantity))

	return &inventorypb.ReleaseInventoryResponse{
		Success: true,
		Message: "释放成功",
	}, nil
}

// ConfirmInventory 确认库存扣减
func (s *InventoryService) ConfirmInventory(ctx context.Context, req *inventorypb.ConfirmInventoryRequest) (*inventorypb.ConfirmInventoryResponse, error) {
	// 开启事务
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询锁定记录
	var lock model.InventoryLock
	if err := tx.Where("lock_id = ? AND status = ?", req.LockId, "LOCKED").First(&lock).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("查询库存锁定记录失败",
			zap.String("lock_id", req.LockId),
			zap.Error(err))
		return nil, fmt.Errorf("查询库存锁定记录失败: %v", err)
	}

	// 更新锁定记录状态
	if err := tx.Model(&model.InventoryLock{}).Where("lock_id = ?", req.LockId).Update("status", "CONFIRMED").Error; err != nil {
		tx.Rollback()
		logger.Log.Error("更新锁定记录状态失败",
			zap.String("lock_id", req.LockId),
			zap.Error(err))
		return nil, fmt.Errorf("更新锁定记录状态失败: %v", err)
	}

	// 更新库存（从预留库存中扣除，总库存减少）
	if err := tx.Model(&model.Inventory{}).Where("product_id = ?", lock.ProductID).Updates(map[string]interface{}{
		"total_stock":  gorm.Expr("total_stock - ?", lock.Quantity),
		"reserved":     gorm.Expr("reserved - ?", lock.Quantity),
		"last_updated": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("更新库存失败",
			zap.String("product_id", lock.ProductID),
			zap.Error(err))
		return nil, fmt.Errorf("更新库存失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Log.Error("提交事务失败",
			zap.String("lock_id", req.LockId),
			zap.Error(err))
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	logger.Log.Info("确认库存扣减成功",
		zap.String("lock_id", req.LockId),
		zap.String("product_id", lock.ProductID),
		zap.Int("quantity", lock.Quantity))

	return &inventorypb.ConfirmInventoryResponse{
		Success: true,
		Message: "确认成功",
	}, nil
}
