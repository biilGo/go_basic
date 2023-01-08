package main

import (
	"fmt"
	"time"
)

func sendData(ch chan string) {
	ch <- "Washington"
	ch <- "Tripoli"
	ch <- "London"
	ch <- "Beijing"
	ch <- "Tokoy"
}

// getData使用了无限循环,它随着sendData的发送完成和ch变空也结束了
func getData(ch chan string) {
	var input string

	// time.Sleep(2e9)

	for {
		input = <-ch
		fmt.Printf("%s", input)
	}
}

// main函数种启动了两个协程,sendData通过通道ch发送了5个字符串,getData按顺序接收它们并打印出来

// 如果2个协程需要通信,你必须给他们同一个通道作为参数才行
func main() {
	ch := make(chan string)

	// 如果我们移除一个或所有go关键字,程序无法运行且抛出panic
	/*
		Program exited with code -2147483645: panic: all goroutines are asleep-deadlock!
	*/
	// 运行时会检查所有的协程是否在等待者什么东西,这意味着程序将无法继续执行,这时死锁的一种形式,而运行时可以为我们检测到这种情况
	go sendData(ch)
	go getData(ch)

	// main等待1秒让两个协程完成,如果不这样,sendData就没有机会输出
	time.Sleep(1e9)
}

// 不要使用打印状态来表明通道的发送和接收顺序:由于打印状态和通道实际发生读写的时间延迟和真实发生的顺序不同
