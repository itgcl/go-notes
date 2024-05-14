package main

import "fmt"

type Stringer interface {
	String() string
}

func Stringify[T Stringer](s []T) (ret []string) {
	for _, v := range s {
		ret = append(ret, v.String())
	}
	return ret
}

// 使用的interface有方法就不能定义其他类型了
//func StringifyV2[T int | Stringer](s []T) (ret []string) {
//	for _, v := range s {
//		fmt.Println(v)
//	}
//	return nil
//}

type S string

func (s S) String() string {
	return string(s)
}

func main() {
	var s1 S = "q"
	var s2 S = "w"
	var s3 S = "e"
	var s4 S = "r"
	s5 := S("t")
	s6 := S("y")
	str := []Stringer{s1, s2, s3, s4, s5, s6}
	fmt.Println(Stringify(str))
}
