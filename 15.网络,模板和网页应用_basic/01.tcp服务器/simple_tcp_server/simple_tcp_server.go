// simple multi-thread/multi-core TCP server

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// 使用25字节的缓冲读取客户端发送的数据
const maxRead = 25

// 得到了服务器地址和端口，返回一个指针类型
func initServer(hostAndPort string) *net.TCPListener {
	serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
	checkError(err, "Resolving address:port falied: '"+hostAndPort+"'")
	listener, err := net.ListenTCP("tcp", serverAddr)
	checkError(err, "ListenTCP:")
	println("Listening to:", listener.Addr().String())
	return listener
}

// 每一个连接都会以协程的方式运行connectionHandler
func connectionHandler(conn net.Conn) {
	// 获取客户端的地址并显示出来
	connFrom := conn.RemoteAddr().String()
	println("Connection from:", connFrom)
	sayHello(conn)

	// 循环读取25字节缓冲读取客户端发送的数据并一一打印出来
	for {
		var ibuf []byte = make([]byte, maxRead+1)
		length, err := conn.Read(ibuf[0:maxRead])
		ibuf[maxRead] = 0 // to prevent overflow

		// 读取错误进入switch语句default分支
		switch err {
		case nil:
			handleMsg(length, err, ibuf)
		//try again
		// 由于版本的更新，会提示undenfined
		case os.EAGAIN:
			continue

		// 退出无限循环并关闭连接
		default:
			goto DISCONNECT
		}
	}
DISCONNECT:
	err := conn.Close()
	println("Closed connection:", connFrom)
	checkError(err, "Close:")
}

func sayHello(to net.Conn) {
	obuf := []byte{'L', 'e', 't', '\'', 's', ' ', 'G', 'O', '!', '\n'}
	wrote, err := to.Write(obuf)
	checkError(err, "Write:wrote"+string(wrote)+"bytes.")
}

func handleMsg(length int, err error, msg []byte) {
	if length > 0 {
		print("<", length, ":")
		for i := 0; ; i++ {
			if msg[i] == 0 {
				break
			}
			fmt.Printf("%c", msg[i])
		}
		print(">")
	}
}

// 所有的错误检查都被重构在独立的函数中，当错误产生时，利用错误上下文触发panic
func checkError(error error, info string) {
	if error != nil {
		panic("ERROR:" + info + " " + error.Error()) // terminate program
	}
}

func main() {
	flag.Parse()

	// 服务器不再是硬编码，而是通过命令行参数传入
	// 通过flag包来读取这些参数，这里使用了flag.NArg()检查是否按照期望传入了2个参数
	if flag.NArg() != 2 {
		panic("usage: host port")
	}

	// 传入的参数通过Sprintf函数格式化成字符串
	hostAndPort := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1))

	listener := initServer(hostAndPort)

	for {
		conn, err := listener.Accept()
		checkError(err, "Accept:")
		go connectionHandler(conn)
	}
}
