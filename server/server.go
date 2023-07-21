package server

import (
	"fmt"

	"github.com/RaphaelL2e/Rasync/broker"
	"github.com/RaphaelL2e/Rasync/broker/redis"
	"github.com/google/uuid"
)

type Server struct {
	groupName string

	broker *broker.Broker
}

func NewServer(groupName string, b *broker.Broker) *Server {
	return &Server{
		groupName: groupName,
		broker:    b,
	}
}

func (s *Server) Send(worker string, data interface{}) (JobId string, err error) {
	id := uuid.New().String()
	return id, s.broker.Send(s.queueName(worker),
		&redis.Message{
			Id:         id,
			WorkerName: worker,
			Data:       data,
		})
}

func (s *Server) queueName(worker string) string {
	return fmt.Sprintf("{%s}_%s", s.groupName, worker)
}

func (s *Server) Cancel(worker, jobId string) error {
	return s.broker.Cancel(s.queueName(worker), jobId)
}
