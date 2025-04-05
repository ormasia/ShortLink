package main

import (
	"fmt"
	"os"

	"shortLink/cache"
	"shortLink/config"
	"shortLink/logger"
	"shortLink/model"
	"shortLink/mq"
	"shortLink/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	// // 初始化配置
	// if err := config.InitConfig("config/config.yaml"); err != nil {
	// 	fmt.Println("config init failed", err)
	// }

	// nacos初始化配置
	if err := config.InitConfigFromNacos(); err != nil {
		fmt.Println("config init failed", err)
	}
	// 初始化日志
	logger.InitLogger(config.GlobalConfig.Kafka.Topic, zapcore.InfoLevel)

	// 初始化 Kafka
	if err := mq.InitKafka([]string{"localhost:9092"}); err != nil {
		logger.Log.Error("❌ Kafka 初始化失败:", zap.Error(err))
		os.Exit(1)
	}
	// 添加日志消费者协程
	go mq.StartLogConsumer([]string{"localhost:9092"}, "log-consumer-group", []string{"shortlink-log"})

	logger.Log.Info("program start")

	// 初始化数据库
	if err := model.InitDB(config.GlobalConfig.MySQL.GetDSN()); err != nil {
		logger.Log.Error("database init failed", zap.Error(err))
	}

	// 初始化Redis
	// TODO:使用函数直接配置
	cache.InitRedis(config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Password, config.GlobalConfig.Redis.Port, config.GlobalConfig.Redis.DB)
	//初始化布隆过滤器
	cache.InitBloom(100000, 0.01)
	// 预热布隆过滤器
	cache.WarmUpBloomFromDB()

	// 设置Gin模式
	gin.SetMode(config.GlobalConfig.App.Mode)
	// 初始化路由
	r := gin.Default()
	router.InitRoutes(r)

	// 获取配置文件中的主机和端口
	host, port := config.GlobalConfig.App.GetHost()
	r.Run(fmt.Sprintf("%s:%d", host, port))
}
