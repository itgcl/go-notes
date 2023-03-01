package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	example1()
}

func example() {
	rand.Seed(time.Now().UnixNano())
	group, _ := errgroup.WithContext(context.Background())
	ch := make(chan int, 1)
	stopCh := make(chan int, 1)
	const NumSenders = 10
	const Max = 1000
	// 多发送者
	for i := 0; i < NumSenders; i++ {
		group.Go(func() error {
			for {
				select {
				// 监听退出信号
				case <-stopCh:
					fmt.Println("send stop...")
					return nil
				case ch <- rand.Intn(Max):
				}
			}
		})
	}

	group.Go(func() error {
		for v := range ch {
			// 达到某个条件触发停止
			if v == Max-1 {
				close(stopCh)
				return nil
			}
			fmt.Println(v)
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}

func example1() {
	ch := make(chan int, 1)
	ch2 := make(chan int, 1)
	group, _ := errgroup.WithContext(context.Background())
	const NumSenders = 10
	const Max = 10000

	group.Go(func() error {
		defer close(ch)
		for i := 0; i < Max; i++ {
			ch <- i
		}
		return nil
	})
	// 发送完关闭ch2，使用waitGroup
	group.Go(func() error {
		var wg sync.WaitGroup
		defer close(ch2)
		for i := 0; i < NumSenders; i++ {
			wg.Add(1)
			group.Go(func() error {
				defer wg.Done()
				for v := range ch {
					ch2 <- v
				}
				return nil
			})
		}
		wg.Wait()
		return nil
	})

	group.Go(func() error {
		for v := range ch2 {
			fmt.Println(v)
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}
