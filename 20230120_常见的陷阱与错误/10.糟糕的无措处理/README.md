# 糟糕的错误处理

## 不要使用布尔值
像下面代码一样,创建一个布尔值变量用于测试错误条件是多余的
```
var good bool
    // 测试一个错误，`good`被赋为`true`或者`false`
    if !good {
        return errors.New("things aren’t good")
    }
```

立即检测一个错误
```
... err1 := api.Func1()
if err1 != nil { … }
```

## 避免错误检测使代码变得混乱
避免写出这样的代码
```
... err1 := api.Func1()
if err1 != nil {
    fmt.Println("err: " + err.Error())
    return
}
err2 := api.Func2()
if err2 != nil {
...
    return
}
```

首先,包括在一个初始化的if语句中对函数的调用,但即使代码中到处都是以if语句的形式通知错误,通过这种方式,很难分辨什么使正常的程序逻辑,什么是错误检测或错误通知.还需注意的是,大部分代码都是致力于错误的检测,通常解决此问题的好办法是尽可能以闭包的形式封装你的错误检测
```
func httpRequestHandler(w http.ResponseWriter, req *http.Request) {
    err := func () error {
        if req.Method != "GET" {
            return errors.New("expected GET")
        }
        if input := parseInput(req); input != "command" {
            return errors.New("malformed command")
        }
        // 可以在此进行其他的错误检测
    } ()
        if err != nil {
            w.WriteHeader(400)
            io.WriteString(w, err)
            return
        }
        doSomething() ...
```
这种方法可以很容易分辨出错误检测,错误通知和正常的程序逻辑