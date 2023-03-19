package service

import (
	"context"
	"fmt"
	"go-notes/notes/grpc/v4/proto"
)

type HelloService struct {
	proto.UnimplementedHelloServiceServer
}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (s *HelloService) SayHello(ctx context.Context, req *proto.SayHelloRequest) (*proto.SayHelloReply, error) {
	return &proto.SayHelloReply{Data: fmt.Sprintf("hello: %s", req.Name)}, nil
}
