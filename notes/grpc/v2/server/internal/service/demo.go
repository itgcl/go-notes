package demo

import (
	"io"
	"log"

	"go-notes/notes/grpc/v2/proto"
	pb "go-notes/notes/grpc/v2/proto"
)

type Service struct {
	pb.UnimplementedStreamDemoServiceServer
	port int
}

func NewService(port int) *Service {
	return &Service{
		port: port,
	}
}

func (s *Service) InputStream(req proto.StreamDemoService_InputStreamServer) error {
	var sum int64
	for {
		param, err := req.Recv()
		if err == io.EOF {
			log.Printf("input stream recv EOF... %s\n", err)
			break
		}
		if err != nil {
			log.Printf("reve error: %s\n", err)
			return err
		}
		sum += param.Value
	}
	log.Printf("input stream for break...")
	if err := req.SendAndClose(&proto.DataReply{Data: sum}); err != nil {
		log.Printf("send and cloer error: %s\n", err)
		return err
	}
	return nil
}

func (s *Service) OutputStream(request *proto.OutputStreamRequest, ss proto.StreamDemoService_OutputStreamServer) error {
	if err := ss.Send(&pb.DataReply{Data: request.X + request.Y}); err != nil {
		log.Printf("outpu stream send error: %s\n", err)
		return err
	}
	if err := ss.Send(&pb.DataReply{Data: request.X - request.Y}); err != nil {
		log.Printf("outpu stream send error: %s\n", err)
		return err
	}
	if err := ss.Send(&pb.DataReply{Data: request.X * request.Y}); err != nil {
		log.Printf("outpu stream send error: %s\n", err)
		return err
	}
	if err := ss.Send(&pb.DataReply{Data: request.X / request.Y}); err != nil {
		log.Printf("outpu stream send error: %s\n", err)
		return err
	}
	return nil
}

func (s *Service) BidirectionalStream(server proto.StreamDemoService_BidirectionalStreamServer) error {
	for {
		request, err := server.Recv()
		if err == io.EOF {
			log.Printf("bidirect stream EOF:%s\n", err)
			return nil
		}
		if err != nil {
			log.Printf("bidirect stream error:%s\n", err)
			return err
		}
		if err := server.Send(&proto.DataReply{Data: request.X + request.Y}); err != nil {
			log.Printf("bidirect send error: %s\n", err)
			return err
		}
	}
}
