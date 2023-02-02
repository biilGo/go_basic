# 使用select切换协程
从不同的并发执行的协程中获取值可以通过关键字select来完成,它和switch控制语句非常相似,也被称作通信开关,它的行为像是"你准备好了么"的轮询机制;select监听进入通道的数据,也可以是用通道发送值的时候
```
select {
case u:= <- ch1:
        ...
case v:= <- ch2:
        ...
        ...
default: // no value ready to be received
        ...
}
```

default语句是可选的,fallthrough行为和普通的switch相似,是不允许的.在任何一个case中执行break或者return,select就结束了

select做的就是:选择处理列出的多个通信情况中的一个
1. 如果阻塞了会等待直到其中一个可以处理
2. 如果多个可以处理,随机选择一个
3. 如果没有通道可以处理并且写了default语句,他就会执行default永远是可运行的

在select中使用发送操作并且有default可以确保发送不被阻塞!如果没有default,select就会一直阻塞

select语句实现了一种监听模式,通常用在无限循环中,在某种情况下,通过break语句循环退出