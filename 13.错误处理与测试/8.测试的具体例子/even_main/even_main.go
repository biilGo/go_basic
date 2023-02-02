package main

import (
	"fmt"
	"mywork/8.测试的具体例子/even"
)

func main() {
	for i := 0; i <= 100; i++ {
		fmt.Printf("is the integer %d even? %v\n", i, even.Even(i))
	}
}
