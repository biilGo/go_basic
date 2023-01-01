package main

import "os"

func main() {
	// 使用os.Stdout.WriteString()可以输出到屏幕
	os.Stdout.WriteString("hello,world\n")

	// 以只写模式创建或打开文件，并且忽略了可能发生的错误
	f, _ := os.OpenFile("D:/git_biilGo/go_basic/20221231_RW_data/fileoutput/filewrite/test", os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()

	// 我们不适用缓冲区直接将内容写入文件
	f.WriteString("hello.world in a file\n")
}
