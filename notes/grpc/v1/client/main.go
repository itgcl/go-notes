package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "go-notes/notes/grpc/v1/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

// serviceName grpc服务端名称
var serviceConfig = `{
	"loadBalancingPolicy": "round_robin",
	"healthCheckConfig": {
		"serviceName": "grpc.article.v1.Article"
	}
}`

func queryArticle(ctx context.Context, c pb.ArticleServiceClient) {
	r, err := c.QueryArticle(ctx, &pb.RequestQueryArticle{ArticleId: 1})
	if err != nil {
		fmt.Printf("queryArticle: %v, %s\n", r, err)
	} else {
		fmt.Printf("queryArticle: %s, %s\n", r.Author, r.Title)
	}
}

func main() {
	ctx := context.Background()
	r := manual.NewBuilderWithScheme("whatever")
	r.InitialState(resolver.State{
		Addresses: []resolver.Address{
			{Addr: "localhost:8181"},
			{Addr: "localhost:8182"},
		},
	})
	address := fmt.Sprintf("%s:///unused", r.Scheme())
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithResolvers(r),
		grpc.WithDefaultServiceConfig(serviceConfig),
	}
	conn, err := grpc.Dial(address, options...)
	if err != nil {
		log.Fatalf("grpc.Dial error: %v", err)
	}
	defer conn.Close()
	c := pb.NewArticleServiceClient(conn)

	for {
		queryArticle(ctx, c)
		time.Sleep(time.Second)
	}
}
