// 如何拷贝一个文件到另一个文件？最简单的方式就是使用io包
package main

import (
	"fmt"
	"io"
	"os"
)

// 注意defer的使用
// 当打开dst文件时发生了错误，那么defer仍然能够确保src.Close()执行
// 如果不那么做，src文件会一直保持打开状态并占用资源

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func main() {
	CopyFile("D:/git_biilGo/go_basic/20221231_RW_data/cp_sourcefile_destfile/filecopy/target.txt", "D:/git_biilGo/go_basic/20221231_RW_data/cp_sourcefile_destfile/filecopy/source.txt")
	fmt.Println("Copy done!")
}
