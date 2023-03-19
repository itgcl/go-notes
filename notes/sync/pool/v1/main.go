package main

import (
	"fmt"
	"sync"
)

type Student struct {
	Name string
	Age  int32
}

func main() {
	studentPool := sync.Pool{
		New: func() interface{} {
			return new(Student)
		},
	}
	studentPool.Put(&Student{
		Name: "a",
		Age:  1,
	})
	studentPool.Put(&Student{
		Name: "b",
		Age:  2,
	})
	studentPool.Put(&Student{
		Name: "c",
		Age:  3,
	})
	stu1 := studentPool.Get().(*Student)
	stu2 := studentPool.Get().(*Student)
	stu3 := studentPool.Get().(*Student)
	stu4 := studentPool.Get().(*Student)
	fmt.Println(stu1, stu2, stu3, stu4)
}
