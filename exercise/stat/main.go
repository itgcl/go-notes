package main

import (
	"fmt"
)

/*
要求统计1-8000的数字中，哪些是素数？将统计素数的任务分给4个协程去完成
*/
func main() {
	goroutineReceiveCount := 4
	nums := make(chan int, 1000)
	results := make(chan int, 1000)
	exitCh := make(chan struct{}, goroutineReceiveCount)
	go putNum(nums)

	for i := 0; i < goroutineReceiveCount; i++ {
		go parseNum(nums, results, exitCh)
	}

	go func() {
		for i := 0; i < goroutineReceiveCount; i++ {
			<-exitCh
		}
		close(results)
	}()
	for {
		value, ok := <-results
		if !ok {
			break
		}
		fmt.Println(value)
	}
	fmt.Println("结束")
}

func putNum(nums chan<- int) {
	for i := 0; i < 8000; i++ {
		nums <- i
	}
	close(nums)
}

func parseNum(nums <-chan int, results chan<- int, exitCh chan<- struct{}) {
	var flag bool
	for {
		value, ok := <-nums
		if !ok {
			break
		}
		flag = true
		for i := 2; i < value; i++ {
			if value%i == 0 {
				flag = false
				break
			}
		}
		if flag {
			results <- value
		}
	}
	exitCh <- struct{}{}
}
