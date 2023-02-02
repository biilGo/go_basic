package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("Get : %v", err)
	}
}

func main() {
	// http.Get()获取并显示网页内容
	// Get返回值中的Body属性包含了网页内容
	res, err := http.Get("http://www.biying.com")

	checkError(err)

	// ioutil.ReadAll读取出来
	data, err := ioutil.ReadAll(res.Body)

	checkError(err)

	fmt.Printf("Got: %q", string(data))
}
