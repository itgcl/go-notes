package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var (
	wg  sync.WaitGroup
	lc  sync.Mutex
	rwl sync.RWMutex
)

func main() {
	some19()
}

func some1() {
	fmt.Println("cpu: ", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("groot:", runtime.GOROOT())
	fmt.Println("goos:", runtime.GOOS)
}

func some2() {
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("goroutine..", i)
		}
	}()
	for i := 0; i < 5; i++ {
		runtime.Gosched()
		fmt.Println("some2..", i)
	}
}

func f() {
	defer fmt.Println("defer...")
	fmt.Println("aaaa")
	runtime.Goexit()
	fmt.Println("Bbbb")
}

func some3() {
	go func() {
		fmt.Println("fun start...")
		f()
		fmt.Println("fun over...")
	}()
	time.Sleep(time.Second)
	fmt.Println("some3 over..")
}

func some4() {
	a := 1
	go func() {
		a = 2
		fmt.Println(a)
	}()

	a = 3
	time.Sleep(1)
	fmt.Println(a)
}

var ticket = 10

func some5() {
	go saleTickets("售票处1")
	go saleTickets("售票处2")
	go saleTickets("售票处3")
	go saleTickets("售票处4")
	go saleTickets("售票处5")
	time.Sleep(time.Second * 5)
	fmt.Println("over")
}

func saleTickets(s string) {
	rand.Seed(time.Now().UnixNano())
	for {
		func() {
			defer lc.Unlock()
			lc.Lock()
			if ticket > 0 {
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				fmt.Println(s, "售出", ticket)
				ticket--
			} else {
				fmt.Println("售罄")
				return
			}
		}()
	}
}

func some6() {
	wg.Add(3)
	go writeData()
	go readData()
	go readData()
	//go readData()
	//go writeData()
	wg.Wait()
	fmt.Println("main over...")
}

func readData() {
	defer wg.Done()
	fmt.Println("准备读...")
	rwl.RLock()
	fmt.Println("正在读取数据...")
	time.Sleep(time.Second * 3)
	fmt.Println("读取结束...")
	rwl.RUnlock()
}
func writeData() {
	defer wg.Done()
	fmt.Println("准备写...")
	rwl.Lock()
	fmt.Println("正在写入数据...")
	time.Sleep(time.Second * 3)
	fmt.Println("写入结束...")
	rwl.Unlock()
}

func some7() {
	var ch = make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
		fmt.Println("go runtime over")
		ch <- 100
	}()
	<-ch
	fmt.Println("some7 over")
}

func some8() {
	ch := make(chan int)
	go func() {
		fmt.Println("goroutine start...")
		time.Sleep(time.Second * 3)
		v := <-ch
		fmt.Println("data:", v)
	}()
	ch <- 10
	fmt.Println("some 8 over")
}

// 携程作为发送者，向channel发送数据，任务发送完成close channel
func some9() {
	var ch = make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= 10; i++ {
			ch <- i
		}
	}()
	for {
		time.Sleep(time.Second)
		v, ok := <-ch
		if !ok {
			fmt.Println("channel closed...")
			break
		}
		fmt.Println(v)
	}
	fmt.Println("some9 over")
}

func some10() {
	var (
		ch   = make(chan int, 5)
		quit = make(chan bool)
	)
	go func() {
		for {
			select {
			case v, ok := <-ch:
				if !ok {
					fmt.Println("channel closed...")
					quit <- true
					return
				}
				// 模拟数据处理
				time.Sleep(time.Second)
				fmt.Println(v)
			}
		}
	}()
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)
	fmt.Println("channel send over...")
	<-quit
	fmt.Println("over")
}

func some11() {
	var ch = make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * 200)
			ch <- i
		}
		fmt.Println("channel send over")
	}()
	for v := range ch {
		fmt.Println("data:", v)
	}
	fmt.Println("over")
}

func some12() {
	ch := make(chan int, 3)
	go func() {
		defer close(ch)
		var wg sync.WaitGroup
		for i := 1; i <= 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 10; i++ {
					ch <- i
					time.Sleep(time.Second)
				}
			}()
		}
		wg.Wait()
		fmt.Println("send channel over...")
	}()
	go func() {
		for v := range ch {
			fmt.Println("data:", v)
		}
	}()
	time.Sleep(time.Second * 100)
}

func some13() {
	ch := make(chan int, 3)
	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
			time.Sleep(time.Second)
		}
		fmt.Println("send channel over...")
	}()
	for v := range ch {
		fmt.Println("data:", v)
	}
	fmt.Println("some13 over")
}

func some14() {
	var (
		ch   = make(chan string)
		done = make(chan bool)
	)
	go sendData(ch, done)
	fmt.Println("some14 read ch start...")
	data := <-ch
	fmt.Println("some14 read ch over")
	fmt.Println(data)

	fmt.Println("some14 write start...")
	ch <- "some14"
	fmt.Println("some14 write over")
	<-done
}

func sendData(ch chan string, done chan bool) {
	fmt.Println("go routine write start...")
	ch <- "go routine"
	fmt.Println("go routine write over")
	data := <-ch
	fmt.Println("goroutine:", data)
	done <- true
	fmt.Println("send data over")
}

func some15() {
	ch := make(chan int)
	// 写入
	go fun1(ch)
	// 读取
	data := <-ch
	fmt.Println("func1 data:", data)
	go fun2(ch)
	ch <- 100
	fmt.Println("some15 over")
}

func fun1(ch chan<- int) {
	ch <- 10
	fmt.Println("func1 over")
}
func fun2(ch <-chan int) {
	v := <-ch
	fmt.Println("func2 data:", v)
	fmt.Println("func2 over")
}

func some16() {
	timer := time.NewTimer(3 * time.Second)
	ch := timer.C
	fmt.Println(time.Now())
	fmt.Println(<-ch)
}

func some17() {
	timer := time.NewTimer(time.Second * 1)
	go func() {
		fmt.Println("timer start")
		<-timer.C
		fmt.Println("timer over")
	}()
	go func() {
		fmt.Println("stop start")
		timer.Stop()
		fmt.Println("stop over")
	}()
	time.Sleep(time.Second * 5)
	fmt.Println("some17 over...")
}

func some18() {
	ch := time.After(3 * time.Second)
	fmt.Println(time.Now())
	fmt.Println(<-ch)
}

func some19() {
	var (
		ch1 = make(chan int)
		ch2 = make(chan int)
	)
	go func() {
		time.Sleep(time.Second)
		ch1 <- 1
	}()
	go func() {
		ch2 <- 2
	}()
	select {
	case v := <-ch1:
		fmt.Println("ch1:", v)
	case v := <-ch2:
		fmt.Println("ch2:", v)
	}
	fmt.Println("some19")
}
