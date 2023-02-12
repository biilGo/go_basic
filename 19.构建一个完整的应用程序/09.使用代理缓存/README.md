# 使用代理缓存
`URLStore`已经成为了有效的RPC服务,现在可以创建另一种代表RPC客户端的类型,它会转发请求到RPC服务器,我们称它为`ProxyStore`
```
type ProxyStore struct {
    client *rpc.Client
}
```

一个RPC客户端必须使用`DialHTTP()`方法连接到服务器,所以我们把这句话加入`NewProxyStore`函数,它用于创建`ProxyStore`对象
```
func NewProxyStore(addr string) *ProxyStore {
    client, err := rpc.DialHTTP("tcp", addr)
    if err != nil {
        log.Println("Error constructing ProxyStore:", err)
    }
    return &ProxyStore{client: client}
}
```

`ProxyStore`有`Get`和`Put`方法,它们利用RPC客户端的`Call`方法,将请求直接传递给服务器:
```
func (s *ProxyStore) Get(key, url *string) error {
    return s.client.Call("Store.Get", key, url)
}
func (s *ProxyStore) Put(url, key *string) error {
    return s.client.Call("Store.Put", url, key)
}
```

## 带缓存的ProxyStore
可是,如果slave进程只是简单地代理所有的工作到master节点,不会得到任务增益,我们打算用slave节点来应对Get请求,要做到这点,它们必须有`URLStore`中的map的一份副本.因此我们对`ProxyStore`的定义进行扩展,将`URLStore`包含在其中:
```
type ProxyStore struct {
    urls *URLStore
    client *rpc.Client
}
```

`NewProxyStore`也必须做修改:
```
func NewProxyStore(addr string) *ProxyStore {
    client, err := rpc.DialHTTP("tcp", addr)
    if err != nil {
        log.Println("ProxyStore:", err)
    }
    return &ProxyStore{urls: NewURLStore(""), client: client}
}
```

还必须修改`NewURLStore`以便给出空文件名时,不会尝试从磁盘写入或读取文件:
```
func NewURLStore(filename string) *URLStore {
    s := &URLStore{urls: make(map[string]string)}
    if filename != "" {
        s.save = make(chan record, saveQueueLength)
        if err := s.load(filename); err != nil {
            log.Println("Error loading URLStore: ", err)
        }
        go s.saveLoop(filename)
    }
    return s
}
```

`ProxyStore`的get方法需要扩展:它应该首先检查缓存中是否有对应的键.如果有,get返回已缓存的结果.否则,应该发起RPC调用,然后用结果更新基本地缓存:
```
func (s *ProxyStore) Get(key, url *string) error {
    if err := s.urls.Get(key, url); err == nil { // url found in local map
        return nil
    }
    // url not found in local map, make rpc-call:
    if err := s.client.Call("Store.Get", key, url); err != nil {
        return err
    }
    s.urls.Set(key, url)
    return nil
}
```

同样地,put方法仅当成功完成了远程RPC`Put`调用,才更新本地缓存:
```
func (s *ProxyStore) Put(url, key *string) error {
    if err := s.client.Call("Store.Put", url, key); err != nil {
        return err
    }
    s.urls.Set(key, url)
    return nil
}
```

## 汇总
slave节点使用`ProxyStore`,只有master使用`URLStore`.有鉴于创造它们的方式,它们看上去十分一致:两者都实现了相同签名的get和put方法,因此我们可以指定一个store接口来概括它们的行为:
```
type Store interface {
    Put(url, key *string) error
    Get(key, url *string) error
}
```

现在全局变量store可以成为Store类型:`var store Store`

最后,我们修改main()函数以便程序只作为master或slave启动.

为此我们添加一个没有默认值的新命令行标志`masterAddr`.

`var masterAddr = flag.String("master", "", "RPC master address")`

如果给出master地址,就启动一个slave进程并创建新的`ProxyStore`;否则启动master进程并创建新的`URLStore`:
```
func main() {
    flag.Parse()
    if *masterAddr != "" { // we are a slave
        store = NewProxyStore(*masterAddr)
    } else { // we are the master
        store = NewURLStore(*dataFile)
    }
    ...
}
```

这样,我们已启用`ProxyStore`作为web前端,以代替`URLStore`.

其余的前端代码继续和之前的一样地工作,它们不必在意Store接口,只有master进程会写数据文件

现在可以加载一个master节点和数个slave节点,对slave进行压力测试.

编译这个版本4或者使用现有的可执行程序.

要进行测试,首先在命令行用以下命令启动master节点:`./goto -http=:8081 -rpc=true    # （Windows 平台用 goto 代替 ./goto）`

这里提供了2个标志:master监听8081端口,已启动RPC

slave节点用以下命令启动:`./goto -master=127.0.0.1:8081`

它获取到master地地址,并在8080端口接受客户端请求.

在源码目录下已包含了以下shell脚本demo.sh用来在类Unix系统下自动启动程序:
```
#!/bin/sh
gomake
./goto -http=:8081 -rpc=true &
master_pid=$!
sleep 1
./goto -master=127.0.0.1:8081 &
slave_pid=$!
echo "Running master on :8081, slave on :8080."
echo "Visit: http://localhost:8080/add"
echo "Press enter to shut down"
read
kill $master_pid
kill $slave_pid
```