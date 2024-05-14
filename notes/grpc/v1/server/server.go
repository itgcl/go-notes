package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"time"

	pb "go-notes/notes/grpc/v1/proto"
	article "go-notes/notes/grpc/v1/server/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

const (
	serverName = "grpc.article.v1.Article"
	AppId      = "zhangsan"
	AppKey     = "123456"
)

var (
	port  = flag.Int("port", 8181, "the port to serve on")
	sleep = flag.Duration("sleep", time.Second*5, "duration between changes in health")
)

func Run() {
	flag.Parse()
	// 注册interceptor
	opts := grpc.UnaryInterceptor(interceptor)
	// 实例化server
	server := grpc.NewServer(opts)
	healthcheck := health.NewServer()
	// healthServer.SetServingStatus("serverName", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(server, healthcheck)
	// 实例化article service
	service := article.NewService(*port)
	// 注册 article
	pb.RegisterArticleServiceServer(server, service)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
	go simulationUnhealthyStatus(healthcheck)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}

// interceptor 拦截器
func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}
	// 继续处理请求
	return handler(ctx, req)
}

// auth 验证Token
func auth(ctx context.Context) error {
	md, has := metadata.FromIncomingContext(ctx)
	if !has {
		return errors.New("params error")
	}
	var (
		appid  string
		appkey string
	)
	if val, ok := md["appId"]; ok {
		appid = val[0]
	}
	if val, ok := md["appKey"]; ok {
		appkey = val[0]
	}
	if appid != AppId || appkey != AppKey {
		return errors.New("username or password error")
	}
	return nil
}

// 模拟服务异常
func simulationUnhealthyStatus(healthcheck *health.Server) {
	go func() {
		// asynchronously inspect dependencies and toggle serving status as needed
		next := healthpb.HealthCheckResponse_SERVING

		for {
			healthcheck.SetServingStatus(serverName, next)

			if next == healthpb.HealthCheckResponse_SERVING {
				next = healthpb.HealthCheckResponse_NOT_SERVING
			} else {
				next = healthpb.HealthCheckResponse_SERVING
			}

			time.Sleep(*sleep)
		}
	}()
}
