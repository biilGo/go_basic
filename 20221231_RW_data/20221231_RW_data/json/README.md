# JSON数据格式
数据结构要在网络中传输或保存到文件，就必须对其编码和解码；目前存在很多编码格式：
    - JSON
    - XML
    - gob
    - Google
    - 缓冲协议等等
    Go语言支持所有这些编码格式

结构可能包含二进制数据，如果将其作为文本打印，那么可读性是很差的，另外结构内部可能包含匿名字段，而不清楚数据的用意

通过把数据转换成纯文本，使用命名的字段来标注，让其具有可读性。这样的数据格式可以通过网络传输，而且是与平台无关的，任何类型的应用能够读取和输出，不与操作系统和编程语言的类型相关

下面是一些术语说明：

- 数据结构 --> 指定格式 = 序列化或编码（传输之前）
- 指定数据 --> 数据格式 = 反序列化或解码（传输之后）

序列化是在内存中把数据转换成指定格式（data --> string），反之亦然（string --> data strcuture）
编码也是一样，只是输出一个数据流；解码是从一个数据流输出到一个数据结构

我们都比较熟悉XML格式，但是有些时候JSON被作为首选，主要是由于其格式上非常简洁
通常JSON被用于web后端和浏览器之间的通讯，但是在其他场景同样的有用

简短的JSON片段：
{
    "Person": {
        "FirstName": "Laura"
        "LastName": "Lynn"
    }
}

尽管XML被广泛的应用，但是JSON更加简洁、轻量和更好的可读性，这也使它很受欢迎

##### JSON与Go类型对应如下：
- bool对应JSON的boolean
- float64对应JSON的number
- string对应JSON的string
- nil对应JSON的null

不是所有的数据都可以编码为JSON类型：只有验证通过的数据结构才能被编码：

- json对象只支持字符串类型的key，要编码一个Go map类型，map必须是map[string]T(T是json包中支持的任何类型)
- channel，复杂类型和函数类型不能被编码
- 不支持循环数据结构，它将引起序列化进入一个无限循环
- 指针可以被编码，实际上是对指针指向的值进行编码（或者指针是nil）

## 反序列化
UnMarshal()的函数签名是:
> func Unmarshal(data []byte, v interface{}) error 

把json解码为数据结构。

虽然反射能够让json字段去尝试匹配目标结构字段,但是只有真正匹配上的字段才会填充数据.字段没有匹配不会报错,而是直接忽略掉

## 解码任意的数据
json包使用`map[string]interface{}`和`[]interface{}`

存储任意的json对象和数组,其可以被反序列化为任何我的json blob存储到接口之中
来看这个json数据,被存储在变量b中:
> b := []byte(`{"Name": "Wednesday", "Age": 6, "Parents": ["Gomez", "Morticia"]}`)

不用理解这个数据的结构,我们可以直接使用Unmarshal把这个数据编码并保存在接口值中:
>var f interface{}

>err := json.Unmarshal(b, &f)

f指向的值是一个map,key是一个字符串,value是自身存储作为空接口类型的值:
```
map[string]interface{} {
    "Name": "Wednesday",
    "Age":  6,
    "Parents": []interface{} {
        "Gomez",
        "Morticia",
    },
}
```

要访问这个数据,我们可以使用类型断言: 
>m := f.(map[string]interface{})

我们可以通过for range语法和type switch来访问其实际类型:
```
for k, v := range m {
    switch vv := v.(type) {
    case string:
        fmt.Println(k, "is string", vv)
    case int:
        fmt.Println(k, "is int", vv)
    case []interface{}:
        fmt.Println(k, "is an array:")
        for i, u := range vv {
            fmt.Println(i, u)
        }
    default:
        fmt.Println(k, "is of a type I don’t know how to handle")
    }
}
```

通过这种方式,你可以处理未知的json数据,同时可以确保类型安全

## 解码数据到结构
如果我们事先知道json数据,我们可以定义一个适当的结构并对json数据反序列化
```
type FamilyMember struct {
    Name    string
    Age     int
    Parents []string
}
```

并对其反序列化
```
var m FamilyMember
err := json.Unmarshal(b, &m)
```

程序实际上是分配了一个新的切片,这是一个典型的反序列化引用类型

## 编码和解码流
json包提供Decoder和Encoder类型来支持常用json数据流读写,NewDecoder和NewEncoder函数分别封装了`io.Reader`和`io.Writer`接口
```
func NewDecoder(r io.Reader) *Decoder
func NewEncoder(w io.Writer) *Encoder
```

要想把json直接写入文件，可以使用json.NewEncoder初始化文件，并调用Encode()；反过来与其对应的是使用`json.NewDecoder`和`Decode()`函数
```
func NewDecoder(r io.Reader) *Decoder
func (dec *Decoder) Decode(v interface{}) error
```

来看下接口是如何对实现进行抽象的，数据结构可以是任何类型，只要其实现了某种接口，目标或数据源能够被编码就必须实现`io.Writer`或`io.Reader`接口.由于Go语言中到处都实现了Reader和Writer因此Encoder和Decoder可被应用的场景非常广泛