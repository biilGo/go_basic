package main

import (
	"fmt"
	"time"
)

var values = [5]int{10, 11, 12, 13, 14}

func main() {

	// version A
	for ix := range values {
		func() {
			fmt.Print(ix, " ")
		}() // 调用闭包打印每个索引值
	}
	fmt.Println()
	// versionA调用闭包5次,打印每个索引值

	// version B
	for ix := range values {
		go func() {
			fmt.Print(ix, " ")
		}()
	}
	fmt.Print()
	time.Sleep(5e9)
	/*
		versionB也做相同的事,但是通过协程调用每个闭包.按理说这将执行得更快,因为闭包是并发执行的.
		如果我们阻塞足够多的时间,让所有协程执行完毕,versionB的输出是`4 4 4 4 4`.
		versionB的循环中,ix变量实际上是一个单变量,表示每个数组元素的索引值,因为这些闭包都只绑定到一个变量,
		这是一个比较好的方式,当你运行这段代码时,你将看见每次循环都打印最后一个索引值4,而不是每个元素的索引值.
		因为协程可能在循环结束后还没有开始执行,而此时ix值是4.
	*/

	// version C 正确的处理方式
	for ix := range values {
		go func(ix interface{}) {
			fmt.Print(ix, " ")
		}(ix)
	}
	fmt.Println()
	time.Sleep(5e9)
	/*
		正确写法,调用每个闭包时将ix作为参数传递给闭包,ix在每次循环时都被重新赋值,并将每个协程的ix放置在栈中,
		所以当协程最终被执行时,每个索引值对协程都是可用的.
	*/

	// version D 输出值
	for ix := range values {
		val := values[ix]
		go func() {
			fmt.Print(val, " ")
		}()
	}
	time.Sleep(1e9)
	/*
		versionD中的变量声明是在循环体内部,所以在每次循环时,这些变量相互之间是不共享的,所以这些变量可以单独的被每个闭包使用
	*/

}

// 输出
/*
0 1 2 3 4
4 4 4 4 4
1 0 3 4 2
10 11 12 13 14
*/
