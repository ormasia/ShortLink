package main

import (
	"log"
	"net"
	"shortLink/userservice/cache"
	"shortLink/userservice/config"
	"shortLink/userservice/logger"
	"shortLink/userservice/model"
	"shortLink/userservice/mq"
	"shortLink/userservice/repository"
	"shortLink/userservice/service/auth"
	"shortLink/userservice/service/rbac"
	"shortLink/userservice/service/user"
	"time"

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
	db := model.GetDB()

	// 初始化Redis
	cache.InitRedis(config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Password, config.GlobalConfig.Redis.Port, config.GlobalConfig.Redis.DB)
	rdb := cache.GetRedis()

	// 创建依赖组件
	// 1. 创建缓存适配器
	redisCache := repository.NewRedisCache(rdb, 24*time.Hour)
	rbacCache := repository.NewRBACRedisCache(rdb, 24*time.Hour)

	// 2. 创建存储库
	userRepository := repository.NewGormUserRepository(db)
	rbacRepository := repository.NewGormRBACRepository(db)

	// 3. 创建认证服务
	authService := auth.NewDefaultAuthService(redisCache)

	// 4. 创建用户服务
	userService := user.NewUserService(userRepository, authService)

	// 5. 创建RBAC服务
	rbacService := rbac.NewRBACService(rbacRepository, rbacCache)

	// 创建gRPC服务器
	grpcServer := grpc.NewServer()

	// 注册服务
	pb.RegisterUserServiceServer(grpcServer, userService)
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
