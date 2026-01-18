package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ouz/goboilerplate/pkg/cache"
	"github.com/ouz/goboilerplate/pkg/errors"
	"github.com/ouz/goboilerplate/pkg/log"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// redisCacheService implements the RedisCacheService interface from the cache package
type redisCacheService struct {
	client *redis.Client
}

// NewRedisCacheService creates a new instance of RedisCacheService
func NewRedisCacheService(client *redis.Client) cache.RedisCacheService {
	return &redisCacheService{
		client: client,
	}
}

// ConnectRedis establishes a connection to Redis
func ConnectRedis(logger *log.Logger, host, port string, monitoringEnabled bool) (*redis.Client, error) {
	opt := &redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	}

	client := redis.NewClient(opt)

	if monitoringEnabled {
		commandFilter := func(cmd redis.Cmder) bool {
			switch cmd.Name() {
			case "xreadgroup", "xadd", "xack", "xgroup", "xinfo", "xlen", "xrange", "xpending":
				return true // true = skip tracing for this command
			default:
				return false // false = trace this command
			}
		}

		if err := redisotel.InstrumentTracing(client, redisotel.WithCommandFilter(commandFilter)); err != nil {
			logger.Warn("Failed to instrument Redis tracing", "error", err)
		}
		if err := redisotel.InstrumentMetrics(client); err != nil {
			logger.Warn("Failed to instrument Redis metrics", "error", err)
		}
	}

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// CloseRedisClient closes the Redis client connection
func CloseRedisClient(client *redis.Client) error {
	return client.Close()
}

func buildRedisFullKey(prefix, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

// Set stores a value in Redis with a specified prefix, key, and TTL
func (r *redisCacheService) Set(ctx context.Context, prefix, key string, ttl time.Duration, value interface{}) error {
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

// Get retrieves a value from Redis by prefix and key
func (r *redisCacheService) Get(ctx context.Context, prefix, key string, result any) (bool, error) {
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

// Exists checks if a key exists in the cache
func (r *redisCacheService) Exists(ctx context.Context, prefix, key string) (bool, error) {
	fullKey := buildRedisFullKey(prefix, key)

	exists := r.client.Exists(ctx, fullKey).Val()
	if exists == 0 {
		return false, nil
	}
	return true, nil
}

// Evict removes a specific key from the cache
func (r *redisCacheService) Evict(ctx context.Context, prefix, key string) error {
	fullKey := buildRedisFullKey(prefix, key)
	err := r.client.Del(ctx, fullKey).Err()
	if err != nil {
		return errors.GenericError("error deleting key from redis", err)
	}
	return nil
}

// EvictByPrefix removes all keys with the given prefix
func (r *redisCacheService) EvictByPrefix(ctx context.Context, prefix string) error {
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

// SAdd adds a member to a Redis set with TTL
func (r *redisCacheService) SAdd(ctx context.Context, prefix, key string, ttl time.Duration, member string) error {
	fullKey := buildRedisFullKey(prefix, key)

	// Add member to set
	err := r.client.SAdd(ctx, fullKey, member).Err()
	if err != nil {
		return errors.GenericError("error adding member to redis set", err)
	}

	// Set expiration if TTL is provided
	if ttl > 0 {
		err = r.client.Expire(ctx, fullKey, ttl).Err()
		if err != nil {
			return errors.GenericError("error setting TTL on redis set", err)
		}
	}

	return nil
}

// SMembers returns all members of a Redis set
func (r *redisCacheService) SMembers(ctx context.Context, prefix, key string) ([]string, error) {
	fullKey := buildRedisFullKey(prefix, key)

	members, err := r.client.SMembers(ctx, fullKey).Result()
	if err == redis.Nil {
		return []string{}, nil
	} else if err != nil {
		return nil, errors.GenericError("error getting members from redis set", err)
	}

	return members, nil
}

// SCard returns the number of members in a Redis set
func (r *redisCacheService) SCard(ctx context.Context, prefix, key string) (int64, error) {
	fullKey := buildRedisFullKey(prefix, key)

	count, err := r.client.SCard(ctx, fullKey).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, errors.GenericError("error getting set cardinality from redis", err)
	}

	return count, nil
}

// Scan returns all keys matching a pattern
func (r *redisCacheService) Scan(ctx context.Context, pattern string) ([]string, error) {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, errors.GenericError("error scanning keys from redis", err)
	}

	return keys, nil
}

// CloseRedisClient closes the Redis client connection
func (r *redisCacheService) CloseRedisClient() error {
	return r.client.Close()
}
