# 探索template包
前面我们使用了`template`对象把数据结构整合到html模板中,这项技术确实对网页应用程序非常有用,然而模板时一项更为通用的技术方案,数据驱动的模板被创建出来,以生成文本输出,html仅是其中一种特定的使用案例

模板通过与数据结构的整合来生成,通常为结构体或切片,当数据项传输给`tmpl.Execute()`它用其中元素进行替换,动态地重写某一小段文本.只有被导出的数据项才可以被整合进模板中.可以在`{{`和`}}`中加入数据求值或控制结构.数据项可以是值或指针,接口隐藏了它们的差异

## 字段替换:`{{.FileName}}`
要在模板中包含某个字段的内容,使用双花括号起以点`.`开头的字段名,假设`Name`是某个结构体的字段,其值要在被模板整合时替换,则在模板中使用文本`{{.Nmae}}`,当Name是map的键时这么做也是可行的,要创建一个新的Template对象,调用`template.New`其字符串参数可以指定模板的名称.

正如前面出现过的`Parse`方法通过解析模板定义字符串,生成模板的内部表示,当使用包含模板定义字符串的文件时,将文件路径传递给`ParseFiles`来解析,解析过程产生错误,这两个函数第二个返回值error != nil最后通过`Execute`方法,数据结构中的内容与模板整合,并将结果写入方法第一个参数中,其类型为`io.Writer`再一次地可能会有error返回,

看程序`template_field.go`输出通`os.Stdout`被写到控制台

## 验证模板格式
为了确保模板定义语法是正确的，使用`Must`函数处理`Parse`的返回结果;看程序`template_validation.go`

## `if-else`
运行Execute产生的结果来自模板的输出，它包含静态文本，以及被`{{}}`包裹的称之为管道的文本
```
t := template.New("template test")
t = template.Must(t.Parse("This is just static text. \n{{\"This is pipeline data - because it is evaluated within the double braces.\"}} {{`So is this, but within reverse quotes.`}}\n"))
t.Execute(os.Stdout, nil)
```

现在我们可以对管道数据的输出结果用`if-else-end`设置条件约束:如果管道是空的,类似于`{{if ``}} Will not print. {{end}}`

那么if条件的求值结果为false,不会有输出内容,但如果是这样:`{{if `anything`}} Print IF part. {{else}} Print ELSE part.{{end}}`
输出Print IF patr.看程序`template_ifelse.go`

## 点号和`with-end`
点号可以在go模板中使用:其值`{{.}}`被设置为当前管道的值
`with`语句将点号设为管道的值,如果管道是空的,那么不管`with-end`块之间有什么,都会被忽略.在被嵌套时,点号根据最近的作用域取得值

看程序`template_with_end.go`

## 模板变量`$`
可以在模板内为管道设置本地变量,变量名以`$`符号作为前缀,变量名只能包含字母,下划线,数字
看程序`template_variables.go`

## `range-end`
`range-end`结构格式`{{range pipeline}} T1 {{else}} T0 {{end}}`

`range`被用于在集合上迭代:管道的值必须时数组,切片或map如果管道的值长度为0,点号的值不受影响,且执行`T0`;否则点号被设置为数组,切片,map内元素的值并执行T1

如果模板为:
```
{{range .}}
{{.}}
{{end}}
```

那么执行代码
```
s := []int{1,2,3,4}
t.Execute(os.Stdout, s)
```

输出
```
1
2
3
4
```

如需更实用的示例
```
{{range .}}
    {{with .Author}}
        <p><b>{{html .}}</b> wrote:</p>
    {{else}}
        <p>An anonymous person wrote:</p>
    {{end}}
    <pre>{{html .Content}}</pre>
    <pre>{{html .Date}}</pre>
{{end}}
```

这里range在结构体切片上迭代,每次都包含author,content和date字段

## 模板预定义函数
也有一些可以在模板代码中使用的预定义函数,如`printf`函数工作方式类似于`fmt.Sprintf`