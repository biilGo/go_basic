// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

// send the sequence 2,3,4,...to returned channel
func generate() chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			ch <- i
		}
	}()
	return ch
}

// filter out input values divisible by prime,send rest to returned channel
func filter(in chan int, prime int) chan int {
	out := make(chan int)
	go func() {
		for {
			if i := <-in; i%prime != 0 {
				out <- i
			}
		}
	}()
	return out
}

// 函数sieve,generate,filter都是工厂,它们创建通道并返回,而且使用了协程的lambda函数使得main函数短小清晰
func sieve() chan int {
	out := make(chan int)

	go func() {
		ch := generate()
		for {
			prime := <-ch
			ch = filter(ch, prime)
			out <- prime
		}
	}()
	return out
}

// 调用sieve返回包含素数的通道,然后通过`println`打印出来
func main() {
	primes := sieve()
	for {
		fmt.Println(<-primes)
	}
}
