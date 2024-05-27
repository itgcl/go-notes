package main

import (
	"context"
	"log"

	"go-notes/notes/temporal/signal"

	"go.temporal.io/sdk/client"
)

func main() {
	// 创建Temporal客户端
	c, err := client.Dial(client.Options{
		HostPort:  "xx",
		Namespace: "xx",
	})
	defer c.Close()
	// 工作流ID和运行ID
	workflowID := "yourWorkflowID" // 你的工作流实例ID
	runID := "yourRunID"           // 可以是空字符串，表示最新的运行实例

	// 构造信号内容
	audioResult := signal.Task{
		// 填充结构体字段
	}
	// 发送信号
	err = c.SignalWorkflow(context.Background(), workflowID, runID, "audioDataResult", audioResult)
	if err != nil {
		log.Fatalln("Error sending signal to workflow", err)
	}

	log.Println("Signal sent successfully")
}
