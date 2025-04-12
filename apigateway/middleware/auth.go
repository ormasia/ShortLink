package middleware

import (
	"net/http"
	"strings"

	"shortLink/apigateway/cache"
	"shortLink/apigateway/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取 Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
			c.Abort()
			return
		}

		// 提取 token 字符串
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析 token
		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "解析-无效或过期的 Token"})
			c.Abort()
			return
		}
		// 校验 token 是否在 Redis 中存在（即是否有效）
		cachedToken := cache.Get(tokenStr)
		if cachedToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "校验-无效或过期的 Token"})
			c.Abort()
			return
		}
		// 注入上下文（便于控制器获取 userID/role）
		c.Set("UserID", claims.UserID)
		c.Set("Role", claims.Role)

		c.Next() // 放行
	}
}
