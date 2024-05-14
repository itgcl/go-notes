package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Product struct {
	CreateTime *timestamppb.Timestamp `json:"createTime"`
}

func main() {
	var product Product
	jsonData := []byte(`{"createTime":"2023-11-15T03:19:50.000Z"}`)
	if err := json.Unmarshal(jsonData, &product); err != nil {
		fmt.Println("解组JSON时出错:", err)
		return
	}

	fmt.Println("产品创建时间:", product.CreateTime.AsTime())
	return
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/websocket"}
	log.Printf("连接到 %s", u.String())
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("连接错误:", err)
	}
	defer c.Close()

	fmt.Println(resp)

	done := make(chan struct{})

	//go func() {
	//	defer close(done)
	//	for {
	//		_, message, err := c.ReadMessage()
	//		if err != nil {
	//			log.Println("读取错误:", err)
	//			return
	//		}
	//		log.Printf("收到消息: %s", message)
	//	}
	//}()

	for i := 2000; i < 3000; i++ {
		err := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("客户端发送消息: %d", i)))
		if err != nil {
			log.Println("写入错误:", err)
			return
		}
		fmt.Printf("send message over  %d\n", i)
	}
	return

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte("客户端发送消息: "+t.String()))
			if err != nil {
				log.Println("写入错误:", err)
				return
			}
			fmt.Println("send message over")
		case <-interrupt:
			log.Println("接收到中断信号，关闭连接")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("关闭连接错误:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
