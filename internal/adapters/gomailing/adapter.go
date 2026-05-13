package gomailing

import (
	"bytes"
	"log"
	"text/template"

	"gopkg.in/gomail.v2"
)

type GoMailer struct {
	host, email, password string
	port                  int
}

func NewGoMail(host, email, password string,
	port int) *GoMailer {
	return &GoMailer{host, email, password, port}
}

func (goMail *GoMailer) SendMail(r Request) error {
	log.Printf("Get Mail Request is: %+v", r)
	var m *gomail.Message
	if r.MailType == HTML {
		tmpl, err := template.ParseFiles(*r.TempFile)
		if err != nil {
			return err
		}
		var body bytes.Buffer
		tmpl.Execute(&body, r.Data)
		m = gomail.NewMessage()
		m.SetHeader("From", goMail.email)
		m.SetHeader("To", r.To...)
		m.SetHeader("Subject", r.Subject)
		m.SetBody("text/html", body.String())

	} else {
		m = gomail.NewMessage()
		m.SetHeader("From", goMail.email)
		m.SetHeader("To", r.To...)
		m.SetHeader("Subject", r.Subject)
		m.SetBody("text/html", *r.Body)
	}
	d := gomail.NewDialer(goMail.host, goMail.port, goMail.email, goMail.password)
	return d.DialAndSend(m)
}
