# 网络和网页应用

## 模板
制作,解析并使模板生效:
`var strTempl = template.Must(template.New("TName").Parse(strTemplateHTML))`

在网页应用中使用HTML过滤器过滤HTML特殊字符:
`{{html .}} 或者通过一个字段 FieldName {{ .FieldName |html }}`

使用缓存模板