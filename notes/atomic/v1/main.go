package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var a int64 = 10
	// 原子操作 对数据操作
	//s := atomic.AddInt64(&a, 2)
	//fmt.Println(s)

	// 读取数据
	v := atomic.LoadInt64(&a)
	fmt.Println(v)

}
