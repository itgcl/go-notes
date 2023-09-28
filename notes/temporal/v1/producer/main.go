package main

import (
	"context"
	buy "go-notes/notes/temporal/v1"
	"log"

	"go.temporal.io/sdk/client"
)

// 启动工作流
func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	// 设置workflowID和队列名
	workflowOptions := client.StartWorkflowOptions{
		ID:        "123456",
		TaskQueue: "test",
	}
	orderID := "10000"
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, buy.Workflow, orderID)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
