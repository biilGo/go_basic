package main

// 外层名字会覆盖内层名字，但是两者的内存空间都保留，这提供了一种重载字段或方法的方式
// 如果相同的名字在同一级别出现两次，如果这个名字被程序使用了，将会引发一个错误，没有办法来解决这种问题引起的二义性，必须有程序自己修正

type A struct{ a int }

type B struct{ a, b int }

type C struct {
	A
	B
}

type D struct {
	B
	b float32
}

// 使用d.b是没有问题的，它是float32，而不是B的b。如果想要内层的b可以通过d.B.b得到

func main() {

}
