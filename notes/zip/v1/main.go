package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mholt/archiver/v4"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := example(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func example(ctx context.Context) error {
	var (
		archiverFileCh = make(chan archiver.File, 1)
		zip            = new(archiver.Zip)
		gz             = new(archiver.Gz)
	)
	// 压缩文件
	zipFile, err := os.Create(fmt.Sprintf("/tmp/%s.zip", time.Now().Format("2006-01-02-15-04-05")))
	if err != nil {
		return err
	}
	defer func() {
		// 关闭文件句柄
		if err := zipFile.Close(); err != nil {
			log.Printf("zip file close error: %v\n", err)
		}
		//// 删除本地文件
		//if err := os.RemoveAll(zipFile.Name()); err != nil {
		//	log.Printf("zip file remove error: %v\n", err)
		//}
	}()
	// 压缩方式
	compress, err := gz.OpenWriter(zipFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := compress.Close(); err != nil {
			log.Printf("compress close error: %v\n", err)
		}
	}()
	group, ctx := errgroup.WithContext(ctx)
	// 解析为archiver.File，写入channel
	group.Go(func() error {
		defer close(archiverFileCh)
		list := []string{
			"576bc1be4c9b48acab543492978b5905.wav",
			"12879445.flac",
			"ATR981156162.mp3",
		}
		for _, item := range list {
			archiverFile, err := getArchiverFile(filepath.Join("/Users/guochenglong/Downloads/", item))
			if err != nil {
				return err
			}
			archiverFileCh <- *archiverFile
		}
		return nil
	})
	group.Go(func() error {
		// 写入文件
		if err := zip.ArchiveAsync(ctx, compress, archiverFileCh); err != nil {
			return err
		}
		// 这个必须关闭，不然文件内容没有写完
		_ = compress.Close()
		return nil
	})
	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func getArchiverFile(filePath string) (*archiver.File, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	return &archiver.File{
		FileInfo:      info,
		NameInArchive: fmt.Sprintf("/musics/%s", filepath.Base(filePath)),
		LinkTarget:    "",
		Open: func() (io.ReadCloser, error) {
			return os.Open(filePath)
		},
	}, nil
}
