package main

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	group, _ := errgroup.WithContext(context.Background())
	ch := make(chan int, 1)
	// 生产者 TODO

	// 消费者
	group.Go(func() error {
		for v := range ch {
			fmt.Println(v)
		}
		return nil
	})
}
