package main

import (
	"fmt"
	"log"
	"time"

	"go-notes/notes/temporal/signal"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	wor "go.temporal.io/sdk/workflow"
)

func main() {
	// 创建Temporal客户端
	c, err := client.Dial(client.Options{
		HostPort:  "localhost:7233",
		Namespace: "default",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	var options worker.Options
	// 工作流ID和运行ID
	//workflowID := "11" // 你的工作流实例ID
	//workflowOptions := client.StartWorkflowOptions{
	//	ID:        strconv.FormatInt(11, 10),
	//	TaskQueue: "test",
	//}
	// 构造信号内容
	//audioResult := signal.Task{
	//	ID: 456,
	//}
	//we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, "workflow", "aaa")
	//if err != nil {
	//	log.Fatalln("Error starting workflow", err)
	//}
	//fmt.Println(we.GetID())
	// 发送信号
	//err = c.SignalWorkflow(context.Background(), workflowID, "", "audioDataResult1", audioResult)
	//if err != nil {
	//	log.Fatalln("Error sending signal to workflow", err)
	//}
	//
	//log.Println("Signal sent successfully")

	w := worker.New(c, "test", options)
	w.RegisterWorkflow(workflow)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func workflow(ctx wor.Context, req interface{}) error {
	// 设置工作超时时间 StartToCloseTimeout或ScheduleToCloseTimeout二者必须设置一个，不然报错
	ao := wor.ActivityOptions{
		StartToCloseTimeout: time.Hour * 24,
		HeartbeatTimeout:    time.Second * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 10,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    3,
			MaximumInterval:    time.Minute * 5,
		}, // 重试策略
	}
	ctx = wor.WithActivityOptions(ctx, ao)
	// 打印请求参数
	log.Printf("req: %v", req)
	YourWorkflow(ctx)
	return nil
}

func YourWorkflow(ctx wor.Context) error {
	fmt.Println(11)
	signalChan := wor.GetSignalChannel(ctx, "audioDataResult")
	audioResult := signal.Task{}
	signalChan.Receive(ctx, &audioResult)
	fmt.Println(audioResult)

	// selector := wor.NewSelector(ctx)

	//selector.AddReceive(signalChan, func(c wor.ReceiveChannel, more bool) {
	//	var signalPayload string
	//	c.Receive(ctx, &signalPayload)
	//	wor.GetLogger(ctx).Info("Received signal", "payload", signalPayload)
	//})
	//
	//// 等待信号
	//selector.Select(ctx)

	return nil
}
