# 新旧模型对比:任务和worker
假设我们需要处理很多任务;一个workfer处理一项任务,任务可以被定义为一个结构体
```
type Task struct {
    // some state
}
```

旧模式:使用共享内存进行同步

由各个任务组成的任务池共享内存;为了同步各个worker以及避免资源竞争,我们需要对任务池进行加锁保护
```
    type Pool struct {
        Mu      sync.Mutex
        Tasks   []*Task
    }
```

sync.Mutex可以进入该灵临界区,如果出现了同一时间多个go协程都进入了该临界区,则会产生竞争:Pool结构就不能保证被正确更新.在传统的模式中,worker代码可能这样写
```
func Worker(pool *Pool) {
    for {
        pool.Mu.Lock()
        // begin critical section:
        task := pool.Tasks[0]        // take the first task
        pool.Tasks = pool.Tasks[1:]  // update the pool of tasks
        // end critical section
        pool.Mu.Unlock()
        process(task)
    }
}
```

这些worker有许多都可以并发执行;他们可以在go协程中启动,一个worker先将pool锁定,从pool获取第一项任务,再解锁和处理任务.加锁保证了同一时间只有一个go协程可以进入到pool中:一项任务有且只能被赋予一个worker.如果不加锁,则工作协程可能会在`task := pool.Task[0]`发生切换,导致`pool.Tasks=pool.Tasks[1:]`结果异常:一些worker获取不到任务,而一些任务可能被多个woker得到.加锁实现同步的方式在工作协程比较少时可以工作的很好,但是当工作协程数量很大,任务量也很多时,处理效率将会因为频繁的加锁\解锁开销而降低,当工作协程数量增加到一个阈值时,程序效率会急剧下降,这就成立瓶颈

新模式:使用通道

使用通道进行同步:使用要给通道接受需要处理的任务,一个通道接受处理完成的任务.worker在协程中启动,其数量N应该根据任务数量进行调整

主线程扮演Master节点角色,可能协程如下形式
```
    func main() {
        pending, done := make(chan *Task), make(chan *Task)
        go sendWork(pending)       // put tasks with work on the channel
        for i := 0; i < N; i++ {   // start N goroutines to do work
            go Worker(pending, done)
        }
        consumeWork(done)          // continue with the processed tasks
    }
```

worker的逻辑比较简单,从pending通道拿任务,处理后放到done通道中
```
    func Worker(in, out chan *Task) {
        for {
            t := <-in
            process(t)
            out <- t
        }
    }
```

这里并不是使用锁:从通道得到新任务的过程没有任务竞争,随着任务数量增加,worker数量也应该相应增加,同时性能并不会像第一种方式那样下降明显,在pending通道中存在一份任务的拷贝.第一个worker从pending通道中获得第一个任务并进行处理,这里并不存在竞争.某一个任务会在哪一个worker中被执行时不可知的,反过来也是,worker数量的增加也会增加通信的开销.这会对性能有轻微的影响

从这个简单的例子中可能很难看出第二种模式的优势,但含有复杂锁运用的程序不仅在编写上显得困难,也不容易编写正确,使用第二种模式的话,就无需考虑这么复杂的东西了

因此第二种模式对比第一种模式而言,不仅性能时一个主要优势,而且还有更大的优势,代码显得更清晰,一个更符合go语言习惯的worker写法
```
    func Worker(in, out chan *Task) {
        for {
            t := <-in
            process(t)
            out <- t
        }
    }
```

对于任何可以建模为master-worker范例的问题,一个类似于worker使用通道进行通信和交互,master进行调整协调的方案都能完美解决,如果系统部署在多台机器上,各个机器上执行worker协程,master和worker之间使用netchan或者RPC进行通信

如何现在该使用锁还是通道?

通道是一个较新的概念,我们着重强调了go协程里通道的使用,但这并不意味着经典的锁方法就不能使用.go语言让你可以根据实际问题进行选择,创建一个优雅,简答,可读性强,在大多数场景性能表现都很好的方案,如果你的问题适合使用锁,也不要忌讳使用它.go语言注重实用,什么方式最能解决你的问题就用什么方式,而不是强迫你使用一种编码风格

一个普遍的经验法则:
1. 使用锁的情况:

- 访问共享数据结构中的缓存信息
- 保存应用程序上下文和状态信息数据

2. 使用通道的情况:

- 以异步操作的结果进行交互
- 分发任务
- 传递数据所有权