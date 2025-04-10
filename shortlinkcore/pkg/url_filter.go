package pkg

import "strings"

var blacklist = []string{"porn", "hack", "malware"}

// IsValidURL 检查URL是否在黑名单中
// 参数：
//   - url: 需要检查的URL
// 返回：
//   - bool: 如果URL在黑名单中，返回false，否则返回true
func IsValidURL(url string) bool {
	for _, word := range blacklist {
		if strings.Contains(url, word) {
			return false
		}
	}
	return true
}
