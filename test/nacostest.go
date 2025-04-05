package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

const (
	base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func GenerateShortURL(length int, checkExists func(string) bool) (string, error) {
	if length <= 0 {
		length = 6
	}

	for range 3 {
		// 生成随机字节
		randomBytes := make([]byte, 8)
		if _, err := rand.Read(randomBytes); err != nil {
			return "", err
		}

		// 将时间戳和随机数组合
		timestamp := time.Now().UnixNano()
		fmt.Println("timestamp:", timestamp)
		combined := make([]byte, 16)
		binary.BigEndian.PutUint64(combined[:8], uint64(timestamp))
		copy(combined[8:], randomBytes)
		fmt.Println("combined:", combined)

		// 将组合后的数据转换为base62
		result := make([]byte, length)
		for i := range length {
			// 使用不同的字节位置来增加随机性
			pos := (int(combined[i%8]) + int(combined[(i+8)%8])) % 62
			result[i] = base62Chars[pos]
		}

		shortURL := string(result)

		// 如果不需要检查存在性，或者检查后不存在，则返回
		if checkExists == nil || !checkExists(shortURL) {
			return shortURL, nil
		}
	}

	// 不返回错误，增加纠错机制：在db中查询是否存在，如果存在则重新生成
	return "shortURL", nil
}
func main() {
	// // nacos初始化配置
	// if err := config.InitConfigFromNacos(); err != nil {
	// 	fmt.Println("config init failed", err)
	// }
	// // if err := config.InitConfig("config/config.yaml"); err != nil {
	// // 	fmt.Println("config init failed", err)
	// // }
	// // 打印配置
	// fmt.Println("nacos config:")
	// fmt.Println(config.GlobalConfig)

	shortURL, err := GenerateShortURL(6, nil)
	if err != nil {
		fmt.Println("generate short url failed", err)
	}
	fmt.Println("short url:", shortURL)
}
