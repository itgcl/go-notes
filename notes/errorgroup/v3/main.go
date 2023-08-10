package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	err := Handler(context.Background())
	log.Fatal(err)
	group, _ := errgroup.WithContext(context.Background())
	group.Go(func() error {
		fmt.Println("aaww")
		group.Go(func() error {
			time.Sleep(time.Second * 3)
			fmt.Println("write")
			return nil
		})
		//time.Sleep(time.Second)
		fmt.Println("defer")
		return nil
	})
	if err := group.Wait(); err != nil {
		panic(err)
	}
	fmt.Println("over")
}

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

	// Enqueue tasks
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
