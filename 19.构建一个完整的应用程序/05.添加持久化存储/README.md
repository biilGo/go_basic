# 添加持久化存储
## 持久化存储：gob
当`goto`终止,这迟早会发生,内存map中缩短的URL就会丢失,要保留这些数据,就得将其保存到磁盘文件中,我们将修改`URLStore`,使它可以保存数据到文件,且在`goto`启动时还原这些数据,为此我们使用go标准库的`encoding/gob`包:它用于序列化和反序列化,将数据结构转换为字节数组,反之亦然.

通过`gob`包的`NewEncoder`和`NewDecoder`函数,可以指定数据要写入或读取的位置.返回`Encoder`和`Decoder`对象提供了`Encode`和`Decode`方法,用于对文件写入和从中读取go数据结构.

提示:`Encoder`实现了`Writer`接口,同样`Decoder`实现了`Reader`接口,我们在`URLStore`上增加了一个新的`file`字段,它是用于读写已打开文件的句柄.
```
type URLStore struct {
    urls map[string]string
    mu sync.RWMutex
    file *os.File
}
```
我们把这个文件命名为`store.gob`,当初始化`URLStore`时将其作为参数传入:`var store = NewURLStore("store.gob")`

接着,调整`NewURLStore`函数:
```
func NewURLStore(filename string) *URLStore {
    s := &URLStore{urls: make(map[string]string)}
    f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal("URLStore:", err)
    }
    s.file = f
    return s
}
```

现在,更新后的`NewURLStore`函数接收一个文件名参数,它会打开该文件,将返回的`*os.File`作为`file`字段的值存储在`URLStore`变量`store`中,即这里的本地变量`s`

对`OpenFile`的调用可能会失败,它会返回一个错误err,注意go是如何处理这种情况的:
```
f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
if err != nil {
    log.Fatal("URLStore:", err)
}
```

当err不为`nil`,表示确实发生了错误,那么输出一条消息并停止程序执行.这是处理错误的一种方式,大多数情况下错误应该返回给调用函数,但这种检测错误的模式在go代码中也很普遍,在`}`之后可以确定文件被成功打开了.

打开该文件时启动了写入标志,更精确地说时"追加模式".每当一对新的短/长URL在程序中创建后,我们通过`gob`把它存储到文件"store.gob"中.

为达到目的,定义一个新的结构体类型`record`:
```
type record struct {
    Key, URL string
}
```

以及更新的`save`方法,将给定的健和URL组成`record`,以`gob`编码的形式写入磁盘.
```
func (s *URLStore) save(key, url string) error {
    e := gob.NewEncoder(s.file)
    return e.Encode(record{key, url})
}
```

`goto`程序启动时,磁盘上存储的数据必须读取到`URLStore`的map中,为此我们编写`load`方法:
```
func (s *URLStore) load() error {
    if _, err := s.file.Seek(0, 0); err != nil {
        return err
    }
    d := gob.NewDecoder(s.file)
    var err error
    for err == nil {
        var r record
        if err = d.Decode(&r); err == nil {
            s.Set(r.Key, r.URL)
        }
    }
    if err == io.EOF {
        return nil
    }
    return err
}
```

这个新的`load`方法会寻址到文件的起始位置,读取并编码每一条记录,然后用`Set`方法将数据存储到map中,再次注意无处不再的错误处理模式.文件的编码由一个无限循环完成,只要没有错误就会一直继续:
```
for err == nil {
    …
}
```

如果得到了一个错误,可能时刚编码了最后一条记录,于是产生了`io.EOF`错误,若并非此种错误,表示产生了解码错误,用`return err`来返回它,对该方法的调用必须假如到`NewURLStore`中:
```
func NewURLStore(filename string) *URLStore {
    s := &URLStore{urls: make(map[string]string)}
    f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal("Error opening URLStore:", err)
    }
    s.file = f
    if err := s.load(); err != nil {
        log.Println("Error loading data in URLStore:", err)
    }
    return s
}
```

同时在`Put`方法中,当新的URL对加入到map中,也应该立即将他们保存到数据文件中:
```
func (s *URLStore) Put(url string) string {
    for {
        key := genKey(s.Count())
        if s.Set(key, url) {
            if err := s.save(key, url); err != nil {
                log.Println("Error saving to URLStore:", err)
            }
            return key
        }
    }
    panic("shouldn’t get here")
}
```

编译并测试这第二个版本的程序,或直接使用现有的可执行程序,验证关闭服务器并重启后,短URL仍然有效,goto程序第一次启动时,文件store.gob还不存在,因此当载入数据时会得到错误:`2011/09/11 11:08:11 Error loading URLStore: open store.gob: The system cannot find the file specified.`

结束进程并重启后,就能正常工作了,或者可以在goto启动前先创建空的store.gob文件

### 备注
当第二次启动goto时,可能会产生错误:`Error loading URLStore: extra data in buffer`

这是由于gob是基于流的协议,它不支持重新开始.