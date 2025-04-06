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

	// 初始化日志
	logger.InitLogger(config.GlobalConfig.Kafka.Topic, zapcore.InfoLevel)

	// 初始化数据库
	model.InitDB(config.GlobalConfig.MySQL.GetDSN())

	// 初始化Redis
	cache.InitRedis(config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Password, config.GlobalConfig.Redis.Port, config.GlobalConfig.Redis.DB)

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &service.UserService{})

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Log.Error("failed to listen", zap.Error(err))
	}
	logger.Log.Info("✅ user-service 启动成功，监听端口 :8081")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Log.Error("failed to serve", zap.Error(err))
	}
}
