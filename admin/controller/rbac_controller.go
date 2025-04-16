package controller

import (
	"context"
	"net/http"
	"shortLink/proto/userpb"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// RBACController RBAC权限管理控制器
type RBACController struct {
	rbacClient userpb.RBACServiceClient
}

// NewRBACController 创建RBAC控制器实例
func NewRBACController(conn *grpc.ClientConn) *RBACController {
	return &RBACController{
		rbacClient: userpb.NewRBACServiceClient(conn),
	}
}

// GetUserRoles 获取用户角色
func (c *RBACController) GetUserRoles(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	req := &userpb.GetUserRolesRequest{
		UserId: uint32(userID),
	}

	resp, err := c.rbacClient.GetUserRoles(context.Background(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp.Roles)
}

// GetRolePermissions 获取角色权限
func (c *RBACController) GetRolePermissions(ctx *gin.Context) {
	roleID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	req := &userpb.GetRolePermissionsRequest{
		RoleId: uint32(roleID),
	}

	resp, err := c.rbacClient.GetRolePermissions(context.Background(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp.Permissions)
}

// CreateRole 创建角色
func (c *RBACController) CreateRole(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &userpb.CreateRoleRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	resp, err := c.rbacClient.CreateRole(context.Background(), grpcReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp.Role)
}

// CreatePermission 创建权限
func (c *RBACController) CreatePermission(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Resource    string `json:"resource" binding:"required"`
		Action      string `json:"action" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &userpb.CreatePermissionRequest{
		Name:        req.Name,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
	}

	resp, err := c.rbacClient.CreatePermission(context.Background(), grpcReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp.Permission)
}

// AssignRoleToUser 为用户分配角色
func (c *RBACController) AssignRoleToUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &userpb.AssignRoleToUserRequest{
		UserId: uint32(userID),
		RoleId: uint32(req.RoleID),
	}

	resp, err := c.rbacClient.AssignRoleToUser(context.Background(), grpcReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

// AssignPermissionToRole 为角色分配权限
func (c *RBACController) AssignPermissionToRole(ctx *gin.Context) {
	roleID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	var req struct {
		PermissionID uint `json:"permission_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &userpb.AssignPermissionToRoleRequest{
		RoleId:       uint32(roleID),
		PermissionId: uint32(req.PermissionID),
	}

	resp, err := c.rbacClient.AssignPermissionToRole(context.Background(), grpcReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

// CheckPermission 检查用户权限
func (c *RBACController) CheckPermission(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Resource string `json:"resource" binding:"required"`
		Action   string `json:"action" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &userpb.CheckPermissionRequest{
		UserId:   uint32(userID),
		Resource: req.Resource,
		Action:   req.Action,
	}

	resp, err := c.rbacClient.CheckPermission(context.Background(), grpcReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"has_permission": resp.HasPermission,
		"message":        resp.Message,
	})
}
