package main

import (
	"fmt"
	"time"
)

func pump1(ch chan int) {
	for i := 0; ; i++ {
		ch <- i * 2
	}
}

func pump2(ch chan int) {
	for i := 0; ; i++ {
		ch <- i + 5
	}
}

func suck(ch1, ch2 chan int) {
	for {
		select {
		case v := <-ch1:
			fmt.Printf("received on channel 1:%d\n", v)
		case v := <-ch2:
			fmt.Printf("received on channel 2:%d\n", v)
		}
	}
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go pump1(ch1)
	go pump2(ch2)
	go suck(ch1, ch2)

	// 一秒内的输出非常惊人
	time.Sleep(1e9)
}
