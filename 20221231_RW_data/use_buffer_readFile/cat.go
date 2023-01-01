// 用buffer读取文件

// 结合使用缓冲读取文件和命令行flag解析这两项技术，如果不加参数，那么你输入什么屏幕就打印什么

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func cat(r *bufio.Reader) {
	for {
		buf, err := r.ReadBytes('\n')
		fmt.Fprintf(os.Stdout, "%s", buf)
		if err == io.EOF {
			break
		}
	}
	return
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		cat(bufio.NewReader(os.Stdin))
	}

	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s:error reading from %s:%s\n", os.Args[0], flag.Arg(i), err.Error())
			continue
		}

		cat(bufio.NewReader(f))
		f.Close()
	}
}
