package middleware

// import (
// 	"fmt"
// 	"net/http"
// 	"shortLink/apigateway/cache"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-redis/redis/v8"
// )

// const (
// 	QueueKey = "request_queue"
// 	MaxWaitTime = 30 * time.Second
// 	CheckInterval = 100 * time.Millisecond
// )

// // QueueMiddleware 实现请求排队的中间件
// func QueueMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		redisClient := cache.GetRedis()

// 		// 获取当前时间戳作为请求的序号
// 		score := float64(time.Now().UnixNano())
// 		member := fmt.Sprintf("%f", score)

// 		// 将请求加入队列
// 		if err := redisClient.ZAdd(c, QueueKey, &redis.Z{Score: score, Member: member}).Err(); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法加入队列"})
// 			c.Abort()
// 			return
// 		}
// 		defer redisClient.ZRem(c, QueueKey, member)

// 		// 等待轮到自己
// 		start := time.Now()
// 		for {
// 			// 检查是否超时
// 			if time.Since(start) > MaxWaitTime {
// 				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "队列等待超时"})
// 				c.Abort()
// 				return
// 			}

// 			// 获取队列中最小的序号
// 			result, err := redisClient.ZRangeWithScores(c, QueueKey, 0, 0).Result()
// 			if err != nil || len(result) == 0 {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "队列错误"})
// 				c.Abort()
// 				return
// 			}

// 			// 如果当前请求是队列中的第一个，则处理请求
// 			if result[0].Member.(string) == member {
// 				break
// 			}

// 			// 等待一段时间后再次检查
// 			time.Sleep(CheckInterval)
// 		}

// 		// 获取当前队列长度
// 		queueLen, err := redisClient.ZCard(c, QueueKey).Result()
// 		if err == nil {
// 			c.Header("X-Queue-Position", strconv.FormatInt(queueLen, 10))
// 		}

// 		// 继续处理请求
// 		c.Next()
// 	}
// }
