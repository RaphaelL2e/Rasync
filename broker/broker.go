package broker

import "github.com/RaphaelL2e/Rasync/broker/redis"

type Broker struct {
	*redis.RedisBroker
}

func NewBroker(host, port, password string) (*Broker, error) {
	broker, err := redis.NewRedisBroker(host, port, password)
	if err != nil {
		return nil, err
	}
	return &Broker{
		RedisBroker: broker,
	}, nil
}
