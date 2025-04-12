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
	"shortLink/shortlinkcore/pkg/gopool"
	"shortLink/shortlinkcore/pkg/locker"
	"shortLink/shortlinkcore/pkg/safebrowsing"
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
	shortUrl, err := Shorten(req.OriginalUrl, req.UserId)
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

	// 2. 异步更新点击量
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

func Shorten(longUrl, userID string) (string, error) {
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
		if err = lock.Unlock(); err != nil {
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
	if err := model.SaveURLMappingWithUserID(shortKey, longUrl, userID); err != nil {
		logger.Log.Error("数据库保存失败", zap.Error(err), zap.String("shortKey", shortKey))
		return "", errors.New("持久化失败")
	}

	// 6.1 异步安全检查
	// 使用协程池进行异步安全检查，避免阻塞主流程
	// 如果发现不安全URL，会更新数据库状态为blocked
	pool := gopool.GetPool()
	pool.Submit(func() {
		logger.Log.Info("开始安全检查", zap.String("url", longUrl))

		isSafe, threatType, err := safebrowsing.CheckURL(longUrl)
		if err != nil {
			logger.Log.Error("安全检查失败",
				zap.String("url", longUrl),
				zap.Error(err))
			return
		}

		//
		// 如果URL不安全，需要：
		// 1. 从缓存中删除该短链接
		// 2. 记录警告日志，包含威胁类型
		// 3. 更新数据库中的URL状态为blocked
		// 4. 记录操作日志
		if !isSafe {
			cache.Del(shortKey)
			logger.Log.Warn("发现不安全URL",
				zap.String("url", longUrl),
				zap.String("threatType", threatType))

			// 更新数据库状态为blocked
			err = model.UpdateStatus(shortKey, "blocked", threatType)
			if err != nil {
				logger.Log.Error("更新URL状态失败",
					zap.String("shortURL", shortKey),
					zap.Error(err))
				return
			}
			logger.Log.Info("已封禁不安全URL",
				zap.String("shortURL", shortKey),
				zap.String("threatType", threatType))
			return
		}

		logger.Log.Info("URL安全检查通过",
			zap.String("url", longUrl))
	})

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

// DeleteUserURLs 删除用户的所有短链接
func (s *ShortlinkService) DeleteUserURLs(ctx context.Context, req *shortlinkpb.DeleteUserURLsRequest) (*shortlinkpb.DeleteUserURLsResponse, error) {
	logger.Log.Info("收到删除用户短链接请求", zap.String("userId", req.UserId))

	// 1. 获取用户的所有短链接
	var mappings []model.URLMapping
	if err := model.GetDB().Where("user_id = ?", req.UserId).Find(&mappings).Error; err != nil {
		logger.Log.Error("获取用户短链接失败", zap.String("userId", req.UserId), zap.Error(err))
		return nil, fmt.Errorf("获取用户短链接失败: %w", err)
	}

	if len(mappings) == 0 {
		logger.Log.Info("用户没有短链接", zap.String("userId", req.UserId))
		return &shortlinkpb.DeleteUserURLsResponse{DeletedCount: 0}, nil
	}

	// 保持正确的删除顺序：先删除数据库，再删除缓存
	/* 	线程A删除缓存
	线程B查询数据库，发现数据还存在
	线程B将数据写入缓存
	线程A删除数据库
	最终导致缓存和数据库不一致*/

	// 2. 先删除数据库记录
	result := model.GetDB().Where("user_id = ?", req.UserId).Delete(&model.URLMapping{})
	if result.Error != nil {
		logger.Log.Error("删除用户短链接失败", zap.String("userId", req.UserId), zap.Error(result.Error))
		return nil, fmt.Errorf("删除用户短链接失败: %w", result.Error)
	}

	deletedCount := int32(result.RowsAffected)

	// 3. 删除Redis缓存和点击量
	redis := cache.GetRedis()
	for _, mapping := range mappings {
		// 删除短链接缓存
		redis.Del(ctx, mapping.ShortURL)
		// 删除点击量
		redis.Del(ctx, fmt.Sprintf("click:%s-%s", mapping.ShortURL, mapping.OriginalURL))
		// 从排行榜中删除
		redis.ZRem(ctx, "shortlink:rank", fmt.Sprintf("%s-%s", mapping.ShortURL, mapping.OriginalURL))
	}

	logger.Log.Info("删除用户短链接成功",
		zap.String("userId", req.UserId),
		zap.Int32("deletedCount", deletedCount))

	return &shortlinkpb.DeleteUserURLsResponse{DeletedCount: deletedCount}, nil
}
