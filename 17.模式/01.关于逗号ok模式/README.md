# 逗号ok模式
我们经常在一个表达式返回2个参数时使用这种模式:`,ok`第一个参数是一个值或者`nil`第二个参数是`true/false`或者一个错误`error`,在一个需要赋值的if条件语句中,使用这种模式去检测第二个参数值会让代码显得优雅简洁,这种模式在go语言编码规范中非常重要

1. 在函数返回时检测错误
```
value, err := pack1.Func1(param1)
if err != nil {
    fmt.Printf("Error %s in pack1.Func1 with parameter %v", err.Error(), param1)
    return err
}
// 函数Func1没有错误:
Process(value)
e.g.: os.Open(file) strconv.Atoi(str)
```

这段代码中函数将错误返回给它的调用者,当函数执行成功时,返回的错误是nil,所以使用这种写法:
```
func SomeFunc() error {
    …
    if value, err := pack1.Func1(param1); err != nil {
        …
        return err
    }
    …
    return nil
}
```

这种模式也常用于defer使程序从panic中恢复执行

要实现简洁的错误检测代码,更好的方式是使用闭包:

2. 检测映射中是否存在一个键值,key1在映射map1中是否有值?
```
if value, isPresent = map1[key1]; isPresent {
        Process(value)
}
// key1不存在
…
```

3. 检测一个接口类型变量,varI是否包含了类型T:类型断言
```
if value, ok := varI.(T); ok {
    Process(value)
}
// 接口类型varI没有包含类型T
```

4. 检测一个通道ch是否关闭
```
    for input := range ch {
        Process(input)
    }
```

或者
```
    for {
        if input, open := <-ch; !open {
            break // 通道是关闭的
        }
        Process(input)
    }
```