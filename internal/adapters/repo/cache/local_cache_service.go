package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/coocood/freecache"
)

type Duration time.Duration

var (
	cacheKeysWithPrefix = make(map[string]map[string]struct{})
	keyPrefixMu         sync.Mutex
)

type LocalCacheService interface {
	Set(prefix, key string, ttl time.Duration, value interface{})
	Get(prefix, key string, result interface{}) (bool, error)
	Evict(prefix, key string)
	EvictByPrefix(prefix string)
}

type cache struct {
	cache *freecache.Cache
}

func NewLocalCacheService() LocalCacheService {
	return &cache{
		cache: freecache.NewCache(100 * 1024 * 1024),
	}
}

func buildFullKey(prefix, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

func (c *cache) Set(prefix, key string, ttl time.Duration, value interface{}) {
	fullKey := buildFullKey(prefix, key)

	jsonData, err := json.Marshal(value)

	if err != nil {
		// TODO LOG
		return
	}

	c.cache.Set([]byte(fullKey), jsonData, int(ttl.Seconds()))
	keyPrefixMu.Lock()
	defer keyPrefixMu.Unlock()

	if _, exists := cacheKeysWithPrefix[prefix]; !exists {
		cacheKeysWithPrefix[prefix] = make(map[string]struct{})
	}

	if _, exists := cacheKeysWithPrefix[prefix][fullKey]; !exists {
		cacheKeysWithPrefix[prefix][fullKey] = struct{}{}
	}
}

func (c *cache) Get(prefix, key string, result interface{}) (bool, error) {
	fullKey := buildFullKey(prefix, key)

	cachedData, err := c.cache.Get([]byte(fullKey))
	if err != nil || cachedData == nil {
		return false, err
	}

	if err := json.Unmarshal(cachedData, result); err != nil {
		return false, err
	}

	return true, nil
}

func (c *cache) Evict(prefix, key string) {
	fullKey := buildFullKey(prefix, key)
	c.cache.Del([]byte(fullKey))

	keyPrefixMu.Lock()
	defer keyPrefixMu.Unlock()

	if keys, exists := cacheKeysWithPrefix[prefix]; exists {
		delete(keys, fullKey)
		if len(keys) == 0 {
			delete(cacheKeysWithPrefix, prefix)
		}
	}
}

func (c *cache) EvictByPrefix(prefix string) {
	keyPrefixMu.Lock()
	defer keyPrefixMu.Unlock()

	if keys, exists := cacheKeysWithPrefix[prefix]; exists {
		for fullKey := range keys {
			c.cache.Del([]byte(fullKey))
		}
		delete(cacheKeysWithPrefix, prefix)
	}
}
