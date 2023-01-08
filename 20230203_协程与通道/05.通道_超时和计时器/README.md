# 通道,超时和计时器
`time`包中有一些有趣的功能可以和通道组合使用

其中就包含了`time.Ticker`结构体,这个对象以指定的时间间隔重复的向通道C发送时间值
```
type Ticker struct {
    C <-chan Time // the channel on which the ticks are delivered.
    // contains filtered or unexported fields
    ...
}
```

时间间隔是ns(纳秒,int64),在工厂函数`time.NewTicker`中以`Duration`类型的参数传入:`func NewTicker(dur) *Ticker`

在协程周期性的执行一些事情的时候非常有用

调用`Stop()`使计时器停止,在`defer`语句中使用,这些都很好的适应`select`语句
```
ticker := time.NewTicker(updateInterval)
defer ticker.Stop()
...
select {
case u:= <-ch1:
    ...
case v:= <-ch2:
    ...
case <-ticker.C:
    logState(status) // call some logging function logState
default: // no value ready to be received
    ...
}
```

`time.Tick()`函数声明为`Tick(d Duration) <-chan Time`当你想返回一个通道而不必关闭它的时候这个函数非常有用;它以d为周期给返回的通道发送时间,d是纳秒.如果需要像下边代码一样限制处理频率
```
import "time"
rate_per_sec := 10
var dur Duration = 1e9 / rate_per_sec
chRate := time.Tick(dur) // a tick every 1/10th of a second
for req := range requests {
    <- chRate // rate limit our Service.Method RPC calls
    go client.Call("Service.Method", req, ...)
}
```

这样只会按照指定频率处理请求:`chRate`阻塞了更高的频率,每秒处理的频率可以根据机器负载资源的情况而增加或减少

定时器结构体看上去和计时器结构体的确很像,但是它只发送一次时间,在`Dration d`之后,还有`time.After(d)`函数声明如下:`func After(d Duration) <-chan Time`

在`Duration d`之后,当前时间被发到返回的通道,所以它和`NewTimer(d).C`是等价的;类似`Tick()`但是`After()`只发送一次时间

示例`timer_goroutine.go`

习惯用法:简单超时模式

要从通道ch中接收数据,但是最多等待1秒,先创建一个信号通道,然后启动一个`lambda`协程,协程在给通道发送数据之前是休眠的:
```
timeout := make(chan bool, 1)
go func() {
        time.Sleep(1e9) // one second
        timeout <- true
}()
```

然后使用`select`语句接收ch或者timeout的数据,如果ch在1秒内没有收到数据,就选择到了time分支并放弃了ch的读取
```
select {
    case <-ch:
        // a read from ch has occured
    case <-timeout:
        // the read from ch has timed out
        break
}
```

第二种形式:取消耗时很长的同步调用

也可以使用`time.After()`函数替换`timeout-channel`可以在select中通过`time.After()`发送的超时信号来停止协程的执行.在`timeoutNs`纳秒后执行select的timeout分支后,执行`client.Call`的协程也随之结束,不会给通道ch返回值
```
ch := make(chan error, 1)
go func() { ch <- client.Call("Service.Method", args, &reply) } ()
select {
case resp := <-ch
    // use resp and reply
case <-time.After(timeoutNs):
    // call timed out
    break
}
```

注意缓冲大小设置为1是必要的,可以避免协程死锁以及确保超时的通道可以被垃圾回收,此外,需要注意在有多个case符合条件时,select对case的选择时伪随机的,如果上面的代码稍作修改如下,则select语句可能不会在定时器超时信号到来时立刻选中`time.After(timeoutNs)`对应的case,因此协程可能不会严格按照定时器设置的时间结束
```
ch := make(chan int, 1)
go func() { for { ch <- 1 } } ()
L:
for {
    select {
    case <-ch:
        // do something
    case <-time.After(timeoutNs):
        // call timed out
        break L
    }
}
```

第三种形式:假设程序从多个复制的数据库同时读取,只需要一个答案,需要接收首先到达的答案,query函数获取数据库的连接切片并请求,并行请求每一个数据库并返回收到的第一个相应
```
func Query(conns []Conn, query string) Result {
    ch := make(chan Result, 1)
    for _, conn := range conns {
        go func(c Conn) {
            select {
            case ch <- c.DoQuery(query):
            default:
            }
        }(conn)
    }
    return <- ch
}
```

再次声明,结果通道ch必须带缓冲,以保证第一个发送进来的数据有地方可以存放,确保放入的首个数据总会成功,所以第一个到达的值会被获取而与执行的顺序无关,正在执行的协程可以总是使用`runtime.Goexit()`来停止

在应用中缓存数据:

应用程序中用到了来自数据库的数据时,经常会把数据缓存到内存中,因为从数据库中获取数据的操作代价很高,如果数据库中的值不会发生变化就没有问题,但是如果值有变化,我们需要一个机制来周期性的从数据库重新读取这些值:缓存的值就不可用了,而且我们也不希望用户看到陈旧的数据