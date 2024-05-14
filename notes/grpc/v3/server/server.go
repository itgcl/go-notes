package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"runtime/debug"

	pb "go-notes/notes/grpc/v3/proto"
	article "go-notes/notes/grpc/v3/server/internal/service"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Run() {
	flag.Parse()
	// 注册interceptor
	// var opts []grpc.ServerOption
	// opts = append(opts, grpc.UnaryInterceptor(interceptor))
	// 实例化server
	// server := grpc.NewServer(opts...)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptor,
			interceptor2,
		)),
	)
	// 实例化article service
	service := article.NewService()
	// 注册 article
	pb.RegisterArticleServiceServer(server, service)
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	fmt.Println("start server")
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}

// interceptor 拦截器
func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

func interceptor2(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("interceptor2 method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("interceptor2 method: %s, %v", info.FullMethod, resp)
	return resp, err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)
}
