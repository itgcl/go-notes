package worker

import (
	"errors"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

func PollingWorkflow(ctx workflow.Context) error {
	var (
		activityErr   error             // 错误
		timeout       bool              // 超时记录
		retryInterval = 1 * time.Minute // 用来保存活动错误
		maxAttempts   = 5               // 最大重试次数
		attempt       = 0               // 当前重试次数
	)
	for {
		attempt++
		if attempt > maxAttempts {
			// 返回一个错误，指出已达到最大重试次数
			return errors.New("reached maximum retry attempts")
		}

		// 设置定时器
		timerFuture := workflow.NewTimer(ctx, retryInterval)

		// 执行活动
		var result string
		activityFuture := workflow.ExecuteActivity(ctx, SomeActivity)

		// 等待定时器或活动完成
		selector := workflow.NewSelector(ctx)
		selector.AddFuture(timerFuture, func(f workflow.Future) {
			// 定时器触发，不执行操作，将继续下一轮循环
			timeout = true
		})
		selector.AddFuture(activityFuture, func(f workflow.Future) {
			err := f.Get(ctx, &result)
			if err != nil {
				workflow.GetLogger(ctx).Error("Activity failed.", "Error", err)
				activityErr = err // 保存错误
			}
		})
		selector.Select(ctx)
		if timeout {
			continue // 如果定时器触发，继续下一轮循环
		}
		// 如果有错误，返回这个错误
		if activityErr != nil {
			return activityErr
		}
		// 检查活动结果
		fmt.Println(result)
	}
}

func SomeActivity() error {
	return nil
}
