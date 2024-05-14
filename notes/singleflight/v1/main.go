package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
)

type Result string

func find(ctx context.Context, query string) (Result, error) {
	return Result(fmt.Sprintf("result for %q", query)), nil
}

func main() {
	var g singleflight.Group
	const n = 5
	// waited := int32(n)
	// done := make(chan struct{})
	key := "https://weibo.com/1227368500/H3GIgngon"
	//for i := 0; i < n; i++ {
	//	go func(j int) {
	//		v, _, shared := g.Do(key, func() (interface{}, error) {
	//			ret, err := find(context.Background(), key)
	//			return ret, err
	//		})
	//		//if atomic.AddInt32(&waited, -1) == 0 {
	//		//	close(done)
	//		//}
	//		fmt.Printf("index: %d, val: %v, shared: %v\n", j, v, shared)
	//	}(i)
	//}
	//
	//select {
	//case <-done:
	//case <-time.After(time.Second):
	//	fmt.Println("Do hangs")
	//}
	//time.Sleep(time.Second * 3)
	//}

	for i := 0; i < n; i++ {
		ch := g.DoChan(key, func() (interface{}, error) {
			ret, err := find(context.Background(), key)
			return ret, err
		})
		// Create our timeout
		timeout := time.After(500 * time.Millisecond)

		var ret singleflight.Result
		select {
		case <-timeout: // Timeout elapsed
			fmt.Println("Timeout")
			return
		case ret = <-ch: // Received result from channel
			fmt.Printf("index: %d, val: %v, shared: %v\n", i, ret.Val, ret.Shared)
		}
	}
}
