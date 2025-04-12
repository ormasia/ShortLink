package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"shortLink/proto/shortlinkpb"
	"shortLink/shortlinkcore/logger"
	"shortLink/shortlinkcore/model"

	"go.uber.org/zap"
)

// BatchShortenResult 表示批量生成短链接的结果
type BatchShortenResult struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	Error       string `json:"error,omitempty"`
}

// BatchShortenURLs 批量生成短链接
// 参数：
//   - ctx: 上下文
//   - urls: 需要转换的原始长URL列表
//   - concurrency: 并发处理的数量，默认为10
//
// 返回：
//   - []BatchShortenResult: 批量生成结果
//   - error: 错误信息
func BatchShortenURLs(ctx context.Context, urls []string, userID string, concurrency int) ([]BatchShortenResult, error) {
	logger.Log.Info("收到批量生成短链接请求",
		zap.Int("urlCount", len(urls)),
		zap.Int("concurrency", concurrency))

	if len(urls) == 0 {
		return nil, errors.New("URL列表为空")
	}

	// 设置默认并发数
	if concurrency <= 0 {
		concurrency = 10
	}

	// 限制最大并发数
	if concurrency > 50 {
		concurrency = 50
	}

	// 预检查数据库中是否已存在这些URL
	results := make([]BatchShortenResult, 0, len(urls))
	var urlsToProcess []string

	for _, url := range urls {
		if shortURL := model.IsOriginalURLExist(url); shortURL != "" {
			// URL已存在，直接使用已有的短链接
			results = append(results, BatchShortenResult{
				OriginalURL: url,
				ShortURL:    shortURL,
			})
			logger.Log.Debug("使用已存在的短链接",
				zap.String("originalUrl", url),
				zap.String("shortUrl", shortURL))
		} else {
			// URL不存在，加入待处理列表
			urlsToProcess = append(urlsToProcess, url)
		}
	}

	// 如果所有URL都已存在，直接返回结果
	if len(urlsToProcess) == 0 {
		return results, nil
	}

	// 创建结果通道，只处理不存在的URL
	resultChan := make(chan BatchShortenResult, len(urlsToProcess))

	// 创建工作池
	var wg sync.WaitGroup
	// 使用信号量控制并发数
	semaphore := make(chan struct{}, concurrency)

	// 处理需要新生成短链接的URL
	for i, url := range urlsToProcess {
		wg.Add(1)
		go func(index int, originalURL string) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 检查上下文是否已取消
			select {
			case <-ctx.Done():
				resultChan <- BatchShortenResult{
					OriginalURL: originalURL,
					Error:       "请求已取消",
				}
				return
			default:
				// 继续处理
			}

			// 生成短链接
			shortURL, err := Shorten(originalURL, userID)
			result := BatchShortenResult{OriginalURL: originalURL}

			if err != nil {
				result.Error = err.Error()
				logger.Log.Warn("批量生成短链接失败",
					zap.String("originalUrl", originalURL),
					zap.Error(err))
			} else {
				result.ShortURL = shortURL
				logger.Log.Debug("批量生成短链接成功",
					zap.String("originalUrl", originalURL),
					zap.String("shortUrl", shortURL))
			}

			resultChan <- result
		}(i, url)
	}

	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集新生成的短链接结果
	for result := range resultChan {
		results = append(results, result)
	}

	logger.Log.Info("批量生成短链接完成",
		zap.Int("totalCount", len(urls)),
		zap.Int("successCount", countSuccesses(results)))

	return results, nil
}

// countSuccesses 计算成功生成的短链接数量
func countSuccesses(results []BatchShortenResult) int {
	count := 0
	for _, result := range results {
		if result.Error == "" {
			count++
		}
	}
	return count
}

// BatchShortenRequest 表示批量生成短链接的请求
type BatchShortenRequest struct {
	URLs        []string `json:"urls"`
	Concurrency int      `json:"concurrency,omitempty"`
}

// BatchShortenResponse 表示批量生成短链接的响应
type BatchShortenResponse struct {
	Results      []BatchShortenResult `json:"results"`
	TotalCount   int                  `json:"total_count"`
	SuccessCount int                  `json:"success_count"`
	ElapsedTime  string               `json:"elapsed_time"`
}

// BatchShortenHandler 处理批量生成短链接的请求
func (s *ShortlinkService) BatchShortenURLs(ctx context.Context, req *shortlinkpb.BatchShortenRequest) (*shortlinkpb.BatchShortenResponse, error) {
	logger.Log.Info("收到批量生成短链接请求", zap.Int("urlCount", len(req.OriginalUrls)))

	startTime := time.Now()

	// 调用批量生成函数
	results, err := BatchShortenURLs(ctx, req.OriginalUrls, req.UserId, int(req.Concurrency))
	if err != nil {
		logger.Log.Error("批量生成短链接失败", zap.Error(err))
		return nil, fmt.Errorf("批量生成短链接失败: %w", err)
	}

	// 转换结果为protobuf格式
	pbResults := make([]*shortlinkpb.BatchShortenResult, 0, len(results))
	for _, result := range results {
		pbResult := &shortlinkpb.BatchShortenResult{
			OriginalUrl: result.OriginalURL,
			ShortUrl:    result.ShortURL,
			Error:       result.Error,
		}
		pbResults = append(pbResults, pbResult)
	}

	elapsedTime := time.Since(startTime)
	logger.Log.Info("批量生成短链接完成",
		zap.Int("totalCount", len(results)),
		zap.Int("successCount", countSuccesses(results)),
		zap.Duration("elapsedTime", elapsedTime))

	return &shortlinkpb.BatchShortenResponse{
		Results:      pbResults,
		TotalCount:   int32(len(results)),
		SuccessCount: int32(countSuccesses(results)),
		ElapsedTime:  elapsedTime.String(),
	}, nil
}
