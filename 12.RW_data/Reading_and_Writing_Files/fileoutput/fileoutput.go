package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// var outputWriter *bufio.Writer
	// var outputFile *os.File
	// var outputError os.Error
	// var outputString string

	// OpenFile函数有三个参数，：文件名、一个或多个标志（使用逻辑运算符"|"连接）、使用的文件权限：
	// os.O_RDONLY只读
	// os.O_WRONLY只写
	// os.O_CREATE创建，如果指定文件不存在，就创建该文件
	// os.O_TRUNC截断，如果指定文件已存在，就将该文件的长度截为0
	outputFile, outputError := os.OpenFile("D:/git_biilGo/go_basic/20221231_RW_data/fileoutput/output.dat", os.O_WRONLY|os.O_CREATE, 0666)
	// 在读文件的时候，文件的权限是被忽略的，所以使用OpenFile时传入第三个参数可以用0。
	// 而在写文件时，不管时Unix还是Windows都需要使用0666

	// 如果文件不存在则自动创建
	if outputError != nil {
		fmt.Printf("An error occurred with file opening or creaton\n")
		return
	}

	defer outputFile.Close()

	// 除了文件句柄，我们还需要bufio的Writer，我们以只读模式打开output.dat

	// 创建一个写入器（缓冲区）对象
	outputWriter := bufio.NewWriter(outputFile)

	outputString := "hello world!\n"

	for i := 0; i < 10; i++ {
		// 使用for循环将字符串写入缓冲区
		outputWriter.WriteString(outputString)
	}

	// 缓冲区的内容紧接着被完全写入文件
	outputWriter.Flush()
}

// 如果写入的东西很简单，我们可以使用fmt.Fprintf(outputFile,"Some test data.\n")直接将内容写入文件
// fmt包里的F开头的Print函数可以直接写入任何io.Writer，包括文件
