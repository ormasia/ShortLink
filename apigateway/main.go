package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"shortLink/apigateway/cache"
	"shortLink/apigateway/config"
	"shortLink/apigateway/middleware"
	"shortLink/apigateway/model"
	pbShortlink "shortLink/proto/shortlinkpb"
	pb "shortLink/proto/userpb"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 处理请求和响应，user服务只要做对应的user操作，其他服务由网关处理
func main() {

	// 初始化配置和数据库
	err := config.InitConfigFromNacos()
	if err != nil {
		log.Fatalf("❌ 初始化配置失败: %v", err)
	}
	// 初始化数据库
	model.InitDB(config.GlobalConfig.MySQL.GetDSN())

	// 初始化Redis
	// TODO:使用函数直接配置
	cache.InitRedis(config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Password, config.GlobalConfig.Redis.Port, config.GlobalConfig.Redis.DB)

	r := gin.Default()
	// 启用跨域支持（允许前端访问）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接 user-service 失败: %v", err)
	}
	client := pb.NewUserServiceClient(conn)

	r.POST("/api/user/register", func(c *gin.Context) {
		var req pb.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		// 调用用户服务注册，已经检验过参数，所以这里有可能的错误是用户名已存在，数据库错误
		res, err := client.Register(ctx, &req)
		if err != nil {
			// TODO: 数据库错误 区分用户名已存在和数据库错误
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": res.Message})
	})

	r.POST("/api/user/login", func(c *gin.Context) {
		var req pb.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		res, err := client.Login(ctx, &req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "登录失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token": res.Token,
			"user":  res.User,
		})
	})

	connShortlink, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接 user-service 失败: %v", err)
	}
	clientShortlink := pbShortlink.NewShortlinkServiceClient(connShortlink)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware()) // JWT 鉴权中间件
	{                                     // 创建短链接
		auth.POST("/api/shorten", func(c *gin.Context) {
			var req pbShortlink.ShortenRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			res, err := clientShortlink.ShortenURL(ctx, &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建短链接失败"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"shortlink": res.ShortUrl})
		})

		auth.GET("/api/shorten/top", func(c *gin.Context) {
			req := &pbShortlink.TopRequest{Count: 10}

			// 超时2秒就返回
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			resp, err := clientShortlink.GetTopLinks(ctx, req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "获取排行榜失败"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"top": resp.Top})
		})
	}
	// 跳转
	r.GET("/:short_url", func(c *gin.Context) {
		var req pbShortlink.ResolveRequest
		req.ShortUrl = c.Param("short_url")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		res, err := clientShortlink.Redierect(ctx, &req)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "短链接无效"})
			return
		}
		c.Redirect(http.StatusFound, res.OriginalUrl)
	})

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
