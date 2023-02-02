package main

import (
	"errors"
	"fmt"
)

// 可以把它用于计算平方根函数的参数测试：
func Sqrt(f float64) (float64, error) {
	if f < 0 {
		return 0, errors.New("math - square root if negative number")
	}

	return Sqrt(1)
	// implementation of sqrt
}

// var errNotFound error = errors.New("Not found error")

func main() {
	// fmt.Printf("error:%v", errNotFound)

	// 可以像下面这样调用Sqrt函数
	if f, err := Sqrt(-1); err != nil {
		fmt.Printf("Error:%s\n", err)
	} else {
		fmt.Println(f)
	}
}
