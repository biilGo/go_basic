// os包中有一个string类型的切片变量os.Args用来处理一些基本的命令行参数
// 它在程序启动后读取命令行输入的参数
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	who := "Alice"

	if len(os.Args) > 1 {
		// 命令行参数会放置在切片os.Args[]中，以空格分隔，从索引1开始
		// 函数strings.Join以空格为间隔连接这些参数
		who += strings.Join(os.Args[1:], " ")
	}
	fmt.Println("Good Morning", who)
}

// 在IDE或编译器中直接运行程序输出：Good Morning Alice
// 在命令行加入参数：./os_args John Bill Marc Luke得到输出：Good Morning Alice John Bill Marc Luke
