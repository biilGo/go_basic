package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
)

func main() {

	// 通过调用sha1.New()创建一个新的hash.Hash对象,用来计算SHA1校验值.hash类型实际上是一个接口,它实现了io.Writer接口

	hasher := sha1.New()

	io.WriteString(hasher, "test")

	b := []byte{}

	fmt.Printf("Result:%x\n", hasher.Sum(b))
	fmt.Printf("Result:%d\n", hasher.Sum(b))

	hasher.Reset()

	data := []byte("We shall overcome!")

	n, err := hasher.Write(data)

	if n != len(data) || err != nil {
		log.Printf("Hash write error: %v / %v", n, err)
	}

	checksum := hasher.Sum(b)

	fmt.Printf("Result:%x\n", checksum)
}
