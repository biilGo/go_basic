// 程序展示启动巨量的go协程是多么容易

// 这些协程已经全部在main函数中for循环里启动

// 当循环完成,一个0被写入到最右边的通道里

// 展示如何通过flag.Int来解析命令行中的参数以指定协程数量

package main

import (
	"flag"
	"fmt"
)

var ngoroutine = flag.Int("n", 1000000, "how many goroutines")

func f(left, right chan int) {
	left <- 1 + <-right
}

func main() {
	flag.Parse()

	leftmost := make(chan int)

	var left, right chan int = nil, leftmost

	for i := 0; i < *ngoroutine; i++ {
		left, right = right, make(chan int)

		// 因为没有发送者一直处于等待状态
		go f(left, right)
	}

	// bang!
	// 主线程,<-0,right不是最初循环的那个right,而是最终循环的right
	// 主线程执行时,类似于递归函数在最内层产生返回值一般
	right <- 0

	// wait for completion
	x := <-leftmost

	fmt.Println(x)
}
