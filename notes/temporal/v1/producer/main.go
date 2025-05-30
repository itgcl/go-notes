package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"
)

// 启动工作流
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
	workflowOptions := client.StartWorkflowOptions{
		ID:                  "20",
		TaskQueue:           "test",
		WorkflowRunTimeout:  time.Minute * 5,
		WorkflowTaskTimeout: time.Hour * 24,
	}
	//err = c.TerminateWorkflow(context.Background(), "ChildWorkflow-10000", "","自定义错误信息")
	//if err != nil {
	//	log.Fatalln("Unable to terminate workflow", err)
	//}
	//return

	//ctx := context.Background()
	// First, let's make the task queue use the build id versioning feature by adding an initial
	// default version to the queue:
	//err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
	//	TaskQueue: "test",
	//	Operation: &client.BuildIDOpAddNewIDInNewDefaultSet{
	//		BuildID: "2.0",
	//	},
	//})
	//if err != nil {
	//	log.Fatalln("Unable to update worker build id compatibility", err)
	//}
	//err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
	//	TaskQueue: "test",
	//	Operation: &client.BuildIDOpAddNewCompatibleVersion{
	//		BuildID:                   "6.0",
	//		ExistingCompatibleBuildID: "5.0",
	//	},
	//})
	//if err != nil {
	//	log.Fatalln("Unable to update worker build id compatibility", err)
	//}
	//build, err := c.GetWorkerBuildIdCompatibility(ctx, &client.GetWorkerBuildIdCompatibilityOptions{TaskQueue: "test"})
	//if err != nil {
	//	log.Fatalln("Unable to update worker build id compatibility", err)
	//}
	//fmt.Println(build.Default())
	//for _, id := range build.Sets {
	//	fmt.Println("sets id :", id)
	//}
	//err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
	//	TaskQueue: "test",
	//	Operation: &client.BuildIDOpAddNewIDInNewDefaultSet{
	//		BuildID: "1.0",
	//	},
	//})
	//if err != nil {
	//	log.Fatalln("Unable to update worker build id compatibility", err)
	//}
	orderID := "33"
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, "Workflow", orderID)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	//	var result string
	//	err = we.Get(context.Background(), &result)
	//	if err != nil {
	//		log.Fatalln("Unable get workflow result", err)
	//	}
	//	log.Println("Workflow result:", result)
}
