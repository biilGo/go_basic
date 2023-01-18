# 一个简单的web服务器
http是比tcp更高层的协议,它描述了网页服务器如何与客户端浏览器进行通信.

引入`http`包并启动了网页服务器,`net.Listen("tcp","localhost:50000")`函数的tcp服务器是类似的,使用`http.ListenAndServer("localhost:8080",nil)`函数,如果成功会返回空,否则返回一个错误

`http.URL`用于表示网页地址,其中字符串属性`Path`用于保存url路径,`http.Request`描述了客户端请求,内含一个`URL`字段

如果`req`来自html表单的POST类型请求,'var1'是该表单中一个输入的域的名称,那么用户输入的值就可以通过go代码`req.FormValue("var1")`获取到.还有一种方法是先执行`request.ParseForm()`,然后再获取`request.Form["var1"]`的第一个返回参数:`var1, found := request.Form["var1"]`

第二个参数`found`为true,如果var1并未出现在表单中,found就是false

表单属性实际上是`map[string][]string`类型,网页服务器发送一个`http.Response`响应,通过`http.ResponseWrite`对象输出的,后者组装了HTTP服务器响应,通过对其写入内容,我们就将数据发送给HTTP客户端

继续编写程序,以实现服务器必须做的事,即如何处理请求.这是通过`http.HandleFunc`函数完成的.当根路径`/`被请求的时候,`HelloServer`函数就被执行了.这个函数是`http.HandlerFunc`类型的,它们通常被命名为Prefhandler和某个路径前缀Pref匹配

`http.HandleFunc`注册了一个处理函数来处理对应`/`的请求

`/`可以被替换为其他更特定的url,比如`/create`或者`/edit`等等,你可以为每一个特定的url定义一个单独的处理函数,这个函数需要两个参数,第一个是`ReponseWrite`类型的`w`;第二个是请求`req`.程序向`w`写入了`Hello`和`r.URL.Path[1:]`组成的字符串;末尾的`[1:]`表示创建一个从索引为1的字符到结尾的子切片,用来丢弃路径开头的`/`,`fmt.Fprintf()`函数完成本次写入;另一种可行的写法是`io.WriteString(w,"hello,world!\n")`

除了`http.HandleFunc("/", Hfunc)`其中`HFunc`是一个处理哈数,签名为
```
func HFunc(w http.ResponseWriter, req *http.Request) {
    ...
}
```

也可以使用这种方式:`http.Handle("/", http.HandlerFunc(HFunc))`

`HandleFunc`只是定义了上述HFunc签名的别名:`type HandlerFunc func(ResponseWriter, *Request)`

它是一个可以把普通的函数当作HTTP处理器的适配器,如果函数f声明的合适,`HandlerFunc(f)`就是一个执行f函数的Handler对象

`http.Handle`的第二个参数也可以是T类型的对象obj:`http.Handle("/", obj)`

如果T有ServeHttp方法,那就实现了http的Handler接口
```
func (obj *Typ) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    ...
}
```

