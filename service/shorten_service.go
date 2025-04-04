package service

import (
	"errors"
	"shortLink/cache"
	"shortLink/config"
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

	// // 生成唯一ID
	// id, err := model.GenerateID()
	// if err != nil {
	// 	return "", err
	// }
	// // 将ID转换为短链接key
	// shortKey := pkg.EncodeID(id)

	// 生成短链接
	shortKey, err := pkg.GenerateShortURL(config.GlobalConfig.App.Base62Length, nil)
	if err != nil {
		return "", err
	}

	// 将短链接加入布隆过滤器,用于后续判断短链接是否存在
	cache.AddToBloom(shortKey)

	// 将短链接与原始URL的映射关系保存到数据库
	err = model.SaveURLMapping(shortKey, url)
	if err != nil {
		return "", err
	}

	// 将映射关系缓存到Redis中，提高访问速度
	cache.Set(shortKey, url)

	// 返回生成的短链接key
	return shortKey, nil
}
