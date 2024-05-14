package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct{}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (h *HelloService) Hello(req string, reply *string) error {
	*reply = "hello: " + req
	return nil
}

/*
启动tcp服务 nc -l 1234
-----------
echo -e '{"method":"HelloService.Hello","params":["hello"],"id":1}' | nc localhost 1234
*/

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
		fmt.Println(33)
		conn, err := lis.Accept()
		fmt.Println(22)
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
