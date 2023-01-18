package main

import (
	"fmt"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	// 在服务端控制台打印状态;每个处理函数被调用时,把请求记录下来也许会更有用
	fmt.Println("Inside HelloServer handler")

	// 用来写入合同谈判.ResponseWriter的函数
	// 比如可以使用下面代码构建一个非常简单的网页并插入title和body的值
	/*
		fmt.Fprintf(w, "<h1>%s<h1><div>%s</div>", title, body)
	*/
	fmt.Fprintf(w, "Hello,"+req.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", HelloServer)

	// 没有错误的处理代码的写法
	/*
		http.ListenAndServe(":8080", http.HandlerFunc(HelloServer))
	*/

	// 如果需要使用安全的https连接
	// 使用:http.ListenAndServeTLS()
	err := http.ListenAndServe("localhost:8080", nil)

	if err != nil {
		log.Fatal("ListenAndServer:", err.Error())
	}
}

// 浏览器请求http://localhost:8080/world返回Hello,world
