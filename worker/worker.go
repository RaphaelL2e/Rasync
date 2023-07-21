package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/RaphaelL2e/Rasync/broker"
	"github.com/RaphaelL2e/Rasync/broker/redis"
)

type Worker struct {
	groupName string

	broker  *broker.Broker
	boolNum int

	handleMessage func(cancel Cancel, m interface{}) error

	cancel map[string]context.CancelFunc
}

func NewWorker(groupName string, b *broker.Broker) *Worker {
	return &Worker{
		groupName: groupName,
		broker:    b,
		boolNum:   5,
		cancel:    make(map[string]context.CancelFunc),
	}
}

func (w *Worker) SetBoolNum(boolNum int) {
	w.boolNum = boolNum
}

func (w *Worker) SetHandleMessage(handleMessage func(ctx Cancel, m interface{}) error) {
	w.handleMessage = handleMessage
}

func (w *Worker) Run(worker string) error {
	c := make(chan int, w.boolNum)
	go func() {
		w.GetStopChan(worker)
	}()
	for {
		c <- 1
		go func() {
			m, err := w.broker.Next(context.Background(), w.queueName(worker))
			if err != nil {
				if err == redis.ErrEmptyQueue {
					<-c
					return
				}
				fmt.Println("next failed:", err)
				<-c
				return
			}
			if m == nil {
				<-c
				return
			}
			ctx, cancel := context.WithCancel(context.Background())
			w.cancel[m.Id] = cancel
			fmt.Println("task start:", m.Id)
			err = w.handleMessage(Cancel{ctx: ctx}, m.Data)
			if err != nil {
				fmt.Println("handle message failed:", err)
			}
			err = w.broker.Ack(w.queueName(worker), m)
			if err != nil {
				fmt.Println("ack failed:", err)
			}
			<-c
			return
		}()
	}
}

func (w *Worker) queueName(worker string) string {
	return fmt.Sprintf("{%s}_%s", w.groupName, worker)
}

func (w *Worker) GetStopChan(worker string) {
	tk := time.Tick(time.Second * 5)
	for {
		select {
		case <-tk:
			jid, err := w.broker.GetCancel(w.queueName(worker))
			if err != nil {
				continue
			}
			fmt.Println("cancel job:", jid)
			err = w.Stop(jid)
			if err != nil {
				fmt.Println("stop failed:", err)
			}
		}
	}

}

func (w *Worker) Stop(jobId string) error {
	cancel, ok := w.cancel[jobId]
	if !ok {
		return fmt.Errorf("jobId %s not found", jobId)
	}
	cancel()
	return nil
}

func (w *Worker) StopAll() {
	for _, cancel := range w.cancel {
		cancel()
	}
}

func (w *Worker) IsCancel(jid string) bool {
	if w.cancel[jid] == nil {
		return false
	}
	return true
}
