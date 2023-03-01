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
	example()
}

func example() {
	rand.Seed(time.Now().UnixNano()) // needed before Go 1.20
	const Max = 100000
	const NumReceivers = 100
	const NumThirdParties = 15

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	dataCh := make(chan int)
	closing := make(chan struct{}) // signal channel
	closed := make(chan struct{})

	group, _ := errgroup.WithContext(context.Background())

	stop := func() {
		select {
		// closing通道会输入多次（NumThirdParties），输出一次（sender携程），
		// 因为多次输入，1次输出，NumThirdParties携程（发送方方）会一直阻塞，导致携程泄漏。
		// 这里加case <-closed, 当 case closing <- struct{}{}阻塞，会执行<-closed（sender会接收退出信号，close closed channel）
		case closing <- struct{}{}:
			fmt.Println("stop")
			//<-closed
		case <-closed:
			fmt.Println("closed")
		}
	}

	for i := 0; i < NumThirdParties; i++ {
		group.Go(func() error {
			r := 1 + rand.Intn(3)
			time.Sleep(time.Duration(r) * time.Second)
			fmt.Println("发出退出信号")
			// 发出退出信号，多个携程会发送多次。
			stop()
			return nil
		})
	}

	// the sender
	group.Go(func() error {
		defer func() {
			fmt.Println("close了")
			close(closed)
			close(dataCh)
		}()
		for {
			select {
			case <-closing:
				fmt.Println("接收到了，要退出")
				return nil
			case dataCh <- rand.Intn(Max):
			}
		}
	})

	// receivers
	for i := 0; i < NumReceivers; i++ {
		group.Go(func() error {
			defer wgReceivers.Done()

			for value := range dataCh {
				log.Println(value)
			}
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}
