# 总结和增强
通过逐步构建goto应用程序,我们遇到了几乎所哟go语言特性:

虽然这个程序按照我们的目标行事,仍然有一些可改进的途径:

- 审美:用户界面可以美化.为此可以使用go的`template`包
- 可靠性:master/slave之间的rpc连接应该可以更可靠,如果客户端到服务器之间的连接中断,客户端应该尝试重连.用一个`dialer`协程可以达成
- 资源减负:由于URL数据库大小不断增长,内存占用可能会成为一个问题,可以通过多台master服务器按照键分片来解决
- 删除:要支持删除短URL,master和slave之间的交互将变得更复杂