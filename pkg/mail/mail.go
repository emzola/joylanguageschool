package mail

import "gopkg.in/gomail.v2"

func SendEmail(name, email, subject, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", "")
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	d := gomail.NewDialer("smtp.mailtrap.io", 587, "", "")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}