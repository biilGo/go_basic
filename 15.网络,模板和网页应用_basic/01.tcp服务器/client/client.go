package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// 打开连接:
	// 客户端使用Dial创建了一个和服务器之间的连接
	// 连接到远程系统成功后，dial函数会返回要给conn类型的接口
	// 我们可以用它发送和接收数据
	// Dial函数简洁地抽象了网络层和传输层，所以不管是ipv4还是ipv6，tcp或udp都可以使用这个公用接口
	conn, err := net.Dial("tcp", "localhost:50000")

	if err != nil {
		// 由于目标计算机积极拒绝而非无法创建连接
		fmt.Println("Error dialing", err.Error())
		return // 终止程序
	}

	// 通过无限循环从stdin接收来自键盘的输入
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("First,what is your name?")

	clientName, _ := inputReader.ReadString('\n')
	// fmt.Printf("CLIENTNAME %s", clientName)

	trimmedClient := strings.Trim(clientName, "\r\n") // windows平台下换行符

	// 给服务器发送信息直到程序退出
	for {
		fmt.Println("What to send to the server? Type Q to quit.")

		input, _ := inputReader.ReadString('\n')

		trimmedInput := strings.Trim(input, "\r\n")

		if trimmedClient == "Q" {
			return
		}

		// Write方法送达到服务器
		_, err = conn.Write([]byte(trimmedClient + "says:" + trimmedInput))
	}
}
