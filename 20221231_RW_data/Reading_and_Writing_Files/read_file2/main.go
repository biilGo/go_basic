// 按列读取文件中的数据
package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("D:/git_biilGo/go_basic/read_file2/products2.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	var col1, col2, col3 []string

	for {
		var v1, v2, v3 string
		_, err := fmt.Fscanln(file, &v1, &v2, &v3)
		// scans until newline

		if err != nil {
			break
		}

		col1 = append(col1, v1)
		col2 = append(col2, v2)
		col3 = append(col3, v3)
	}

	fmt.Println(col1)
	fmt.Println(col2)
	fmt.Println(col3)
}

// path包里包含一个子包叫filepath,这个子包提供了跨平台的函数，用于处理文件名和路径
// filename := filepath.Base(path)
