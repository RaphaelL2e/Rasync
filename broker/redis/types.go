package redis

import "errors"

type Message struct {
	Id         string `json:"id"`
	WorkerName string `json:"worker_name"`
	Data       interface{}
}

var ErrEmptyQueue = errors.New("empty queue")
