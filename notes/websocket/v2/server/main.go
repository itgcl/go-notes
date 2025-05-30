package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kbinani/screenshot"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // 允许跨域
}

// 捕获屏幕并转换为JPEG
func captureScreen() ([]byte, error) {
	img, err := screenshot.CaptureDisplay(0) // 仅截取主屏幕
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, &jpeg.Options{Quality: 50}) // 50% 质量压缩
	return buf.Bytes(), nil
}

// 处理WebSocket连接
func handleWS(conn *websocket.Conn) {
	defer conn.Close()

	for {
		time.Sleep(100 * time.Millisecond) // 每100ms发送一次（10FPS）
		imgData, err := captureScreen()
		if err != nil {
			fmt.Println("截图失败:", err)
			continue
		}

		// 发送图片数据
		if err = conn.WriteMessage(websocket.BinaryMessage, imgData); err != nil {
			fmt.Println("发送失败:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("WebSocket连接失败:", err)
			return
		}
		go handleWS(conn) // 处理WebSocket连接
	})

	fmt.Println("WebSocket 服务器启动，访问 ws://localhost:9527/ws")
	http.ListenAndServe(":9527", nil)
}
