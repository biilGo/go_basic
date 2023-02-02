package main

import (
	"fmt"
	"os"
	"text/template"
)

type Person struct {
	Name                string
	nonExportedAgeField string
}

func main() {
	t := template.New("hello")

	// 数据结构中包含一个未导出的字段,当我们尝试把它整合到类似这样的定义字符串:
	// t, _ = t.Parse("your age is {{.nonExportedAgeField}}!")会产生错误
	t, _ = t.Parse("hello {{.Name}}!")

	p := Person{Name: "Mary", nonExportedAgeField: "31"}

	// 如果只是想把Execute()方法的第二个参数用于替换使用{{.}}
	// 当再浏览器环境中进行这些步骤,应首先使用html过滤器来过滤内容;如{{html .}}或者对FieldName过滤:{{.FieldName | html}}
	// |html时请求模板引擎再输出FieldName的结果前把值传递给html格式化器，它会执行html字符转义，避免用户输入数据破坏html文档结构
	if err := t.Execute(os.Stdout, p); err != nil {
		fmt.Println("There was an error:", err.Error())
	}
}
