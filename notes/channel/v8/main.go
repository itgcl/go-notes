package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Data struct {
	Value chan int
	Group chan int
}

func main() {
	group, _ := errgroup.WithContext(context.Background())
	archiverFileMap := map[int]Data{}
	wg := new(sync.WaitGroup)

	group.Go(func() error {
		for i := 0; i <= 4; i++ {
			wg.Add(1)
			archiverFileMap[i] = Data{make(chan int, 1), make(chan int, 1)}
			group.Go(func() error {
				s(archiverFileMap[i], wg)
				return nil
			})
		}
		wg.Wait()
		return nil
	})
	group.Go(func() error {
		for i := 1; i <= 100; i++ {
			v, exists := archiverFileMap[i]
			if !exists {
				continue
			}
			v.Value <- i
		}
		return nil
	})
	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}

func s(data Data, wg *sync.WaitGroup) {
	defer wg.Done()
	for g := range data.Group {
		v(g, data.Value)
	}
}

func v(groupID int, value <-chan int) {
	for v := range value {
		fmt.Println(groupID, v)
	}
}
