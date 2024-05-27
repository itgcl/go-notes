```go
package main

import (
	"go.temporal.io/sdk/workflow"
	"context"
	"time"
)

func YourWorkflow(ctx workflow.Context, param string) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Workflow started", "param", param)

	// 使用版本控制处理活动的变更
	var result string
	v := workflow.GetVersion(ctx, "YourActivityChange", workflow.DefaultVersion, 1)
	// v := workflow.GetVersion(ctx, "YourActivityChange", 1, 2) // 后续内容变更修改版本号
	// v := workflow.GetVersion(ctx, "YourActivityChange", 2, 3) // 后续内容变更修改版本号
	if v == workflow.DefaultVersion {
		err := workflow.ExecuteActivity(ctx, YourOldActivity, param).Get(ctx, &result)
		if err != nil {
			logger.Error("Activity failed.", "Error", err)
			return "", err
		}
	} else {
		err := workflow.ExecuteActivity(ctx, YourNewActivity, param).Get(ctx, &result)
		if err != nil {
			logger.Error("Activity failed.", "Error", err)
			return "", err
		}
	}

	logger.Info("Workflow completed.", "result", result)
	return result, nil
}

// 旧版本活动
func YourOldActivity(ctx context.Context, param string) (string, error) {
	// 旧的活动逻辑
	return "result from old activity", nil
}

// 新版本活动
func YourNewActivity(ctx context.Context, param string) (string, error) {
	// 新的活动逻辑
	return "result from new activity", nil
}
```

## 发送端没有设置版本号，为什么？
在Temporal中，版本号是在工作流代码中控制的，而不是在发送端设置的。发送端（客户端代码）负责启动工作流，但不涉及工作流内部的逻辑细节，如版本控制。版本控制是在工作流逻辑中实现的，以便在工作流的不同执行阶段根据需要应用不同的逻辑或活动版本。
## workflow.GetVersion 中 "YourActivityChange" 是怎么来的？
"YourActivityChange" 是一个开发者定义的字符串，用于标识特定的代码更改点。这个字符串可以是任何值，但它应该在工作流中唯一地标识一个版本控制点。当你调用 workflow.GetVersion 时，Temporal工作流服务会查看这个标识点的历史版本，并返回适当的版本号。这样，你就可以在工作流中使用这个版本号来决定执行哪个代码分支。
## 接收端是否需要同时注册 YourOldActivity 和 YourNewActivity？
是的，如果你的工作流计划支持旧版本和新版本的活动，你需要在工作流启动时注册所有相关的活动。这确保了无论工作流实例是使用旧版本还是新版本启动的，相应的活动都可用。这是确保向后兼容性和新功能并存的关键。