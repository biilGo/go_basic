package main

import (
	"io"
	"net/http"
)

// 当使用字符串常量表示html文本的时候包含<html><body>...</body></html>对于浏览器将它识别为html文档非常重要
// 更安全的做法是处理函数中,在写入返回内容之前将头部的content-type 设置为text/html:
/* w.Header().Set("Content-Type", "text/html") */
const form = `
	<html><body>
		<form action="#" method="post" name="bar">
			<input type="text" name="in" />
			<input type="submit" value="submit"/>
		</form>
	</body></html>
`

// handle a simple get request
// 处理url的/test1使在浏览器里输出hello world
func SimpleServer(w http.ResponseWriter, request *http.Request) {
	io.WriteString(w, "<h1>hello, world</h1>")
}

// 处理/test2:
// 如果url最初由浏览器请求,那么它是一个GET请求返回一个form常量,包含简单的input表单
// 表单里有一个文本框和一个提交按钮,在文本框里输入一些东西并点击提交按钮的时候,会发起一个POST请求
func FormServer(w http.ResponseWriter, request *http.Request) {

	// content-type 会让浏览器认为它可以使用函数 http.DetectContentType([]byte(form)) 来处理收到的数据
	w.Header().Set("Content-Type", "text/html")

	switch request.Method {
	case "GET":
		// display the form to the user
		io.WriteString(w, form)

	// 请求为POST类型时:
	// name属性为inp的文本框内容可以获取
	case "POST":
		// handle the form data,note that parseForm must
		// be called before we can extract form data
		// request.ParseForm();
		// io.WriteString(w,request.Form["in"][0])
		// 获取request.FormValue("inp"),然后将其写回浏览器页面中
		io.WriteString(w, request.FormValue("in"))
	}
}

func main() {
	http.HandleFunc("/test1", SimpleServer)

	http.HandleFunc("/test2", FormServer)

	// 8080端口启动一个网页服务器
	if err := http.ListenAndServe(":8088", nil); err != nil {
		panic(err)
	}
}
