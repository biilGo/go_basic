package main

import "fmt"

// 接口Shaper
type Shaper interface {
	Area() float32
}

// 结构体Square
type Square struct {
	side float32
}

// 接口的方法Area()
func (sq *Square) Area() float32 {
	return sq.side * sq.side
}

func main() {
	sq1 := new(Square) // Square实例
	sq1.side = 5

	// 结构体Square实现了接口Shaper
	var areaIntf Shaper
	areaIntf = sq1
	fmt.Printf("The square has area:%f\n", areaIntf.Area())
}
