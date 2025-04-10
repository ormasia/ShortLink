package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"shortLink/proto/shortlinkpb"
	"shortLink/shortlinkcore/cache"
	"shortLink/shortlinkcore/config"
	"shortLink/shortlinkcore/logger"
	"shortLink/shortlinkcore/model"
	"shortLink/shortlinkcore/pkg"
	"shortLink/shortlinkcore/pkg/locker"
	"shortLink/shortlinkcore/service/click"
	"time"

	"go.uber.org/zap"
)

type ShortlinkService struct {
	shortlinkpb.UnimplementedShortlinkServiceServer
}

// 生成短链接
func (s *ShortlinkService) ShortenURL(ctx context.Context, req *shortlinkpb.ShortenRequest) (*shortlinkpb.ShortenResponse, error) {
	logger.Log.Info("收到生成短链接请求", zap.String("originalUrl", req.OriginalUrl))

	// 1. 检查数据库是否存在该长链接
	ShortUrlDB := model.IsOriginalURLExist(req.OriginalUrl)
	if ShortUrlDB != "" {
		logger.Log.Info("找到已存在的短链接",
			zap.String("originalUrl", req.OriginalUrl),
			zap.String("shortUrl", ShortUrlDB))
		return &shortlinkpb.ShortenResponse{ShortUrl: ShortUrlDB}, nil
	}

	// 2. 生成短链接
	shortUrl, err := Shorten(req.OriginalUrl)
	if err != nil {
		logger.Log.Error("生成短链接失败",
			zap.String("originalUrl", req.OriginalUrl),
			zap.Error(err))
		return nil, fmt.Errorf("生成短链接失败: %w", err)
	}

	// 3. 返回短链接
	logger.Log.Info("短链接生成成功",
		zap.String("originalUrl", req.OriginalUrl),
		zap.String("shortUrl", shortUrl))
	return &shortlinkpb.ShortenResponse{ShortUrl: shortUrl}, nil
}

// 解析短链接
func (s *ShortlinkService) Redierect(ctx context.Context, req *shortlinkpb.ResolveRequest) (*shortlinkpb.ResolveResponse, error) {
	logger.Log.Info("收到解析短链接请求", zap.String("shortUrl", req.ShortUrl))

	// 1. 解析短链接
	originalURL, err := Resolve(req.ShortUrl)
	if err != nil {
		logger.Log.Error("短链接解析失败",
			zap.String("shortUrl", req.ShortUrl),
			zap.Error(err))
		return nil, fmt.Errorf("短链接不存在: %w", err)
	}

	// 2. 更新点击量
	go click.IncrClickCount(req.ShortUrl, originalURL)

	// 3. 返回原始链接
	logger.Log.Info("短链接解析成功",
		zap.String("shortUrl", req.ShortUrl),
		zap.String("originalUrl", originalURL))
	return &shortlinkpb.ResolveResponse{OriginalUrl: originalURL}, nil
}

func (s *ShortlinkService) GetTopLinks(ctx context.Context, req *shortlinkpb.TopRequest) (*shortlinkpb.TopResponse, error) {
	logger.Log.Info("收到获取热门短链接请求", zap.Int64("count", req.Count))

	// 1. 获取点击量最高的短链接
	rankList, err := click.GetTopShortLinks(req.Count)
	if err != nil {
		logger.Log.Error("获取热门短链接失败", zap.Error(err))
		return nil, err
	}

	// 2. 返回点击量最高的短链接
	items := make([]*shortlinkpb.ShortLinkItem, 0)
	for _, r := range rankList {
		items = append(items, &shortlinkpb.ShortLinkItem{
			ShortUrl: r.ShortUrl,
			Clicks:   r.Clicks,
		})
	}

	// 3. 返回点击量最高的短链接
	logger.Log.Info("获取热门短链接成功",
		zap.Int64("count", req.Count),
		zap.Int("resultCount", len(items)))
	return &shortlinkpb.TopResponse{Top: items}, nil
}

// Shorten 将长URL转换为短链接
// 参数：
//   - url: 需要转换的原始长URL

// 返回：
//   - string: 生成的短链接key
//   - error: 错误信息，如果转换成功则为nil
// func Shorten(url string) (string, error) {
// 	// 验证URL是否合法
// 	if !pkg.IsValidURL(url) {
// 		return "", errors.New("链接非法")
// 	}
// 	// 使用分布式锁
// 	lock := locker.NewRedisLock(cache.GetRedis(), "lock:shorten:"+url, 3*time.Second)

// 	ok, err := lock.TryLock()
// 	if err != nil {
// 		logger.Log.Error("加锁失败", zap.Error(err))
// 		return "", err
// 	}
// 	if !ok {
// 		return "", errors.New("请稍后再试（已被处理）")
// 	}

// 	defer lock.Unlock()
// 	// 生成短链接
// 	shortKey, err := pkg.GenerateShortURL(config.GlobalConfig.App.Base62Length, nil)
// 	if err != nil {
// 		return "", err
// 	}

// 	// 将短链接加入布隆过滤器,用于后续判断短链接是否存在
// 	cache.AddToBloom(shortKey)

// 	// 将短链接与原始URL的映射关系保存到数据库
// 	err = model.SaveURLMapping(shortKey, url)
// 	if err != nil {
// 		return "", err
// 	}

// 	// 将映射关系缓存到Redis中，提高访问速度
// 	cache.Set(shortKey, url)

// 	// 返回生成的短链接key
// 	return shortKey, nil
// }

func Shorten(longUrl string) (string, error) {
	// 1. 校验 URL 合法性
	if !pkg.IsValidURL(longUrl) {
		return "", errors.New("链接非法")
	}
	// 2. 分布式锁（对 URL 做哈希防止 key 过长）防止并发过程中生成重复短链
	urlHash := fmt.Sprintf("%x", sha256.Sum256([]byte(longUrl)))
	lockKey := "lock:shorten:" + urlHash
	lock := locker.NewRedisLock(cache.GetRedis(), lockKey, 3*time.Second)
	ok, err := lock.TryLock()
	if err != nil {
		logger.Log.Error("加锁失败", zap.Error(err), zap.String("url", longUrl))
		return "", errors.New("系统繁忙，请稍后重试")
	}
	if !ok {
		return "", errors.New("操作频繁，请稍后重试")
	}
	defer func() {
		logger.Log.Debug("释放分布式锁", zap.String("lockKey", lockKey))
		if err := lock.Unlock(); err != nil {
			logger.Log.Warn("解锁失败", zap.Error(err), zap.String("url", longUrl))
		}
	}()

	// // 3. 加锁后再次检查缓存或数据库（幂等）
	// ShortUrlDB := model.IsOriginalURLExist(longUrl)
	// if ShortUrlDB != "" {
	// 	logger.Log.Info("加锁后发现已存在短链接",
	// 		zap.String("originalUrl", longUrl),
	// 		zap.String("shortUrl", ShortUrlDB))
	// 	cache.Set(ShortUrlDB, longUrl)
	// 	return ShortUrlDB, nil
	// }

	// 4. 生成短链 Key（Base62）
	shortKey, err := pkg.GenerateShortURL(config.GlobalConfig.App.Base62Length, cache.MightContain)
	if err != nil {
		logger.Log.Error("短链生成失败", zap.Error(err))
		return "", errors.New("生成失败")
	}

	// 5. 更新布隆过滤器
	cache.AddToBloom(shortKey)

	// 6. 持久化数据库
	if err := model.SaveURLMapping(shortKey, longUrl); err != nil {
		logger.Log.Error("数据库保存失败", zap.Error(err), zap.String("shortKey", shortKey))
		return "", errors.New("持久化失败")
	}

	// 7. 写入 Redis 缓存
	cache.Set(shortKey, longUrl)

	logger.Log.Info("短链生成成功",
		zap.String("shortKey", shortKey),
		zap.String("url", longUrl),
	)

	return shortKey, nil
}

// Resolve 解析短链接
// 参数：
//   - short: 短链接
//
// 返回：
//   - string: 原始URL
//   - error: 错误信息，如果解析成功则为nil
func Resolve(short string) (string, error) {
	logger.Log.Debug("开始解析短链接", zap.String("shortUrl", short))

	// 查缓存
	if url := cache.Get(short); url != "" {
		logger.Log.Debug("从缓存中获取到原始链接",
			zap.String("shortUrl", short),
			zap.String("originalUrl", url))
		return url, nil
	}

	// 使用布隆过滤器检查短链接是否存在
	// TODO 怎么持久化布隆过滤器，
	//有三个方案：
	// 1. 使用redis持久化布隆过滤器，需要使用redisbloom
	// 2. 使用文件持久化布隆过滤器，需要使用bloom
	// 3. 程序每次启动时全量从 DB 预热布隆过滤器

	if !cache.MightContain(short) {
		logger.Log.Warn("布隆过滤器不存在该值", zap.String("shortUrl", short))
		return "", errors.New("数据不存在")
	}

	// 使用 singleflight 防止缓存击穿
	logger.Log.Debug("使用singleflight从数据库获取原始链接", zap.String("shortUrl", short))
	v, err, _ := pkg.Group.Do(short, func() (any, error) {
		return model.GetOriginalURL(short)
	})
	if err != nil {
		logger.Log.Error("从数据库获取原始链接失败",
			zap.String("shortUrl", short),
			zap.Error(err))
		return "", err
	}
	original := v.(string)

	// 缓存结果
	cache.Set(short, original)
	logger.Log.Debug("解析短链接成功并更新缓存",
		zap.String("shortUrl", short),
		zap.String("originalUrl", original))
	return original, nil
}
