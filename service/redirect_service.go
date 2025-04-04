package service

import (
	"errors"
	"fmt"
	"shortLink/cache"
	"shortLink/model"
	"shortLink/pkg"
)

// Resolve 解析短链接
// 参数：
//   - short: 短链接
//
// 返回：
//   - string: 原始URL
//   - error: 错误信息，如果解析成功则为nil
func Resolve(short string) (string, error) {
	// 查缓存
	if url := cache.Get(short); url != "" {
		return url, nil
	}

	// 使用布隆过滤器检查短链接是否存在
	// TODO 怎么持久化布隆过滤器，
	//有三个方案：
	// 1. 使用redis持久化布隆过滤器，需要使用redisbloom
	// 2. 使用文件持久化布隆过滤器，需要使用bloom

	// 3. 程序每次启动时全量从 DB 预热布隆过滤器
	if !cache.MightContain(short) {
		fmt.Println("布隆过滤器不存在该值")
		return "", errors.New("数据不存在")
	}

	// 使用 singleflight 防止缓存击穿
	v, err, _ := pkg.Group.Do(short, func() (any, error) {
		return model.GetOriginalURL(short)
	})
	if err != nil {
		return "", err
	}
	original := v.(string)
	// 缓存结果
	cache.Set(short, original)
	return original, nil
}
