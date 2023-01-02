# 从panic中恢复
正如名字一样,这个内建函数被用于从panic或错误场景中恢复,让程序可以从panicking重新获得控制权,停止终止过程进而恢复正常执行

`recover`只能在defer修饰的函数中使用,用于取得panic调用中传递过来的错误值,如果是正常执行,调用`recover`会返回nil,且没有其他效果

总结:panic会导致栈被展开直到defer修饰的recover()被调用或者程序中止

该例子中protect函数调用参数g来保护调用者防止从g中抛出的运行时panic,并展示panic中的信息
```
func protect(g func()) {
    defer func() {
        log.Println("done")
        // Println executes normally even if there is a panic
        if err := recover(); err != nil {
            log.Printf("run time panic: %v", err)
        }
    }()
    log.Println("start")
    g() //   possible runtime-error
}
```

log包实现了简单的日志功能,默认的log对象向标准错误输出中写入并打印每条日志信息的日期和时间,除了`Println`和`Printf`函数,其他的致命性都会写完日志信息后调用`os.Exit(1)`那些退出函数也是如此,而panic效果的函数会在写完日志信息后调用panic,可以在程序必须中止或发生临界错误时使用它们,就像当web服务器不能启动时那样

log包用那些方法定义了一个Logger接口类型.

......

`defer-panic-recover`在某种意义上也是一种像if和for这样的控制流机制

go标准库中许多地方都用了这个机制,json包中的解码和regexp包中的Complie函数,go库的原则时即使在包的内部使用了panic,在它的对外接口中也必须用recover处理成返回显示的错误