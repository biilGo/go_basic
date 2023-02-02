# 性能调试:分析并优化go程序
## 时间和内存消耗
可以用这个编写脚本xtime来测量:
```
#!/bin/sh
/usr/bin/time -f '%Uu %Ss %er %MkB %C' "$@"
```

## 用go test调试
如果代码使用了Go中的testing包的基准测试功能,我们可以用gotest标准的`-cpuprofile`和`-memprofile`标志向指定文件写入CPU或内存使用情况报告

使用方式:
> go test -x -v -cpuprofile=prof.out -file x_test.go

编译执行x_test.go中的测试,并向prof.out文件中写入cpu性能分析信息

## 用pprof调试
你可以在单机程序progexec中引入`runtime/pprof`包,这个包以pprof可视化工具需要的格式写入运行时报告数据,对于cpu性能分析来说你需要添加一些代码
```
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
...
```

代码定义了一个名为cpuprofile的flag,调用go flag库来解析命令行flag,如果命令行设置了cpuprofile flag,则开始cpu性能分析并把结果重定向到那个文件.这个分析程序最后在程序推出之前调用StopCPUProfile来刷新挂起的写操作到文件中;我们用defer来保证这一切会在main返回时触发

现在用这个flag运行程序: `progexec -cpuprofile=progexec.prof`

然后可以像这样用gopprof工具: `gopprof progexec progexec.prof`

gopprof程序时google pprof C++分析器的一个轻微变种,关于此工具更多信息:https://github.com/gperftools/gperftools

如果开启CPU性能分析,go程序会以大约每秒100次的频率阻塞,并记录当前执行的goroutine栈上的程序计数器样本

此工具一些有趣的命令

1. topN

用来展示分析结果中最开头的N份样本,top5会展示在程序运行期间调用最频繁的5个函数
```
Total: 3099 samples
626 20.2% 20.2% 626 20.2% scanblock
309 10.0% 30.2% 2839 91.6% main.FindLoops
...
```

第五列表示函数的调用频度

2. web或web函数名

该命令生成一份SVG格式的分析数据图表,并在网络浏览器中打开它,函数被表示成不同的矩形,箭头指示函数调用链.

3. list函数名或weblist函数名

展示对应函数名的代码行列表,第2列表示当前执行消耗的时间,这样就很好的指出了运行过程中消耗最大的代码

如果发现函数`runtime.mallocgc`调用频繁,那么时应该进行内存分析的时候了,找出垃圾回收频繁执行的原因,和内存大量分配的根源

为了做到这一点必须在合适的地方添加下面的代码
```
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
...
CallToFunctionWhichAllocatesLotsOfMemory()
if *memprofile != "" {
    f, err := os.Create(*memprofile)
    if err != nil {
        log.Fatal(err)
    }
    pprof.WriteHeapProfile(f)
    f.Close()
    return
}
```

用`-memprofile flag`运行这个程序: `progexec -memprofile=progexec.mprof`

然后可以像这样再次使用gopprof工具: `gopprof progexec progexec.mprof`

top5,list函数名,等命令同样适用,只不过现在是以Mb为单位测量内存分配情况,这是top命令输出的例子
```
Total: 118.3 MB
    66.1 55.8% 55.8% 103.7 87.7% main.FindLoops
    30.5 25.8% 81.6% 30.5 25.8% main.*LSG·NewLoop
    ...
```

从第一列可以看出,最上面的函数占用了最多的内存

同样有一个报告内存分配计数的有趣工具`gopprof --inuse_objects progexec progexec.mprof`

对于web应用来说,有标准的http接口可以分析数据,在http服务中添加`import _ "http/pprof"`

会在` /debug/pprof/`下的一些url安装处理器,然后你可以用一个唯一的参数,你服务中的分析数据的url来执行gopprof命令,它会下载并执行在线分析
```
gopprof http://localhost:6060/debug/pprof/profile # 30-second CPU profile
gopprof http://localhost:6060/debug/pprof/heap # heap profile
```