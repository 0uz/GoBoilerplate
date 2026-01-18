package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/coocood/freecache"
	"github.com/ouz/goboilerplate/pkg/log"
)

var (
	cacheKeysWithPrefix = make(map[string]map[string]struct{})
	keyPrefixMu         sync.Mutex
)

type localCache struct {
	cache  *freecache.Cache
	logger *log.Logger
}

func NewLocalCacheService(logger *log.Logger, sizeMB int) LocalCacheService {
	return &localCache{
		cache:  freecache.NewCache(sizeMB * 1024 * 1024),
		logger: logger,
	}
}

func buildFullKey(prefix, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

func (c *localCache) Set(prefix, key string, ttl time.Duration, value interface{}) {
	fullKey := buildFullKey(prefix, key)

	jsonData, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("Failed to marshal cache value", "error", err, "prefix", prefix, "key", key)
		return
	}

	if err := c.cache.Set([]byte(fullKey), jsonData, int(ttl.Seconds())); err != nil {
		c.logger.Error("Failed to set cache value", "error", err, "prefix", prefix, "key", key)
	}
	keyPrefixMu.Lock()
	defer keyPrefixMu.Unlock()

	if _, exists := cacheKeysWithPrefix[prefix]; !exists {
		cacheKeysWithPrefix[prefix] = make(map[string]struct{})
	}

	if _, exists := cacheKeysWithPrefix[prefix][fullKey]; !exists {
		cacheKeysWithPrefix[prefix][fullKey] = struct{}{}
	}
}

func (c *localCache) Get(prefix, key string, result interface{}) (bool, error) {
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

func (c *localCache) Evict(prefix, key string) {
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

func (c *localCache) EvictByPrefix(prefix string) {
	keyPrefixMu.Lock()
	defer keyPrefixMu.Unlock()

	if keys, exists := cacheKeysWithPrefix[prefix]; exists {
		for fullKey := range keys {
			c.cache.Del([]byte(fullKey))
		}
		delete(cacheKeysWithPrefix, prefix)
	}
}
