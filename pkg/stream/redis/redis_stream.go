package redis

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/ouz/goauthboilerplate/pkg/errors"
	"github.com/ouz/goauthboilerplate/pkg/log"
	"github.com/ouz/goauthboilerplate/pkg/stream"
	"github.com/redis/go-redis/v9"
)

type redisStreamService struct {
	client *redis.Client
	logger *log.Logger
}

func NewRedisStreamService(logger *log.Logger, client *redis.Client) stream.StreamService {
	return &redisStreamService{
		client: client,
		logger: logger,
	}
}

func (r *redisStreamService) Publish(ctx context.Context, streamKey string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return errors.GenericError("failed to marshal event", err)
	}

	err = r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		Values: map[string]interface{}{
			"data": data,
		},
	}).Err()

	if err != nil {
		return errors.GenericError("failed to publish to stream", err)
	}

	return nil
}

func (r *redisStreamService) CreateGroup(ctx context.Context, streamKey, group string) error {
	err := r.client.XGroupCreateMkStream(ctx, streamKey, group, "0").Err()
	if err != nil {
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			return nil // Group already exists, ignore
		}
		return errors.GenericError("failed to create consumer group", err)
	}
	return nil
}

func (r *redisStreamService) Consume(ctx context.Context, streamKey, group, consumer string, handler stream.HandlerFunc) error {
	if err := r.CreateGroup(ctx, streamKey, group); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Read from consumer group
			streams, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    group,
				Consumer: consumer,
				Streams:  []string{streamKey, ">"}, // ">" means new messages
				Count:    10,                       // Batch size
				Block:    2 * time.Second,          // Block for 2 seconds if no messages
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue // No messages, retry
				}
				if strings.Contains(err.Error(), "client is closed") {
					r.logger.Info("Redis client closed, stopping consumer")
					return nil
				}
				r.logger.Error("Error reading from stream", "error", err)
				time.Sleep(1 * time.Second)
				continue
			}

			for _, xStream := range streams {
				for _, msg := range xStream.Messages {
					dataStr, ok := msg.Values["data"].(string)
					if !ok {
						r.logger.Error("Invalid message format", "msg_id", msg.ID)
						r.Ack(ctx, streamKey, group, msg.ID)
						continue
					}

					if err := handler(ctx, msg.ID, []byte(dataStr)); err != nil {
						r.logger.Error("Failed to process message", "msg_id", msg.ID, "error", err)
					} else {
						if err := r.Ack(ctx, streamKey, group, msg.ID); err != nil {
							r.logger.Error("Failed to ack message", "msg_id", msg.ID, "error", err)
						}
					}
				}
			}
		}
	}
}

func (r *redisStreamService) Ack(ctx context.Context, streamKey, group string, ids ...string) error {
	return r.client.XAck(ctx, streamKey, group, ids...).Err()
}
