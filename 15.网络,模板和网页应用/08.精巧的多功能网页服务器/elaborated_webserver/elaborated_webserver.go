package main

import (
	"bytes"
	"expvar"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// hello workld the web server
var helloRequests = expvar.NewInt("hello-requests")

// flags:
var webroot = flag.String("root", "/home/user", "web root directory")

// simple flag server
var booleanflag = flag.Bool("boolean", true, "another flag for testing")

// Simple counter server. POSTing to it will set the value
type Counter struct {
	n int
}

// a channel
type Chan chan int

// Logger初六函数w.Write Header(404)输出404头部
// 此项技术通常很有用,无论合适服务器执行代码产生错误,都可以应用类似这样的代码:
/* if err != nil {
    w.WriteHeader(400)
    return
} */
// logger包的函数,针对每个请求在服务器端命令行打印日期,时间,URL
func Logger(w http.ResponseWriter, req *http.Request) {
	log.Print(req.URL.String())
	w.WriteHeader(404)
	w.Write([]byte("oops"))
}

// expvar可以创建int,float,string类型变量,并将它们发布为公共变量,在URL/debug/vars上以JSON格式公布,通常被用于服务器操作计数
// 处理函数
func HelloServer(w http.ResponseWriter, req *http.Request) {
	// 该处理函数对其+1,然后写入"hello world"到浏览器
	helloRequests.Add(1)
	io.WriteString(w, "hello, world!\n")
}

// this makes Counter satisfy the expvar,.Var interface, so we can export
// it directly.
// String()方法实现了expvar.Var接口,使其可以被发布,它是一个结构体.ServeHTTP函数使ctr成为处理器,因为它的签名正确实现了http.Handler接口
func (ctr *Counter) String() string {
	return fmt.Sprintf("%d", ctr.n)
}

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET": // increment n
		ctr.n++
	case "POST":
		buf := new(bytes.Buffer)
		io.Copy(buf, req.Body)
		body := buf.String()
		if n, err := strconv.Atoi(body); err != nil {
			fmt.Fprintf(w, "bad POST:%v\nbody:[%v]\n", err, body)
		} else {
			ctr.n = n
			fmt.Fprint(w, "count reset\n")
		}
	}
	fmt.Fprintf(w, "counter = %d\n", ctr.n)
}

// 该函数使用了flag包VisitAll函数跌倒所有的标签，打印它们的名称，值和默认值

func FlagServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Flag:\n")
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() != f.DefValue {
			fmt.Fprintf(w, "%s = %s [default = %s]\n", f.Name, f.Value.String(), f.DefValue)
		} else {
			fmt.Fprintf(w, "%s = %s\n", f.Name, f.Value.String())
		}
	})
}

// simple argument server
// 该函数迭代os.Args以打印出所有的命令行参数,如果没有指定则只有程序名称会被打印出来
// ArgServer; http://localhost:12345/args; ./elaborated_webserver.exe
func ArgServer(w http.ResponseWriter, req *http.Request) {
	for _, s := range os.Args {
		fmt.Fprint(w, s, " ")
	}
}

func ChanCreate() Chan {
	c := make(Chan)
	go func(c Chan) {
		for x := 0; ; x++ {
			c <- x
		}
	}(c)
	return c
}

// Channel;http://localhost:12345/chan;channel send #1
func (ch Chan) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, fmt.Sprintf("channel send #%d\n", <-ch))
}

// exec a program,redirecting output
// 显示当前时间
// os.Pipe返回一对相关联的File从r读取数据,返回已读取的字节数来自于w的写入.函数返回这2个文件和错误:
// func Pipe() (r *File, w *File, err error)
func DateServer(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-type", "text/plain; charset=utf-8")
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(rw, "pipe:%s\n", err)
		return
	}

	p, err := os.StartProcess("/bin/date", []string{"date"}, &os.ProcAttr{Files: []*os.File{nil, w, w}})

	defer r.Close()

	w.Close()

	if err != nil {
		fmt.Fprintf(rw, "fork/exec:%s\n", err)
		return
	}

	defer p.Release()

	io.Copy(rw, r)

	wait, err := p.Wait()

	if err != nil {
		fmt.Fprintf(rw, "wait:%s\n", err)
		return
	}

	if !wait.Exited() {
		fmt.Fprintf(rw, "date:%v\n", wait)
		return
	}
}

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(Logger))
	http.Handle("/go/hello", http.HandlerFunc(HelloServer))

	// the counter is published as a variable directly
	ctr := new(Counter)
	expvar.Publish("counter", ctr)
	http.Handle("/counter", ctr)

	// http.Handle("/go/",http.FileServer(http.Dir("/tmp"))) // uses the OS Filesystem

	// FileServer(root FileSystem) Handler返回一个处理器，它以root作为根，用文件系统的内容相应HTTP请求
	// 获得操作系统的文件系统，用http.Dir:
	// http.Handle("/go/", http.FileServer(http.Dir("/tmp")))
	http.Handle("/go/", http.StripPrefix("/go/", http.FileServer(http.Dir(*webroot))))

	http.Handle("flags", http.HandlerFunc(FlagServer))

	http.Handle("args", http.HandlerFunc(ArgServer))

	// http.Handle("/chan", ChanCreate())
	http.Handle("/chan", ChanCreate())

	http.Handle("date", http.HandlerFunc(DateServer))

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Panicln("ListenAndServer:", err)
	}
}

// 每当有刷新请求到达,通道的ServeHTTP方法从通道获取下一个整数并显示,由此可见,网页服务器可以从通道中获取要发送的相应,它可以由另一个函数产生
/* func ChanResponse(w http.ResponseWriter, req *http.Request) {
    timeout := make (chan bool)
    go func () {
        time.Sleep(30e9)
        timeout <- true
    }()
    select {
    case msg := <-messages:
        io.WriteString(w, msg)
    case stop := <-timeout:
        return
    }
} */
