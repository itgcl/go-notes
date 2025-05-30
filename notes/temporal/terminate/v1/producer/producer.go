package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	//c1, err := client.NewNamespaceClient(client.Options{})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//retentionPeriod := time.Hour * 24 * 30
	//if err := c1.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
	//	Namespace:                        "dev",
	//	WorkflowExecutionRetentionPeriod: &retentionPeriod, // 执行记录的保留时间
	//}); err != nil {
	//	fmt.Println(err)
	//}
	//
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	// 设置workflowID和队列名
	//workflowOptions := client.StartWorkflowOptions{
	//	ID:                  "19",
	//	TaskQueue:           "test",
	//	WorkflowRunTimeout:  time.Minute * 5,
	//	WorkflowTaskTimeout: time.Hour * 24,
	//}
	err = c.TerminateWorkflow(context.Background(), "ChildWorkflow-33", "", "自定义错误信息")
	if err != nil {
		log.Fatalln("Unable to terminate workflow", err)
	}
	return
}
