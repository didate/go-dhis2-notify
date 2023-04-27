package mail

import (
	"crypto/tls"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

func Send(mSubject string, mBody string) error {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", os.Getenv("MAIL_USERNAME"))

	// Set E-Mail receivers
	m.SetHeader("To", "didate224@gmail.com", "lyve.diallo@gmail.com")

	// Set E-Mail subject
	m.SetHeader("Subject", mSubject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", mBody)

	// Settings for SMTP server
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic(err)
	}
	d := gomail.NewDialer(os.Getenv("SMTP_SERVER"), port, os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	return d.DialAndSend(m)

}
