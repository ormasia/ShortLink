package pkg

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"math"
	"time"

	"shortLink/shortlinkcore/config"
	"shortLink/shortlinkcore/logger"

	"go.uber.org/zap"
)

const (
	base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

const (
	// 默认短链接长度
	defaultLength = 6
)

// GenerateShortURL 生成随机短链接
// 参数：
//   - length: 短链接长度，默认为6
//   - checkExists: 检查短链接是否存在的函数，如果为nil则不检查
//
// 返回：
//   - string: 生成的短链接
//   - error: 错误信息
func GenerateShortURL(length int, checkExists func(string) bool) (string, error) {
	logger.Log.Debug("开始生成短链接", zap.Int("length", length))

	if length <= 0 {
		length = defaultLength
		logger.Log.Debug("使用默认短链接长度", zap.Int("defaultLength", defaultLength))
	}

	retryCount := 0
	for range config.GlobalConfig.App.MaxRetries {
		retryCount++
		// 生成随机字节
		randomBytes := make([]byte, 8)
		if _, err := rand.Read(randomBytes); err != nil {
			logger.Log.Error("生成随机字节失败", zap.Error(err))
			return "", err
		}

		// 将时间戳和随机数组合
		timestamp := time.Now().UnixNano()
		combined := make([]byte, 16)
		binary.BigEndian.PutUint64(combined[:8], uint64(timestamp))
		copy(combined[8:], randomBytes)

		// 将组合后的数据转换为base62
		result := make([]byte, length)
		for i := range length {
			// 使用不同的字节位置来增加随机性
			pos := (int(combined[i%8]) + int(combined[(i+8)%16])) % 62
			result[i] = base62Chars[pos]
		}

		shortURL := string(result)
		logger.Log.Debug("生成短链接候选值", zap.String("shortURL", shortURL))

		// 如果不需要检查存在性，或者检查后不存在，则返回
		if checkExists == nil {
			logger.Log.Debug("无需检查短链接是否存在，直接返回")
			return shortURL, nil
		} else if !checkExists(shortURL) {
			logger.Log.Debug("短链接不存在，可以使用")
			return shortURL, nil
		} else {
			logger.Log.Debug("短链接已存在，需要重新生成",
				zap.String("shortURL", shortURL),
				zap.Int("retryCount", retryCount))
		}
	}

	logger.Log.Warn("生成短链接失败，已达到最大重试次数",
		zap.Int("maxRetries", config.GlobalConfig.App.MaxRetries))
	// 不返回错误，增加纠错机制：在db中查询是否存在，如果存在则重新生成
	return "", errors.New("生成短链接失败，请重试")
}

// EncodeID 将ID转换为短链接key
// 参数：
//   - id: 需要转换的ID
//
// 返回：
//   - string: 生成的短链接key
func EncodeID(id int64) string {
	logger.Log.Debug("开始将ID编码为短链接", zap.Int64("id", id))

	if id == 0 {
		logger.Log.Debug("ID为0，返回base62第一个字符")
		return string(base62Chars[0])
	}

	var result []byte
	for id > 0 {
		// 将id对62取余，得到一个0-61之间的数字
		// 将这个数字转换为base62Chars中的一个字符
		// 将这个字符添加到result中
		// 将id除以62
		result = append([]byte{base62Chars[id%62]}, result...)
		id /= 62
	}

	encoded := string(result)
	logger.Log.Debug("ID编码完成", zap.String("encoded", encoded))
	return encoded
}

// CalculatePossibleCombinations 计算可能的组合数
func CalculatePossibleCombinations(length int) int64 {
	return int64(math.Pow(62, float64(length)))
}
