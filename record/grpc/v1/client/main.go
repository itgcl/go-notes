package main

import (
	"context"
	"fmt"
	pb "go-notes/record/grpc/v1/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
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
		fmt.Printf("queryArticle: %s\n", r.Author)
	}
}

func main() {
	var (
		ctx     = context.Background()
		options = []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithDefaultServiceConfig(serviceConfig),
		}
	)
	conn, err := grpc.Dial("127.0.0.1:8181", options...)
	if err != nil {
		log.Fatalf("grpc.Dial error: %v", err)
	}
	defer conn.Close()
	c := pb.NewArticleServiceClient(conn)

	//healthClient := healthpb.NewHealthClient(conn)
	//ir := &healthpb.HealthCheckRequest{
	//	Service: "aaaaaaa",
	//}

	for {
		queryArticle(ctx, c)
		//healthCheckResponse, err := healthClient.Check(context.Background(), ir)
		//fmt.Println(healthCheckResponse, err)
		time.Sleep(time.Second)
	}
}
   