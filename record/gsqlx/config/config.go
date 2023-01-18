package main

import (
	"fmt"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

type Student struct {
	Name        string
	FirstLetter string // 首字母
}

func main() {
	studentList := []*Student{
		{Name: "艾欧尼亚"}, {Name: "abc"}, {Name: "安全网"}, {Name: "茴香"},
		{Name: "回家"}, {Name: "直接"}, {Name: "不去玩"}, {Name: "bbd"},
		{Name: "123"}, {Name: "234"}, {Name: "！@#￥"}, {Name: "!@#"},
		{Name: ""},
	}
	for index := range studentList {
		v := studentList[index]
		// 获取首字母
		v.FirstLetter = GetFirstLetter(v.Name)
	}
	// 对首字母排序
	sortStudentList := firstLetterSort(studentList)
	for _, sortStudent := range sortStudentList {
		fmt.Printf("%+v\n", sortStudent)
	}
}

func GetFirstLetter(str string) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.FIRST_LETTER // 首字母风格，只返回拼音的首字母部分
	// 替换空格
	str = strings.Replace(str, " ", "", -1)
	// 处理空字符串
	if str == "" {
		return str
	}
	a := pinyin.Pinyin(str, args)
	if len(a) > 0 {
		return strings.ToUpper(a[0][0])
	} else {
		return strings.ToUpper(str[:1])
	}
}

func firstLetterSort(arr []*Student) []*Student {
	length := len(arr)
	for i := 0; i < length-1; i++ {
		flag := false
		for j := 0; j < length-1-i; j++ {
			if arr[j].FirstLetter > arr[j+1].FirstLetter {
				// 元素交换
				arr[j+1], arr[j] = arr[j], arr[j+1]
				flag = true
			}
		}
		if !flag {
			break
		}
	}
	return arr
}
