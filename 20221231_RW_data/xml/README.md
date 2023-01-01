# XML数据格式
如同json包一样,也有`Marshal()`和`UnMarshal()`从xml中编码和解码数据,但这个更通用,可以从文件中读取和写入
和json的方式一样,xml数据可以被序列化为结构,或者从结构反序列化为xml数据

程序`xml.go`利用encoding/xml包实现了一个简单的xml解析器,用来解析xml数据内容

程序中定义了若干xml标签类型:`StartElement`,`Chardata`,`EndElement`,`Comment`,`Directive`,`Proclnst`
保重同样定义了一个结构解析器:`NewParse`方法持有一个`io.Reader`并生成一个解析器类型的对象.还有一个`Token()`方法返回输入流里的下一个xml token.在输入流的结尾处,会返回`nil`,`io.EOF`

xml文本被循环处理知道`Token()`返回一个错误,因为已经达到文件尾部,再没有内容可供处理.通过一个`type-switch`可以根据一些xml标签进一步处理.Chardata中的内容只是一个`[]byte`通过字符串转换让其变得可读性强一些.