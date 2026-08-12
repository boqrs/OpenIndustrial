package redis

import (
	"context"
	"encoding/json"

	"github.com/OpenIndustrial/cloud/internal/event"
	"github.com/go-redis/redis/v8"
)

// StreamPublisher 实现了 event.Publisher 接口，使用 Redis Streams 作为消息中间件。
type StreamPublisher struct {
	client *redis.Client
}

// NewStreamPublisher 创建一个新的 Redis Stream 发布者实例。
func NewStreamPublisher(client *redis.Client) event.Publisher {
	return &StreamPublisher{client: client}
}

// Publish 将事件发布到指定的 Redis Stream。
// 它将整个事件信封序列化为 JSON，并作为单个字段发布。
func (p *StreamPublisher) Publish(ctx context.Context, streamName string, env *event.Envelope) error {
	// 将整个 Envelope 序列化为 JSON 字符串
	eventJSON, err := json.Marshal(env)
	if err != nil {
		return err // 如果序列化失败，这是一个严重的内部错误
	}

	// 使用 XADD 命令将事件添加到 Stream
	args := &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"event": string(eventJSON), // 将整个事件作为单个字段的值
		},
	}

	// 执行 XADD 命令
	return p.client.XAdd(ctx, args).Err()
}