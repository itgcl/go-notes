package main

import "fmt"

type sumType interface {
	int32 | float32
}

func sum[T float32 | sumType](n []T) T {
	var s T
	for _, item := range n {
		s += item
	}
	return s
}
func main() {
	data1 := []int32{10, 20, 30, 40, 50}
	data2 := []float32{10.1, 20.2, 30.3, 40.4, 50.5}
	//data3 := []string{"a", "b", "C"}
	sum1 := sum(data1)
	sum2 := sum(data2)
	//sum3 := sum(data3)
	fmt.Printf("sum1: %v (%T)\n", sum1, sum1)
	fmt.Printf("sum2: %v (%T)\n", sum2, sum2)
	//fmt.Printf("sum2: %v (%T)\n", sum3, sum3)
}
