package reflect_test

import (
	"fmt"
	"reflect"
	"testing"
)

type Student struct {
	Name        string
	Age         int
	Score       int
	Address     string
	IsGood      bool
	PocketMoney float64
}

func TestReflect(t *testing.T) {
	{
		student := Student{
			Name:        "张三",
			Age:         20,
			Score:       99,
			Address:     "地球村",
			IsGood:      true,
			PocketMoney: 20.20,
		}
		reflectValue := reflect.TypeOf(student)
		for i := 0; i < reflectValue.NumField(); i++ {
			key := reflectValue.Field(0)
			fmt.Println(key.Name, key.Type)
			// fmt.Printf("key:%s, value:%s", reflectValue.Field(i))
			// fmt.Printf("%d: %")
		}
	}
}
