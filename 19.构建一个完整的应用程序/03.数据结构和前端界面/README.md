# 版本1-数据结构和前端界面
## 数据结构
当程序运行在生产环境时,会收到很多短网址的请求,同时会有一些将长URL转换成短URL的请求,我们的程序要以什么样的结构存储这些数据?(A)和(B)两种URL都是字符串,此外它们相互关联:给定键(B)能获取到值(A),他们互相映射(map).要将这些数据存储在内存中,我们需要这种结构,它们几乎存在于所有的编程语言中,只是名称有所不同,例如"哈希表"或"字典"

go语言就有这种内建的映射(map):`map[string]string`

键的类型写在`[`和`]`之间,紧接着是值的类型.为特定类型指定一个别名在严谨的程序中非常实用.go语言中通过关键字`type`来定义,因此有定义:`type URLStore map[string]string`

它从短URL映射到长URL,两者都是字符串.

要创建那种类型的变量,并命名为m,使用:`m := make(URLStore)`

假设http://goto/a映射到http://google.com,我们要把它们存储到m中,可以用如下语句:`m["a"] = "http://google.com/"`.(键只是http://goto/的后缀,其前缀是不变的)

要获得给定a对应的长URL,可以这么写:`url := m["a"]`

此时`url`的值等于http://google.com

注意:使用了`:=`就不需要指明url的类型为string,编译器会从右侧的值中推断出来

## 使用程序线程安全
这里,变量`URLStroe`是中心化的内存存储,当收到网络流量时,会有很多`Redirect`服务的请求,这些请求其实只涉及读操作:已给定的短URL作为键,返回对应长URL的值.然而,对Add服务的请求则大不相同,它们会更改`URLStroe`,添加新的键值对,当在瞬间收到大量更新请求时,可能会产生如下问题:
1. 添加操作可能被另一个同类请求打断,写入的长URL值可能会丢失;
2. 读取和更改同时进行时,导致可能读到脏数据.代码中的map并不保证当开始更新数据时,会彻底阻止另一个更新操作的启动.
3. 也就是说,map不是线程安全的,goto会并发地为很多请求提供服务.

因此必须使用`URLStore`是线程安全的,以便可以从不同的线程访问它,最简单和经典的方法是为其增加一个锁,它是go标准库`sync`包中的`Mutex`类型,必须导入到我们的代码中

我们把`URLStore`类型的定义更改为一个结构体,它包含了两个字段`map`和`sync`包的`RWNUtex`
```
import "sync"
type URLStore struct {
    urls map[string]string        // map from short to long URLs
    mu sync.RWMutex
}
```

`RWMutex`有两种锁:分别对应读和写,多个客户端可以同时设置读锁,但只有一个客户端可以设置写锁,有效的串行化变更,使他们按顺序生效.

我们将在`get`函数中实现`Redirect`服务的读请求,在set函数中实现`Add`服务的写请求,`Get`函数类似下面这样:
```
func (s *URLStore) Get(key string) string {
    s.mu.RLock()
    url := s.urls[key]
    s.mu.RUnlock()
    return url
}
```

函数按照键返回对应映射后的URL,它所处理的变量是指针类型,指向`URLStore`,但在读取值之前,先用`s.mu.Rlock()`放置一个读锁,这样就不会有更新操作妨碍读取.数据读取后撤销锁定,以便挂起的更新操作可以开始,如果键不存在于map中会怎样?会返回字符串的零值(空字符串).注意点号类似面向对象的语言:在s的mu字段上调用方法RLock()

Set函数同时需要URL的键值对,且必须放置写锁`Lock()`来排除同一时刻任何其他更新操作.函数返回布尔值true或false来表示Set操作是否成功:
```
func (s *URLStore) Set(key, url string) bool {
    s.mu.Lock()
    _, present := s.urls[key]
    if present {
        s.mu.Unlock()
        return false
    }
    s.urls[key] = url
    s.mu.Unlock()
    return true
}
```

形式`_,present := s.urls[key]`可以测试map中是否已经包含该键,包含则`present`为true,否则为`false`这种形式称为`逗号ok模式`,在go代码中会频繁出现,如果键已存在,Set函数直接返回布尔值`false`,map不会被更新,如果键不存在,把它加入map中并返回true,左侧`_`是一个值的占位符,赋值给`_`来表明我们不会使用它,注意在更新后尽早调用`Unlock()`来释放对`URLStroe`的锁定

## 使用defer简化代码
目前代码还比较简单,容易记得操作完成后调用`Unlock()`解锁,然而在代码更复杂时容易忘记解锁,或者放置在错误的位置,往往导致问题很难追踪,对于这种情况Go提供了一个特殊关键字`defer`.可以在Lock之后立即示意`Unlock`不过其效果时`Unlock()`只会在函数返回之前被调用

Get可以简化成以下代码
```
func (s *URLStore) Get(key string) string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.urls[key]
}
```

Set的逻辑在某种程度上也变得清晰
```
func (s *URLStore) Set(key, url string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    _, present := s.urls[key]
    if present {
        return false
    }
    s.urls[key] = url
    return true
}
```

## URLStore工厂函数
`URLStore`结构体中包含map类型的字段,使用前必须先用make初始化,在Go中创建一个结构体实例,一般是通过定义一个前缀为`New`,能返回该类型已初始化实例的函数
```
func NewURLStore() *URLStore {
    return &URLStore{ urls: make(map[string]string) }
}
```

在`return`语句中,创建了`URLStore`字面量实例,其中包含初始化了的map映射,锁无需特别指明初始化,这是go创建结构体实例的惯例.`&`是取址运算符,它将我们要返回的内容变成指针,因为`NewURLStore`返回类型是`*URLStore`,然后调用该函数来创建`URLSDtore`变量:`var store = NewURLStore()`

## 使用URLStore
要新增一对短\长URL到map中,我们只需调用s上的Set方法,由于返回布尔值,可以把它包裹在if语句中:
```
if s.Set("a", "http://google.com") {
    // 成功
}
```

要获取给定短URL对应的长URL,调用s上的Get方法,将返回值放入变量url:
```
if url := s.Get("a"); url != "" {
    // 重定向到 url
} else {
    // 键未找到
}
```

这里我们利用go语言if语句的特性,可以在起始部分,条件判断前放置初始化语句,另外还需要一个Count方法以获取map中键值对的数量,可以使用内建的len函数:
```
func (s *URLStore) Count() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return len(s.urls)
}
```

如何根据给定的长URL计算出短URl?为此我们创建一个函数`genKey(n int) string{...}`,将`s.Count()`的当前值作为其整型参数传入.

现在,我们可以创建一个Put方法,接收一个长URL用genKey生成其短URL键,调用Set方法在此键下存储长URL数据,然会返回这个键:
```
func (s *URLStore) Put(url string) string {
    for {
        key := genKey(s.Count())
        if s.Set(key, url) {
            return key
        }
    }
    // shouldn’t get here
    return ""
}
```

for循环会已知尝试调用Set直到成功为止,现在我们定义好了数据存储,以及配套的可工作函数,但这本身并不能完成任务,我们还需要开发web服务器以交付Add和Redirect服务