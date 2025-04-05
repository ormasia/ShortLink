package router

import (
	"shortLink/controller"

	"shortLink/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 路由注册
func InitRoutes(r *gin.Engine) {
	r.POST("/shorten", controller.ShortenURL)   // 将长URL转换为短链接
	r.GET("/:shortUrl", controller.RedirectURL) // 重定向短链接
}
func InitRoutesWithAuth() *gin.Engine {
	r := gin.Default()

	// ✅ 启用跨域支持（允许前端访问）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 你的前端地址
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// ✅ 公共接口（注册 & 登录）
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)

	// ✅ 登录后才能访问的接口（如生成短链）
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware()) // JWT 鉴权中间件
	{
		auth.POST("/shorten", controller.ShortenURL)
	}

	// ✅ 公开访问的短链接跳转
	r.GET("/:shortUrl", controller.RedirectURL)

	return r
}
