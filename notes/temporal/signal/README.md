## 使用信号
使用信号（Signals）在Temporal中实现异步通信是一种有效的方法，特别是当你的工作流需要从外部系统接收事件或数据时。以下是如何在Temporal工作流中实现信号接收方和信号通知方的示例代码。
1. 信号接收方
   在工作流定义中，你首先需要定义一个可以接受信号的方法。这个方法通过 workflow.RegisterSignalHandler 注册信号处理器。这里我们以音频语种检测结果为例：
```go
package main

func AudioLanguageWorkflow(ctx workflow.Context, taskID int64) error {
// 设置信号通道
signalChan := workflow.GetSignalChannel(ctx, "startActivitySignal")

    // 初始化变量来接收信号数据
    var audioResult TaskReply

    // 注册信号处理器
    signalChan.Receive(ctx, &audioResult)
    workflow.GetLogger(ctx).Info("Received signal with audio language result.", "result", audioResult)

    // 根据接收到的结果继续处理
    // ...

    return nil
}
```
在这个例子中，signalChan.Receive 会阻塞等待直到信号到达。你也可以使用 workflow.NewSelector 来同时等待多个输入，例如定时器和信号。

2. 信号通知方
   信号的发送方通常是外部的系统或服务。在Temporal中，你可以使用客户端API来向工作流发送信号。这里是如何从外部应用发送信号到工作流的示例：

```go
package main

import (
   "errors"
   "fmt"

   "go.temporal.io/sdk/client"
   "go.temporal.io/sdk/workflow"
)

func AudioLanguageWorkflow(ctx workflow.Context, taskID int64) error {
   // 创建一个接收信号的通道
   sigChan := workflow.GetSignalChannel(ctx, "startActivitySignal")

   // 设置一个1小时的超时定时器
   timeoutDuration := time.Hour
   timer := workflow.NewTimer(ctx, timeoutDuration)

   // 使用Selector来等待信号或超时
   selector := workflow.NewSelector(ctx)
   var audioResult TaskReply
   var timeout bool

   selector.AddReceive(sigChan, func(c workflow.ReceiveChannel, more bool) {
      c.Receive(ctx, &audioResult)
      workflow.GetLogger(ctx).Info("Received signal to start activity", "content", audioResult)
   })

   selector.AddFuture(timer, func(f workflow.Future) {
      // Timer完成，但不需要在这里处理任何东西
      timeout = true
   })

   // 阻塞等待信号或定时器
   selector.Select(ctx)
   if timeout {
      return errors.New("timeout")
   }
   // 接收到数据，进行后续处理
   fmt.Println(audioResult)
   return nil
}
```
在这个示例中，我们首先创建了一个Temporal客户端，然后使用该客户端的 SignalWorkflow 方法发送信号。你需要指定工作流ID和RunID（如果你只想发送到最近的运行实例，RunID可以是空字符串），信号名称，以及要发送的数据。

通过这种方式，你可以从任何可以访问Temporal客户端的地方发送信号到工作流，实现异步的交互。确保信号名在发送方和接收方保持一致。
