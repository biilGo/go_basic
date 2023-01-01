// 读文件

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// 此程序我们以只读模式打开input.dat文件

func main() {
	// 变量inputFile是*os.File类型的
	inputFile, inputError := os.Open("D:/git_biilGo/go_basic/fileinput/input.dat")
	// 该类型是一个结构，表示一个打开文件的描述符，然后使用os包里Open函数打开一个文件
	// 该函数的参数是文件名，类型为string

	// 如果文件不存在或者程序没有足够的权限打开这个文件，Open函数会返回一个错误：
	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputfile\n" +
			"Does the file exist?\n" +
			"Have you got acces to it?\n")
		return //exit the function on error
	}

	// 如果打开正常我们就使用defer inputFile.Close()语句确保在程序退出前关闭该文件
	defer inputFile.Close()

	// 然后我们使用bufio.NewReader来获得一个读取器变量,使用bufio包提供的读取器，我们可以很方便的操作相对高层的string对象，而避免去操作比较底层的字节
	inputReader := bufio.NewReader(inputFile)

	// 在一个无限循环中使用ReadString('\n')或ReadBytes('\n')将文件内容逐行读取出来
	// Unix和Linux行结束符是\n，Windows行结束符是\r\n
	// ReadString和ReadBytes方法，我们不需要关心操作系统的类型，直接使用\n就可以了，另外我们也可以使用ReadLine()方法来实现相同的功能
	for {
		inputString, readerError := inputReader.ReadString('\n')
		fmt.Printf("The input was:%s", inputString)

		// 一旦读取到文件末尾，变量readError的值将变成非空且其值为常量io.EOF
		// 我们就会执行return语句从而退出循环
		if readerError == io.EOF {
			return
		}
	}
}
