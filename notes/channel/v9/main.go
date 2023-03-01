package main

import (
	"fmt"
	"sync"
)

type TestStruct struct {
	Test string
}

func main() {
	var stringsSlice = []string{"a", "b", "c"}
	channelsMap := make(map[string]chan *TestStruct)

	//for i := 1; i <= 10; i++ {
	//	stringsSlice = append(stringsSlice, fmt.Sprintf("value%d", i))
	//}

	var wg sync.WaitGroup

	wg.Add(len(stringsSlice))
	for _, value := range stringsSlice {
		channelsMap[value] = make(chan *TestStruct, 1)
		go watchChannel(value, channelsMap[value], &wg)
	}

	for _, value := range stringsSlice {
		for i := 0; i < 100; i++ {
			testStruct := new(TestStruct)
			testStruct.Test = fmt.Sprint("Hello! ", i)
			channelsMap[value] <- testStruct
		}
	}

	for _, ch := range channelsMap {
		close(ch)
	}

	wg.Wait()
	fmt.Println("Program ended")
}

func watchChannel(channelMapKey string, ch <-chan *TestStruct, wg *sync.WaitGroup) {
	defer wg.Done()
	for channelValue := range ch {
		fmt.Printf("Channel '%s' used. Passed value: '%s'\n", channelMapKey, channelValue.Test)
	}
}
