package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type ctxReader struct {
	ctx context.Context
	r   io.Reader
}

func (cr *ctxReader) Read(p []byte) (int, error) {
	select {
	case <-cr.ctx.Done():
		return 0, cr.ctx.Err() // 如果context已取消，返回错误
	default:
		fmt.Println(p)
		return cr.r.Read(p) // 否则，正常读取数据
	}
}

func copyWithContext(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	// 给定的 context 和源阅读器包装在自定义 ctxReader 中
	cr := &ctxReader{ctx, src}
	// 使用 io.Copy 像往常一样复制，但会通过 ctxReader 检测 context 的取消
	return io.Copy(dst, cr)
}

func main() {
	// 假设我们有两个文件 src.txt 和 dst.txt
	srcFile, err := os.Open("/Users/guochenglong/workspace/go-notes/ef079b2234d740d69096ef7f781b3b76.wav")
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create("aa.wav")
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()

	// 创建一个将在一段时间后取消的 context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// 开始复制并传递 context
	written, err := copyWithContext(ctx, dstFile, srcFile)
	if err != nil {
		fmt.Printf("Copy error: %v\n", err)
	} else {
		fmt.Printf("Copied %d bytes\n", written)
	}

	// 让 main 函数不立即退出，以便观察 context 处理
	time.Sleep(1 * time.Second)
}
