### 今日会议

#### 接口

通过第三方库文件库read，write操作引出接口

```go
package demo_io_interface

import (
	gconv "github.com/og/x/conv"
	glist "github.com/og/x/list"
	l "github.com/og/x/log"
	gtime "github.com/og/x/time"
	"testing"
	"time"
)
type Writer interface {
	Write(s string)
}
func downloadSome(w Writer){
	glist.Run(10, func(i int) (_break bool) {
		w.Write(gconv.IntString(i))
		return
	})
}
func saveSome(s string) {
	l.V("save:", s)
}
type Memory struct {
	v []string
}
func (m *Memory)Write(s string) {
	m.v = append(m.v, s)
}
func TestBasic(t *testing.T) {
	 memory := Memory{}
	 downloadSome(&memory)
	 l.V(memory)

}
```

当函数downloadSome()传入参数是接口类型。

在结构体类型和接口类型不一致，只要结构体实现接口，就是接口类型，结构体就可以传入这个函数downloadSome()中。

#### 问题

> 1.  函数downloadSome(&memory) 为什么要传指针，不传引用类型就编辑器报错?

当接口有方法的时候，肯定是需要使用接口里的方法。使用方法可能会修改自身数据，所以需要传指针。



> 2. 当  Writer  是一个空接口，使用函数downloadSome(memory)为什么可以不用传指针了？

因为空接口没有方法，也就不需要做什么操作，所以可以不传指针。



**在写代码的时候遇到这种问题，都是修改好不报错就行了，直接跳过了，并没有脑子里想**

为什么要传指针？ 为什么空接口不需要指针？为什么go这样设计？



#### 需要提高的点

**通过上面的问题，学习到了就算是语法等方面问题，也要多想进行累计，多想了以后，在后面遇到问题就可以发散思维想到解决方法或者别人想不到的层面。例如（上面的接口指针问题）**



> 接口没有属性怎么解决属性的问题？

一般对属性的操作基本上是读取和写入修改。通过两个方法实现属性问题。

```go
type Order interface {
   GetAmount()	//获取
   SetAmount(amount float64) //设置
}
```



> 反射获取结构体标签如何把`json：“xxx”`变成[]string？

```go
type Test struct {
    A string `json:test("a","b")`
   Number string `json:test(10)`
    Abcde string `json:test("abcde")`
    AB string `json:test("a,b")`
}
```

例如：获取标签 `json:test("a","b")` 两个参数，反射获取到 "a“,”b"，然后strings.Split( `"a","b"`, `,`)以逗号分割[]string{"a", "b"}。

但是`json:test("a,b")`这种特殊格式的不行，反射获取"a,b"是一个参数，再以逗号分割会两个参数。是错误的。[]string{"a","b"}。

平时多想累计

使用json_decode 先把"a,b"字符串拼接成 [  "a,b"  ]，在json_decode获得数组[]string{"a,b"}。

**如果实在想不到，可以把这种特殊用法记下来，以后这种问题可以使用。**