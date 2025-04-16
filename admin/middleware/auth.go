package middleware

import (
	"context"
	"net/http"
	"shortLink/proto/userpb"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your_jwt_secret")

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RequirePermission 权限检查中间件
func RequirePermission(rbacClient userpb.RBACServiceClient, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
			c.Abort()
			return
		}

		// 提取token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析token
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 获取claims
		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token格式"})
			c.Abort()
			return
		}

		// 设置用户ID到上下文
		c.Set("user_id", claims.UserID)

		// 检查权限
		req := &userpb.CheckPermissionRequest{
			UserId:   uint32(claims.UserID),
			Resource: resource,
			Action:   action,
		}

		resp, err := rbacClient.CheckPermission(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !resp.HasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": resp.Message})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminRequired 管理员权限检查中间件
func AdminRequired(rbacClient userpb.RBACServiceClient) gin.HandlerFunc {
	return RequirePermission(rbacClient, "admin", "access")
}

// UserRequired 用户检查中间件
func UserRequired(rbacClient userpb.RBACServiceClient) gin.HandlerFunc {
	return RequirePermission(rbacClient, "user", "manage")
}

// LinkRequired 链接检查中间件
func LinkRequired(rbacClient userpb.RBACServiceClient) gin.HandlerFunc {
	return RequirePermission(rbacClient, "link", "manage")
}

// RoleRequired 角色检查中间件
func RoleRequired(rbacClient userpb.RBACServiceClient) gin.HandlerFunc {
	return RequirePermission(rbacClient, "role", "manage")
}

// PermissionRequired 权限检查中间件
func PermissionRequired(rbacClient userpb.RBACServiceClient) gin.HandlerFunc {
	return RequirePermission(rbacClient, "permission", "manage")
}
