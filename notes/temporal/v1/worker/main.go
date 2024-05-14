package main

import (
	"log"

	"go-notes/notes/temporal/v1"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	// 监听队列名
	// w := worker.New(c, "test", worker.Options{})
	w := worker.New(c, "test",
		worker.Options{
			// BuildID:                 "6.0",
			// UseBuildIDForVersioning: true,
		},
	)
	// w := worker.New(c, "test", worker.Options{BuildID: "1.1", UseBuildIDForVersioning: true})
	// 注册 workflow和activity
	w.RegisterWorkflow(buy.Workflow)
	w.RegisterActivity(buy.Activity)
	w.RegisterActivity(buy.ActivityV1)
	w.RegisterActivity(buy.ActivityV2)
	w.RegisterActivity(buy.ActivityV3)
	w.RegisterActivity(buy.ActivityV4)
	w.RegisterActivity(buy.ActivityV5)
	// 执行消费
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("Temporal Worker 运行失败：%v", err)
	}
	// 停止 Worker
	w.Stop()
	// 关闭 Temporal 客户端
	c.Close()
	// 关闭日志
	log.Println("Temporal Worker 已退出")
}
