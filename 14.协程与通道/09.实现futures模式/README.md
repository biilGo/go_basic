# 实现futures模式
所谓futures就是指:有时候在你使用某一个值之前需要先对其进行计算,这种情况下,你就可以在另一个处理器上进行该值的计算,到使用时该值就已经计算完毕了

futures模式通过闭包和通道可以很容易实现,类似于生成器,不同地方在于futures需要返回一个值

假设我们有一个矩阵类型,我们需要计算俩个矩阵A和B乘积的逆,首先我们通过函数`Inverse(M)`分别对其进行求逆运算,再将结果相乘,如下函数`InverseProduct()`实现了如上过程
```
func InverseProduct(a Matrix, b Matrix) {
    a_inv := Inverse(a)
    b_inv := Inverse(b)
    return Product(a_inv, b_inv)
}
```

例子中,a和b的求逆矩阵需要先被计算,那么为什么再计算b的逆矩阵时,需要等待a的逆计算完成呢?显然不必要,这两个求逆运算可以并行执行,换句话说,调用`Product`函数只需要等到`a_inv`和`b_inv`的计算完成.

如下代码实现了并行计算方式:
```
func InverseProduct(a Matrix, b Matrix) {
    a_inv_future := InverseFuture(a)   // start as a goroutine
    b_inv_future := InverseFuture(b)   // start as a goroutine
    a_inv := <-a_inv_future
    b_inv := <-b_inv_future
    return Product(a_inv, b_inv)
}
```

`InversFuture`函数以`goroutine`的形式起了一个闭包,该闭包会将矩阵求逆结果放入到future通道中
```
func InverseFuture(a Matrix) chan Matrix {
    future := make(chan Matrix)
    go func() {
        future <- Inverse(a)
    }()
    return future
}
```

当开发一个计算密集型库时,使用Futures模式设计API接口是很有意义的,在你的包使用Futures模式,且能保持友好的API接口,此外Futures可以通过要给异步的APPI暴露出来,这样你可以以最小的成本将包中的并行计算移到用户代码中