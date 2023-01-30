package main

import (
	"context"
	"fmt"
	pb "go-notes/notes/grpc/v2/proto"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
)

func main() {
	var (
		ctx = context.Background()
	)
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}
	conn, err := grpc.Dial("localhost:8181", options...)
	if err != nil {
		log.Fatalf("grpc.Dial error: %v", err)
	}
	defer conn.Close()
	sc := pb.NewStreamDemoServiceClient(conn)
	c := NewClient(sc)
	// -------请求测试---------
	{
		fmt.Println("------------------------")
		// 输入流
		reply, err := c.InputStream(ctx)
		if err != nil {
			return
		}
		log.Printf("input stream reply: %v\n", reply.Data)
		fmt.Println("------------------------")
		// 输出流
		_, err = c.OutputStream(ctx)
		if err != nil {
			return
		}
		fmt.Println("------------------------")
		// 双向流
		err = c.BidirectionalStream(ctx)
		fmt.Println(err)
	}

}

type Client struct {
	sc pb.StreamDemoServiceClient
}

func NewClient(sc pb.StreamDemoServiceClient) *Client {
	return &Client{sc: sc}
}

// InputStream 输入流
func (c *Client) InputStream(ctx context.Context) (*pb.DataReply, error) {
	var list = []*pb.InputStreamRequest{
		{Value: 1},
		{Value: 2},
		{Value: 3},
		{Value: 4},
		{Value: 5},
	}
	streamClient, err := c.sc.InputStream(ctx)
	if err != nil {
		log.Printf("input stream client error: %s\n", err)
		return nil, err
	}
	for _, req := range list {
		if err := streamClient.Send(req); err != nil {
			log.Printf("input stream send error: %s\n", err)
			return nil, err
		}
	}
	resp, err := streamClient.CloseAndRecv()
	if err != nil {
		log.Printf("input stream close and recv error: %s\n", err)
		return nil, err
	}
	return &pb.DataReply{Data: resp.Data}, nil
}

// OutputStream 输出流
func (c *Client) OutputStream(ctx context.Context) (*pb.DataReply, error) {
	// 设置超时时间
	newCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	sc, err := c.sc.OutputStream(newCtx, &pb.OutputStreamRequest{X: 10, Y: 2})
	if err != nil {
		log.Printf("output stream error: %s\n", err)
		return nil, err
	}
	var sum int64
	for {
		resp, err := sc.Recv()
		if err == io.EOF {
			log.Printf("output stream recv EOF...")
			break
		}
		if err != nil {
			log.Printf("output stream recv error: %s\n", err)
			return nil, err
		}
		sum += resp.Data
		log.Printf("output stream reply:%d\n", resp.Data)
	}
	return &pb.DataReply{Data: sum}, nil
}

// BidirectionalStream 双向流
func (c *Client) BidirectionalStream(ctx context.Context) error {
	sc, err := c.sc.BidirectionalStream(ctx)
	if err != nil {
		log.Printf("bidirectional stream error:%s\n", err)
		return err
	}
	ch := make(chan int64, 10)
	go func() {
		for {
			resp, err := sc.Recv()
			if err == io.EOF {
				log.Printf("bidirectional recv EOF..\n")
				close(ch)
				return
			}
			if err != nil {
				log.Printf("bidirectional  recv error:%s\n", err)
				return
			}
			ch <- resp.Data
		}
	}()
	// 发送消息
	for i := 1; i <= 10; i++ {
		v := int64(i)
		if err := sc.Send(&pb.BidirectionalStreamRequest{
			X: v,
			Y: v,
		}); err != nil {
			log.Printf("bidirectional stream error:%s\n", err)
			return err
		}
		log.Printf("bidirectional send success: %d\n", i)
	}
	// 关闭stream，会eof，让channel关闭
	if err := sc.CloseSend(); err != nil {
		log.Printf("bidirectional close send error:%s\n", err)
	}
	for v := range ch {
		fmt.Printf("bidirectional resp: %d\n", v)
	}
	fmt.Println("bidirectional over")
	return nil
}
