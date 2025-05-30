package main

import (
	"fmt"
	"time"
)

func timer(funcToMeasure func(int) int) func(int) int {
	return func(n int) int {
		start := time.Now()
		result := funcToMeasure(n)
		duration := time.Since(start)
		fmt.Printf("timer函数执行时间: %v\n", duration)
		return result
	}
}

func myFunc(n int) int {
	return n * 2
}

func timer2(n int, funcToMeasure func(n int) int) (int, error) {
	start := time.Now()
	result := funcToMeasure(n)
	duration := time.Since(start)
	fmt.Printf("timer2函数执行时间: %v\n", duration)
	return result, nil
}

func myFunc2(n int) int {
	return n * 2
}

func timer3(funcToMeasure func() int) (int, error) {
	start := time.Now()
	result := funcToMeasure()
	duration := time.Since(start)
	fmt.Printf("timer3函数执行时间: %v\n", duration)
	return result, nil
}

func myFunc3(n int) int {
	return n * 2
}

func main() {
	timedMyFunc := timer(myFunc) // 传递函数
	result := timedMyFunc(2)     // 调用被包装的函数
	fmt.Printf("函数1 result: %d\n", result)

	result, err := timer2(3, myFunc2) // 传递闭包
	fmt.Printf("函数2 result: %v, error: %v\n", result, err)

	result3, err := timer3(func() int {
		return myFunc3(4)
	}) // 传递匿名函数
	fmt.Printf("函数3 result: %v, error: %v\n", result3, err)
}
