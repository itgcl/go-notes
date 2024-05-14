package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	group, _ := errgroup.WithContext(context.Background())
	group.Go(func() error {
		fmt.Println("aaww")
		group.Go(func() error {
			time.Sleep(time.Second * 3)
			fmt.Println("write")
			return nil
		})
		// time.Sleep(time.Second)
		fmt.Println("defer")
		return nil
	})
	if err := group.Wait(); err != nil {
		panic(err)
	}
	fmt.Println("over")
}
