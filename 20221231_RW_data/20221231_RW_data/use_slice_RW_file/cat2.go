// 切片提供Go中数量I/O缓冲的标准方式，这版cat函数中，在一个切片缓冲内使用无限for循环读取文件

// 使用os包中os.File和Read方法

package main

import (
	"flag"
	"fmt"
	"os"
)

func cat(f *os.File) {
	const NBUF = 512
	var buf [NBUF]byte
	for {
		switch nr, err := f.Read(buf[:]); true {
		case nr < 0:
			fmt.Fprintf(os.Stderr, "cat:error reading:%s\n", err.Error())
		case nr == 0: // EOF
			return
		case nr > 0:
			if nw, ew := os.Stdout.Write(buf[0:nr]); nw != nr {
				fmt.Fprintf(os.Stderr, "cat error writing:%s\n", ew.Error())
			}
		}
	}
}

func main() {
	flag.Parse() // Scans the arg list and sets up flags
	if flag.NArg() == 0 {
		cat(os.Stdin)
	}

	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if f == nil {
			fmt.Fprintf(os.Stderr, "cat:can't open %s:error %s\n", flag.Arg(i), err)
			os.Exit(1)
		}
		cat(f)
		f.Close()
	}
}
