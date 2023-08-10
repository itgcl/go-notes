package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	err := Handler2(context.Background())
	fmt.Println(err)
}

// Handler 控制携程数
func Handler(ctx context.Context) error {
	// 把idList转换成channel，可以控制携程数
	var (
		idList      = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		workerCount = 2 // 携程数
		taskCh      = make(chan int, len(idList))
	)
	group, ctx := errgroup.WithContext(ctx)

	for i := 0; i < workerCount; i++ {
		group.Go(func() error {
			return worker(taskCh, ctx)
		})
	}

	group.Go(func() error {
		defer close(taskCh)
		for _, id := range idList {
			taskCh <- id
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("over")
	return nil
}

func worker(taskCh <-chan int, ctx context.Context) error {
	for id := range taskCh {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Println(id)
			time.Sleep(time.Second * 2)
		}
	}
	return nil
}

func Handler2(ctx context.Context) error {
	var (
		idList     = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		workerChan = make(chan struct{}, 2)
	)
	group, ctx := errgroup.WithContext(ctx)
	for _, id := range idList {
		workerChan <- struct{}{}
		group.Go(func() error {
			defer func() { <-workerChan }()
			fmt.Println(id)
			time.Sleep(time.Second * 2)
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("over")
	return nil
}
