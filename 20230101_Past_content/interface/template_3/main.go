package main

import "fmt"

type stockPosition struct {
	ticker     string
	sharePrice float32
	count      float32
}

// method ti determine the value of a stock posoition
func (s stockPosition) getValue() float32 {
	return s.sharePrice * s.count
}

type car struct {
	make  string
	model string
	price float32
}

// method to determine the value of a car
func (c car) getValue() float32 {
	return c.price
}

// contract that defines different things that have value
type valueable interface {
	getValue() float32
}

func showValue(asset valueable) {
	fmt.Printf("Value of the asset is %f\n", asset.getValue())
}

func main() {
	var o valueable = stockPosition{"GOOD", 577.20, 4}
	showValue(o)
	o = car{"BMW", "M3", 66500}
	showValue(o)
}
