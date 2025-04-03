package router

import (
	"shortLink/controller"

	"github.com/gin-gonic/gin"
)

// 路由注册
func InitRoutes(r *gin.Engine) {
	r.POST("/shorten", controller.ShortenURL)   // 将长URL转换为短链接
	r.GET("/:shortUrl", controller.RedirectURL) // 重定向短链接
}
