# 在多核心上并行计算
假设我们有`NCPU`个CPU核心:`const NCPU = 4 //对应一个四核处理器`然后我们想把计算量分成NCPU个部分,每一部分都和其他部分并行运行.

```
func DoAll(){
    sem := make(chan int, NCPU) // Buffering optional but sensible
    for i := 0; i < NCPU; i++ {
        go DoPart(sem)
    }
    // Drain the channel sem, waiting for NCPU tasks to complete
    for i := 0; i < NCPU; i++ {
        <-sem // wait for one task to complete
    }
    // All done.
}
func DoPart(sem chan int) {
    // do the part of the computation
    sem <-1 // signal that this piece is done
}
func main() {
    runtime.GOMAXPROCS(NCPU) // runtime.GOMAXPROCS = NCPU
    DoAll()
}
```

1. `DoAll()`函数创建了一个`sem`通道,每个并行计算都将在对其发送完成信号;在一个for循环中NCPU个协程被启动了,每个协程会承担1/NCPU的工作量,每一个`DoPart()`协程都会向sem通道发送完成信号
2. `DoAll()`会在for循环中等待NCPU个协程完成:sem通道就像一个信号量,这份代码展示了一个经典的信号量模式