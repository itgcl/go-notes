package main

import "fmt"

/*
开启一个writeData协程，向管道intChan中写入50个整数
开启一个readChan协程，从管道中读取writeData写入的数据
注意：writeData和readChan操作的是同一个管道
主线程需要等待writeData和readData协程都完成工作才能退出【管道】
*/

/*
	发送者关闭通道 接收者处理通道关闭情况
	不会出现发送者任务发送一半程序退出的情况, 因为主线程等待exitChannel数据 exitChannel在接收者完成后才会产生
	接收者完成的情况是依据发送者发送完成并关闭通道 所以确保了发送者一定全部发送成功
*/

// 1.
//func main()  {
//	ch := make(chan int, 20)
//	exitCh := make(chan struct{})
//	go send(ch)
//	go receive(ch, exitCh)
//	for range exitCh{
//		fmt.Println("收到程序结束")
//		break
//	}
//
//}
//
//func send(ch chan<- int)  {
//	for i := 0; i< 50; i++ {
//		ch <- i
//	}
//	fmt.Println("发送完成 关闭通道")
//	close(ch)
//}
//
//func receive(ch <-chan int, exitCh chan<- struct{})  {
//	for {
//		value, ok := <-ch
//		if ok == false {
//			break
//		}
//		fmt.Printf("value:%v\n", value)
//	}
//	fmt.Println("接收完成信号")
//	exitCh <- struct{}{}
//}

// 2. 扩展 使用线程池(接收者)
func main() {
	routineCount := 5
	ch := make(chan int, 20)
	receivePassCh := make(chan struct{}, routineCount)
	go send(ch)
	for i := 0; i < routineCount; i++ {
		go receive(ch, receivePassCh)
	}
	total := 0
	for range receivePassCh {
		total++
		if total == routineCount {
			break
		}
	}
	fmt.Println("收到程序结束")
}

func send(ch chan<- int) {
	for i := 0; i < 50; i++ {
		ch <- i
	}
	close(ch)
}

// 线程数完成进行报告 完成报告数等于线程数 结束
func receive(ch <-chan int, receivePassCh chan<- struct{}) {
	for {
		value, ok := <-ch
		if !ok {
			break
		}
		fmt.Printf("value:%v\n", value)
	}
	fmt.Println("接收完成信号")
	receivePassCh <- struct{}{}
}

// TODO 目前使用了排它锁达到效果 不是很好需要修改
// 3. 扩展 使用线程池(发送者)
//var once sync.Once
//func main()  {
//	ch := make(chan int, 20)
//	exitCh := make(chan struct{})
//	counter := new(int)
//	// 开启5个goroutine发送者
//	for i := 0; i < 5; i++{
//		go send(ch, counter)
//	}
//	go receive(ch, exitCh)
//	for range exitCh{
//		fmt.Println("收到程序结束")
//		break
//	}
//
//}
//
//var lock  sync.Mutex
//func send(ch chan<- int, counter *int)  {
//	for i := 0; i < 50; i++ {
//		lock.Lock()
//		*counter++
//		if *counter <= 50 {
//			ch <- i
//		}else{
//			once.Do(func() {
//				close(ch)
//				fmt.Println("发送完成 关闭通道")
//			})
//			break
//		}
//		lock.Unlock()
//	}
//
//}
//
//func receive(ch <-chan int, exitCh chan<- struct{})  {
//	for {
//		value, ok := <-ch
//		if ok == false {
//			break
//		}
//		fmt.Printf("value:%v\n", value)
//	}
//	fmt.Println("接收完成信号")
//	exitCh <- struct{}{}
//}
