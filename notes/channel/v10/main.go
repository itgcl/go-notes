package main

import (
	"fmt"
	"sync"
)

type TestStruct struct {
	Test string
}
type Data struct {
	Test  chan *TestStruct
	Group string
}

func main() {
	var stringsSlice = []string{"a", "b", "c"}
	channelsMap := make(map[string]Data)
	var wg sync.WaitGroup
	wg.Add(len(stringsSlice))
	for _, value := range stringsSlice {
		//channelsMap[value] = make(chan *TestStruct, 1)
		channelsMap[value] = Data{
			Test:  make(chan *TestStruct, 1),
			Group: value,
		}
		go watchChannel(channelsMap[value].Group, channelsMap[value].Test, &wg)
	}

	for _, value := range stringsSlice {
		for i := 0; i < 100; i++ {
			testStruct := new(TestStruct)
			testStruct.Test = fmt.Sprint("Hello! ", i)
			channelsMap[value].Test <- testStruct
		}
	}

	for _, ch := range channelsMap {
		close(ch.Test)
	}

	wg.Wait()
	fmt.Println("Program ended")
}

func watchChannel(group string, ch <-chan *TestStruct, wg *sync.WaitGroup) {
	defer wg.Done()
	for channelValue := range ch {
		fmt.Printf("Channel '%s' used. Passed value: '%s'\n", group, channelValue.Test)
	}
}
