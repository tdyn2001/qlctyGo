package gRedis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(address string) *RedisClient {
	redisClient := &RedisClient{}
	redisClient.rdb = redis.NewClient(&redis.Options{
		Addr: address,
	})
	return redisClient
}

func (rc *RedisClient) Subscribe(topic string, handler func([]byte)) {
	subscribeTopic := rc.rdb.Subscribe(context.Background(), topic)
	channel := subscribeTopic.Channel()
	for msg := range channel {
		handler([]byte(msg.Payload))
	}
}

func (rc *RedisClient) Publish(topic string, message []byte) {
	err := rc.rdb.Publish(context.Background(), topic, message).Err()
	if err != nil {
		panic(err)
	}
}
