package main

import (
	"fmt"
	"sync"
)

var (
	once sync.Once
	Conf *Config
)

type Config struct {
	Mysql string
	Redis string
}

func coreInitConf() {
	fmt.Printf("init config...\n")
	Conf = &Config{
		Mysql: "aa",
		Redis: "bb",
	}
}

func InitConf(i int) {
	fmt.Printf("counter：%d...\n", i)
	once.Do(coreInitConf)
}

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			InitConf(i)
		}(i)
	}
	wg.Wait()
	fmt.Printf("config；%+v", Conf)
}
