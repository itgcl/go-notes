package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// c581e5fa56b3000d96b0df458e787651-2

const blockSize = 1024 * 1024 * 16 // 16MB
// GetFileETag 计算文件的 ETag 值，基于文件内容的 MD5 哈希值
func GetFileETag(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Save the file stat.
	fileStat, err := file.Stat()
	if err != nil {
		return "", err
	}

	if fileStat.Size() > blockSize {
		return fileTagMultipart(file)
	} else {
		return fileTag(file)
	}
}

func fileTagMultipart(reader io.Reader) (string, error) {
	// 创建一个 MD5 哈希对象
	combinedHash := md5.New()
	buf := make([]byte, blockSize)
	shards := 0
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			// 每个块计算 MD5 哈希
			hash := md5.New()
			hash.Write(buf[:n]) // 将当前块的数据写入哈希对象
			blockHash := hash.Sum(nil)

			// 将当前块的哈希合并到最终的合并哈希
			combinedHash.Write(blockHash)
			shards += 1
		}

		if err == io.EOF {
			// 文件读取完毕
			break
		}
		if err != nil {
			// 如果遇到其他错误
			return "", err
		}
	}
	// 计算哈希值并转换为十六进制字符串
	return fmt.Sprintf("\"%s-%d\"", hex.EncodeToString(combinedHash.Sum(nil)), shards), nil
}

func fileTag(reader io.Reader) (string, error) {
	// 创建一个 MD5 哈希对象
	hash := md5.New()

	// 使用 io.Copy 将文件内容写入哈希对象
	_, err := io.Copy(hash, reader)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func main() {
	var b []byte
	b = append(b, 'a', 'b', 'c', 'd', 'e', 'f', 'g', 100, 254)
	fmt.Println(b)
	var r []rune
	r = append(r, 97, 256, 30000, '你', '好', 'a')
	fmt.Println(r)
	return
	// 50b08b4f016b8a61e5aac071ccd9cb38
	// 示例文件路径
	filePath := "/Users/guochenglong/Downloads/016.flac"
	// 81b7c8c85c5e5a1ed19b4110ef5e56ba
	etag, err := GetFileETag(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("ETag:", etag)
}
