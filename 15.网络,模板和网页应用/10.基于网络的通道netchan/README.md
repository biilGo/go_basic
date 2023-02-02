# 基于网络的通道netchan
go团队决定改进并重新打造`netchan`包的现有版本,已被移至`old/netchan`

`old`目录用于存放过时的包代码,它们不会成为go1.x的一部分

一项可`rpc`密切相关的技术是基于网络的通道,它们仅存在于被执行的机器内容空间中.netchan包实现了类型安全的网络化通道;它允许一个通道两端出现由网络连接的不同计算机.其实现原理是,在其中一台机器上将传输数据发送到通道中,那么就可以被另一台计算机上同类型的通道接收.一个导出器会按名称发布通道,导入器连接到导出的机器,并按名称导入这些通道,之后2台机器就可按通常的方式来使用通道.网络通道不是同步的,它们类似于带缓存的通道.

发送端示例代码如下:
```
exp, err := netchan.NewExporter("tcp", "netchanserver.mydomain.com:1234")
if err != nil {
    log.Fatalf("Error making Exporter: %v", err)
}
ch := make(chan myType)
err := exp.Export("sendmyType", ch, netchan.Send)
if err != nil {
    log.Fatalf("Send Error: %v", err)
}
```

接收端示例代码如下
```
imp, err := netchan.NewImporter("tcp", "netchanserver.mydomain.com:1234")
if err != nil {
    log.Fatalf("Error making Importer: %v", err)
}
ch := make(chan myType)
err = imp.Import("sendmyType", ch, netchan.Receive)
if err != nil {
    log.Fatalf("Receive Error: %v", err)
}
```