package main

import (
	"fmt"
	"os"
)

var user = os.Getenv("USER")

func check() {
	if user == "" {
		panic("unknown user:no value if $USER")
	}
}

func main() {
	fmt.Println("Starting the program")
	panic("A severe error occurrede:stopping the program")
	fmt.Println("Ending the program")

	check()
}
