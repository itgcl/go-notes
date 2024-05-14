package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	example3()
}

func example() {
	// 作用域，example方法无法对输入数据本身做改变，因为不是传址
	var v int64 = 10
	v1 := add(v)
	fmt.Printf("原数据：%d, add:%d\n", v, v1)
}

func add(v int64) int64 {
	return atomic.AddInt64(&v, 2)
}

func example1() {
	type Data struct {
		a, b, c int
	}
	var (
		data = Data{1, 2, 3}
		v    atomic.Value
	)
	// 设置值
	v.Store(data)
	// 读取值
	d1 := v.Load().(Data)
	fmt.Println(d1)
}

func example2() {
	var v int64 = 10
	// 对数据操作
	// 加
	v1 := atomic.AddInt64(&v, 2)
	fmt.Printf("原数据：%d, 相加后：%d\n", v, v1)
	// 减
	v2 := atomic.AddInt64(&v, -3)
	fmt.Printf("原数据：%d, 相减后：%d\n", v, v2)
	// 读取值
	fmt.Printf("load: %d\n", atomic.LoadInt64(&v))
	// 设置值
	atomic.StoreInt64(&v, 33)
	fmt.Printf("store: %d\n", v)
	// 替换新值，返回旧值
	oldValue := atomic.SwapInt64(&v, 111)
	fmt.Printf("swap new: %d, old: %d\n", v, oldValue)
	// 比较并交换 v的值是111，则替换成22，否则不变
	ok := atomic.CompareAndSwapInt64(&v, 111, 22)
	fmt.Printf("CompareAndSwapInt64：%v，old value: %d,  current value: %d, swap value: %d\n", ok, 111, v, 22)
	ok = atomic.CompareAndSwapInt64(&v, 23, 11)
	fmt.Printf("CompareAndSwapInt64：%v, value: %d, swap value: %d\n", ok, v, 11)
}

func example3() {
	type Data struct {
		a, b, c int
	}
	var (
		data  = Data{1, 2, 3}
		data2 = Data{4, 5, 6}
		data3 = Data{7, 8, 9}
		v     atomic.Value
	)
	// 设置值
	v.Store(data)
	// 读取值
	v.Swap(data2)
	fmt.Println(v.Load().(Data))
	v.CompareAndSwap(data2, data3)
	fmt.Println(v)
}
