# 自定义包中的错误处理和panicking
这是所有内部包实现者应该遵守的最佳实践

- 在包内部,总是应该从panic和recover;不允许显示的超出包范围的panic()
- 向包的调用者返回错误值

在包内部,特别时非导出函数中有很深层次的嵌套调用时,将panic转换成error来告诉调用方为何出错,是很使用的

当没有东西需要转换或者转换成整数失败时,这个包会panic,但是可导出的Parse函数会从panic中recover并用所有这些信息返回一个错误给调用者,在`panic_recover.go`中调用了parse包.不可解析的字符串导致错误并被打印出来

......

