package controller

import (
	"fmt"
	"net/http"
	"shortLink/service"

	"github.com/gin-gonic/gin"
)

// RedirectURL 重定向短链接
// 参数：
//   - c: Gin的上下文
//
// 返回：
//   - 无
func RedirectURL(c *gin.Context) {
	short := c.Param("shortUrl")
	fmt.Println(short)
	origin, err := service.Resolve(short)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "链接不存在"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, origin) //重定向,301
}
