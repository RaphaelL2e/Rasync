package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	rp "runtime/pprof"
	"time"

	"github.com/RaphaelL2e/Rasync/broker"
	"github.com/RaphaelL2e/Rasync/worker"
)

func main() {
	// pprof
	// 启动 HTTP 服务器，监听地址为 localhost:6060
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// cpu profiling
	f, _ := os.Create("cpu.prof")
	err := rp.StartCPUProfile(f)
	if err != nil {
		panic(err)
	}
	defer rp.StopCPUProfile()

	// memory profiling
	f1, _ := os.Create("mem.prof")
	rp.WriteHeapProfile(f1)
	defer f1.Close()

	// main
	broker1, err := broker.NewBroker("localhost", "6379", "")
	if err != nil {
		panic(err)
	}

	// worker run
	w := worker.NewWorker("test", broker1)
	w.SetHandleMessage(handleMessage)
	err = w.Run("worker1")
	if err != nil {
		panic(err)
	}
}

func handleMessage(cancel worker.Cancel, m interface{}) error {
	// do something
	sec := rand.Intn(100)
	fmt.Println(fmt.Sprintf("value: %v sleep %d s", m, sec))
	time.Sleep(time.Duration(sec) * time.Second)
	if cancel.IsCancel() {
		fmt.Println("cancel")
		return fmt.Errorf("cancel")
	}
	sec = rand.Intn(100)
	fmt.Println(fmt.Sprintf("value: %v sleep %d s", m, sec))
	time.Sleep(time.Duration(sec) * time.Second)
	if cancel.IsCancel() {
		fmt.Println("cancel")
		return fmt.Errorf("cancel")
	}
	fmt.Println(fmt.Sprintf("value: %v sleep %d s done", m, sec))
	return nil
}
