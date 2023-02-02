package main

import "fmt"

type specialString string

var whatIsThis specialString = "hello"

func TypeSwitch() {
	testFunc := func(any interface{}) {
		switch v := any.(type) {
		case bool:
			fmt.Printf("any %v is a bool type", v)
		case int:
			fmt.Printf("any %v is an int type", v)
		case float32:
			fmt.Printf("any %v is a float32 type", v)
		case string:
			fmt.Printf("any %v is string type", v)
		case specialString:
			fmt.Printf("any %v is special String", v)
		default:
			fmt.Printf("unknow type")
		}
	}

	testFunc(whatIsThis)
}

func main() {
	TypeSwitch()
}
