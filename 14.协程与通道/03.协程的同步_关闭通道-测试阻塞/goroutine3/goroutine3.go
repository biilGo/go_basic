package main

import "fmt"

// 函数的最后,关闭了通道
func sendData(ch chan string) {
	ch <- "Washington"
	ch <- "Tripoli"
	ch <- "London"
	ch <- "Beijing"
	ch <- "Tokio"
	close(ch)
}

func getData(ch chan string) {
	// for循环getData中,在每次接收通道的数据之前都使用if !open来检测
	for {
		input, open := <-ch
		if !open {
			break
		}
		fmt.Printf("%s", input)
	}
}

func main() {
	ch := make(chan string)

	// sendData是协程
	go sendData(ch)

	// 和main在同一个线程里
	getData(ch)
}

// 使用for-range语句来读取通道是更好的办法,因为这回自动检测通道是否关闭
/*
for input := range ch {
      process(input)
}
*/
