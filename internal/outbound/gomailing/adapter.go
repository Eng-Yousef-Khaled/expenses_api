package gomailing

import (
	"context"
	"log/slog"

	"github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
	"gopkg.in/gomail.v2"
)

type GoMailer interface {
	SendVerification(ctx context.Context, mail auth.EmailAddress, message string) error
}

type goMailer struct {
	host, email, password string
	port                  int
}

func NewGoMail(host, email, password string,
	port int) GoMailer {
	return &goMailer{host, email, password, port}
}

func (g *goMailer) SendVerification(ctx context.Context, mail auth.EmailAddress, message string) error {
	slog.Info("Get Mail Request mail is: %v , message is %v", string(mail), message)
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
	m.SetHeader("To", string(mail))
	m.SetHeader("Subject", "")
	m.SetBody("text/html", message)
	// }
	d := gomail.NewDialer(g.host, g.port, g.email, g.password)
	return d.DialAndSend(m)
}
