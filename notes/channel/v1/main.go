package main

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"
)

func main() {
	example()
}

func example() {
	group, ctx := errgroup.WithContext(context.Background())
	ch := make(chan int, 1)
	// 生产
	group.Go(func() error {
		defer close(ch)
		for i := 1; i <= 100; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case ch <- i:
			}
		}
		return nil
	})
	// 消费
	group.Go(func() error {
		for v := range ch {
			fmt.Println(v)
		}
		return nil
	})
	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}
