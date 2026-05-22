package gomailing

import (
	"context"
	"log"

	"github.com/eng-yousef-khaled/expenses_api/internal/domain/auth"
	"gopkg.in/gomail.v2"
)

type GoMailer interface {
	Send(ctx context.Context, r auth.EmailMessage) error
}

type goMailer struct {
	host, email, password string
	port                  int
}

func NewGoMail(host, email, password string,
	port int) GoMailer {
	return &goMailer{host, email, password, port}
}

func (g *goMailer) Send(ctx context.Context, r auth.EmailMessage) error {
	log.Printf("Get Mail Request is: %+v", r)
	var m *gomail.Message
	// if r.MailType == auth.HTML {
	// 	tmpl, err := template.ParseFiles(*r.TempFile)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	var body bytes.Buffer
	// 	tmpl.Execute(&body, r.Data)
	// 	m = gomail.NewMessage()
	// 	m.SetHeader("From", g.email)
	// 	m.SetHeader("To", r.To...)
	// 	m.SetHeader("Subject", r.Subject)
	// 	m.SetBody("text/html", body.String())

	// } else {
	m = gomail.NewMessage()
	m.SetHeader("From", g.email)
	m.SetHeader("To", r.To...)
	m.SetHeader("Subject", r.Subject)
	m.SetBody("text/html", *r.Body)
	// }
	d := gomail.NewDialer(g.host, g.port, g.email, g.password)
	return d.DialAndSend(m)
}
