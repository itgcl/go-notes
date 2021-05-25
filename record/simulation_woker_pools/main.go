package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
使用goroutine和channel实现一个计算int64随机数各位数和的程序。
	开启一个goroutine循环生成int64类型的随机数，发送到jobChan
	开启24个goroutine从jobChan中取出随机数计算各位数的和，将结果发送到resultChan
	主goroutine从resultChan取出结果并打印到终端输出
 */
type Results struct {
	Sum int64
	Original int64
}
func main()  {
	jobs := make(chan int64, 11)
	results := make(chan Results, 11)
	// 启动线程池
	for i := 0; i < 24; i++ {
		go workerPools(jobs, results)
	}
	// 启动生成任务线程
	go generateInt64(jobs)

	for v := range results{
		fmt.Println(v)
	}
}

func generateInt64(jobs chan int64)  {
	for  {
		time.Sleep(time.Millisecond * 20)
		jobs <- rand.Int63()
	}
}

func workerPools(jobs <-chan int64, results chan<- Results)  {
	for {
		time.Sleep(time.Millisecond * 20)
		value := <- jobs
		res := Results{
			Sum:      0,
			Original: value,
		}
		for value > 0 {
			res.Sum += value % 10
			value /= 10
		}
		results <- res
	}
}
