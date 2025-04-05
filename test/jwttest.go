package main

import (
	"fmt"
	"shortLink/controller"
	"shortLink/pkg/jwt"
	"time"

	"github.com/gin-gonic/gin"
)

func jwtTest() {
	token, err := jwt.GenerateToken(1, "admin", time.Hour*24)
	if err != nil {
		fmt.Println("生成token失败", err)
	}
	fmt.Println("token", token)
	claims, err := jwt.ParseToken(token)
	if err != nil {
		fmt.Println("解析token失败", err)
	}
	fmt.Println("claims", claims)
}

func Logintest() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.Run(":8080")
}

func main() {
	jwtTest()
	Logintest()
}
