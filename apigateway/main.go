package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"shortLink/apigateway/cache"
	"shortLink/apigateway/config"
	"shortLink/apigateway/middleware"
	pbShortlink "shortLink/proto/shortlinkpb"
	pb "shortLink/proto/userpb"
	"strconv"
	"time"

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

	r.POST("/api/v1/users", func(c *gin.Context) {
		var req pb.RegisterRequest
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		// 调用用户服务注册，已经检验过参数，所以这里有可能的错误是用户名已存在，数据库错误
		res, registerErr := client.Register(ctx, &req)
		if registerErr != nil {
			// TODO: 数据库错误 区分用户名已存在和数据库错误--不用区分，用户存在不返回错误
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册失败", "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": res.Message, "data": res.Message})
	})

	r.POST("/api/v1/users/login", func(c *gin.Context) {
		var req pb.LoginRequest
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		res, loginErr := client.Login(ctx, &req)
		if loginErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "登录失败", "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "登录成功",
			"data": gin.H{
				"token": res.Token,
				"user":  res.User,
			},
		})
	})

	r.POST("/api/v1/users/logout", func(c *gin.Context) {
		var req pb.LogoutRequest
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		res, logoutErr := client.Logout(ctx, &req)
		if logoutErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "登出失败", "data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "登出成功",
			"data": gin.H{
				"message": res.Message,
			},
		})
	})

	connShortlink, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接 user-service 失败: %v", err)
	}
	clientShortlink := pbShortlink.NewShortlinkServiceClient(connShortlink)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware()) // JWT 鉴权中间件
	{
		// 创建短链接
		auth.POST("/api/v1/links", func(c *gin.Context) {
			var req pbShortlink.ShortenRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
				return
			}
			req.UserId = strconv.Itoa(int(c.GetUint("UserID")))
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			res, err := clientShortlink.ShortenURL(ctx, &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建短链接失败", "data": nil})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": gin.H{"shortlink": res.ShortUrl}})
		})

		// 批量生成短链接 - 添加特殊的批量限流中间件
		auth.POST("/api/v1/links/batch", middleware.BatchRateLimitMiddleware(), func(c *gin.Context) {
			var req pbShortlink.BatchShortenRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
				return
			}

			// 检查URL列表是否为空
			if len(req.OriginalUrls) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "URL列表不能为空", "data": nil})
				return
			}

			// 设置超时时间，批量处理可能需要更长时间
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
			defer cancel()

			req.UserId = strconv.Itoa(int(c.GetUint("UserID")))
			req.Concurrency = 10
			res, err := clientShortlink.BatchShortenURLs(ctx, &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "批量生成短链接失败", "data": nil})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "批量创建成功",
				"data": gin.H{
					"results":       res.Results,
					"total_count":   res.TotalCount,
					"success_count": res.SuccessCount,
					"elapsed_time":  res.ElapsedTime,
				},
			})
		})

		auth.GET("/api/v1/links/top", func(c *gin.Context) {
			req := &pbShortlink.TopRequest{Count: 10}

			// 超时2秒就返回
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			resp, err := clientShortlink.GetTopLinks(ctx, req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取排行榜失败", "data": nil})
				return
			}

			c.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取成功", "data": gin.H{"top": resp.Top}})
		})

		// 删除用户的所有短链接
		auth.DELETE("/api/v1/links", func(c *gin.Context) {
			userID := strconv.Itoa(int(c.GetUint("UserID")))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req := &pbShortlink.DeleteUserURLsRequest{
				UserId: userID,
			}

			res, err := clientShortlink.DeleteUserURLs(ctx, req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "删除短链接失败",
					"data":    nil,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "删除成功",
				"data": gin.H{
					"deleted_count": res.DeletedCount,
				},
			})
		})
	}
	// 跳转接口也添加限流
	r.GET("/api/v1/links/:short_url", middleware.RateLimitMiddleware(), func(c *gin.Context) {
		var req pbShortlink.ResolveRequest
		req.ShortUrl = c.Param("short_url")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		res, err := clientShortlink.Redierect(ctx, &req)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "短链接无效", "data": nil})
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
