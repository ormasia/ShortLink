package service

import (
	"errors"
	"shortLink/cache"
	"shortLink/model"
	"shortLink/pkg"
)

// Shorten 将长URL转换为短链接
// 参数：
//   - url: 需要转换的原始长URL
//
// 返回：
//   - string: 生成的短链接key
//   - error: 错误信息，如果转换成功则为nil
func Shorten(url string) (string, error) {
	// 验证URL是否合法
	if !pkg.IsValidURL(url) {
		return "", errors.New("链接非法")
	}

	// 生成唯一ID
	id := model.GenerateID()
	// 将ID转换为短链接key
	shortKey := pkg.EncodeID(id)

	// 将短链接与原始URL的映射关系保存到数据库
	err := model.SaveURLMapping(shortKey, url)
	if err != nil {
		return "", err
	}

	// 将映射关系缓存到Redis中，提高访问速度
	cache.Set(shortKey, url)

	// 返回生成的短链接key
	return shortKey, nil
}
