package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisBroker struct {
	client *RedisClient
}

func NewRedisBroker(host, port, password string) (*RedisBroker, error) {
	client, err := NewRedisClient(host, port, password)
	if err != nil {
		return nil, err
	}
	return &RedisBroker{
		client: client,
	}, err
}

func (r *RedisBroker) Next(ctx context.Context, queue string) (msg *Message, err error) {
	tk := time.After(2 * time.Second)
	for {
		select {
		case <-tk:
			// timeout
			return nil, nil
		default:

			value, err := r.client.RPopLPush(ctx, queue)
			if err != nil {
				if err == redis.Nil {
					return msg, ErrEmptyQueue
				}
			}

			err = json.Unmarshal([]byte(value), &msg)
			if err != nil {
				return msg, err
			}
			return msg, nil
		}
	}
}

func (r *RedisBroker) Send(queue string, msg *Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.client.LPush(context.Background(), queue, b)
}

func (r *RedisBroker) Ack(queue string, msg *Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.client.LRem(context.Background(), queue, 1, b)
}

func (r *RedisBroker) Cancel(queue string, Jid string) error {
	return r.client.Push(context.Background(), queue, Jid)
}

func (r *RedisBroker) GetCancel(queue string) (string, error) {
	return r.client.CancelPop(context.Background(), queue)
}
