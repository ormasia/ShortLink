package main

import (
	"fmt"
	"shortLink/config"
)

func main() {
	// nacos初始化配置
	if err := config.InitConfigFromNacos(); err != nil {
		fmt.Println("config init failed", err)
	}
	// if err := config.InitConfig("config/config.yaml"); err != nil {
	// 	fmt.Println("config init failed", err)
	// }
	// 打印配置
	fmt.Println("nacos config:")
	fmt.Println(config.GlobalConfig)
}
