package main

import (
	"fmt"
	"log"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	connect, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		log.Fatalf("connect error: %s", err)
	}
	l := zk.NewLock(connect, "/look1", zk.WorldACL(zk.PermAll))
	if err := l.Lock(); err != nil {
		log.Fatalf("lock error: %s", err)
	}
	println("lock succ, do your business logic")
	go func() {
		l := zk.NewLock(connect, "/look1", zk.WorldACL(zk.PermAll))
		fmt.Println("goroutine start lock....")
		fmt.Println(22, l.Lock())
		println("goroutine lock success")
		time.Sleep(time.Second * 5)
		if err := l.Unlock(); err != nil {
			log.Fatalf("goroutine unlock error: %s", err)
		}
		println("goroutine unlock success")
	}()
	time.Sleep(time.Second * 10)
	// do some thing
	if err := l.Unlock(); err != nil {
		log.Fatalf("unlock error: %s", err)
	}
	println("unlock succ, finish business logic")
	time.Sleep(time.Second * 20)
}
