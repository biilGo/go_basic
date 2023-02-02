// flag包有一个扩展功能用来解析命令行选项

// 通常被用来替换基本常量

// flag包中有一个Flag被定义为一个含有如下字段的结构体
// type Flag struct {
// 	Name     string
// 	Usage    string
// 	Value    Value
// 	DefValue string
// }

package main

import (
	"flag"
	"os"
)

// flag定义其他类型可以使用flag.Int()\flag.Float64()\flag.String()

// flag.Bool() 定义了一个默认值是 false 的 flag

// 当在命令行出现第一个参数(n)，flag被设置成true(NewLine是*bool类型)
var NewLine = flag.Bool("n", false, "print newline")

// 对于flag.Bool可以设置布尔型flag来测试你的代码

// var processedFlag = flag.Bool("proc", false, "nothing processed yet")

// 用如下代码来测试：
// if *processedFlag { // found flag -proc
//     r = process()
// }

const (
	Space   = " "
	Newline = "\n"
)

func main() {

	// flag.VisitAll(fn func(*Flag))是另外一个有用的功能

	// 按照字典顺序遍历flag并且对每个标签调用fn

	// 打印flag的使用帮助信息
	flag.PrintDefaults()

	// 扫描参数列表并设置flag，flag.Arg(i)表示第i个参数

	// Parse()之后，flag.Arg(i)全部可用，flag.Arg(0)就是第一个真是的flag而不是像os.Args(0)放置程序的名字
	flag.Parse()

	var s string = ""

	// flag.Narg()返回参数的数量，解析后flag或常量就可用了
	for i := 0; i < flag.NArg(); i++ {
		if i > 0 {
			s += " "
			// flag被解引用到*NewLine所以当值是treu时将添加一个NewLine("\n")
			if *NewLine {
				s += Newline
			}
		}
		s += flag.Arg(i)
	}
	os.Stdout.WriteString(s)
}
