package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"shortLink/userservice/cache"
	"shortLink/userservice/config"
	"shortLink/userservice/logger"
	"shortLink/userservice/model"
	"shortLink/userservice/mq"
	"shortLink/userservice/pkg/discovery"
	"shortLink/userservice/repository"
	"shortLink/userservice/scripts"
	"shortLink/userservice/service/auth"
	"shortLink/userservice/service/rbac"
	"shortLink/userservice/service/user"
	"syscall"
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

	// 初始化角色和权限
	if err := scripts.InitRolesAndPermissions(db); err != nil {
		logger.Log.Error("初始化角色和权限失败", zap.Error(err))
	}

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

	// 初始化服务发现客户端
	if err := discovery.InitNamingClient(); err != nil {
		logger.Log.Error("初始化服务发现客户端失败", zap.Error(err))
		return
	}

	// 注册服务
	if err := discovery.RegisterService(); err != nil {
		logger.Log.Error("注册服务失败", zap.Error(err))
		return
	}

	// 优雅退出
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM) //注册信号处理器的函数，信号进入sigCh，原本阻塞的channel就恢复了
		<-sigCh

		logger.Log.Info("正在关闭服务...")

		// 注销服务
		if err := discovery.DeregisterService(); err != nil {
			logger.Log.Error("注销服务失败", zap.Error(err))
		}

		// 停止gRPC服务器
		grpcServer.GracefulStop()
	}()

	logger.Log.Info("gRPC server started on :8081")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Log.Error("failed to serve", zap.Error(err))
	}
}
