# 用json持久化存储
## 版本4-用json格式存储
第四个版本:goto_v4

## 以json格式存储
如果你是个敏锐的测试者也许已经注意到,当goto程序启动2次,第二次启动后能读取短URL且完美地工作,然而从第三次开始,会得到错误:`Error loading URLStore: extra data in buffer`

这是由于gob是基于流的协议,它不支持重新开始,为补救该问题,这里我们使用json作为存储协议,它以纯文本形式存储数据,因此也可以被非go语言编写的进程读取,同时也显示了更换一种不同的持久化协议是多么简单,因为与存储打交道的代码被清晰地隔离在2个方法中,即`load`和`saveLoop`

从创建新的空文件store.json开始,更改main.go中声明文件名变量的那一行:`var dataFile = flag.String("file", "store.json", "data store file name")`

在store.go中导入json取代gob,然后再saveLoop中唯一需压被修改的行:`e := gob.NewEncoder(f)`

更改为:
`e := json.NewEncoder(f)`

类似的,在`load`方法中:`d := gob.NewDecoder(f)`

修改为:`d := json.NewDecoder(f)`

这就是所有要改动的地方!编译,启动并测试,你会发现之前的错误不会在发生