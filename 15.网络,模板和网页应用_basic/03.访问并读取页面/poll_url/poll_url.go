package main

import (
	"fmt"
	"net/http"
)

// 数组中的url将被访问
var urls = []string{
	"http://www.biying.com/",
	"http://www.bookstack.cn/",
	"http://mail.sina.com",
}

func main() {
	// Execute an http head request for all url's
	// and returns the http status string or an error string

	for _, url := range urls {

		// 发送一个http.Head()请求查看返回值
		resp, err := http.Head(url)

		if err != nil {
			fmt.Println("Error:", url, err)
		}

		// 打印响应的Response状态码
		fmt.Println(url, ": ", resp.Status)
	}
}
