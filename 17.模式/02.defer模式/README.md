# defer模式
使用defer可以确保资源不再需要时,都会被恰当地关闭或归还到"池子"中,更重要的一点是,它可以恢复panic

1. 关闭一个文件流:
```
// 先打开一个文件f
defer f.Close()
```

2. 解锁一个被锁定的资源:
```
mu.Lock()
defer mu.Unclock()
```

3. 关闭一个通道
```
ch := make(chan float64)
defer close(ch)
```

也可以是两个通道
```
answera, answerb := make(chan int), make(chan int)
defer func() {
    close(answera);
    close(answerB)
}()
```

1. 从panic恢复
```
defer func() {
    if err := recover();
    err != nil {
        log.Printf("run time panic:%v",err)
    }
}()
```

2. 停止一个计时器
```
tick1 := time.NewTicker(updateInterval)
defer tick1.Stop()
```

3. 释放一个进程p:
```
p, err := os.StartProcess(..., ..., ...)
defer p.Release()
```

4. 停止CPU性能分析并立即写入:
```
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()
```

当然defer也可以在打印报表时避免忘记输出页脚