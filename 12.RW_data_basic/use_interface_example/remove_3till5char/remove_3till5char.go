// 代码有一个输入文件goprogram然后以每一行为单位读取，从读取的当前行中截取第3到第5字节写入另一个文件

// 运行程序输出的文件是空文件，找出逻辑中的bug

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	inputFile, _ := os.Open("goprogram")

	outputFile, _ := os.OpenFile("goprogramT", os.O_WRONLY|os.O_CREATE, 0666)

	defer inputFile.Close()

	defer outputFile.Close()

	inputReader := bufio.NewReader(inputFile)

	outputWriter := bufio.NewWriter(outputFile)

	for {
		intputString, _, readerError := inputReader.ReadLine()

		if readerError == io.EOF {
			fmt.Println("EOF")
			return
		}

		outputString := string(intputString[2:5]) + "\r\n"

		_, err := outputWriter.WriteString(outputString)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Conversion done")
}
