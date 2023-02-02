package main

import (
	"fmt"
	"testing"
)

func BenchmarkChannelSync(b *testing.B) {
	ch := make(chan int)
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- i
		}
		close(ch)
	}()
	for range ch {
	}
}

func BenchmarkChannelBuffered(b *testing.B) {
	ch := make(chan int, 128)
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- i
		}
		close(ch)
	}()
	for range ch {
	}
}

// N的值将由gotest来判断并取得一个足够大的数字,以获得合理的基准测试结果
// 当然同样的基准测试方法也适用于普通函数

func main() {
	// 通过testing.Benchmark调用N次
	// Benchmark有一个String()方法来输出其结果

	fmt.Println("sync", testing.Benchmark(BenchmarkChannelSync).String())

	fmt.Println("buffered", testing.Benchmark(BenchmarkChannelBuffered).String())
}
