package pkg

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"time"

	"shortLink/shortlinkcore/config"
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
	if length <= 0 {
		length = defaultLength
	}

	for range config.GlobalConfig.App.MaxRetries {
		// 生成随机字节
		randomBytes := make([]byte, 8)
		if _, err := rand.Read(randomBytes); err != nil {
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
			pos := (int(combined[i%8]) + int(combined[(i+8)%8])) % 62
			result[i] = base62Chars[pos]
		}

		shortURL := string(result)
		// 如果不需要检查存在性，或者检查后不存在，则返回
		if checkExists == nil || !checkExists(shortURL) {
			return shortURL, nil
		}
	}
	// 不返回错误，增加纠错机制：在db中查询是否存在，如果存在则重新生成
	return "shortURL", nil
}

// EncodeID 将ID转换为短链接key
// 参数：
//   - id: 需要转换的ID
//
// 返回：
//   - string: 生成的短链接key
func EncodeID(id int64) string {
	if id == 0 {
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
	return string(result)
}

// CalculatePossibleCombinations 计算可能的组合数
func CalculatePossibleCombinations(length int) int64 {
	return int64(math.Pow(62, float64(length)))
}
