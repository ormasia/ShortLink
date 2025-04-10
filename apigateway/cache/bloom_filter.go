package cache

import (
	"shortLink/apigateway/model"
	"sync"

	"github.com/bits-and-blooms/bloom/v3"
)

var (
	bloomFilter *bloom.BloomFilter
	once        sync.Once
)

// 初始化布隆过滤器
// n:预计插入的数据量 fp:误判率
func InitBloom(n uint, fp float64) { //单例模式
	once.Do(func() {
		bloomFilter = bloom.NewWithEstimates(n, fp)
	})
}

// 添加数据到布隆过滤器
func AddToBloom(data string) {
	if bloomFilter != nil {
		bloomFilter.Add([]byte(data))
	}
}

// 判断数据是否在布隆过滤器中
// Test returns true if the data is in the BloomFilter, false otherwise.
// If true, the result might be a false positive.
// If false, the data is definitely not in the set.
func MightContain(data string) bool {
	if bloomFilter == nil {
		return false
	}
	return bloomFilter.Test([]byte(data))
}

func WarmUpBloomFromDB() {
	urls := model.GetAllShortUrls()
	for _, short := range urls {
		AddToBloom(short)
	}
}
