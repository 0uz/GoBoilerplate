package cache

import (
	"context"
	"time"
)

type LocalCacheService interface {
	Set(prefix, key string, ttl time.Duration, value any)
	Get(prefix, key string, result any) (bool, error)
	Evict(prefix, key string)
	EvictByPrefix(prefix string)
}

type RedisCacheService interface {
	Set(ctx context.Context, prefix, key string, ttl time.Duration, value any) error
	Get(ctx context.Context, prefix, key string, result any) (bool, error)
	Exists(ctx context.Context, prefix, key string) (bool, error)
	Evict(ctx context.Context, prefix, key string) error
	EvictByPrefix(ctx context.Context, prefix string) error
	SAdd(ctx context.Context, prefix, key string, ttl time.Duration, member string) error
	SMembers(ctx context.Context, prefix, key string) ([]string, error)
	SCard(ctx context.Context, prefix, key string) (int64, error)
	Scan(ctx context.Context, pattern string) ([]string, error)
	CloseRedisClient() error
}
