package main

import (
	"fmt"
	"log"

	"shortLink/cache"
	"shortLink/config"
	"shortLink/model"
	"shortLink/router"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("program start")

	if err := config.InitConfig("config/config.yaml"); err != nil {
		log.Fatal(err)
	}

	// 初始化数据库
	if err := model.InitDB(config.GlobalConfig.MySQL.GetDSN()); err != nil {
		log.Fatal(err)
	}

	// 初始化Redis
	// TODO:使用函数直接配置
	cache.InitRedis(config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Password, config.GlobalConfig.Redis.DB)
	//初始化布隆过滤器
	cache.InitBloom(2, 2.2)

	// 4. 设置Gin模式
	gin.SetMode(config.GlobalConfig.App.Mode)

	r := gin.Default()
	router.InitRoutes(r)
	// 获取配置文件中的主机和端口
	host, port := config.GlobalConfig.App.GetHost()
	r.Run(fmt.Sprintf("%s:%d", host, port))
}
