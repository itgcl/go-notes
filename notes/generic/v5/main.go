package main

import "fmt"

type valueType interface {
	int32 | float32
}

func channelGen[T valueType]() chan T {
	ch := make(chan T)
	return ch
}

func sum[T valueType](ch chan T, values []T) {
	var result T
	for _, v := range values {
		result += v
	}
	ch <- result
}

func main() {
	data1 := []int32{10, 20, 30, 40, 50}
	data2 := []float32{10.1, 20.2, 30.3, 40.4, 50.5}
	ch1 := channelGen[int32]()
	ch2 := channelGen[float32]()
	go sum(ch1, data1) // 讓函式推論泛型型別
	go sum(ch2, data2)
	fmt.Println("sum1:", <-ch1)
	fmt.Println("sum2:", <-ch2)
}
