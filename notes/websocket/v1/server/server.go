package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // 创建一个WebSocket升级器

func main() {
	http.HandleFunc("/websocket", handleWebSocket) // 定义WebSocket处理函数
	log.Fatal(http.ListenAndServe(":8080", nil))   // 启动HTTP服务器
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // 升级HTTP连接为WebSocket连接
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	for {
		// 读取客户端发送的消息
		t, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		//message := "ww"
		//t := 1
		//// 处理消息并回复
		log.Printf("Received message: %s, type: %d", message, t)
		//err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Server received message: %s", message)))
		//if err != nil {
		//	log.Println("Error writing message:", err)
		//	break
		//}
		time.Sleep(time.Millisecond * 100)
	}
}
