package v1

import (
	"fmt"
	"testing"

	"go.uber.org/goleak"
)

func generateInteger() func() int {
	ch := make(chan int)
	count := 0
	go func() {
		for {
			ch <- count
			count++
		}
	}()

	return func() int {
		return <-ch
	}
}

func TestGO(t *testing.T) {
	defer goleak.VerifyNone(t)
	generate := generateInteger()
	fmt.Println(generate()) // 0
	// fmt.Println(generate()) // 1
	// fmt.Println(generate()) // 2
}
