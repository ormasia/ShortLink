package router

import (
	"shortLink/admin/controller"
	"shortLink/admin/middleware"
	"shortLink/proto/userpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// SetupAdminRouter 设置后台管理路由
func SetupAdminRouter(r *gin.Engine, conn *grpc.ClientConn) {
	// 创建控制器实例
	rbacController := controller.NewRBACController(conn)

	// 后台管理路由组
	admin := r.Group("/admin")
	admin.Use(middleware.AdminRequired(userpb.NewRBACServiceClient(conn)))

	// 用户角色管理
	users := admin.Group("/users")
	{
		users.GET("/:id/roles", rbacController.GetUserRoles)                // 获取用户角色
		users.POST("/:id/roles", rbacController.AssignRoleToUser)           // 为用户分配角色
		users.POST("/:id/check-permission", rbacController.CheckPermission) // 检查用户权限
	}

	// 角色管理
	roles := admin.Group("/roles")
	{
		roles.POST("", rbacController.CreateRole)                             // 创建角色
		roles.GET("/:id/permissions", rbacController.GetRolePermissions)      // 获取角色权限
		roles.POST("/:id/permissions", rbacController.AssignPermissionToRole) // 为角色分配权限
	}

	// 权限管理
	permissions := admin.Group("/permissions")
	{
		permissions.POST("", rbacController.CreatePermission) // 创建权限
	}
}
