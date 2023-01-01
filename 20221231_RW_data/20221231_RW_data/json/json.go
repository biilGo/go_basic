// json.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Address struct {
	Type    string
	City    string
	Country string
}

type VCard struct {
	FirstName string
	LastName  string
	Address   []*Address
	Remark    string
}

func main() {
	pa := &Address{"private", "Aartselaar", "Belgium"}

	wa := &Address{"work", "Boom", "Belgium"}

	// vc编码后的数据为js，对其解码时，我们首先创建VCard用来保存解码的数据
	// var v VCard并调用json.Unmarshal(js,&v)解析[]byte中的json数据并将结果存入指针&v指向的值
	vc := VCard{"Jan", "Kersschot", []*Address{pa, wa}, "none"}

	// fmt.Printf("%v: \n", vc) // {Jan Kersschot [0x126d2b80 0x126d2be0] none}:

	// JSON format:

	// json.Marshal()的函数签名是
	// func Marshal(v insterface{}) ([]byte,error),数据编码后JSON文本实际上是一个[]byte

	// 处于安全考虑在web应用中最好使用json.MarshalforHTML()函数，其对数据执行HTML转码。所以文件可以被安全的嵌在HTML标签中

	js, _ := json.Marshal(vc)

	fmt.Printf("JSON format:%s", js)

	// using an encoder:
	file, error := os.OpenFile("D:/git_biilGo/go_basic/20221231_RW_data/format_JSON_Data/json/vcard.json", os.O_CREATE|os.O_WRONLY, 0666)
	if error != nil {
		fmt.Println("err")
	}

	defer file.Close()

	// json.NewEncoder()的函数签名是：func NewEncoder(w io.Writer) *Encoder
	// 返回的Encoder类型的指针可调用方法Encode(v interface{}),将数据对象v的json编码写入io.Write的w中

	enc := json.NewEncoder(file)
	err := enc.Encode(vc)

	if err != nil {
		log.Println("Error in encoding json")
	}
}
