# 协程间的信道
## 概念
在第一个例子中,协程是独立运行的,他们之间没有通信,他们必须通信才会变得更有用;彼此之间大宋和接收信息并且协调/同步他们的工作.

协程可以使用共享变量来通信,但是很不提倡这样做,因为这种方式给所有的共享内存的多线程都带来困难

而go有一种特殊的类型,通道`channel`就像一个可以用于发送类型化数据的管道,由其负责协程之间的通信,从而避开所有共享内存导致的陷阱;这种通过通道进行通信的方式保证了同步性.

数据在通道中进行传递:在任何给定时间,一个数据被设计为只有一个协程可以对齐访问,所以不会发生数据竞争.数据的所有权也因此被传递

工厂的传送带是个很有用的例子,一个机器在传送带上放置物品,另外一个机器拿到物品并打包.通道服务于通信的两个目的:值的交换,同步的,保证了两个计算任务时候都是可知状态.

通常使用这样的格式来声明通道:`var identifier chan datatype`

未初始化的通道的值是nil

所以通道只能传输一种类型的数据,比如`chan int`或者`chan string`,所有的类型都可以用于通道,空接口`interface{}`也可以,甚至可以(有时非常有用)创建通道的通道.

通道实际上是类型化消息的队列:使数据得以传输,它使先进先出(FIFO)的结构所以可以保证发送给他们的元素的顺序(有些人知道,通道可以比作UnixShells中的双向管道(two-way pipe))

通道也是引用类型,所以使用`make`函数来给他们分配内存
```
var ch1 chan string
ch1 = make(chan string)
```

当然可以更短`ch1 := make(chan string)`

我们构建一个int通道的通道:`chanOfChans := make(chan int)`或者函数通道:`funcChan := make(chan func())`

所以通道是第一类对象,可以存储在变量中,作为函数的参数传递,从函数返回以及通过通道发送它们自身,另外它们是类型化的,允许类型检查,比如尝试使用整数通道发送一个指针

## 通信操作符<-
这个操作符直观的标示了数据的传输:信息按照箭头的方向流动

流向通道(发送)
`ch <- int1`表示:用通道ch发送变量int1(双目运算符,中缀=发送)

从通道流出(接收),3种方式:
`int2 = <- ch`表示:变量int2从通道ch(一元运算的前缀操作符,前缀=接收)接收数据(获取新值):假设int2已经声明过了;如果没有的化可以写成`int2 := <- ch`
`<- ch`可以单独调用获取通道的值,当前值会被丢弃,但是可以用来验证,所以以下代码是合法的
```
if <- ch != 1000 {
    ...
}
```

同一个操作符`<-`即用于发送也用于接收,但go会根据对象弄明白该干什么.虽非强制要求,但未了可读性通道的命名通常以`ch`开头或包含`chan`通道的发送和接收都是原子操作;它们总是互不干扰完成的.示例`goroutine2.go`

## 通道阻塞
默认情况下,通信时同步且无缓冲的,在由接受者接收数据之前,发送不会结束,可以想象一个无缓冲的通道在没有空间来保存数据的时候,必须要由一个接收者准备好接收通道的数据然后发送者可以直接把数据发送给接收者.所以通道的发送/接收操作在对方准备好之前是阻塞的:
1. 对于同一个通道,发送操作(协程或函数中)在接收者准好之前是阻塞的,如果ch中的数据无人接收,就无法再给通道传入其他数据;新的输入无法在通道非空的情况下传入.所以发送操作会等待ch再次变为可用状态,就是通道值被接收时(可以传入变量)
2. 对于同一个通道,接收操作是阻塞的,知道发送者可用,如果通道中没有数据,接收者就阻塞了

尽管这看上去是非常严格是约束,实际在大部分情况下工作很不错,程序`channel_block.go`验证了以上理论,一个协程在无限循环中给通道发送整数数据,不过因为没有接收者,只输出了一个数字0

## 通过一个或多个通道交换数据进行协程同步
通信时一种同步形式,通过通道,两个协程在通信中某刻同步交换数据.无缓冲通道成为了多个协程同步的完美工具

甚至可以在通道两端互相阻塞对方,形成了叫做死锁的状态,go运行时会检查panic停止程序,死锁几乎完全时由糟糕的设计导致的

无缓冲通道会被阻塞,设计无阻塞的程序可以避免这种情况,或者使用带缓冲的通道,程序`blocking.go`解释为什么程序会导致panic,所有的协程都休眠 - 死锁!

## 同步通道-使用带缓冲的通道
一个无缓冲通道只能包含1个元素,有时明显很局限.我们给通道提供了一个缓存,可以在扩展的make命令设置它的容量如下:
```
buf := 100
ch1 := make(chan string, buf)
```

buf是通道可以同时容纳的元素(这里是string)个数

在缓冲满载之前,给一个带缓冲的通道发送数据是不会阻塞的,而从通道读取数据也不会阻塞,知道缓冲空了.

缓冲容量和类型无关,所以可以给一些通道设置不同的容量,只要他们拥有同样的元素类型,内置的`cap`函数可以返回缓冲区的容量.

如果容量大于0,通道就是异步的了,缓冲满载或变空之前通信不会阻塞,元素会按照发送的顺序被接收.如果容量是0或者未设置,通信仅在收发双方准备好的情况下才会成功

同步`ch := make(chan type, value)`
1. value == 0 -> synchronous, unbuffered(阻塞)
2. value > 0 -> asyncchronous, buffered(非阻塞)取决于value元素

若使用通道的缓冲,你的程序会在请求激增的时候表现更好:更具弹性,专业术语叫:更具伸缩性.在设计算法时首先考虑使用无缓冲通道,只在不确定的情况下使用缓冲

## 协程中用通道输出结果
为了知道计算何时完成,可以通过信道回报这样写
```
ch := make(chan int)
go sum(bigArray, ch) // bigArray puts the calculated sum on ch
// .. do something else for a while
sum := <- ch // wait for, and retrieve the sum
```

也可以使用通道来达到同步的目的,这个很有效的用法在传统计算中称为信号量,或者换个方式:通过通道发送信号告知处理已经完成

在其他协程运行时让main程序无限阻塞的通常做法是在main函数的最后放置一个`select{}`

也可以使用通道让main程序等待协程完成,就是所谓的信号量模式

## 信号量模式
协程通过在通道ch中放置一个值来处理结束的信号,main协程等待<-ch直到从中获取到值
我们期望从这个通道中获取返回的结果
```
func compute(ch chan int){
    ch <- someComputation() // when it completes, signal on the channel.
}
func main(){
    ch := make(chan int)     // allocate a channel.
    go compute(ch)        // start something in a goroutines
    doSomethingElseForAWhile()
    result := <- ch
}
```

这个信号也可以是其他的,不返回结果,比如下面这个协程中的匿名函数协程:
```
ch := make(chan int)
go func(){
    // doSomething
    ch <- 1 // Send a signal; value does not matter
}()
doSomethingElseForAWhile()
<- ch    // Wait for goroutine to finish; discard sent value.
```

或者等待两个协程完成,每一个都会对切片s的一部分进行排序,片段如下
```
done := make(chan bool)
// doSort is a lambda function, so a closure which knows the channel done:
doSort := func(s []int){
    sort(s)
    done <- true
}
i := pivot(s)
go doSort(s[:i])
go doSort(s[i:])
<-done
<-done
```

下边的代码用完整的信号量模式对长度为N的float64切片进行N个`soSomething()`计算并同时完成,通道sem分配了相同的长度,待所有的计算都完成后,发送信号.在循环中从通道sem不停的接收数据来等待所有的协程完成
```
type Empty interface {}
var empty Empty
...
data := make([]float64, N)
res := make([]float64, N)
sem := make(chan Empty, N)
...
for i, xi := range data {
    go func (i int, xi float64) {
        res[i] = doSomething(i, xi)
        sem <- empty
    } (i, xi)
}
// wait for goroutines to finish
for i := 0; i < N; i++ { <-sem }
```

注意上述代码中闭合函数的用法:`i`和`ix`都是作为参数传入闭合函数的,这一做法使得每个协程获得一份`i`和`ix`的单独拷贝,从而闭合函数内部屏蔽了外层循环中的`i`和`ix`变量,否则for循环的下一迭代会更新所有协程中`i`和`ix`的值,另一方面,切片res没有传入闭合函数,因为协程不需要res的单独拷贝,切片res也在闭合函数中但并不是参数

## 实现并行的for循环
for循环的每一个迭代是并行完成的
```
for i, v := range data {
    go func (i int, v float64) {
        doSomething(i, v)
        ...
    } (i, v)
}
```

在for循环中并行计算迭代可能带来很好的性能提升,不过所有的迭代都必须是独立完成的,有些语言比如Fortress或者其他并行框架以不同的结构实现了这种方式,在go中用协程实现起来非常容易

## 用带缓冲通道实现一个信号量
信号量是实现互斥锁常见的同步机制,限制对资源的访问,解决读写问题;比如没有实现信号量的`sync`的go包,使用带缓冲的通道可以轻松实现
1. 带缓冲通道的容量和要同步的资源容量相同
2. 通道长度与当前资源被使用的数量相同
3. 容量减去通道的长度就是未处理的资源个数

不用管通道中存放的是什么,只关注长度;因此我们创建一个长度可变容量为0的通道
```
type Empty interface {}
type semaphore chan Empty
```

将可用资源的数量N来初始化信号量`semaphore：sem = make(semaphore, N)`
然后直接对信号量进行操作
```
// acquire n resources
func (s semaphore) P(n int) {
    e := new(Empty)
    for i := 0; i < n; i++ {
        s <- e
    }
}
// release n resources
func (s semaphore) V(n int) {
    for i:= 0; i < n; i++{
        <- s
    }
}
```

可以用来实现一个互斥的例子
```
/* mutexes */
func (s semaphore) Lock() {
    s.P(1)
}
func (s semaphore) Unlock(){
    s.V(1)
}
/* signal-wait */
func (s semaphore) Wait(n int) {
    s.P(n)
}
func (s semaphore) Signal() {
    s.V(1)
}
```

习惯用法:通道工厂模式,编程中常见的另外一种模式:不将通道作为参数传递给协程,而用函数来生成一个通道并返回(工厂角色),函数内有个匿名函数被协程调用.示例`channel_idiom.go`

## 给通道使用for循环
for循环的range语句可以用在通道上,便可以从通道中获取值
```
for v := range ch {
    fmt.Printf("The value is %v\n", v)
}
```

它从指定通道中读取数据直到通道关闭,才继续执行下边的代码.很明显,另外一个协程必须写入ch而且必须在写入完成后才关闭.suck函数可以这样写,且在协程中调用这个动作,程序变成了这样

示例:`channel_idiom2.go`

习惯用法:通道迭代模式

这个模式用到了`producter_consumer.go`的生产者-消费者,通常需要从包含地址索引字段items的容器给通道填入元素,为容器的类型定义一个方法`Iter()`返回一个只读的通道items
```
func (c *container) Iter () <- chan item {
    ch := make(chan item)
    go func () {
        for i:= 0; i < c.Len(); i++{    // or use a for-range loop
            ch <- c.items[i]
        }
    } ()
    return ch
}
```

在协程里,一个for循环迭代容器c中的元素.调用这个方法的代码可以这样迭代容器:`for x := range container.Iter() { ... }`

其运行在自己启动的协程中,所以上边的迭代用到了一个通道和两个协程.这样我们就有了一个典型的生产者-消费者模式.如果在程序结束之前,向通道写值的协程未完成工作,则这个协程不会被垃圾回收;这是设计使然,这种看起来并不符合预期的行为正是由通道这种线程安全的通信方式所导致的.如此一来,一个协程为了写入一个永远无人读取的通道而被挂起就变成了一个bug,而并非你预想中的那样被悄悄回收掉

习惯用法:生产者消费者模式

假设你有Producte()函数来产生Consume函数需要的值,它们都可以运行在独立的协程中,生产者在通道中放入给消费者读取的值,整个处理过程可以替换为无线循环
```
for {
    Consume(Produce())
}
```

## 通道的方向
通道类型可以用注解来表示它只发送或者只接收
```
var send_only chan<- int         // channel can only receive data
var recv_only <-chan int        // channel can only send data
```

只接收的通道无法关闭,因为关闭通道是发送者用来表示不再给通道发送值了,所以只接收通道是没有意义的,通道创建的时候都是双向的,但也可以分配有方向的通道变量,就像这样
```
var c = make(chan int) // bidirectional
go source(c)
go sink(c)

func source(ch chan<- int){
    for { ch <- 1 }
}

func sink(ch <-chan int) {
    for { <-ch }
}
```

习惯用法:管道和选择器模式

更具体的例子还有协程处理它从通道接收的数据并发送给输出的通道
```
sendChan := make(chan int)
receiveChan := make(chan string)
go processChannel(sendChan, receiveChan)

func processChannel(in <-chan int, out chan<- string) {
    for inValue := range in {
        result := ... /// processing inValue
        out <- result
    }
}
```

通过使用方向注解来限制协程对通道的操作

示例`sieve1.go`和`sieve2.go`