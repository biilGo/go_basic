# 用户界面:web服务端
我们的web服务器由它启动,例如用如下命令在本地8080端口启动web服务器:`http.ListenAndServe(":8080", nil)`.
web服务器会在一个无限循环中监听到来的请求,但我们必须定义针对这些请求,服务器如何相应,可以用被称为HTTP处理器的`HandleFunc`函数来办到,例如代码:`http.HandleFunc("/add", Add)`

如此,每个以`/add`结尾的亲求都会调用`Add`函数.程序有两个HTTP处理器

- Redirect,用于对短URL重定向
- Add,用于处理新提交的URL

示意图:
![19.4示意图](./assets/19.4%E7%A4%BA%E6%84%8F%E5%9B%BE.jpg)

最简单的`main()`函数类似这样:
```
func main() {
    http.HandleFunc("/", Redirect)
    http.HandleFunc("/add", Add)
    http.ListenAndServe(":8080", nil)
}
```

对`/add`的请求由`Add`处理器处理,所有其他请求会被`Redirect`处理器处理,处理函数从到来的请求中获取信息,然后产生相应并写入`http.ResponseWriter`类型变量`w`

`Add`函数必须做的事有:
1. 读取长URL,即:用`r.FormValue("url")`从HTML表单中提交的HTTP请求中读取URL
2. 使用`stroe`上的Put方法存储长URL
3. 将对应的短URL发送给用户

每个需求都转化为一行代码:
```
func Add(w http.ResponseWriter, r *http.Request) {
    url := r.FormValue("url")
    key := store.Put(url)
    fmt.Fprintf(w, "http://localhost:8080/%s", key)
}
```

这里的`fmt`包的`Fprintf`函数用来替换字符串中的关键字`%s`,然后将结果作为相应发送回客户端,注意`Fprintf`把数据写到了`ResponseWrite`中,其实`Fpritnf`可以将数据写到任何实现了`io.Writer`的数据结构,即该结构实现了`Write`方法,Go中`io.Writer`称为接口,可见`Fprintf`利用接口变得十分通用,可以对很多不同的类型写入数据.Go中接口的使用十分普遍,它使代码更通用.

还需要一个表单,仍然可以用`Fprintf`来输出,这次将常量写入w,让我们来修改`Add`当未指定URL时显示HTML表单
```
func Add(w http.ResponseWriter, r *http.Request) {
    url := r.FormValue("url")
    if url == "" {
        fmt.Fprint(w, AddForm)
        return
    }
    key := store.Put(url)
    fmt.Fprintf(w, "http://localhost:8080/%s", key)
}
const AddForm = `
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
`
```

在那种情况下,发送字符串常量`AddForm`到客户端,它是html表单,包含一个url输入域和一个提交按钮,点击发送POST请求到`/add`这样`Add`处理函数被再次调用,此时url的值来自文本域.

`Redirect`函数在http请求路径中找到键,用get函数从store检索到对应的长URL,对用户发送http重定向,如果没找到URL,发送404 Not Found错误取而代之:
```
func Redirect(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Path[1:]
    url := store.Get(key)
    if url == "" {
        http.NotFound(w, r)
        return
    }
    http.Redirect(w, r, url, http.StatusFound)
}
```

## 编译和运行
- Linux和OSX平台,在终端窗口源码目录下启动make命令,或者LitelIDE中构建项目
- Windows平台,启动MINGW环境,步骤为:开始菜单,所有程序,MinGW,MinGW shell,在命令行窗口输入make并回车,源代码被编译并链接为原生exe可执行程序

生成内容为可执行程序,Linux/OSX下为goto,Windows下为goto.exe

要启动并运行web服务器,那么:
- linux和OSX平台,输入命令./goto
- Windows平台,从GoIDE启动程序

## 测试该程序
打开浏览器并请求url:http://localhost:8080/add
这会激活Add处理函数,请求还未包含url变量,所以相应会输出html表单询问输入:
![19.4相应输出html表单.jpg](./assets/19.4%E7%9B%B8%E5%BA%94%E8%BE%93%E5%87%BAhtml%E8%A1%A8%E5%8D%95.jpg)

条件一个长URL以获取等价的缩短版本,例如`http://golang.org/pkg/bufio/#Writer`,然后单击按钮,应用会为你产生一个短URL并打印出来,例如:`http:// localhost:8080/2`
![19.4获取等价缩短版本.jpg](./assets/19.4%E8%8E%B7%E5%8F%96%E7%AD%89%E4%BB%B7%E7%BC%A9%E7%9F%AD%E7%89%88%E6%9C%AC.jpg)

复制该URL并在浏览器地址栏粘贴以发出请求,现在轮到Redirect处理函数上场了,对应长URL的页面被显示了出来
![19.4Redirect函数.jpg](./assets/19.4Redirect%E5%87%BD%E6%95%B0.jpg)