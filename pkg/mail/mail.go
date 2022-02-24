package mail

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(name, email, subject, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", "joylanguagesschool@gmail.com")
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	d := gomail.NewDialer("smtp.mailtrap.io", 465, os.Getenv("SMTPUSERNAME"), os.Getenv("SMTPPASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}
