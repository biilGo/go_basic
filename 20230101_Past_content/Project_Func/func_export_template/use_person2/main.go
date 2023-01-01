package main

import (
	"fmt"
	person "mywork/Project_Func/func_export_template/person2"
)

func main() {
	p := new(person.Person)
	// p.firstName undefind
	// cannot refer to unexported field or method firstName
	// p.firstName = "Eric"
	p.SetFirstName("Eric")
	fmt.Println(p.FirstName())
}
