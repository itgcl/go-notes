package main

import (
	"context"
	"flag"
	"fmt"
	"go-notes/notes/grpc/v4/internal/service"
	"go-notes/notes/grpc/v4/proto"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// Update
)

var grpcPort = flag.Int("grpc_port", 9001, "the port to serve on")
var httpPort = flag.Int("http_port", 8080, "the port to restful serve on")

func main() {
	server := grpc.NewServer()
	helloService := service.NewHelloService()
	proto.RegisterHelloServiceServer(server, helloService)
	// Serve gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Serving gRPC on 0.0.0.0" + fmt.Sprintf(":%d", *grpcPort))
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 2. 启动 HTTP 服务
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", *grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = proto.RegisterHelloServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *httpPort),
		Handler: gwmux,
	}
	log.Println("Serving gRPC-Gateway on http://127.0.0.1" + fmt.Sprintf(":%d", *httpPort))
	log.Fatalln(gwServer.ListenAndServe())
}
