// Copyright 2023 The go Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package main

import "fmt"

// Send the sequence 2,3,4... to channel `ch`\
func generate(ch chan int) {
	for i := 2; ; i++ {
		// send i to channel
		ch <- i
	}
}

// copy the values from channel in to channel out
// removing those divisible by prime

// 协程filter拷贝整数到输出通道，丢弃掉可以被prime整除的数字
// 然后每个prime又开启了一个新的协程，生成器和选择器并发请求
func filter(in, out chan int, prime int) {
	for {
		// receive value of new variable i from in
		i := <-in
		if i%prime != 0 {
			// send i to channel out
			out <- i
		}
	}
}

// the prime sieve:daisy-chain filter processes together
func main() {
	// create a new channel
	ch := make(chan int)

	// start generate() as a goroutine
	go generate(ch)

	for {
		prime := <-ch
		fmt.Print(prime, " ")
		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
	}
}
