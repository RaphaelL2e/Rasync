package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	cli         *redis.Client
	WaitQueue   string
	WorkQueue   string
	CancelQueue string
}

func NewRedisClient(host, port, password string) (*RedisClient, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RedisClient{
		cli:         cli,
		WaitQueue:   "_wait_queue",
		WorkQueue:   "_work_queue",
		CancelQueue: "_cancel_queue",
	}, nil
}

func (r *RedisClient) LPush(ctx context.Context, prefix string, values ...interface{}) error {
	return r.cli.LPush(ctx, prefix+r.WaitQueue, values).Err()
}

func (r *RedisClient) RPopLPush(ctx context.Context, prefix string) (string, error) {
	return r.cli.RPopLPush(ctx, prefix+r.WaitQueue, prefix+r.WorkQueue).Result()
}

func (r *RedisClient) LRem(ctx context.Context, prefix string, count int64, value interface{}) error {
	return r.cli.LRem(ctx, prefix+r.WorkQueue, count, value).Err()
}

func (r *RedisClient) Push(ctx context.Context, prefix string, values ...interface{}) error {
	return r.cli.LPush(ctx, prefix+r.CancelQueue, values).Err()
}

func (r *RedisClient) CancelPop(ctx context.Context, prefix string) (string, error) {
	return r.cli.RPop(ctx, prefix+r.CancelQueue).Result()
}
