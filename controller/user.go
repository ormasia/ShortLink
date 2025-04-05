package controller

import (
	"shortLink/model"
	"shortLink/pkg/jwt"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	var user model.User
	if err := model.GetDB().Where("username = ?", req.Username).First(&user).Error; err != gorm.ErrRecordNotFound {
		c.JSON(409, gin.H{"error": "用户名已存在"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user = model.User{
		Username: req.Username,
		Password: string(hash),
		Nickname: req.Nickname,
		Email:    req.Email,
	}
	if err := model.GetDB().Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "注册失败"})
		return
	}

	c.JSON(200, gin.H{"message": "注册成功"})
}

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	var user model.User
	if err := model.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "用户名不存在"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "密码错误"})
		return
	}

	token, _ := jwt.GenerateToken(user.ID, user.Role, 24*time.Hour)
	c.JSON(200, gin.H{"token": token, "user": gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"email":    user.Email,
		"role":     user.Role,
	}})
}
