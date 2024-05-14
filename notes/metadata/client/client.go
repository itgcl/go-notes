package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/metadata"
	mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	gg "google.golang.org/grpc"
	mmmd "google.golang.org/grpc/metadata"
)

func main() {
	// callHTTP()
	callGRPC()
}

func callHTTP() {
	conn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			mmd.Client(),
		),
		http.WithEndpoint("127.0.0.1:8000"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterHTTPClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToClientContext(ctx, "x-md-global-extra", "2233")
	reply, err := client.SayHello(ctx, &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply)
}

func callGRPC() {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMiddleware(
			mmd.Client(),
		),
		grpc.WithUnaryInterceptor(
			// metadataInterceptor,
			metadataInterceptor1,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToClientContext(ctx, "x-md-global-extra", "2233")
	reply, err := client.SayHello(ctx, &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v \n", reply)
}

func metadataInterceptor1(ctx context.Context, method string, req, reply interface{}, cc *gg.ClientConn, invoker gg.UnaryInvoker, opts ...gg.CallOption) error {
	// 修改metadata
	// ctx = metadata.AppendToClientContext(ctx, "x-md-global-eeeeee", "rrrrr")
	// ctx = metadata.MergeToClientContext(ctx, metadata.Metadata{"x-md-global-a": []string{"eee"}})
	// 继续执行请求
	// m, ok := metadata.FromClientContext(ctx)
	// fmt.Println(m, ok)
	ctx = mmmd.AppendToOutgoingContext(ctx, "x-md-global-tt", "99")
	mm, ok := mmmd.FromOutgoingContext(ctx)
	fmt.Println(mm, ok, 222)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func AppendToClientContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToClientContext got an odd number of input pairs for metadata: %d", len(kv)))
	}
	md, _ := metadata.FromClientContext(ctx)
	md = md.Clone()
	for i := 0; i < len(kv); i += 2 {
		md.Set(kv[i], kv[i+1])
	}
	return metadata.NewClientContext(ctx, md)
}
