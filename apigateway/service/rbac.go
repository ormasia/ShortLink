package service

import (
	"net/http"
	"shortLink/apigateway/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RBACService 处理RBAC相关的请求
type RBACService struct {
	DB *gorm.DB
}

// CreateRole 创建新角色
func (s *RBACService) CreateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if err := s.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建角色失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": role})
}

// UpdateRole 更新角色信息
func (s *RBACService) UpdateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if err := s.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新角色失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": role})
}

// DeleteRole 删除角色
func (s *RBACService) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := s.DB.Delete(&model.Role{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除角色失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功", "data": nil})
}

// ListRoles 获取角色列表
func (s *RBACService) ListRoles(c *gin.Context) {
	var roles []model.Role
	if err := s.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取角色列表失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取成功", "data": roles})
}

// CreatePermission 创建新权限
func (s *RBACService) CreatePermission(c *gin.Context) {
	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if err := s.DB.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建权限失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": permission})
}

// AssignRoleToUser 为用户分配角色
func (s *RBACService) AssignRoleToUser(c *gin.Context) {
	var userRole model.UserRole
	if err := c.ShouldBindJSON(&userRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if err := s.DB.Create(&userRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "分配角色失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "分配成功", "data": nil})
}

// AssignPermissionToRole 为角色分配权限
func (s *RBACService) AssignPermissionToRole(c *gin.Context) {
	var rolePermission model.RolePermission
	if err := c.ShouldBindJSON(&rolePermission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if err := s.DB.Create(&rolePermission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "分配权限失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "分配成功", "data": nil})
}
