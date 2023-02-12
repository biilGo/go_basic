# 分布式程序
## 多服务器处理架构
目前为止`goto`以单线程运行,但即使用协程,在一台机器上运行的单一进程,也只能为一定数量的并发请求提供服务.一个缩短网址服务,相对于`Add`用`Put()`写入,通常`Redirect`服务要多得多.因此我们应该可以创建任意数量的只读的从服务器,提供服务并缓方法调用的结果,将`Put`方法调用的结果,将`Put`请求转发给主服务器,类似如下架构:
![19.8_主从架构](./assets/19.8_%E4%B8%BB%E4%BB%8E%E6%9E%B6%E6%9E%84.jpg)

对于slave进程,要在网络上运行`goto`应用的一个master节点实例,它们必须能相互通信.go的`rpc`包为跨越网络发起函数调用提供了便捷的途径,这里将把`URLStore`变为RPC服务.slave进程将应对get请求以交付长URL.当一个长URL要被转换为缩短版本时,它们通过rpc连接把任务委托给master进程,因此只有master节点会写入数据文件.

截至目前,`URLStore`上基本的`Get()`和`Put()`方法具有如下签名:
```
func (s *URLStore) Get(key string) string
func (s *URLStore) Put(url string) string
```

而RPC调用仅能使用如下形式的方法(t是T类型的值):`func (t T) Name(args *ArgType, reply *ReplyType) error`

要使`URLStore`成为RPC服务,需要修改put和get方法使它们符合上述函数签名,以下是修改后的签名:
```
func (s *URLStore) Get(key, url *string) error
func (s *URLStore) Put(url, key *string) error
```

`Get()`代码变更为:
```
func (s *URLStore) Get(key, url *string) error {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if u, ok := s.urls[*key]; ok {
        *url = u
        return nil
    }
    return errors.New("key not found")
}
```

现在,键和长URL都变成了指针,必须加上前缀`*`来取得它们的值,例如`*key`这种形式,`u`是一个值,可以用`*url=u`来将其赋值给指针.

接着对`Put()`代码做同样的改动:
```
func (s *URLStore) Put(url, key *string) error {
    for {
        *key = genKey(s.Count())
            if err := s.Set(key, url); err == nil {
            break
        }
    }
    if s.save != nil {
        s.save <- record{*key, *url}
    }
    return nil
}
```

`Put()`调用`Set()`,由于后者也要做调整,`key`和`url`参数现在是指针类型,还必须返回`error`取代`boolean`:
```
func (s *URLStore) Set(key, url *string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    if _, present := s.urls[*key]; present {
        return errors.New("key already exists")
    }
    s.urls[*key] = *url
    return nil
}
```

同样,当从`load()`调用`Set()`时,也必须做调整:`s.Set(&r.Key, &r.URL)`

还必须修改HTTP处理函数以适应`URLStore`上的更改,`Redirect`处理函数现在返回`URLStore`给出错误的字符串形式:
```
func Redirect(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Path[1:]
    var url string
    if err := store.Get(&key, &url); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, url, http.StatusFound)
}
```

`Add`处理函数也以基本相同的方式修改:`
```
func Add(w http.ResponseWriter, r *http.Request) {
    url := r.FormValue("url")
    if url == "" {
        fmt.Fprint(w, AddForm)
        return
    }
    var key string
    if err := store.Put(&url, &key); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "http://%s/%s", *hostname, key)
}
```

要使应用程序更灵活,正如之前章节所为,可以添加一个命令行标志flag来决定是否在main()函数中启动RPC服务器:`var rpcEnabled = flag.Bool("rpc", false, "enable RPC server")`

要使RPC工作,还要用rpc包来注册URLStore,并用`HandleHTTP`创建基于HTTP的RPC处理器:
```
func main() {
    flag.Parse()
    store = NewURLStore(*dataFile)
    if *rpcEnabled { // flag has been set
        rpc.RegisterName("Store", store)
        rpc.HandleHTTP()
    }
    ... (set up http like before)
}
```