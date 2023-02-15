package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Data struct {
	AID int
	BID int
	CID int
}

func main() {
	example()
}

func example() {
	group, ctx := errgroup.WithContext(context.Background())
	ch := make(chan int, 1)
	someBCh := make(chan *Data, 1)
	someCCh := make(chan *Data, 1)
	wg := &sync.WaitGroup{}
	// 获取aid列表
	group.Go(func() error {
		return ListAID(ctx, ch)
	})
	// 通过aid获取bid，可能是查询数据库，为了提高性能，按批次查询
	group.Go(func() error {
		defer close(someBCh)
		var (
			data      []int
			threshold = 10
		)
		for v := range ch {
			data = append(data, v)
			// 达到阈值，异步获取数据
			if len(data) == threshold {
				wg.Add(1)
				var sendData = make([]int, threshold)
				copy(sendData, data)
				data = data[:0]
				// 异步获取cid数据
				group.Go(func() error {
					return ListBID(ctx, sendData, someBCh, wg)
				})
			}
		}
		if len(data) > 0 {
			group.Go(func() error {
				wg.Add(1)
				return ListBID(ctx, data, someBCh, wg)
			})
		}
		wg.Wait()
		return nil
	})
	// 通过bid获取cid
	group.Go(func() error {
		defer close(someCCh)
		return listCID(ctx, someBCh, someCCh)
	})
	// 入库
	group.Go(func() error {
		var count int
		for v := range someCCh {
			count++
			fmt.Println(v)
		}
		fmt.Println("count:", count)
		return nil
	})
	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}

// ListAID 获取aid列表
func ListAID(ctx context.Context, ch chan<- int) error {
	defer close(ch)
	for i := 1; i <= 100; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- i:
		}
	}
	return nil
}

func ListBID(ctx context.Context, intList []int, ach chan<- *Data, wg *sync.WaitGroup) error {
	defer wg.Done()
	for _, i := range intList {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ach <- &Data{AID: i, BID: i * 2, CID: 0}:
		}
	}
	return nil
}

func listCID(ctx context.Context, ach <-chan *Data, bch chan<- *Data) error {
	for v := range ach {
		v.CID = v.BID * 2
		select {
		case <-ctx.Done():
			return ctx.Err()
		case bch <- v:
		}
	}
	return nil
}

func example1() {
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
