package main

import "fmt"

type Some struct{}

func sum(n []interface{}) interface{} {
	var s float32 // 用 float32 來加總
	// 將值轉為 float32 來加總
	for _, item := range n {
		// switch 型別斷言
		switch t := item.(type) {
		case int32:
			s += float32(t)
		case float32:
			s += t
		case int:
			s += float32(t)
		}
	}
	// 檢查切片第一個元素，並用對應的型別傳回加總值
	if len(n) > 0 {
		// 型別斷言
		if _, ok := n[0].(int32); ok {
			return int32(s)
		}
	}
	return s
}

func sumAny[T int | int32 | float32 | string](t []T) T {
	var sum T
	for _, v := range t {
		sum += v
	}
	return sum
}

// TODO 返回float32类型，调用sumAnyV2[int32]有入参则是明确定义入参类型，没有入参则定义出参类型
func sumAnyV2[T int | int32 | float32](t []T) float32 {
	var sum float32
	for _, v := range t {
		sum += float32(v)
	}
	return sum
}

func sumAnyV22[T int | int32 | float32]() float32 {
	var sum float32
	return sum
}

func sumAnyV222[T int | int32 | float32 | string]() T {
	var sum T
	return sum
}

// 不能做加法运算的类型编辑器会检测报错
//func sumAnyV3[T int | int32 | Some](t []T) float32 {
//	var sum float32
//	for _, v := range t {
//		sum += float32(v)
//	}
//	return sum
//}

// T只是别名
func sumAnyV4[P int | float32 | string](t []P) P {
	var sum P
	for _, v := range t {
		sum += v
	}
	return sum
}

func sumAnyV5[T int | float32 | string](t []T, extra T) T {
	var sum T
	for _, v := range t {
		sum += v
	}
	return sum + extra
}

// 都是T只有int类型，也不能和int相加
//func sumAnyV6[T int](t []T, extra int) T {
//	var sum T
//	for _, v := range t {
//		sum += v + extra
//	}
//	return sum
//}

func main() {
	data1 := []int32{10, 20, 30, 40, 50}
	data2 := []float32{10.1, 20.2, 30.3, 40.4, 50.5}
	data3 := []string{"a", "b", "c"}
	sum1 := sumAny(data1)
	sum2 := sumAny(data2)
	sum3 := sumAny(data3)
	// 用 fmt 套件檢視傳回值的動態型別 (目前仍是 interface{})
	fmt.Printf("sum1: %v (%T)\n", sum1, sum1)
	fmt.Printf("sum2: %v (%T)\n", sum2, sum2)
	fmt.Printf("sum3: %v (%T)\n", sum3, sum3)

	sum11 := sumAnyV2[int32](data1)
	sum22 := sumAnyV2[float32](data2)
	// 用 fmt 套件檢視傳回值的動態型別 (目前仍是 interface{})
	fmt.Printf("sum11: %v (%T)\n", sum11, sum11)
	fmt.Printf("sum22: %v (%T)\n", sum22, sum22)

	sum111 := sumAnyV22[int]()
	sum222 := sumAnyV22[float32]()
	// 用 fmt 套件檢視傳回值的動態型別 (目前仍是 interface{})
	fmt.Printf("sum111: %v (%T)\n", sum111, sum111)
	fmt.Printf("sum222: %v (%T)\n", sum222, sum222)

	sum1111 := sumAnyV222[int32]()
	sum2222 := sumAnyV222[float32]()
	sum3333 := sumAnyV222[string]()
	// 用 fmt 套件檢視傳回值的動態型別 (目前仍是 interface{})
	fmt.Printf("sum1111: %v (%T)\n", sum1111, sum1111)
	fmt.Printf("sum2222: %v (%T)\n", sum2222, sum2222)
	fmt.Printf("sum3333: %v (%T)\n", sum3333, sum3333)
}
