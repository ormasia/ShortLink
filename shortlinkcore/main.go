package main

import (
	"log"
	"net"
	"shortLink/proto/shortlinkpb"
	"shortLink/shortlinkcore/cache"
	"shortLink/shortlinkcore/config"
	"shortLink/shortlinkcore/logger"
	"shortLink/shortlinkcore/model"
	"shortLink/shortlinkcore/mq"
	"shortLink/shortlinkcore/service"

	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func main() {
	// 初始化配置和数据库
	err := config.InitConfigFromNacos()
	if err != nil {
		log.Fatalf("❌ 初始化配置失败: %v", err)
	}

	// 初始化mq
	mq.InitKafka(config.GlobalConfig.Kafka.Brokers)

	// 初始化日志
	logger.InitLogger(config.GlobalConfig.Kafka.Topic, zapcore.InfoLevel)

	// 初始化数据库
	model.InitDB(config.GlobalConfig.MySQL.GetDSN())

	// 初始化Redis
	// TODO:使用函数直接配置
	cache.InitRedis(config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Password, config.GlobalConfig.Redis.Port, config.GlobalConfig.Redis.DB)
	//初始化布隆过滤器
	cache.InitBloom(100000, 0.01)
	// 预热布隆过滤器
	cache.WarmUpBloomFromDB()

	// 创建 gRPC 服务器并注册服务
	grpcServer := grpc.NewServer()
	shortlinkpb.RegisterShortlinkServiceServer(grpcServer, &service.ShortlinkService{}) //DB: model.GetDB(), Cache: cache.GetRedis()

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("❌ 监听端口失败: %v", err)
	}
	log.Println("✅ shortlink-core gRPC 服务已启动，监听端口 :8082")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ 启动 gRPC 服务失败: %v", err)
	}
}
