# 误用短声明导致变量覆盖
```
var remember bool = false
if something {
    remember := true //错误
}
// 使用remember
```

在此代码段中,`remember`变量永远不会在if语句外面编程true,如果something为true,由于使用了短声明`:=`,if语句内部的新变量remember将覆盖外面的remember变量,并且该变量的值为true,但是id语句外面,变量remember的值变成了false,所以正确的写法应该是:
```
if something {
    remember = true
}
```

此类错误也容易在for循环中出现,尤其当函数返回一个具名变量时难于察觉,如下代码段:
```
func shadow() (err error) {
    x, err := check1() // x是新创建变量，err是被赋值
    if err != nil {
        return // 正确返回err
    }
    if y, err := check2(x); err != nil { // y和if语句中err被创建
        return // if语句中的err覆盖外面的err，所以错误的返回nil！
    } else {
        fmt.Println(y)
    }
    return
}
```