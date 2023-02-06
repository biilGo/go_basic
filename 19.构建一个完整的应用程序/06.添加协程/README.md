# 添加协程
## 用协程优化性能
如果由太多客户端同时尝试添加URL,第二个版本依旧存在性能问题.得益于锁机制,我们的map可以在并发访问环境下安全的更新,但每条新产生的记录都要立即写入磁盘,这种机制成为了瓶颈,写入操作可能同时发生,根据不同操作系统特性,可能会产生数据损坏.就算不产生写入冲突,每个客户端在`Put`函数返回前,必须等待数据写入磁盘.因此,在一个I/O负载很高的系统中,客户端为了完成`Add`请求,将等待更长的不必要的实践.

为缓解该问题,必须对`Put`和存储进程解耦:我们将使用go的并发机制.我们不再将记录直接写入磁盘,而是发送到一个通道中,它是某种形式的缓冲区,因而发送函数不必等待它完成.

保存进程会从该通道读取数据并写入磁盘,它是以`saveLoop`协程启动的独立线程.现在`main`和`saveLoop`并行的执行,不会再发生阻塞.

将`URLStore`的`file`字段替换为`record`类型的通道:`save chan record`
```
type URLStore struct {
    urls map[string]string
    mu sync.RWMutex
    save chan record
}
```

通道和map一样,必须用`make`创建,我摸会以此修改`NewURLStore`工厂函数,并给定缓存的`save`通道:
```
func (s *URLStore) Put(url string) string {
    for {
        key := genKey(s.Count())
        if s.Set(key, url) {
            s.save <- record{key, url}
            return key
        }
    }
    panic("shouldn't get here")
}
```

`save`通道的另一端必须由一个接收者:新的`saveLoop`方法在独立的协程中运行,它接收record值并将它们写入到文件.`saveLoop`是在`NewURLStore()`函数中用`go`关键字启动的,可以移除不必要的打开文件的代码,以下是修改后的`NewURLStore()`
```
const saveQueueLength = 1000
func NewURLStore(filename string) *URLStore {
    s := &URLStore{
        urls: make(map[string]string),
        save: make(chan record, saveQueueLength),
    }
    if err := s.load(filename); err != nil {
        log.Println("Error loading URLStore:", err)
    }
    go s.saveLoop(filename)
    return s
}
```

以下是`saveLoop`方法的代码:
```
func (s *URLStore) saveLoop(filename string) {
    f, err := os.Open(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal("URLStore:", err)
    }
    defer f.Close()
    e := gob.NewEncoder(f)
    for {
        // taking a record from the channel and encoding it
        r := <-s.save
        if err := e.Encode(r); err != nil {
            log.Println("URLStore:", err)
        }
    }
}
```

在无限循环中,记录从save通道读取,然后编码到文件中.

14章深入学习了协程和通道,但在这里我们见到了实用的案例,更好的管理程序的不同部分,注意现在`Encoder`对象只被创建一次,而不是每次保存时都创建,这也可以节省了一些内存和运算处理.

还有一个改进可以使goto更灵活:我们可以将文件名,监听地址和主机名定义为标志,来代替程序中硬编码或定义常量,这样当程序启动时,可以在命令行中指定它们的新值,如果没有指定,将采用flag的默认值,该功能来自另一个包,所以需要`import "flag"`

先创建一些全局变量来保存flag的值:
```
var (
    listenAddr = flag.String("http", ":8080", "http listen address")
    dataFile = flag.String("file", "store.gob", "data store file name")
    hostname = flag.String("host", "localhost:8080", "host name and port")
)
```

为了处理命令行参数,必须把`flag.Parse()`添加到`main`函数中,在flag解析后才能实例化`URLStore`,一旦得知了`dataFile`的值:
```
var store *URLStore
func main() {
    flag.Parse()
    store = NewURLStore(*dataFile)
    http.HandleFunc("/", Redirect)
    http.HandleFunc("/add", Add)
    http.ListenAndServe(*listenAddr, nil)
}
```

现在`Add`处理函数中必须用`*hostname`提花`localhost:8080`:
`fmt.Fprintf(w, "http://%s/%s", *hostname, key)`

编译或直接使用现有的可执行程序测试第三版本