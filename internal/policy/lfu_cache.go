package policy

import "github.com/Code-Hex/go-generics-cache/policy/lfu"

const invalidCacheCapacity = "cache capacity must be greater than zero"

type LFUGraphCache struct {
	cache *lfu.Cache[string, *Graph]
}

func NewLFUGraphCache(capacity int) *LFUGraphCache {
	if capacity <= 0 {
		panic(invalidCacheCapacity)
	}
	return &LFUGraphCache{
		cache: lfu.NewCache[string, *Graph](lfu.WithCapacity(capacity)),
	}
}

func (c *LFUGraphCache) Get(key string) (*Graph, bool) {
	return c.cache.Get(key)
}

func (c *LFUGraphCache) Set(key string, value *Graph) {
	c.cache.Set(key, value)
}
