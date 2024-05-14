package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// 1个接收者和n个发送者 通过close通知接收者通道关闭
	jobs := make(chan int)

	go receive(jobs)

	// 5个 发送者
	for i := 0; i < 5; i++ {
		send(jobs)
	}
}

func receive(jobs <-chan int) {
	// 只处理50个任务
	limit := 10
	var completeCount int

	for limit > completeCount {
		completeCount++
		value := <-jobs
		fmt.Println(value)
	}
}

func send(jobs chan<- int) {
	for {
		jobs <- rand.Intn(10)
	}
}
