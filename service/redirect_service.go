package service

import (
	"shortLink/cache"
	"shortLink/model"

	"golang.org/x/sync/singleflight"
)

var g singleflight.Group

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
	// 使用 singleflight 防止缓存击穿
	v, err, _ := g.Do(short, func() (any, error) {
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
