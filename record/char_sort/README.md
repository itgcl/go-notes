## 业务场景
需求：对用户名称按首字母进行排序展示。

**会遇到的问题**
>如果用户名称全部是英文，实现是比较简单，有中文名称的情况进行排序，实现起来比较复杂。

### 排序方式
目前有2种排序方式。`（这里是从前到后的顺序）`

1. 数字 < 英文 < 中文
2. 数字 < 中英文混合

## examples
### 排序方式一
这里只提一下数据库查询时按第一种方式排序，修改字段字符格式为gbk。
```sql
select name from student order by CONVERT( name USING gbk ) COLLATE gbk_chinese_ci ASC
```
### 排序方式二

*实现思路*
>获取第一个中文，把中文转成拼音，取拼音第一个字符。

安装第三方包
```go
go get github.com/mozillazg/go-pinyin
```

结构体
```go 
type Student struct {  
   Name        string  
   FirstLetter string  // 首字母
}
```
获取首字母
```go
func GetFirstLetter(str string) string {  
   args := pinyin.NewArgs()  
   args.Style = pinyin.FIRST_LETTER  
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
```

完整代码
```go
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
   sortStudentList := FirstLetterSort(studentList)  
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
  
func FirstLetterSort(arr []*Student) []*Student {  
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
```
