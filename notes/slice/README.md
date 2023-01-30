### 切片（slice）

> 切片和数组的区别

数组：定义的时候就要求说明长度，而且是固定的。

切片：长度和容量可以不固定。

#### 使用：

> 使用make定义切片时，长度大于容量会报错。

```go
package main

import "fmt"

func main(){
   data := make([]int, 11, 10)
   fmt.Println(data)
}
//# command-line-arguments
//.\test.go:6:14: len larger than cap in make([]int)

```

### 切片截取 

**新的切片是引用类型。修改新切片的值，之前的切片值也会修改。**

假设 **slice[1:3:8] **用 **[low:high:max]**表示

**max不能大于初始化切片的最大容量**。

不指定max默认使用被切片的最大容量。

len 为长度，计算方式  **len = high - low**

新切片的容量计算方式 **newSliceCap = max - low**

```go
package main

import (
	"fmt"
)

func main(){
	intList := []int{10,20,30,40,50,60,70,80,90,100}
    q := intList[1:3:8]  //从下标1（包含）截取到下标3（不包含）
    fmt.Println(q, len(q), cap(q)) [20,30]  len=2 (3-1) cap=7(8-1)
}

```

#### 切片多次截取

q的值是 [20 30]，w截取q[1:6]，q[0]=20 q[1]=30 

理论上是 [30]  为什么会是[30 40 50 60 70]呢

因为切片截取是引用类型，他们使用的是同一个内存空间。所以当w找下标6时q只有0和1，他会找到刚开始的intList进行获取。

也就是从30开始往后推 30 40 50 60 70。

```go
package main

import (
   "fmt"
)

func main(){
   intList := []int{10,20,30,40,50,60,70,80,90,100}
   q := intList[1:3]  //从下标1（包含）截取到下标3（不包含）[20 30]
   //此时在新的切片再进行截取
   w := q[1:6]
   fmt.Println(w) //[30 40 50 60 70]
}
```