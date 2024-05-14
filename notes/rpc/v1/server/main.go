package main

import (
	"log"
	"net"
	"net/rpc"
)

type HelloService struct{}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (h *HelloService) Hello(req string, reply *string) error {
	*reply = "hello: " + req
	return nil
}

func main() {
	err := rpc.RegisterName("HelloService", NewHelloService())
	if err != nil {
		log.Fatalf("register error: %v", err)
	}
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go rpc.ServeConn(conn)
	}
}
