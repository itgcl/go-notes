package main

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"
)

type TestStruct struct {
	Test string
}
type Data struct {
	Test  chan *TestStruct
	Group string
}

var count int

func main() {
	stringsSlice := []string{"a", "b", "c", "d", "a", "b"}
	channelsMap := make(map[string]Data)
	group, _ := errgroup.WithContext(context.Background())
	group.Go(func() error {
		defer func() {
			for _, v := range channelsMap {
				close(v.Test)
			}
		}()
		for _, value := range stringsSlice {
			// wg.Add(1)
			func() {
				_, exists := channelsMap[value]
				if !exists {
					channelsMap[value] = Data{
						Test:  make(chan *TestStruct, 1),
						Group: value,
					}
					group.Go(func() error {
						watchChannel(channelsMap[value].Group, channelsMap[value].Test)
						return nil
					})
				}
				for i := 0; i < 10000; i++ {
					testStruct := new(TestStruct)
					testStruct.Test = fmt.Sprint("Hello! ", i)
					channelsMap[value].Test <- testStruct
				}
			}()
		}
		return nil
	})
	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)
	fmt.Println("Program ended")
}

func watchChannel(group string, ch <-chan *TestStruct) {
	for channelValue := range ch {
		//_ = channelValue
		fmt.Printf("Channel '%s' used. Passed value: '%s'\n", group, channelValue.Test)
		count++

	}
}
