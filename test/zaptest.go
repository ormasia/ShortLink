package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"shortLink/logger"
	"shortLink/mq"
)

func main() {
	// 初始化 Kafka
	if err := mq.InitKafka([]string{"localhost:9092"}); err != nil {
		fmt.Println("❌ Kafka 初始化失败:", err)
		os.Exit(1)
	}

	// 初始化 zap（控制台 + Kafka）
	logger.InitLogger("shortlink-log", zapcore.InfoLevel)

	// 异步启动日志消费者
	go func() {
		err := mq.StartLogConsumer([]string{"localhost:9092"}, "log-consumer-group", []string{"shortlink-log"})
		if err != nil {
			fmt.Println("❌ Kafka 消费者启动失败:", err)
			os.Exit(1)
		}
	}()

	// 模拟结构化日志写入
	for i := range 3 {
		go func() {
			logger.Log.Info("📦 日志测试", zap.Int("index", i), zap.String("type", "info"))
			logger.Log.Warn("⚠️ 警告日志", zap.Int("index", i))
			logger.Log.Error("🔥 错误日志", zap.Int("index", i), zap.String("reason", "测试错误"))
			time.Sleep(1 * time.Second)
		}()
	}

	fmt.Println("✅ 测试日志已发送，Kafka 消费者已在后台运行...")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println("👋 程序已手动终止")
}
