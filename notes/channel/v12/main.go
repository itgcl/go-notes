package main

import (
	"fmt"
	"sync"
)

func WordCount(text string) {
	numberOfGroups := 4
	wordChannel := make(chan map[string]int)
	wg := new(sync.WaitGroup)
	wg.Add(numberOfGroups)

	for i := 0; i < numberOfGroups; i++ {
		go processToCounting("", wordChannel, wg)
	}

	wg.Wait()
	fmt.Print(<-wordChannel)

	finalMap := make(map[string]int)
	close(wordChannel)

	for i := 0; i < numberOfGroups; i++ {
		for k, v := range <-wordChannel {
			finalMap[k] += v
		}
	}
}

func processToCounting(textSlice []string, wordChannel chan map[string]int, wg *sync.WaitGroup) {
	freq := make(map[string]int)
	for _, v := range textSlice {
		freq[v]++
	}
	wg.Done()
	wordChannel <- freq
}
