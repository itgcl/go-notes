package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// n个接收者和一个发送者 通过close通知接收者通道关闭

var wg sync.WaitGroup

func main() {
	jobs := make(chan int, 5)
	// 开启3个接收者执行任务
	for i := 0; i < 3; i++ {
		go receive(jobs)
	}

	// 1个发送者 发送十个任务
	for i := 0; i < 10; i++ {
		wg.Add(1)
		jobs <- rand.Intn(10)
	}
	close(jobs)
	wg.Wait()

}

func receive(jobs <-chan int) {
	for {
		value, ok := <-jobs
		if !ok {
			break
		}
		fmt.Println(value)
		wg.Done()
		// DB insert update delete ...
	}
}
