package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	err := server(ctx)
	fmt.Println(err)
	fmt.Println("some...")
	time.Sleep(time.Second * 3)
	fmt.Println("wwww...")
}

func server(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	group, gctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				fmt.Printf("ctx done......\n")
				return gctx.Err()
			default:
				fmt.Printf("worker run ..\n")
				break
			}
			time.Sleep(time.Second)
		}
	})
	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func worker() {

}
