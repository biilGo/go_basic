package main

import (
	"fmt"
	"time"
)

// pump函数为通道提供数值,也被叫做生产者
func pump(ch chan int) {
	for i := 0; ; i++ {
		ch <- i
	}
}

// 为通道接触阻塞定义了suck函数在无限循环中读取通道
func suck(ch chan int) {
	for {
		fmt.Println(<-ch)
	}
}

func main() {
	ch1 := make(chan int)
	go pump(ch1) // pump hangs
	go suck(ch1)
	time.Sleep(1e9)
	fmt.Println(<-ch1)
}
