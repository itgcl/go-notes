package main

import (
	"io"
	"log"
	"net/http"
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
curl localhost:1234/ping -X POST \
    --data '{"method":"HelloService.Hello","params":["hello"],"id":0}'
*/

func main() {
	err := rpc.RegisterName("HelloService", NewHelloService())
	if err != nil {
		log.Fatalf("register error: %v", err)
	}
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}
		err := rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
		if err != nil {
			log.Println(err)
			return
		}
	})
	err = http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
