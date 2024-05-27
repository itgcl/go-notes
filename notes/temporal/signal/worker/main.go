package producer

import (
	"errors"
	"fmt"
	"time"

	"go-notes/notes/temporal/signal"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func AudioLanguageWorkflow(ctx workflow.Context, taskID int64) error {
	// 设置信号通道
	signalChan := workflow.GetSignalChannel(ctx, "audioLanguageResult")

	// 初始化变量来接收信号数据
	var audioResult *signal.Task

	// 注册信号处理器
	signalChan.Receive(ctx, &audioResult)
	workflow.GetLogger(ctx).Info("Received signal with audio language result.", "result", audioResult)

	// 根据接收到的结果继续处理
	// ...

	return nil
}

func SignalSelect(ctx workflow.Context) error {
	// 初始化变量来接收信号数据
	var (
		audioResult    *signal.Task
		signalReceived bool
		timeout        bool
	)
	// 创建Temporal客户端
	c, err := client.Dial(client.Options{
		HostPort:  "xx",
		Namespace: "xx",
	})
	if err != nil {
		return err
	}
	defer c.Close()
	// 创建一个接收信号的通道
	sigChan := workflow.GetSignalChannel(ctx, "audioLanguageResult")

	// 设置一个1小时的超时定时器
	timeoutDuration := time.Hour
	timer := workflow.NewTimer(ctx, timeoutDuration)

	// 使用Selector来等待信号或超时
	selector := workflow.NewSelector(ctx)

	// 添加信号通道
	selector.AddReceive(sigChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &audioResult) // 接收数据
		signalReceived = true        // 标记接收到信号
		workflow.GetLogger(ctx).Info("Received signal to start activity", "content", audioResult)
	})

	// 添加定时器
	selector.AddFuture(timer, func(f workflow.Future) {
		timeout = true
		workflow.GetLogger(ctx).Info("Timer fired, timeout occurred")
	})

	// 阻塞等待信号或定时器
	selector.Select(ctx)

	// 检查是否超时
	if timeout {
		// 如果定时器触发，处理超时逻辑
		return errors.New("timeout error: did not receive signal within 1 hour")
	}

	// 检查是否因为信号接收而退出
	if signalReceived {
		// 处理信号内容
		fmt.Println(audioResult)
	}
	return nil
}
