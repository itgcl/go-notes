package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	example()
}

func example() {
	rand.Seed(time.Now().UnixNano())
	const Max = 10000
	const NumReceivers = 10
	const NumSenders = 10

	dataCh := make(chan int)
	stopCh := make(chan struct{})
	toStop := make(chan string)
	group, _ := errgroup.WithContext(context.Background())
	// 裁判员携程，接收退出信号，关闭stopCh，发送者和接收者在发送或接收数据时，检查stopCh是否关闭，关闭的话则退出
	group.Go(func() error {
		v := <-toStop
		fmt.Println(v)
		close(stopCh)
		return nil
	})
	// 发送者
	for i := 0; i < NumSenders; i++ {
		i := i
		group.Go(func() error {
			for {
				v := rand.Intn(Max)
				// 模拟触发停止
				if v == 0 {
					// 为什么这里使用select？防止toStop通道阻塞
					select {
					case toStop <- fmt.Sprintf("sender#%d", i):
						fmt.Println("sender over")
					default:
						fmt.Println("sender default")
					}
					return nil
				}
				// 发送时，检查通道是否关闭，
				select {
				case <-stopCh:
					fmt.Println("sender quit...")
					return nil
				case dataCh <- v:
				}
			}
		})
	}

	// 接收者
	for i := 0; i < NumReceivers; i++ {
		//i := i
		group.Go(func() error {
			for {
				select {
				case <-stopCh:
					fmt.Println("receiver quit...")
					return nil
				case v := <-dataCh:
					if v == Max-1 {
						select {
						case toStop <- fmt.Sprintf("receiver #%d", i):
							fmt.Println("receiver over")
						default:
							fmt.Println("receiver default")
						}
						return nil
					}
					fmt.Println(v)
				}
			}
		})
	}
	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}
