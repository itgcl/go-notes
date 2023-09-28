package buy

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// 流程活动函数
func CreateOrder(ctx context.Context, orderId string) error {
	// 调用订单服务API创建订单
	fmt.Println("create order over", orderId)
	return nil
}

func Payment(ctx context.Context, orderId string) error {
	// 调用支付服务API进行支付
	fmt.Println("payment over", orderId)
	return nil
}

func CancelOrder(ctx context.Context, orderId string) error {
	// 调用订单服务API取消订单
	fmt.Println("cancel order over", orderId)
	return nil
}

func DeductInventory(ctx context.Context, orderId string) error {
	// 调用库存服务API扣减库存
	fmt.Println("deduct inventory over", orderId)
	return nil
}

// Workflow 工作流定义
func Workflow(ctx workflow.Context, orderId string) (string, error) {
	// 设置工作超时时间 TODO StartToCloseTimeout或ScheduleToCloseTimeout二者必须设置一个，不然报错
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	//logger := workflow.GetLogger(ctx)
	//logger.Info("HelloWorld workflow started", "orderId", orderId)

	var result string
	err := workflow.ExecuteActivity(ctx, Activity, orderId).Get(ctx, &result)
	if err != nil {
		fmt.Println("execution error:", err)
		return "", err
	}
	fmt.Println("workflow completed.", "result", result)
	return result, nil
}

func Activity(ctx context.Context, orderId string) (string, error) {
	fmt.Println(111111)
	err := CreateOrder(ctx, orderId)
	if err != nil {
		fmt.Println("create order error", err)
		return "", err
	}
	err = Payment(ctx, orderId)
	if err != nil {
		err := CancelOrder(ctx, orderId)
		if err != nil {
			return "", err
		}
	}

	err = DeductInventory(ctx, orderId)
	if err != nil {
		return "", err
	}

	return "xxx", nil
}
