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
	w := worker.New(c, "test", worker.Options{})
	// 注册 workflow和activity
	w.RegisterWorkflow(buy.Workflow)
	w.RegisterActivity(buy.Activity)
	// 执行消费
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
