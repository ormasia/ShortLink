package controller

import (
	"net/http"
	"shortLink/logger"
	"shortLink/model"
	"shortLink/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ShortenRequest struct {
	URL string `json:"url" binding:"required"` // binding是gin框架的标签，用于验证请求参数
}

// ShortenURL 将长URL转换为短链接
// 参数：
//   - c: Gin的上下文
//
// 返回：
//   - 无
func ShortenURL(c *gin.Context) {
	var req ShortenRequest
	if err := c.BindJSON(&req); err != nil {
		logger.Log.Error("参数错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	//判断是否存在 有必要放到redis中吗？
	// shorturlCache := cache.Get(req.URL)
	// if shorturlCache != "" {
	// 	fmt.Println("cache yes")
	// 	c.JSON(http.StatusOK, gin.H{"short_url": shorturlCache})
	// 	return
	// }
	ShortUrlDB := model.IsOriginalURLExist(req.URL)
	if ShortUrlDB != "" {
		logger.Log.Info("DB yes")
		c.JSON(http.StatusOK, gin.H{"short_url": ShortUrlDB})
		return
	}

	//不存在，生成短链接
	shortUrl, err := service.Shorten(req.URL)
	if err != nil {
		logger.Log.Error("生成短链接失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"short_url": shortUrl})
}
