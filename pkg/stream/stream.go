package stream

import "context"

type HandlerFunc func(ctx context.Context, msgID string, payload []byte) error

type StreamService interface {
	Publish(ctx context.Context, stream string, event interface{}) error
	Consume(ctx context.Context, stream, group, consumer string, handler HandlerFunc) error
	CreateGroup(ctx context.Context, stream, group string) error
	Ack(ctx context.Context, stream, group string, ids ...string) error
}
