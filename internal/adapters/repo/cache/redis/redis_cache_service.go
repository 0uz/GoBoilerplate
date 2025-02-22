package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

type RedisCacheService interface {
	Set(ctx context.Context, prefix, key string, ttl time.Duration, value interface{}) error
	Get(ctx context.Context, prefix, key string, result interface{}) (bool, error)
	Exists(ctx context.Context, prefix, key string) (bool, error)
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
	return fmt.Sprintf("redis://%s:%s", config.Get().Valkey.Host, config.Get().Valkey.Port)
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
		return errors.GenericError("error marshaling value", err)
	}

	err = r.client.Set(ctx, fullKey, jsonData, ttl).Err()
	if err != nil {
		return errors.GenericError("error setting value to redis", err)
	}

	return nil
}

func (r *redisCache) Get(ctx context.Context, prefix, key string, result any) (bool, error) {
	fullKey := buildRedisFullKey(prefix, key)
	
	cachedData, err := r.client.Get(ctx, fullKey).Bytes()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, errors.GenericError("error getting value from redis", err)
	}

	if err := json.Unmarshal(cachedData, result); err != nil {
		return false, errors.GenericError("error unmarshaling value", err)
	}

	return true, nil
}

func (r *redisCache) Exists(ctx context.Context, prefix, key string) (bool, error) {
	fullKey := buildRedisFullKey(prefix, key)

	exists := r.client.Exists(ctx, fullKey).Val()
	if exists == 0 {
		return false, nil
	}
	return true, nil
}

func (r *redisCache) Evict(ctx context.Context, prefix, key string) error {
	fullKey := buildRedisFullKey(prefix, key)
	err := r.client.Del(ctx, fullKey).Err()
	if err != nil {
		return errors.GenericError("error deleting key from redis", err)
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
		return errors.GenericError("error scanning keys from redis", err)
	}

	if len(keys) > 0 {
		err := r.client.Del(ctx, keys...).Err()
		if err != nil {
			return errors.GenericError("error deleting keys from redis", err)
		}
	}
	return nil
}
