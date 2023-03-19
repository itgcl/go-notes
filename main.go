package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Download(w http.ResponseWriter, r *http.Request) {
	filename := "/Users/guochenglong/Downloads/sqweqweqeq.zip"
	file, _ := os.OpenFile(filename, os.O_CREATE, 0777)
	defer file.Close()
	fileStat, _ := file.Stat()
	resp, err := http.Get("https://baidu.com")
	w.Header().Set("Content-Disposition", "attachment; filename=ww22w.zip")
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	file.Seek(0, 0)
	_, err = io.Copy(file, resp.Body)
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
	fmt.Println(333333333)
}

type Data struct {
	t  sql.NullTime
	t1 time.Time
	t2 *time.Time
}

func main() {
	defer func() {
		var err error
		fmt.Println("defer", err)
	}()
	err := demo()
	fmt.Println("demo", err)
}

func demo() error {
	return errors.New("aaa")
}
