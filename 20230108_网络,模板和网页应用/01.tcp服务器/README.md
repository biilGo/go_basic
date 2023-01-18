# tcp服务器
这部分我们将使用tcp协议和之前讲到的协程范式编写一个简单的客户端-服务器应用，一个web服务器应用需要响应众多客户端的并发请求：go会为每一个客户端产生一个协程用来处理请求。我们需要使用net包中网络通信的功能，它包含了处理tcp/ip以及udp协议、域名解析等方法

`net.Error`:

`net`包返回的错误类型遵循惯惯例`error`,但有些错误实现包含额外的方法,他们被定义为`net.Error`接口
```
package net
type Error interface {
    Timeout() bool // 错误是否超时
    Temporary() bool // 是否是临时错误
}
```

通过类型断言,客户端代码可以测试`net.Error`从而区分是临时发生的还是必然会出现的错误.举例来说,一个网络爬虫程序在遇到临时发生的错误时可能会休眠或者重试,如果一个必然发生的错误,则会放弃继续执行
```
// in a loop - some function returns an error err
if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
    time.Sleep(1e9)
    continue // try again
}
if err != nil {
    log.Fatal(err)
}
```