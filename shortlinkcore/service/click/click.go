// 点击计数模块
package click

import (
	"context"
	"fmt"

	"shortLink/shortlinkcore/cache"
	"shortLink/shortlinkcore/logger"

	"go.uber.org/zap"
)

func IncrClickCount(shortUrl, originalUrl string) {
	ctx := context.Background()

	logger.Log.Info("记录短链接点击",
		zap.String("shortUrl", shortUrl),
		zap.String("originalUrl", originalUrl))

	// 计数（用于单个点击展示）「记录某个短链总共被点击了多少次」，以便展示或查询，不用于排行。也可以不记录；
	_, err := cache.GetRedis().Incr(ctx, fmt.Sprintf("click:%s-%s", shortUrl, originalUrl)).Result()
	if err != nil {
		logger.Log.Error("增加点击计数失败",
			zap.String("shortUrl", shortUrl),
			zap.String("originalUrl", originalUrl),
			zap.Error(err))
	}

	// ✅ 更新排行榜（ZSet 自增）ZIncrBy 原子操作
	_, err = cache.GetRedis().ZIncrBy(ctx, "shortlink:rank", 1, fmt.Sprintf("%s-%s", shortUrl, originalUrl)).Result()
	if err != nil {
		logger.Log.Error("更新排行榜失败",
			zap.String("shortUrl", shortUrl),
			zap.String("originalUrl", originalUrl),
			zap.Error(err))
	} else {
		logger.Log.Debug("更新排行榜成功",
			zap.String("shortUrl", shortUrl),
			zap.String("originalUrl", originalUrl))
	}

	// // 设置点击量 key 的过期（排行榜不需要）
	// cache.GetRedis().Expire(ctx, fmt.Sprintf("click:%s", shortUrl), 7*24*time.Hour)
}

type ShortLinkRank struct {
	ShortUrl string  `json:"short_url"`
	Clicks   float64 `json:"clicks"`
}

// 获取 Top N 热门短链
func GetTopShortLinks(n int64) ([]ShortLinkRank, error) {
	logger.Log.Info("获取热门短链接排行", zap.Int64("count", n))

	ctx := context.Background()
	raw, err := cache.GetRedis().ZRevRangeWithScores(ctx, "shortlink:rank", 0, n-1).Result()
	if err != nil {
		logger.Log.Error("获取热门短链接排行失败",
			zap.Int64("count", n),
			zap.Error(err))
		return nil, err
	}

	result := make([]ShortLinkRank, 0, len(raw))
	for _, z := range raw {
		result = append(result, ShortLinkRank{
			ShortUrl: fmt.Sprintf("%v", z.Member),
			Clicks:   z.Score,
		})
	}

	logger.Log.Info("获取热门短链接排行成功",
		zap.Int64("count", n),
		zap.Int("resultCount", len(result)))
	return result, nil
}
