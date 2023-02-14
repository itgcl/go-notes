package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

func main() {
	flag.Parse()

	http.HandleFunc("/", Download)
	err := http.ListenAndServe(":8081", nil)
	if nil != err {
		log.Fatalln("Get Dir Err", err.Error())
	}
}
