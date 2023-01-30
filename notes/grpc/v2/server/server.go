package server

import (
	"flag"
	"fmt"
	pb "go-notes/notes/grpc/v2/proto"
	demo "go-notes/notes/grpc/v2/server/internal/service"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8181, "the port to serve on")
)

func Run() {
	flag.Parse()
	// 实例化server
	server := grpc.NewServer()
	service := demo.NewService(*port)
	// 注册
	pb.RegisterStreamDemoServiceServer(server, service)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
