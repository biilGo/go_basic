package main

import (
	"fmt"
	"time"
)

// main(),longWait(),shortWait()三个函数作为独立的处理单元按顺序启动,然后开始并行运行

func longWait() {

	// 每一个函数都在运行的开始和结束阶段输出了消息

	fmt.Println("Beginning longWait()")

	// time包sleep函数按照指定的时间来暂停函数或协程的执行
	// 单位是纳秒ns(符号1e9表示1乘10的9次方,e=指数)
	time.Sleep(5 * 1e9) // sleep for 5 seconds

	fmt.Println("End of longWait()")
}

func shortWait() {

	// 每一个函数都在运行的开始和结束阶段输出了消息

	fmt.Println("Beginning shortWait()")

	// time包sleep函数按照指定的时间来暂停函数或协程的执行
	// 单位是纳秒ns(符号1e9表示1乘10的9次方,e=指数)
	time.Sleep(2 * 1e9) //sleep for 2 seconds

	fmt.Println("End of shortWait()")
}

func main() {
	fmt.Println("In main()")

	// 以并行的方式执行
	go longWait()

	// 以并行的方式执行
	go shortWait()

	// sleep works with a Duration in nanoseconds(ns)!
	// main函数暂停10s从而确定它会在另外两个协程之后结束
	// 如果不这样,main会提前结束,longWait则无法完成
	fmt.Println("About to sleep in main()")
	time.Sleep(10 * 1e9)
	// 如果不在main中等待,协程会随着程序的结束而消亡

	// 当main函数返回的时候,程序退出,它不会等待任何其他非main协程的结束
	fmt.Println("At the end of main()")
}
