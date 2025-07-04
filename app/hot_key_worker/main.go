package main

import (
	"hot_key/consumer"
	"hot_key/producer"
	"hot_key/redis"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	workerStarter()
}

func workerStarter() {
	redis.Init()
	go producer.Listener()
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU - 1)
	for i := 0; i < numCPU; i++ {
		go consumer.Consume()
	}

	// 创建一个通道接收信号
	sigCh := make(chan os.Signal, 1)
	// 注册要捕获的信号（Ctrl+C、kill 命令等）
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 启动一个 goroutine 等待信号
	go func() {
		sig := <-sigCh
		log.Printf("捕获到信号: %v，执行退出回调\n", sig)
		producer.Checkout()
		os.Exit(0)
	}()

	select {}
}
