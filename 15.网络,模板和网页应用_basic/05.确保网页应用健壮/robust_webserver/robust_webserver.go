package main

import (
	"io"
	"log"
	"net/http"
)

const form = `<html><body><form action="#" method="post" name="bar">
<input type="text" name="in"/>
<input type="submit" value="Submit"/>
</form></html></body>`

// 为增强代码可读性，为页面处理函数创建一个类型
type HandleFnc func(http.ResponseWriter, *http.Request)

// handle a simple get request
func SimpleServer(w http.ResponseWriter, request *http.Request) {
	io.WriteString(w, "<h1>hello,world</h1>")
}

// handle a form ,both the get which display the form and the POST which processes it
func FormServer(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch request.Method {
	case "GET":
		// display the form to the user
		io.WriteString(w, form)
	case "POST":
		/*
			handle the form data, note that ParseForm must
			be called before we can extract form data
		*/
		//request.ParseForm();
		//io.WriteString(w, request.Form["in"][0])
		io.WriteString(w, request.FormValue("in"))
	}
}

// 错误处理函数logPanics
func logPanics(function HandleFnc) HandleFnc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Printf("[%v] caught panic:%v", request.RemoteAddr, x)
			}
		}()
		function(writer, request)
	}
}

func main() {

	// 用logPanics来包装对处理函数的调用
	http.HandleFunc("/test1", logPanics(SimpleServer))
	http.HandleFunc("/test2", logPanics(FormServer))

	if err := http.ListenAndServe(":8088", nil); err != nil {
		panic(err)
	}
}
