// smtp包实现了用于发送邮件的"简单邮件传输协议"
// 它有一个client类型,代表一个连接到SMTP服务器的客户端
// Dial:方法返回一个已连接到SMTP服务器的客户端Client
// 设置Mail(from发件人)和Rcpt(to接收人)
// Data方法返回一个用于写入数据的Writer,这里利用buf.WriteTo(wc)写入

package main

import (
	"bytes"
	"log"
	"net/smtp"
)

func main() {
	// Connect to the remote SMTP server

	client, err := smtp.Dial("smtp.sina.com:25")

	if err != nil {
		log.Fatal(err)
	}

	// set the sender adn recipient
	client.Mail("zrq7002@sina.com")
	client.Rcpt("zhao.ruqing@datablau.com.cn")

	// send the email body
	wc, err := client.Data()
	if err != nil {
		log.Fatal(err)
	}

	defer wc.Close()

	buf := bytes.NewBufferString("this is the email body")
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
}
