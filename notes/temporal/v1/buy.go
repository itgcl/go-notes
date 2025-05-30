package buy

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
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
		StartToCloseTimeout: 60 * time.Second,
		HeartbeatTimeout:    time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 2,
			BackoffCoefficient: 1.0,
			MaximumAttempts:    3,
			MaximumInterval:    time.Minute * 30,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	//cwo := workflow.ChildWorkflowOptions{
	//	// TaskQueue: "compliance-detection",
	//	RetryPolicy: &temporal.RetryPolicy{
	//		InitialInterval:    time.Second * 30, // 间隔时间30秒
	//		BackoffCoefficient: 2.0,              // 退避系数2
	//		MaximumAttempts:    5,                // 最大重试次数5次
	//		MaximumInterval:    time.Hour,        // 最大执行时间1小时
	//	},
	//	WorkflowID: fmt.Sprintf("ChildWorkflow-%s", orderId),
	//}

	// cctx = workflow.WithChildOptions(ctx, cwo)
	//if err := workflow.ExecuteChildWorkflow(cctx, ChildWorkflow, orderId).Get(ctx, nil); err != nil {
	//	if temporal.IsTerminatedError(err) {
	//		fmt.Println("child workflow terminated:", err)
	//	} else {
	//		fmt.Println("child workflow exec  error", err)
	//		return "", err
	//	}
	//}

	ctx = workflow.WithValue(ctx, "test", "qwe")
	logger := workflow.GetLogger(ctx)

	logger.Info("HelloWorld workflow started", "orderId", orderId)
	result, err := v1(ctx, orderId)
	fmt.Printf("result: %s, err: %v\n", result, err)
	if err != nil {
		return "", err
	}
	//var result string
	//err := workflow.ExecuteActivity(ctx, ActivityV1, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}
	//
	//err = workflow.ExecuteActivity(ctx, ActivityV2, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}

	//err = workflow.ExecuteActivity(ctx, ActivityV3, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}
	//
	//err = workflow.ExecuteActivity(ctx, ActivityV4, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}
	//err = workflow.ExecuteActivity(ctx, ActivityV5, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}
	fmt.Println("workflow completed.", "result", result)
	return "", nil
}

func ChildWorkflow(ctx workflow.Context, orderId string) error {
	// 设置工作超时时间 TODO StartToCloseTimeout或ScheduleToCloseTimeout二者必须设置一个，不然报错
	//ao := workflow.ActivityOptions{
	//	StartToCloseTimeout: 60 * time.Second,
	//	HeartbeatTimeout:    time.Second * 5,
	//	RetryPolicy: &temporal.RetryPolicy{
	//		InitialInterval:    time.Second * 2,
	//		BackoffCoefficient: 1.0,
	//		MaximumAttempts:    3,
	//		MaximumInterval:    time.Minute * 30,
	//	},
	//}
	//ctx = workflow.WithActivityOptions(ctx, ao)
	// logger := workflow.GetLogger(ctx)
	// logger.Info("HelloWorld workflow started", "orderId", orderId)
	//err := workflow.ExecuteActivity(ctx, ActivityV1, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}
	//
	//err = workflow.ExecuteActivity(ctx, ActivityV2, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}

	//err = workflow.ExecuteActivity(ctx, ActivityV3, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}
	//
	//err = workflow.ExecuteActivity(ctx, ActivityV4, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}
	//err = workflow.ExecuteActivity(ctx, ActivityV5, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}
	fmt.Println("time sleep start")
	err := workflow.Sleep(ctx, time.Second*120)
	if err != nil {
		fmt.Println("sleep error", err)
		return err
	}
	fmt.Println("child workflow completed.", orderId)
	return nil
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

func v1(ctx workflow.Context, orderId string) (string, error) {
	var result string
	fmt.Println("v1 value: ", ctx.Value("test"))
	err := workflow.ExecuteActivity(ctx, ActivityV1, orderId).Get(ctx, &result)
	if err != nil {
		fmt.Println("execution error:", err)
		return "", err
	}
	//err = workflow.ExecuteActivity(ctx, ActivityV2, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}
	//err = workflow.ExecuteActivity(ctx, ActivityV3, orderId).Get(ctx, &result)
	//if err != nil {
	//	fmt.Println("execution error:", err)
	//	return "", err
	//}
	//fmt.Println("start error")
	return result, nil
}

func ActivityV1(ctx context.Context, name string) (string, error) {
	fmt.Println("v1 order", name)
	//if err := Some11(ctx); err != nil {
	//	return "", err
	//}
	//time.Sleep(time.Second * 2)
	//if err := Some22(ctx); err != nil {
	//	return "", err
	//}
	fmt.Println(ctx.Value("test"), "get ctx value")
	fmt.Println("v1 over")
	return "v1111111", nil
}

func Some11(ctx context.Context) error {
	fmt.Println("11 some")
	return nil
}

func Some22(ctx context.Context) error {
	fmt.Println("22 some")
	return errors.New("some2 error")
}

var count = 0

func ActivityV2(ctx context.Context, name string) (string, error) {
	fmt.Println("v2 orderid", name)
	time.Sleep(time.Second * 2)
	fmt.Println("v2 over")
	//count++
	//if count > 3 {
	//	return "v222222", nil
	//}
	return "v222222222", errors.New("xxxxxxxxx")
}

func ActivityV3(ctx context.Context, name string) (string, error) {
	fmt.Println(" v3 orderid", name)
	return "", errors.New("111111111111")
	time.Sleep(time.Second * 2)
	fmt.Println("v3 over")
	return "v3333333", nil
}

func ActivityV4(ctx context.Context, name string) (string, error) {
	fmt.Println(" v4 orderid", name)
	time.Sleep(time.Second * 8)
	fmt.Println(" v4 over")
	return "v44444", nil
}

func ActivityV5(ctx context.Context, name string) (string, error) {
	fmt.Println(" v5 orderid", name)
	time.Sleep(time.Second * 10)
	fmt.Println(" v5 over")
	return "v555555", nil
}
