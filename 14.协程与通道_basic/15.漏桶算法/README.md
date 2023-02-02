# 漏桶算法
考虑以下客户端-服务器结构:客户端协程执行一个无线循环从某个源头接收数据;数据读取到Buffer类型的缓冲区,为了避免分配过多的缓冲区以及释放缓冲区,它保留了一份空闲缓冲区列表,并且使用一个缓冲通道来表示这个列表:`var freeList = make(chan *Buffer,100)`

这个可重用的缓冲区队列与服务器是共享的,当接收数据时,客户端尝试从`freeList`获取缓冲区;但如果此时通道为空,则会分配新的缓冲区,一旦消息被加在后,它将被发送到服务器上的`serverChan`通道:`var serverChan = make(chan *Buffer)`

以下是客户端算法代码:
```
 func client() {
    for {
        var b *Buffer
        // Grab a buffer if available; allocate if not 
        select {
            case b = <-freeList:
                // Got one; nothing more to do
            default:
                // None free, so allocate a new one
                b = new(Buffer)
        }
        loadInto(b)         // Read next message from the network
        serverChan <- b     // Send to server
    }
 }
```

服务器的循环则接收每一条来自客户端的消息并处理它,之后尝试将缓冲区返回给共享的空闲缓冲区:
```
func server() {
    for {
        b := <-serverChan       // Wait for work.
        process(b)
        // Reuse buffer if there's room.
        select {
            case freeList <- b:
                // Reuse buffer if free slot on freeList; nothing more to do
            default:
                // Free list full, just carry on: the buffer is 'dropped'
        }
    }
}
```

但是这种方法在`freeList`通道已满的时候是行不通的,因为无法放入空闲`freeList`通道的缓冲区会被丢到地上由垃圾收集器回收故名:漏桶算法