package cache

import (
	"sync"

	"github.com/bits-and-blooms/bloom/v3"
)

var (
	bloomFilter *bloom.BloomFilter
	once        sync.Once
)

func InitBloom(n uint, fp float64) {
	once.Do(func() {
		bloomFilter = bloom.NewWithEstimates(n, fp)
	})
}

func AddToBloom(data string) {
	if bloomFilter != nil {
		bloomFilter.Add([]byte(data))
	}
}

func MightContain(data string) bool {
	if bloomFilter == nil {
		return false
	}
	return bloomFilter.Test([]byte(data))
}
