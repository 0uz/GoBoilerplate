package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ouz/goauthboilerplate/internal/config"
)

type RedisCacheService interface {
	Set(ctx context.Context, prefix, key string, ttl time.Duration, value interface{}) error
	Get(ctx context.Context, prefix, key string, result interface{}) (error, bool)
	Evict(ctx context.Context, prefix, key string) error
	EvictByPrefix(ctx context.Context, prefix string) error
	CloseRedisClient() error
}

type redisCache struct {
	client *redis.Client
}

func ConnectRedis() (*redis.Client, error) {
	opt, err := redis.ParseURL(prepareURL())
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CloseRedisClient(client *redis.Client) error {
	return client.Close()
}

func NewRedisCacheService(client *redis.Client) RedisCacheService {
	return &redisCache{
		client: client,
	}
}

func prepareURL() string {
	return fmt.Sprintf("redis://%s:%s", config.Get().Redis.Host, config.Get().Redis.Port)
}

func buildRedisFullKey(prefix, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

func (r *redisCache) CloseRedisClient() error {
	return r.client.Close()
}

func (r *redisCache) Set(ctx context.Context, prefix, key string, ttl time.Duration, value interface{}) error {
	fullKey := buildRedisFullKey(prefix, key)

	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling value: %w", err)
	}

	err = r.client.Set(ctx, fullKey, jsonData, ttl).Err()
	if err != nil {
		return fmt.Errorf("error setting value to redis: %w", err)
	}

	return nil
}

func (r *redisCache) Get(ctx context.Context, prefix, key string, result interface{}) (error, bool) {
	fullKey := buildRedisFullKey(prefix, key)

	cachedData, err := r.client.Get(ctx, fullKey).Bytes()
	if err == redis.Nil {
		return nil, false
	} else if err != nil {
		return fmt.Errorf("error getting value from redis: %w", err), false
	}

	if err := json.Unmarshal(cachedData, result); err != nil {
		return fmt.Errorf("error unmarshaling value: %w", err), false
	}

	return nil, true
}

func (r *redisCache) Evict(ctx context.Context, prefix, key string) error {
	fullKey := buildRedisFullKey(prefix, key)
	err := r.client.Del(ctx, fullKey).Err()
	if err != nil {
		return fmt.Errorf("error deleting key from redis: %w", err)
	}
	return nil
}

func (r *redisCache) EvictByPrefix(ctx context.Context, prefix string) error {
	pattern := fmt.Sprintf("%s:*", prefix)
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("error scanning keys from redis: %w", err)
	}

	if len(keys) > 0 {
		err := r.client.Del(ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("error deleting keys from redis: %w", err)
		}
	}
	return nil
}
