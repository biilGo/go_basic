package main

import (
	"log"
	"net/smtp"
)

func main() {
	// set up authentication information
	auth := smtp.PlainAuth(
		"",
		"user@example.com",
		"password",
		"mail.example.com",
	)

	// connect to the server,authenticate,set the sender and recipient
	// and send the email all in one step

	err := smtp.SendMail(
		"mail.example.com:25",
		auth,
		"sender@example.org",
		[]string{"recipient@example.net"},
		[]byte("this is the email body"),
	)

	if err != nil {
		log.Fatal(err)
	}
}
