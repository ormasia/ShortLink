package main

import (
	"log"
	"shortLink/admin/router"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接userservice的gRPC服务
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接到userservice: %v", err)
	}
	defer conn.Close()

	// 创建Gin引擎
	r := gin.Default()

	// 设置后台管理路由
	router.SetupAdminRouter(r, conn)

	// 启动服务器
	if err := r.Run(":8083"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
