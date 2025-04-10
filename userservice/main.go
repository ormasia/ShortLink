package main

import (
	"log"
	"net"
	"shortLink/userservice/cache"
	"shortLink/userservice/config"
	"shortLink/userservice/logger"
	"shortLink/userservice/model"
	"shortLink/userservice/mq"
	"shortLink/userservice/service"

	pb "shortLink/proto/userpb"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func main() {

	// 初始化配置
	if err := config.InitConfigFromNacos(); err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	// 初始化mq
	mq.InitKafka(config.GlobalConfig.Kafka.Brokers)

	// 初始化日志 TODO热更新logLevel
	logger.InitLogger(config.GlobalConfig.Kafka.Topic, zapcore.InfoLevel)

	// 初始化数据库
	model.InitDB(config.GlobalConfig.MySQL.GetDSN())

	// 初始化Redis
	cache.InitRedis(config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Password, config.GlobalConfig.Redis.Port, config.GlobalConfig.Redis.DB)

	// 创建gRPC服务器
	grpcServer := grpc.NewServer()

	// 注册用户服务
	pb.RegisterUserServiceServer(grpcServer, &service.UserService{})

	// 注册RBAC服务
	rbacService := service.NewRBACService(model.GetDB())
	pb.RegisterRBACServiceServer(grpcServer, rbacService)

	// 启动gRPC服务器
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Log.Error("failed to listen", zap.Error(err))
		return
	}

	logger.Log.Info("gRPC server started on :8081")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Log.Error("failed to serve", zap.Error(err))
	}
}
