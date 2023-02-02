# 协程和恢复
一个用到`recover`的程序停掉了服务器内部一个失败的协程而不影响其他协程的工作
```
func server(workChan <-chan *Work) {
    for work := range workChan {
        go safelyDo(work)   // start the goroutine for that work
    }
}
func safelyDo(work *Work) {
    defer func() {
        if err := recover(); err != nil {
            log.Printf("Work failed with %s in %v", err, work)
        }
    }()
    do(work)
}
```

如果`do(work)`发生panic,错误会被记录且协程会推出并释放,而其他协程不受影响

因为recover总是返回nil,除非直接在defer修饰的函数中调用,defer修饰的代码可以调用那些自身可以使用`panic`和`recover`避免失败的库例程.

`safeyDo()`中`defer`修饰的函数可能在调用`recover`之前就调用了要给`logging`函数,`pancking`状态不会影响`logging`代码的运行,因为加入了恢复模式,函数`do`可以通过调用panic来摆脱不好的情况.但是恢复在panicking的协程内部,不能被另外一个协程恢复