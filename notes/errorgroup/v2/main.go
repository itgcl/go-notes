package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	worker(context.Background())
	// f, err := os.Create("/tmp/ww.zip")
	// fmt.Println(f, err)
}

func worker(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error {
		for i := 1; i <= 10; i++ {
			f, err := os.Create("/tmp/ww.txt")
			if err != nil {
				return err
			}

			// defer func() {
			defer func() {
				fmt.Println("defer")
				f.Close()
			}()
			//}()
			group.Go(func() error {
				time.Sleep(time.Second)
				fmt.Println("write")
				f.WriteString("a")
				// fmt.Println("go start", i)
				// time.Sleep(time.Second * 2)
				return nil
			})
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}
