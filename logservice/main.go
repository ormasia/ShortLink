package main

import (
	"log"
	"os"
	"os/signal"
	"shortLink/logservice/consumer"
	"syscall"
)

func main() {
	brokers := []string{"localhost:9092"} // Kafka 地址，可通过配置管理引入
	topics := []string{"shortlink-log"}
	groupID := "log-service-group"

	// 创建日志目录
	_ = os.MkdirAll("./logs", os.ModePerm)
	logFile, err := os.OpenFile("./logs/aggregated.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("❌ 无法打开日志文件: %v", err)
	}
	defer logFile.Close()

	// 设置日志输出路径
	log.SetOutput(logFile)
	log.Println("✅ log-service 启动，开始消费日志...")

	// 启动 Kafka 消费者，日志落地
	go func() {
		err = consumer.StartLogConsumer(brokers, groupID, topics)
		if err != nil {
			log.Fatalf("❌ Kafka 日志消费者启动失败: %v", err)
		}
	}()
	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 log-service 退出")
}
