package controller

import (
	"net/http"
	"shortLink/logger"
	"shortLink/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RedirectURL 重定向短链接
// 参数：
//   - c: Gin的上下文
//
// 返回：
//   - 无
func RedirectURL(c *gin.Context) {
	short := c.Param("shortUrl")

	logger.Log.Info("重定向短链接", zap.String("shortUrl", short))

	origin, err := service.Resolve(short)
	if err != nil {
		logger.Log.Error("重定向短链接失败", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "链接不存在"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, origin) //重定向,301
}
