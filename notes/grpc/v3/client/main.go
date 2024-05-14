package main

import (
	"context"
	"fmt"
	"log"

	pb "go-notes/notes/grpc/v1/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
)

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
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}
	conn, err := grpc.Dial(":8081", options...)
	if err != nil {
		log.Fatalf("grpc.Dial error: %v", err)
	}
	defer conn.Close()
	c := pb.NewArticleServiceClient(conn)

	queryArticle(ctx, c)
}
