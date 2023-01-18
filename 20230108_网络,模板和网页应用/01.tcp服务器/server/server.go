package main

import (
	"fmt"
	"net"
)

// 客户端的请求将产生一个net.conn类型的连接变量
func doServerStuff(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return //终止程序
		}

		// len获取客户端发送的数据字节数
		fmt.Printf("Received data:%v", string(buf[:len]))
	}
}

func main() {
	fmt.Println("Starting th server ...")

	// 创建listener
	// net.Listener类型的变量listener
	// 实现服务器基本功能:
	/*
		用来监听和接收来自客户端的请求,在localhost端口为50000基于TCP协议
	*/
	listener, err := net.Listen("tcp", "localhost:50000")

	if err != nil {
		fmt.Println("Error listening", err.Error())
		return // 终止程序
	}

	// 监听并接收来自客户端的连接
	for {
		// 等待客户端的请求
		conn, err := listener.Accept()

		// 当客户端发送的数据都被读取完成时,协程就结束
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return //终止程序
		}

		// 一个独立的协程使用这个连接执行,开始使用一个512字节的缓冲data来读取客户端发送来的数据
		// 并且把他们打印到服务器的终端
		go doServerStuff(conn)
	}
}
