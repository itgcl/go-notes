package main

import (
	"fmt"
	"sync"
	"time"
)

var done bool

func main() {
	c := sync.NewCond(&sync.Mutex{})
	go read("read 1", c)
	go read("read 2", c)
	go read("read 3", c)
	go write("write", c)
	time.Sleep(time.Second * 3)
}

func read(s string, c *sync.Cond) {
	c.L.Lock()
	for !done {
		fmt.Println(s, "携程开始等待")
		c.Wait()
		fmt.Println(s, "携程接收到通知")
	}
	fmt.Println("read over")
	c.L.Unlock()
}

func write(s string, c *sync.Cond) {
	fmt.Println(s, "start")
	done = true
	time.Sleep(time.Second)
	fmt.Println(s, "over")
	c.Broadcast()
	fmt.Println("broadcast over")
}
