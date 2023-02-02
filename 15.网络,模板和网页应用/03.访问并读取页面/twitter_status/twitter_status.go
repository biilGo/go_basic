// 通过xml包解析成为一个结构

package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// 这个结构会保存解析后的返回数据
// 他们会形成有层级的xml,可以忽略一些无用的数据
type Status struct {
	Text string
}

type User struct {
	XMLName xml.Name
	Status  Status
}

func main() {
	// 发起请求查询推特goodlnad用户的状态
	response, _ := http.Get("http://twitter.com/users/Googland.xml")

	// 初始化XML返回值的结构
	user := User{xml.Name{"", "user"}, Status{""}}

	// 将XML解析为我们的结构
	// 由于go版本的更新,Unmarshal函数第一个参数需要是[]byte类型,而无法传入Body
	xml.Unmarshal(response.Body, &user)
	fmt.Printf("status:%s", user.Status.Text)
}
