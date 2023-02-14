package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	group, _ := errgroup.WithContext(context.Background())
	ch := make(chan int, 100)
	// 发送者
	group.Go(func() error {
		defer close(ch)
		for i := 1; i <= 100; i++ {
			ch <- i
		}
		return nil
	})
	// 消费者
	// 开启五个worker
	for i := 0; i < 5; i++ {
		i := i
		group.Go(func() error {
			for v := range ch {
				// 模拟业务处理，500毫秒处理一个
				fmt.Printf("goroutine%d customer, data: %d \n", i, v)
				time.Sleep(time.Millisecond * 500)
			}
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}
