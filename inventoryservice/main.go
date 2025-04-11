package main

import (
	"fmt"
	"net"
	"shortLink/inventoryservice/config"
	"shortLink/inventoryservice/logger"
	"shortLink/inventoryservice/model"
	"shortLink/inventoryservice/service"
	"shortLink/proto/inventorypb"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal("加载配置失败", zap.Error(err))
	}

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("连接数据库失败", zap.Error(err))
	}

	// 自动迁移数据库表
	if errAutoMigrate := db.AutoMigrate(
		&model.Product{},
		&model.Inventory{},
		&model.InventoryLock{},
	); err != nil {
		logger.Log.Fatal("数据库迁移失败", zap.Error(errAutoMigrate))
	}

	// 创建gRPC服务器
	server := grpc.NewServer()

	// 注册库存服务
	inventoryService := service.NewInventoryService(db)
	inventorypb.RegisterInventoryServiceServer(server, inventoryService)

	// 启动gRPC服务器
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		logger.Log.Fatal("启动服务失败", zap.Error(err))
	}

	logger.Log.Info("库存服务启动成功", zap.Int("port", cfg.Server.Port))

	if err := server.Serve(lis); err != nil {
		logger.Log.Fatal("服务运行失败", zap.Error(err))
	}
}
