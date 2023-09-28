package main

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	demo()
}

type Some struct {
	ID int
}

// 变量定义在循环外，在循环中每次使用变量的同一块内存存入slice
func demo() {
	var (
		some     Some
		someList []*Some
	)
	for i := 1; i <= 3; i++ {
		some.ID = i
		someList = append(someList, &some)
		fmt.Printf("%p\n", &some)
	}
	fmt.Printf("%p, %p, %p\n", someList[0], someList[1], someList[2])
	fmt.Println(someList[0], someList[1], someList[2])
}

// for range 元素始终使用同一块内存，goroutine没执行，for就执行了下一次循环，导致some的值从1->2->3
func demo2() {
	var someList = []Some{
		{1},
		{2},
		{3},
	}
	group, _ := errgroup.WithContext(context.Background())
	for _, some := range someList {
		group.Go(func() error {
			fmt.Printf("%p, %+v\n", &some, some)
			return nil
		})
		//func() {
		//	fmt.Printf("%p, %+v\n", &some, some)
		//}()
	}
	if err := group.Wait(); err != nil {
		panic(err)
	}
}
