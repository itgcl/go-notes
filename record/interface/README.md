### golang 函数不确定参数使用

> 使用...interface{} 默认转换为[]interface{}

> 两个不定参数的函数使用出现 [[xx,xx,xx]] 问题

```go
func test1(data ...interface{}){
	test2(data) // 打印 []interface {}{1, 2, 3}
}

//test2也接收不定参数
func test2(data ...interface{}){
	fmt.Println(data) //打印 []interface {}{[]interface {}{1, 2, 3}}
}
```

**解决:**

```go
func test1(data ...interface{}){
   test2(data...)   //data...
}

func test2(data ...interface{}){
   l.V(data)  //打印 []interface {}{1, 2, 3}
}
```

**函数不定参数必须是最后一个参数不然报错**

![image-20200429223424665](C:\Users\dell\AppData\Roaming\Typora\typora-user-images\image-20200429223424665.png)