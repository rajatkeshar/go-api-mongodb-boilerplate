package mailer

import (
    "os"
	"log"
	"net/smtp"
)

func SendMail(to string, subject string, body string) {
	from := os.Getenv("EMAIL")
	pass := os.Getenv("PASSWORD")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject +"\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Println("smtp error: %s", err)
		return
	}

	log.Println("email sent to: ", to)
}
