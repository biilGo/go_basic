package main

import (
	// 方便读写文件
	"io/ioutil"
	"log"
	"net/http"

	// 用于验证输入标题
	"regexp"

	// 动态创建html文档
	"text/template"
)

const lenPath = len("/view/")

// 避免黑客构造特殊输入攻击服务器,用如下正则表达式检查用户在浏览器上输入的URL
// 同时也是wiki页面标题,makeHandler会用它对请求管控
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

// 在程序运行时仅做一次解析,在init()函数中处理可以方便的达到目的,所有模板对象都被保持在内存中,存放在以html文件名为索引的map中
// 这被称为模板缓存,推荐的最佳实践
var templates = make(map[string]*template.Template)

// 为真正从模板和结构体构建出页面，必须使用：
// templates[tmpl].Execute(w, p)

var err error

// 需要一种机制把Page结构体数据插入网页的标题和内容中,可以利用template包通过如下步骤完成

/* <h1>{{.Title |html}}</h1>
<p>[<a href="/edit/{{.Title |html}}">edit</a>]</p>
<div>{{printf "%s" .Body |html}}</div> */

type Page struct {
	Title string
	Body  []byte
}

// template.Must(template.ParseFiles(tmpl + ".html"))把模板文件转换为*template.Template类型对象
func init() {
	for _, tmpl := range []string{"edit", "view"} {
		templates[tmpl] = template.Must(template.ParseFiles(tmpl + ".html"))
	}
}

// 高阶函数,其参数是一个函数,返回一个新的闭包函数
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[lenPath:]
		// 先验证输入标题title的有效性,如果标题包含了字母和数字以外的符号,就触发NotFound错误
		if !titleValidator.MatchString(title) {
			http.NotFound(w, r)
			return
		}
		// 闭包封闭了函数变量fn来构造其返回值
		fn(w, r, title)
	}
}

// 尝试按标题读取文本文件,通过调用同load()函数完成
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := load(title)
	if err != nil { // page not found
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// 尝试读取文件,如果存在则用editHanler模板来渲染
// 万一发生错误,创建一个新的包含指定标题的Page对象并渲染
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := load(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// 在编辑页面点击保存按钮时,触发保存页面内容的动作,按钮须放在html表单中,开头如下
// <form action="/save/{{.Title |html}}" method="POST">
// 这意味着，当提交表单到类似 http://localhost/save/{Title} 这样的 URL 格式时，
// 一个 POST 请求被发往网页服务器。针对这样的 URL 我们已经定义好了处理函数：saveHandler()。
// 在 request 对象上调用 FormValue() 方法，可以提取名称为 body 的文本域内容，用这些信息构造一个 Page 对象，然后尝试通过调用 save() 方法保存其内容。
// 万一运行失败，执行 http.Error 以将错误显示到浏览器。如果保存成功，重定向浏览器到该页的阅读页面。
// save() 函数非常简单，利用 ioutil.WriteFile()，写入 Page 结构体的 Body 字段到文件 filename 中，之后会被用于模板替换占位符 {{printf "%s" .Body |html}}。
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// viewHandler另外该值和表示没有error的nil值一起返回给调用者,然后再此方法中将该结构体与模板对象整合
// 万一发生错误，wiki网页再磁盘上不存在，错误会返回给viewHanler,此时会自动重定向,跳转请求对应标题的编辑页面
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates[tmpl].Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func load(title string) (*Page, error) {
	filename := title + ".txt"
	// 构建文件名并用ioutil.ReadFile读取文件内容
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// 如果文件存在,其内容会存入字符串中
	return &Page{Title: title, Body: body}, nil
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err.Error())
	}
}
