FROM golang:1.22.10-alpine3.19 AS builder
# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 安装 golangci-lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

# 构建一个小的可执行文件
# RUN go build -o main .  # 如果你的项目需要的话

FROM alpine:3.19

# 安装必要的依赖 (musl libc)
RUN apk update && apk add --no-cache musl

# 从 builder 阶段复制 golangci-lint
COPY --from=builder /go/bin/golangci-lint /usr/local/bin/golangci-lint

# 从 builder 阶段复制 Go 运行时依赖
COPY --from=builder /usr/local/go/ /usr/local/go/

# 设置环境变量
ENV PATH="/usr/local/go/bin:${PATH}"

# 设置入口点
ENTRYPOINT ["golangci-lint"]
