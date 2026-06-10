package gomailing

import (
	"bytes"
	"context"
	"log/slog"
	"path/filepath"
	"text/template"

	"github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
	"gopkg.in/gomail.v2"
)

type GoMailer interface {
	SendVerification(ctx context.Context, mail auth.EmailAddress, subject string, content string, code string, name string) error
}

type goMailer struct {
	host, email, password string
	port                  int
}

func NewGoMail(host, email, password string,
	port int) GoMailer {
	return &goMailer{host, email, password, port}
}

type VerificationTemplateData struct {
	Name    string
	Code    string
	Content string
}

func (g *goMailer) SendVerification(ctx context.Context, mail auth.EmailAddress, subject string, content string, code string, name string) error {
	slog.Info("Get Mail Request mail is: %v , message is %v", string(mail), content)
	templatePath := filepath.Join("templates", "verification.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		slog.Error("Failed to parse email template file", "path", templatePath, "error", err)
		return err
	}
	data := VerificationTemplateData{
		Name:    name,
		Code:    code,
		Content: content,
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		return err
	}

	// 5. Construct and send email
	m := gomail.NewMessage()
	m.SetHeader("From", g.email)
	m.SetHeader("To", string(mail))
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(g.host, g.port, g.email, g.password)
	return d.DialAndSend(m)
}
