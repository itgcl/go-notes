package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

func main() {
	err := UploadFilePairV2(context.Background(), []int64{1, 2, 3, 4, 5, 6})
	if err != nil {
		return
	}
}

func UploadFilePair(ctx context.Context, jobs []int64) error {
	outputCh := make(chan error, 1)
	defer func() {
		close(outputCh)
		fmt.Println("close over")
	}()
	for _, item := range jobs {
		item := item
		go func() {
			fmt.Println("start goroutine")
			outputCh <- UploadFile(item)
		}()
		select {
		case err := <-outputCh:
			if err != nil {
				fmt.Println("sftp upload file error")
				return err
			}
			fmt.Println("sftp upload success")
		// 超时时间半小时
		case <-time.After(time.Second * 3):
			log.Warn(fmt.Sprintf("sftp upload file timeout..."))
			return errors.New("FtpTimeout")
		}
	}
	return nil
}

func UploadFilePairV2(ctx context.Context, jobs []int64) error {
	for _, item := range jobs {
		err := UploadFile(item)
		if err != nil {
			fmt.Println("sftp upload file error")
			return err
		}
		fmt.Println("sftp upload success")
		return nil
	}
	return nil
}

func UploadFile(id int64) error {
	fmt.Println("req id: ", id)
	time.Sleep(time.Second * 5)
	fmt.Println("reply id: ", id)
	return nil
}
