# 复用
## 典型的客户端/服务器模式
客户端-服务器应用正是goroutine和channels的亮点所在

客户端可以是运行在任意设备上的任意程序,它会按需发送请求至服务器,服务器接收到这个请求后开始相应的工作,然后再将响应返回给客户端,典型情况下一般是多个客户端对应一个服务器

例如我们入场使用的浏览器客户端,其功能就是向服务器请求网页,而web服务器则会向浏览器响应网页数据

使用go的服务器通常会在协程中执行向客户端的响应,故而会对每一个客户端请求启动一个协程,一个常用的操作方法是客户端请求自身中包含一个通道,而服务器则向这个通道发送响应

例如`request`结构,其中内嵌了一个`replyc`通道
```
type Request struct {
    a, b      int    
    replyc    chan int // reply channel inside the Request
}
```

或者通俗的
```
type Reply struct{...}
type Request struct{
    arg1, arg2, arg3 some_type
    replyc chan *Reply
}
```

接下来先使用简单形式,服务器会为每一个请求启动一个协程并在其中执行`run()`函数,此举会将类型为`binOp`的`op`操作返回的int值发送到`replyc`通道
```
type binOp func(a, b int) int
func run(op binOp, req *Request) {
    req.replyc <- op(req.a, req.b)
}
```

`server`协程会无限循环以从`chan *Request`接收请求,并且为了避免长时间操作所堵塞,它将为每一个请求启动一个协程来做具体的工作
```
func server(op binOp, service chan *Request) {
    for {
        req := <-service; // requests arrive here  
        // start goroutine for request:        
        go run(op, req);  // don’t wait for op to complete    
    }
}
```

`server`本身则是以协程的方式在`startServer`函数中启动
```
func startServer(op binOp) chan *Request {
    reqChan := make(chan *Request);
    go server(op, reqChan);
    return reqChan;
}
```

`startServer`则会在`main`协程中被调用

在以下测试例子中,100个请求会被发送到服务器,只有它们全部被送达后我们才会按相反的顺序检查响应
```
func main() {
    adder := startServer(func(a, b int) int { return a + b })
    const N = 100
    var reqs [N]Request
    for i := 0; i < N; i++ {
        req := &reqs[i]
        req.a = i
        req.b = i + N
        req.replyc = make(chan int)
        adder <- req  // adder is a channel of requests
    }
    // checks:
    for i := N - 1; i >= 0; i-- {
        // doesn’t matter what order
        if <-reqs[i].replyc != N+2*i {
            fmt.Println(“fail at”, i)
        } else {
            fmt.Println(“Request “, i, “is ok!”)
        }
    }
    fmt.Println(“done”)
}
```

这个程序仅启动了100个协程,然而即使执行了100,000个协程我们也能在数秒内看到它完成,这说明了go的协程是如何的轻量,如果我们启动了相同数量的真实的线程,程序早就崩溃了

## 卸载:通过信号通道关闭服务器
在一个版本中`server`在`main`函数返回后并没有完全关闭,而被强制结束了,为了改进这一点,我们可以提供一个退出的通道给`server`
```
func startServer(op binOp) (service chan *Request, quit chan bool) {
    service = make(chan *Request)
    quit = make(chan bool)
    go server(op, service, quit)
    return service, quit
}
```

`server`函数现在则使用`select`在`service`通道和`quit`通道之间做出选择
```
func server(op binOp, service chan *request, quit chan bool) {
    for {
        select {
            case req := <-service:
                go run(op, req) 
            case <-quit:
                return   
        }
    }
}
```

当`quit`通道接收到一个`true`值`server`就会返回并结束

在`main`函数中我们做出如下更改:`adder, quit := startServer(func(a, b int) int { return a + b })`

在`main`函数的结尾处我们放一行`quit <- true`