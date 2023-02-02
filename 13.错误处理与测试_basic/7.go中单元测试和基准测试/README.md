# go中的单元测试和基准测试
首先所有的包都应该有一定的必要文档,然后同样重要的是对包的测试.

名为`testing`的包被专门用来进行自动化测试,日志和错误报告,并且还包含一些基准测试函数的功能

gotest是unix bash脚本,所以在windows下需要配置MINGGW在windows环境下把所有的`pkg/linux_amd64`替换成`pkg/windows`

对一个包错测试,需要写一些可以频繁执行的小块测试单元来检查代码的正确性,我们必须写一些go源文件来测试代码,测试程序必须属于被测试的包,并且文件名满足这种形式`*_test.go`所以测试代码和包中的业务代码是分开的

_test程序不会被普通的go编译器编译,所以当放应用部署到生产环境时它们不会被部署,只有gotest会编译所有的程序,普通程序和测试程序

测试文件中必须导入"tesing"包,并写一些名字以`TestZzz`打头的全局函数,这里的Zzz是被测试函数的字母描述.

测试函数必须有这种形式的头部
> func TestAbcde(t *testing.T)

T是传给测试函数的结构类型,用来管理测试状态,支持格式化测试日志,在函数的结尾把输出跟想要的结果对比,如果不等就打印一个错误,成功的测试则直接返回

用下面的这些函数来通知测试失败
```
func (t *T) Fail()
// 标记测试函数为失败，然后继续执行（剩下的测试）
```

```
func (t *T) FailNow()
// args 被用默认的格式格式化并打印到错误日志中
```

```
func (t *T) Fatal(args ...interface{})
// 结合 先执行 3），然后执行 2）的效果
```

运行go test来编译测试程序,并执行程序中所有的TestZZZ函数,如果所有的测试都通过会打印除PASS

gotest可以接收一个或多个函数程序作为参数,并指定一些选项

结合`chatty`或-v选项,每个执行的测试函数以及测试状态会被打印

```
go test fmt_test.go --chatty
=== RUN fmt.TestFlagParser
--- PASS: fmt.TestFlagParser
=== RUN fmt.TestArrayPrinter
--- PASS: fmt.TestArrayPrinter
...
```

testing包中有一些类型和函数可以用来做简单的基准测试;测试代码中必须包含以`BenchmarkZzz`打头的函数并接收一个`*testing.B`类型的参数
```
func BenchmarkReverse(b *testing.B) {
    ...
}
```

命令`go test –test.bench=.*`会运行所有的基准测试函数,代码中的函数会被调用N次并战术N的值和函数执行的平均时间.如果是用testing.Benchmark调用这些函数,直接运行程序即可