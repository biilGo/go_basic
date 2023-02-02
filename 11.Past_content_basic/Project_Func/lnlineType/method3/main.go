package main

import (
	"fmt"
	"math"
)

type Point struct {
	x, y float64
}

type NamedPoint struct {
	Point
	name string
}

func (p *Point) Abs() float64 {
	return math.Sqrt(p.x*p.x + p.y*p.y)
}

func (n *NamedPoint) Abs() float64 {
	return n.Point.Abs() * 100.
} // fmt.Println(n.Abs()) = 500

func main() {
	n := &NamedPoint{Point{3, 4}, "Pythagoras"}
	fmt.Println(n.Abs())
}
