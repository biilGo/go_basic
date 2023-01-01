// 读取压缩文件
// compress包提供了读取压缩文件的功能，支持的压缩文件格式为bzip2/flate/gzip/lzw/zlib

package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
)

func main() {
	fName := "D:/git_biilGo/go_basic/20221231_RW_data/gzipped/gzippedMyFile.txt.gz"

	var r *bufio.Reader

	fi, err := os.Open(fName)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v,Can't open %s:error:%s\n", os.Args[0], fName, err)
		os.Exit(1)
	}

	defer fi.Close()

	fz, err := gzip.NewReader(fi)

	if err != nil {
		r = bufio.NewReader(fi)
	} else {
		r = bufio.NewReader(fz)
	}

	for {
		line, err := r.ReadSlice('\n')

		if err != nil {
			fmt.Println("Done reading file")
			os.Exit(0)
		}
		fmt.Println(line)
	}
}
