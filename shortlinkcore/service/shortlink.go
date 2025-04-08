package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"shortLink/proto/shortlinkpb"
	"shortLink/shortlinkcore/cache"
	"shortLink/shortlinkcore/config"
	"shortLink/shortlinkcore/model"
	"shortLink/shortlinkcore/pkg"
	"shortLink/shortlinkcore/service/click"
	"shortLink/userservice/logger"
	"shortLink/userservice/pkg/locker"
	"time"

	"go.uber.org/zap"
)

type ShortlinkService struct {
	shortlinkpb.UnimplementedShortlinkServiceServer
	// DB    *gorm.DB
	// Cache *redis.Client
}

// 生成短链接
func (s *ShortlinkService) ShortenURL(ctx context.Context, req *shortlinkpb.ShortenRequest) (*shortlinkpb.ShortenResponse, error) {
	ShortUrlDB := model.IsOriginalURLExist(req.OriginalUrl)
	if ShortUrlDB != "" {
		return &shortlinkpb.ShortenResponse{ShortUrl: ShortUrlDB}, nil
	}
	shortUrl, err := Shorten(req.OriginalUrl)
	if err != nil {
		return nil, fmt.Errorf("生成短链接失败: %w", err)
	}

	return &shortlinkpb.ShortenResponse{ShortUrl: shortUrl}, nil
}

// 解析短链接
func (s *ShortlinkService) Redierect(ctx context.Context, req *shortlinkpb.ResolveRequest) (*shortlinkpb.ResolveResponse, error) {
	originalURL, err := Resolve(req.ShortUrl)
	if err != nil {
		return nil, fmt.Errorf("短链接不存在: %w", err)
	}

	// 更新点击量
	go click.IncrClickCount(req.ShortUrl, originalURL)

	return &shortlinkpb.ResolveResponse{OriginalUrl: originalURL}, nil
}

func (s *ShortlinkService) GetTopLinks(ctx context.Context, req *shortlinkpb.TopRequest) (*shortlinkpb.TopResponse, error) {
	rankList, err := click.GetTopShortLinks(req.Count)
	if err != nil {
		return nil, err
	}

	items := make([]*shortlinkpb.ShortLinkItem, 0)
	for _, r := range rankList {
		items = append(items, &shortlinkpb.ShortLinkItem{
			ShortUrl: r.ShortUrl,
			Clicks:   r.Clicks,
		})
	}
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

	// 3. 分布式锁（对 URL 做哈希防止 key 过长）防止并发过程中生成重复短链
	urlHash := fmt.Sprintf("%x", sha256.Sum256([]byte(longUrl)))
	lockKey := "lock:shorten:" + urlHash
	lock := locker.NewRedisLock(cache.GetRedis(), lockKey, 3*time.Second)
	fmt.Println(lockKey, "lockKey")
	ok, err := lock.TryLock()
	if err != nil {
		logger.Log.Error("加锁失败", zap.Error(err), zap.String("url", longUrl))
		return "", errors.New("系统繁忙，请稍后重试")
	}
	if !ok {
		return "", errors.New("操作频繁，请稍后重试")
	}
	defer func() {
		fmt.Println("解锁")
		if err := lock.Unlock(); err != nil {
			logger.Log.Warn("解锁失败", zap.Error(err), zap.String("url", longUrl))
		}
	}()

	// 4. 加锁后再次检查缓存或数据库（幂等）
	ShortUrlDB := model.IsOriginalURLExist(longUrl)
	if ShortUrlDB != "" {
		cache.Set(ShortUrlDB, longUrl)
		return ShortUrlDB, nil
	}

	// 5. 生成短链 Key（Base62）
	shortKey, err := pkg.GenerateShortURL(config.GlobalConfig.App.Base62Length, cache.MightContain)
	if err != nil {
		// logger.Log.Error("短链生成失败", zap.Error(err))
		return "", errors.New("生成失败")
	}

	// 6. 更新布隆过滤器
	cache.AddToBloom(shortKey)

	// 7. 持久化数据库
	if err := model.SaveURLMapping(shortKey, longUrl); err != nil {
		// logger.Log.Error("数据库保存失败", zap.Error(err))
		return "", errors.New("持久化失败")
	}

	// 8. 写入 Redis 缓存
	cache.Set(shortKey, longUrl)

	// logger.Log.Info("短链生成成功",
	// 	zap.String("shortKey", shortKey),
	// 	zap.String("url", longUrl),
	// )

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
