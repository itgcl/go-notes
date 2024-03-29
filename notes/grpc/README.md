# grpc
grpc服务简单实现。

## protobuf生成go文件
>cd record/grpc目录

执行protoc
```go
protoc --go_out=. --go_opt=paths=source_relative \                                                                                                                      ✔  11:07:53 AM  
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./v1/proto/*.proto
```
**protoc如果不存在的话，查看安装教程`https://grpc.io/docs/protoc-installation/`**

## 健康检查

健康检查的作用：https://github.com/grpc/grpc/blob/master/doc/health-checking.md

参考官方示例：https://github.com/grpc/grpc-go/tree/master/examples/features/health

### 项目启动
#### server
```go
go run server/main.go -port=8181 -sleep=5s
go run server/main.go -port=8182 -sleep=10s
```
#### client
```go
go run client/main.go
```

## metadata


## gateway

查看v4目录

grpc-gateway github: `https://github.com/grpc-ecosystem/grpc-gateway`

项目根目录执行下面代码
### 生成go文件
`
protoc -I . -I=./third_party \
--go_out . --go_opt paths=source_relative \
--go-grpc_out . --go-grpc_opt paths=source_relative \
./notes/grpc/v4/proto/hello.proto
`
### 生成gateway源文件
`protoc -I . --proto_path=./third_party --grpc-gateway_out . \
--grpc-gateway_opt logtostderr=true \
--grpc-gateway_opt paths=source_relative \
./notes/grpc/v4/proto/hello.proto`