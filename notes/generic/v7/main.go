package some

import (
	"fmt"
	"strings"
)

func Keys[T comparable](m map[T]interface{}) []T {
	keys := make([]T, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

func Split[T any](s, sep string, converter func(string) (T, error)) ([]T, error) {
	var data []T
	list := strings.Split(s, sep)
	for _, item := range list {
		value, err := converter(item)
		if err != nil {
			return nil, err
		}
		data = append(data, value)
	}
	return data, nil
}

func main() {
	keys := Keys(map[string]interface{}{"a": nil, "b": nil, "c": nil})
	fmt.Println(keys)
	keys2 := Keys(map[int]interface{}{4: nil, 3: nil, 1: nil})
	fmt.Println(keys2)

	originalMap := map[string]int{"a": 1, "b": 2, "c": 3}

	s := Values[string](originalMap)
	fmt.Println(s[0] == 1)
}

type Person struct {
	Name string
	Sort int
}
